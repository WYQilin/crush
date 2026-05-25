package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/charmbracelet/crush/internal/config"
	crushlog "github.com/charmbracelet/crush/internal/log"
	"github.com/charmbracelet/crush/internal/server"
	"github.com/charmbracelet/crush/internal/web"
	"github.com/charmbracelet/x/term"
	"github.com/spf13/cobra"
)

func init() {
	webCmd.Flags().String("addr", ":8080", "TCP address to listen on")
	webCmd.Flags().String("workspace-dir", "", "Workspace root (must contain crush.json); defaults to current working directory")
	webCmd.Flags().String("pages-dir", "", "Directory exposed under /preview/; defaults to <workspace-dir>/pages")
	webCmd.Flags().String("user", "crush", "HTTP Basic auth username")
	webCmd.Flags().String("password", "", "HTTP Basic auth password (env CRUSH_WEB_PASSWORD if empty)")
	webCmd.Flags().Bool("no-auth", false, "Disable HTTP Basic auth (NOT recommended for remote)")
	webCmd.Flags().Bool("require-permissions", false, "Keep per-tool permission prompts on (currently has no web UI; agent will hang)")
	rootCmd.AddCommand(webCmd)
}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start the Crush web frontend (single-user)",
	Long: `Run a browser-based Crush coding agent. The workspace root
(--workspace-dir, defaults to the current working directory) must contain a
crush.json with at least one configured provider/model. The pages directory
(--pages-dir) is exposed read-only under /preview/ so the agent can write
generated artifacts there and the UI can render them.`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		addr, _ := cmd.Flags().GetString("addr")
		workspaceDir, _ := cmd.Flags().GetString("workspace-dir")
		pagesDir, _ := cmd.Flags().GetString("pages-dir")
		user, _ := cmd.Flags().GetString("user")
		pass, _ := cmd.Flags().GetString("password")
		noAuth, _ := cmd.Flags().GetBool("no-auth")
		requirePerms, _ := cmd.Flags().GetBool("require-permissions")
		dataDir, _ := cmd.Flags().GetString("data-dir")
		debug, _ := cmd.Flags().GetBool("debug")

		if pass == "" {
			pass = os.Getenv("CRUSH_WEB_PASSWORD")
		}
		if !noAuth && pass == "" {
			return fmt.Errorf("password required (--password or $CRUSH_WEB_PASSWORD); use --no-auth to disable")
		}
		if noAuth {
			user, pass = "", ""
		}

		if workspaceDir == "" {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve working dir: %w", err)
			}
			workspaceDir = cwd
		}
		absWS, err := filepath.Abs(workspaceDir)
		if err != nil {
			return fmt.Errorf("resolve workspace dir: %w", err)
		}

		// Default pages-dir to <workspace-dir>/pages so generated
		// content always lands inside the project, not next to the
		// caller's cwd.
		if pagesDir == "" {
			pagesDir = filepath.Join(absWS, "pages")
		}
		absPages, err := filepath.Abs(pagesDir)
		if err != nil {
			return fmt.Errorf("resolve pages dir: %w", err)
		}
		if err := os.MkdirAll(absPages, 0o755); err != nil {
			return fmt.Errorf("create pages dir: %w", err)
		}

		// Load the workspace config so the API server can discover
		// providers/models/agents from crush.json. The backend
		// re-loads it per workspace, but loading here surfaces
		// configuration errors before the server starts.
		cfg, err := config.Load(absWS, dataDir, debug)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		if !cfg.Config().IsConfigured() {
			return fmt.Errorf("no providers configured; place a crush.json with providers/models in %s", absWS)
		}

		logFile := filepath.Join(config.GlobalCacheDir(), "web", "crush.log")
		if term.IsTerminal(os.Stderr.Fd()) {
			crushlog.Setup(logFile, debug, os.Stderr)
		} else {
			crushlog.Setup(logFile, debug)
		}

		api := server.NewServer(cfg, "tcp", addr)
		api.SetLogger(slog.Default())

		webSrv, err := web.New(api, web.Options{
			Addr:               addr,
			WorkspaceDir:       absWS,
			PagesDir:           absPages,
			Username:           user,
			Password:           pass,
			RequirePermissions: requirePerms,
		})
		if err != nil {
			return err
		}

		errch := make(chan error, 1)
		sigch := make(chan os.Signal, 1)
		sigs := []os.Signal{os.Interrupt}
		sigs = append(sigs, addSignals(sigs)...)
		signal.Notify(sigch, sigs...)

		go func() { errch <- webSrv.ListenAndServe(cmd.Context()) }()

		select {
		case <-sigch:
			slog.Info("Received interrupt, shutting down...")
		case err := <-errch:
			if err != nil && !errors.Is(err, server.ErrServerClosed) {
				return fmt.Errorf("web server error: %w", err)
			}
			return nil
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return webSrv.Shutdown(ctx)
	},
}
