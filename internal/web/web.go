// Package web wraps the Crush HTTP API with a browser-friendly
// frontend: an embedded chat SPA, a per-page preview server backed by
// the agent's working directory, and a styled cookie-based login
// suitable for single-user remote deployment.
package web

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/subtle"
	"embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/crush/internal/proto"
	"github.com/charmbracelet/crush/internal/server"
	"github.com/charmbracelet/crush/internal/version"
)

//go:embed static
var staticFS embed.FS

const sessionCookie = "crush_session"

// Options configure the web frontend.
type Options struct {
	// Addr is the TCP address to listen on, e.g. ":8080".
	Addr string

	// WorkspaceDir is the directory used as the agent's working
	// directory and as the source of crush.json (providers, models,
	// agents). Must be configured before the server starts.
	WorkspaceDir string

	// PagesDir is the directory exposed read-only under /preview/.
	// Each generated page lives at PagesDir/<name>/ and is served at
	// /preview/<name>/. PagesDir may live anywhere; it does not need
	// to be inside WorkspaceDir, but most users will want it to be a
	// subdirectory so the agent can write to it directly.
	PagesDir string

	// Username and Password gate every route via a styled login form
	// backed by an HttpOnly session cookie. When both are empty,
	// auth is disabled (not recommended for remote deployment).
	Username string
	Password string

	// RequirePermissions, when true, leaves Crush's per-tool
	// permission prompts on. The current web UI has no UI for
	// granting permissions, so the default (false) opts the
	// workspace into YOLO mode where every tool call is implicitly
	// approved. Single-user remote deployments rely on the login
	// above for access control.
	RequirePermissions bool
}

// Server is the web frontend wrapping a [*server.Server].
type Server struct {
	opts        Options
	api         *server.Server
	apiHandler  http.Handler
	workspaceID string
	handler     http.Handler

	mu       sync.RWMutex
	sessions map[string]time.Time // token -> issued-at.
}

// New builds a Server. The underlying [*server.Server] is constructed
// from the given API server (which the caller is responsible for
// configuring with TCP host + crush config).
func New(api *server.Server, opts Options) (*Server, error) {
	if opts.WorkspaceDir == "" {
		return nil, fmt.Errorf("workspace dir is required")
	}
	if opts.PagesDir == "" {
		return nil, fmt.Errorf("pages dir is required")
	}
	absWS, err := filepath.Abs(opts.WorkspaceDir)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve workspace dir: %w", err)
	}
	if st, err := os.Stat(absWS); err != nil || !st.IsDir() {
		return nil, fmt.Errorf("workspace dir %s is not a directory", absWS)
	}
	opts.WorkspaceDir = absWS

	abs, err := filepath.Abs(opts.PagesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve pages dir: %w", err)
	}
	if err := os.MkdirAll(abs, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create pages dir: %w", err)
	}
	opts.PagesDir = abs

	s := &Server{
		opts:       opts,
		api:        api,
		apiHandler: api.Handler(),
		sessions:   make(map[string]time.Time),
	}
	s.handler = s.buildHandler()
	return s, nil
}

// ListenAndServe bootstraps the agent workspace pointing at the pages
// directory, then serves HTTP on opts.Addr.
func (s *Server) ListenAndServe(ctx context.Context) error {
	ln, err := net.Listen("tcp", s.opts.Addr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	if err := s.bootstrapWorkspace(); err != nil {
		_ = ln.Close()
		return err
	}
	slog.Info("Crush web frontend listening",
		"addr", s.opts.Addr,
		"workspace_dir", s.opts.WorkspaceDir,
		"pages_dir", s.opts.PagesDir,
		"workspace_id", s.workspaceID,
	)
	return s.api.ServeHandler(ln, s.handler)
}

// Shutdown forwards to the underlying API server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.api.Shutdown(ctx)
}

// bootstrapWorkspace creates a single agent workspace rooted at the
// configured workspace directory (which must contain crush.json). The
// web frontend treats this workspace as the implicit one for all chat
// traffic.
func (s *Server) bootstrapWorkspace() error {
	_, ws, err := s.api.Backend().CreateWorkspace(proto.Workspace{
		Path:    s.opts.WorkspaceDir,
		Version: version.Version,
		Env:     os.Environ(),
		YOLO:    !s.opts.RequirePermissions,
	})
	if err != nil {
		return fmt.Errorf("create workspace: %w", err)
	}
	s.workspaceID = ws.ID
	return nil
}

func (s *Server) buildHandler() http.Handler {
	uiFS, err := fs.Sub(staticFS, "static")
	if err != nil {
		panic(err)
	}
	uiHandler := http.FileServer(http.FS(uiFS))
	previewHandler := http.StripPrefix("/preview/",
		http.FileServer(http.Dir(s.opts.PagesDir)))

	root := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/login":
			s.handleLogin(w, r)
			return
		case r.URL.Path == "/logout":
			s.handleLogout(w, r)
			return
		}
		if !s.authorized(r) {
			s.unauthorized(w, r)
			return
		}
		switch {
		case strings.HasPrefix(r.URL.Path, "/v1/"):
			s.interceptAgentPost(r)
			s.apiHandler.ServeHTTP(w, r)
		case strings.HasPrefix(r.URL.Path, "/preview/"):
			previewHandler.ServeHTTP(w, r)
		case r.URL.Path == "/api/info":
			s.handleInfo(w, r)
		case r.URL.Path == "/api/pages":
			s.handlePages(w, r)
		default:
			uiHandler.ServeHTTP(w, r)
		}
	})

	return root
}

func (s *Server) handleInfo(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"workspace_id":  s.workspaceID,
		"workspace_dir": s.opts.WorkspaceDir,
		"pages_dir":     s.opts.PagesDir,
		"version":       version.Version,
		"auth_enabled":  s.authEnabled(),
	})
}

// handlePages lists the immediate subdirectories of PagesDir as page
// slugs. A page is a directory containing at least an index.html, but
// we tolerate empty/in-progress dirs to keep the UI responsive while
// the agent is mid-write.
func (s *Server) handlePages(w http.ResponseWriter, _ *http.Request) {
	entries, err := os.ReadDir(s.opts.PagesDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type pageEntry struct {
		name    string
		modTime time.Time
	}
	pages := make([]pageEntry, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		mt := time.Time{}
		if info, err := e.Info(); err == nil {
			mt = info.ModTime()
		}
		// Prefer the directory's index.html mtime when present, since
		// the directory itself may not be touched on file rewrites.
		if fi, err := os.Stat(filepath.Join(s.opts.PagesDir, name, "index.html")); err == nil {
			if t := fi.ModTime(); t.After(mt) {
				mt = t
			}
		}
		pages = append(pages, pageEntry{name: name, modTime: mt})
	}
	sort.SliceStable(pages, func(i, j int) bool {
		return pages[i].modTime.After(pages[j].modTime)
	})
	out := make([]string, 0, len(pages))
	for _, p := range pages {
		out = append(out, p.name)
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}

func (s *Server) authEnabled() bool {
	return s.opts.Username != "" || s.opts.Password != ""
}

// agentPostPathRe matches POST /v1/workspaces/{id}/agent (and not its
// sub-paths like /agent/init or /agent/sessions/...).
var agentPostSuffix = "/agent"

// Markers wrap the injected pages preamble so the SPA can strip it
// from rendered user messages (see static/index.html). Keep both
// constants in sync with the frontend.
const (
	pagesPreambleStart = "<!--__crush_web_preamble__-->"
	pagesPreambleEnd   = "<!--/__crush_web_preamble__-->\n\n"
)

// interceptAgentPost rewrites the body of POST /v1/workspaces/{id}/agent
// to prepend a workspace-specific guidance preamble that pins generated
// pages under the configured pages directory. Without it the agent
// happily writes index.html into the workspace root, escaping the
// preview server.
func (s *Server) interceptAgentPost(r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	if !strings.HasPrefix(r.URL.Path, "/v1/workspaces/") {
		return
	}
	if !strings.HasSuffix(r.URL.Path, agentPostSuffix) {
		return
	}
	if r.Body == nil {
		return
	}
	raw, err := io.ReadAll(r.Body)
	_ = r.Body.Close()
	if err != nil {
		return
	}
	var msg map[string]any
	if err := json.Unmarshal(raw, &msg); err != nil {
		// Restore body so the inner handler can produce its own
		// 400 error.
		r.Body = io.NopCloser(bytes.NewReader(raw))
		return
	}
	if p, ok := msg["prompt"].(string); ok && p != "" {
		msg["prompt"] = pagesPreambleStart + s.pagesPreamble() +
			pagesPreambleEnd + p
	}
	out, err := json.Marshal(msg)
	if err != nil {
		r.Body = io.NopCloser(bytes.NewReader(raw))
		return
	}
	r.Body = io.NopCloser(bytes.NewReader(out))
	r.ContentLength = int64(len(out))
	r.Header.Set("Content-Length", fmt.Sprintf("%d", len(out)))
}

// pagesPreamble returns the workspace convention text injected ahead
// of every user prompt for the web frontend.
func (s *Server) pagesPreamble() string {
	rel, err := filepath.Rel(s.opts.WorkspaceDir, s.opts.PagesDir)
	if err != nil || strings.HasPrefix(rel, "..") {
		rel = s.opts.PagesDir
	}
	return fmt.Sprintf(
		`[系统约束 - 由 Crush Web 注入，必须遵守]
你正在 Crush Web 工作台中工作。所有“生成/修改网页”类任务的产物**必须**位于 `+
			"`%s/<slug>/`"+` 目录下：
- 路径根目录：`+"`%s`"+`
- 每个页面是一个独立子目录，目录名用英文小写短横线 slug（如 sakura-landing）
- 入口文件必须是该子目录下的 `+"`index.html`"+`
- 资源（CSS/JS/图片/manifest 等）一律放在该子目录内，使用相对路径引用
- 严禁把 index.html 或任何页面资源写到工作目录根、或写到 `+"`%s`"+` 之外的任意位置
- 如果用户没说 slug，请基于主题自取一个简洁英文 slug
- 用户说“继续/接着写/再加一段”等，请在已有的 slug 子目录上修改而非新建
- 不涉及网页生成的请求（如纯问答/解释代码）忽略本约束即可
完成后用一句话告知用户预览路径 /preview/<slug>/。`,
		rel, s.opts.PagesDir, rel,
	)
}

// authorized reports whether the request carries a valid session
// cookie (or whether auth is disabled altogether).
func (s *Server) authorized(r *http.Request) bool {
	if !s.authEnabled() {
		return true
	}
	c, err := r.Cookie(sessionCookie)
	if err != nil || c.Value == "" {
		return false
	}
	s.mu.RLock()
	_, ok := s.sessions[c.Value]
	s.mu.RUnlock()
	return ok
}

// unauthorized redirects browser navigations to /login and returns 401
// JSON for API/SSE/XHR clients so the SPA can react.
func (s *Server) unauthorized(w http.ResponseWriter, r *http.Request) {
	if isAPIRequest(r) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"unauthorized"}`))
		return
	}
	target := "/login"
	if r.URL.Path != "/" && r.URL.Path != "" {
		target += "?next=" + r.URL.Path
	}
	http.Redirect(w, r, target, http.StatusSeeOther)
}

func isAPIRequest(r *http.Request) bool {
	if strings.HasPrefix(r.URL.Path, "/v1/") ||
		strings.HasPrefix(r.URL.Path, "/api/") {
		return true
	}
	accept := r.Header.Get("Accept")
	return strings.Contains(accept, "application/json") ||
		strings.Contains(accept, "text/event-stream") ||
		r.Header.Get("X-Requested-With") == "XMLHttpRequest"
}

// handleLogin serves the login page on GET and validates credentials
// on POST. On success it sets an HttpOnly session cookie and returns
// the post-login redirect target as JSON.
func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// If already authorized, bounce home.
		if s.authorized(r) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		serveStaticFile(w, r, "login.html")
	case http.MethodPost:
		if !s.authEnabled() {
			writeJSON(w, http.StatusOK, map[string]string{"redirect": "/"})
			return
		}
		var body struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Next     string `json:"next"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest,
				map[string]string{"error": "invalid request"})
			return
		}
		uOK := subtle.ConstantTimeCompare(
			[]byte(body.Username), []byte(s.opts.Username)) == 1
		pOK := subtle.ConstantTimeCompare(
			[]byte(body.Password), []byte(s.opts.Password)) == 1
		if !uOK || !pOK {
			writeJSON(w, http.StatusUnauthorized,
				map[string]string{"error": "用户名或密码错误"})
			return
		}
		token, err := newToken()
		if err != nil {
			writeJSON(w, http.StatusInternalServerError,
				map[string]string{"error": "internal error"})
			return
		}
		s.mu.Lock()
		s.sessions[token] = time.Now()
		s.mu.Unlock()
		http.SetCookie(w, &http.Cookie{
			Name:     sessionCookie,
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Secure:   r.TLS != nil,
			MaxAge:   60 * 60 * 24 * 30,
		})
		next := body.Next
		if next == "" || !strings.HasPrefix(next, "/") ||
			strings.HasPrefix(next, "//") {
			next = "/"
		}
		writeJSON(w, http.StatusOK, map[string]string{"redirect": next})
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(sessionCookie); err == nil {
		s.mu.Lock()
		delete(s.sessions, c.Value)
		s.mu.Unlock()
	}
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
	if r.Method == http.MethodPost && isAPIRequest(r) {
		writeJSON(w, http.StatusOK, map[string]string{"redirect": "/login"})
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func newToken() (string, error) {
	var b [32]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func serveStaticFile(w http.ResponseWriter, r *http.Request, name string) {
	data, err := staticFS.ReadFile("static/" + name)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	_, _ = w.Write(data)
}
