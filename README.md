# Crush

<p align="center">
    <a href="https://stuff.charm.sh/crush/charm-crush.png"><img width="450" alt="Charm Crush Logo" src="https://github.com/user-attachments/assets/cf8ca3ce-8b02-43f0-9d0f-5a331488da4b" /></a><br />
    <a href="https://github.com/charmbracelet/crush/releases"><img src="https://img.shields.io/github/release/charmbracelet/crush" alt="Latest Release"></a>
    <a href="https://github.com/charmbracelet/crush/actions"><img src="https://github.com/charmbracelet/crush/actions/workflows/build.yml/badge.svg" alt="Build Status"></a>
</p>

<p align="center">Your new coding bestie, now available in your favourite terminal.<br />Your tools, your code, and your workflows, wired into your LLM of choice.</p>
<p align="center">终端里的编程新搭档，<br />无缝接入你的工具、代码与工作流，全面兼容主流 LLM 模型。</p>

<p align="center"><img width="800" alt="Crush Demo" src="https://github.com/user-attachments/assets/58280caf-851b-470a-b6f7-d5c4ea8a1968" /></p>

## Features

- **Multi-Model:** choose from a wide range of LLMs or add your own via OpenAI- or Anthropic-compatible APIs
- **Flexible:** switch LLMs mid-session while preserving context
- **Session-Based:** maintain multiple work sessions and contexts per project
- **LSP-Enhanced:** Crush uses LSPs for additional context, just like you do
- **Extensible:** add capabilities via MCPs (`http`, `stdio`, and `sse`)
- **Works Everywhere:** first-class support in every terminal on macOS, Linux, Windows (PowerShell and WSL), Android, FreeBSD, OpenBSD, and NetBSD
- **Industrial Grade:** built on the Charm ecosystem, powering 25k+ applications, from leading open source projects to business-critical infrastructure

## Installation

Use a package manager:

```bash
# Homebrew
brew install charmbracelet/tap/crush

# NPM
npm install -g @charmland/crush

# Arch Linux (btw)
yay -S crush-bin

# Nix
nix run github:numtide/nix-ai-tools#crush

# FreeBSD
pkg install crush
```

Windows users:

```bash
# Winget
winget install charmbracelet.crush

# Scoop
scoop bucket add charm https://github.com/charmbracelet/scoop-bucket.git
scoop install crush
```

<details>
<summary><strong>Nix (NUR)</strong></summary>

Crush is available via the official Charm [NUR](https://github.com/nix-community/NUR) in `nur.repos.charmbracelet.crush`, which is the most up-to-date way to get Crush in Nix.

You can also try out Crush via the NUR with `nix-shell`:

```bash
# Add the NUR channel.
nix-channel --add https://github.com/nix-community/NUR/archive/main.tar.gz nur
nix-channel --update

# Get Crush in a Nix shell.
nix-shell -p '(import <nur> { pkgs = import <nixpkgs> {}; }).repos.charmbracelet.crush'
```

### NixOS & Home Manager Module Usage via NUR

Crush provides NixOS and Home Manager modules via NUR.
You can use these modules directly in your flake by importing them from NUR. Since it auto detects whether its a home manager or nixos context you can use the import the exact same way :)

```nix
{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    nur.url = "github:nix-community/NUR";
  };

  outputs = { self, nixpkgs, nur, ... }: {
    nixosConfigurations.your-hostname = nixpkgs.lib.nixosSystem {
      system = "x86_64-linux";
      modules = [
        nur.modules.nixos.default
        nur.repos.charmbracelet.modules.crush
        {
          programs.crush = {
            enable = true;
            settings = {
              providers = {
                openai = {
                  id = "openai";
                  name = "OpenAI";
                  base_url = "https://api.openai.com/v1";
                  type = "openai";
                  api_key = "sk-fake123456789abcdef...";
                  models = [
                    {
                      id = "gpt-4";
                      name = "GPT-4";
                    }
                  ];
                };
              };
              lsp = {
                go = { command = "gopls"; enabled = true; };
                nix = { command = "nil"; enabled = true; };
              };
              options = {
                context_paths = [ "/etc/nixos/configuration.nix" ];
                tui = { compact_mode = true; };
                debug = false;
              };
            };
          };
        }
      ];
    };
  };
}
```

</details>

<details>
<summary><strong>Debian/Ubuntu</strong></summary>

```bash
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://repo.charm.sh/apt/gpg.key | sudo gpg --dearmor -o /etc/apt/keyrings/charm.gpg
echo "deb [signed-by=/etc/apt/keyrings/charm.gpg] https://repo.charm.sh/apt/ * *" | sudo tee /etc/apt/sources.list.d/charm.list
sudo apt update && sudo apt install crush
```

</details>

<details>
<summary><strong>Fedora/RHEL</strong></summary>

```bash
echo '[charm]
name=Charm
baseurl=https://repo.charm.sh/yum/
enabled=1
gpgcheck=1
gpgkey=https://repo.charm.sh/yum/gpg.key' | sudo tee /etc/yum.repos.d/charm.repo
sudo yum install crush
```

</details>

Or, download it:

- [Packages][releases] are available in Debian and RPM formats
- [Binaries][releases] are available for Linux, macOS, Windows, FreeBSD, OpenBSD, and NetBSD

[releases]: https://github.com/charmbracelet/crush/releases

Or just install it with Go:

```
go install github.com/charmbracelet/crush@latest
```

> [!WARNING]
> Productivity may increase when using Crush and you may find yourself nerd
> sniped when first using the application. If the symptoms persist, join the
> [Slack][slack] or [Discord][discord] and nerd snipe the rest of us.

## Getting Started

The quickest way to get started is to grab an API key for your preferred
provider such as Anthropic, OpenAI, Groq, OpenRouter, or Vercel AI Gateway and just start
Crush. You'll be prompted to enter your API key.

That said, you can also set environment variables for preferred providers.

| Environment Variable        | Provider                                           |
| --------------------------- | -------------------------------------------------- |
| `HYPER_API_KEY`             | Charm Hyper                                        |
| `ANTHROPIC_API_KEY`         | Anthropic                                          |
| `OPENAI_API_KEY`            | OpenAI                                             |
| `VERCEL_API_KEY`            | Vercel AI Gateway                                  |
| `GEMINI_API_KEY`            | Google Gemini                                      |
| `SYNTHETIC_API_KEY`         | Synthetic                                          |
| `ZAI_API_KEY`               | Z.ai                                               |
| `MINIMAX_API_KEY`           | MiniMax                                            |
| `HF_TOKEN`                  | Hugging Face Inference                             |
| `CEREBRAS_API_KEY`          | Cerebras                                           |
| `OPENROUTER_API_KEY`        | OpenRouter                                         |
| `IONET_API_KEY`             | io.net                                             |
| `ALIBABA_SINGAPORE_API_KEY` | Alibaba (Singapore)                                |
| `GROQ_API_KEY`              | Groq                                               |
| `AVIAN_API_KEY`             | Avian                                              |
| `OPENCODE_API_KEY`          | OpenCode Zen & Go                                  |
| `VERTEXAI_PROJECT`          | Google Cloud VertexAI (Gemini)                     |
| `VERTEXAI_LOCATION`         | Google Cloud VertexAI (Gemini)                     |
| `AWS_ACCESS_KEY_ID`         | Amazon Bedrock (Claude)                            |
| `AWS_SECRET_ACCESS_KEY`     | Amazon Bedrock (Claude)                            |
| `AWS_REGION`                | Amazon Bedrock (Claude)                            |
| `AWS_PROFILE`               | Amazon Bedrock (Custom Profile)                    |
| `AWS_BEARER_TOKEN_BEDROCK`  | Amazon Bedrock                                     |
| `AZURE_OPENAI_API_ENDPOINT` | Azure OpenAI models                                |
| `AZURE_OPENAI_API_KEY`      | Azure OpenAI models (optional when using Entra ID) |
| `AZURE_OPENAI_API_VERSION`  | Azure OpenAI models                                |

### Subscriptions

If you prefer subscription-based usage, here are some plans that work well in
Crush:

- [Synthetic](https://synthetic.new/pricing)
- [GLM Coding Plan](https://z.ai/subscribe)
- [Kimi Code](https://www.kimi.com/membership/pricing)
- [MiniMax Coding Plan](https://platform.minimax.io/subscribe/coding-plan)

### By the Way

Is there a provider you’d like to see in Crush? Is there an existing model that needs an update?

Crush’s default model listing is managed in [Catwalk](https://github.com/charmbracelet/catwalk), a community-supported, open source repository of Crush-compatible models, and you’re welcome to contribute.

<a href="https://github.com/charmbracelet/catwalk"><img width="174" height="174" alt="Catwalk Badge" src="https://github.com/user-attachments/assets/95b49515-fe82-4409-b10d-5beb0873787d" /></a>

## Web UI (单用户远程部署)

除了终端 TUI，Crush 还提供一个轻量的浏览器前端，适合在远程 VPS 上自部署后通过浏览器访问。该模式专为**个人单用户**场景设计：聊天 + 简单单页面生成与预览，不面向多用户协作。

### 启动

```bash
# 设置 HTTP Basic Auth 密码（也可以用 --password 直接传）
export CRUSH_WEB_PASSWORD='your-strong-password'

# 在含 crush.json 的工作区根目录下启动；--pages-dir 仅用作预览根目录
crush web --addr :8080 --pages-dir ./pages --user crush
```

启动后浏览器打开 `http://<host>:8080/`，登录用户名/密码即可进入聊天界面。Web 模式会以当前 `--workspace-dir`（默认 = 启动时的工作目录）为 agent 工作区，从这里加载 `crush.json` 中的 provider/model/agent 配置；agent 写入的产物建议放到 `--pages-dir` 子目录，前端会通过 `/preview/<page>/` 展示。

### 命令行参数

- `--addr`：监听地址，默认 `:8080`
- `--workspace-dir`：工作区根目录，必须包含已配置 provider 的 `crush.json`；默认取启动时的当前目录
- `--pages-dir`：暴露在 `/preview/` 下的预览根目录，**默认 `<workspace-dir>/pages`**（不存在会自动创建），即项目目录下的 `pages/` 子目录
- `--user`：HTTP Basic Auth 用户名，默认 `crush`
- `--password`：HTTP Basic Auth 密码；为空时读取环境变量 `CRUSH_WEB_PASSWORD`
- `--no-auth`：关闭鉴权（**仅本地调试**，远程部署严禁开启）
- `--require-permissions`：保留 Crush 的工具调用权限确认。当前 Web 前端**没有**权限弹窗 UI，开启后 agent 会卡在第一次写文件/执行命令时；默认关闭（YOLO 模式），所有工具调用自动放行，由前面的 Basic Auth 把关
- `--debug`：开启调试日志（日志路径见 [Logging](#logging) 一节，文件位于 `web/crush.log`）

### 配置

Web 模式从 `--workspace-dir` 起按 Crush 标准规则查找 `crush.json`（见 [Configuration](#configuration)），未配置任何 provider 时会直接报错并拒绝启动。首次使用前请在工作区根目录放置一份带有 provider/model 的 `crush.json`（或在该目录内运行 `crush logs in` 等命令完成登录），无需为 Web 单独维护一份配置。

### 路由概览

- `/` 内嵌的聊天 SPA
- `/v1/*` 复用 Crush 的 REST API（sessions、messages、SSE 事件流等）
- `/preview/<page>/` 直接以静态资源方式预览 `pages-dir/<page>/` 下生成的页面
- `/api/info`、`/api/pages` 给前端使用的轻量辅助接口

所有路由都被 HTTP Basic Auth 包裹，浏览器登录一次后 iframe 预览也会自动带上鉴权。

### 前端行为

- **即时连接状态**：SSE 流在响应头返回时立即 flush，连接成功后右上角状态点立刻变绿；中间过 nginx/SLB 反代时也会带上 `X-Accel-Buffering: no` 防缓冲
- **预览自动跟随**：`/api/pages` 按 mtime 倒序返回；当 agent 生成或修改页面（`file` / `agent_event=response` 事件触发）时，预览面板会自动切到本次最新生成的页面，无需手动刷新
- **新会话重置预览**：点击「新会话」会清空消息列表与预览 iframe，下一轮生成完成后自动加载新页面
- **手机自适应**：`<meta viewport>` 强制 1280px 等比缩放，整页按桌面布局展示并按屏宽缩小；触屏设备额外加 `max-height` 限制，避免页面被拉得过长产生「样式错乱」错觉

### 远程部署

仓库根目录附带两份开箱即用的部署模板：

- [`nginx.conf`](./nginx.conf)：HTTPS 反代到 `127.0.0.1:8085`，按 location 分档调超时与缓冲（SSE/agent 长连接关 buffering 与 1h 读超时，静态资源短缓存，预览页面 `no-store`），含 upstream keepalive
- [`ecosystem.config.js`](./ecosystem.config.js)：PM2 进程守护（`fork` 单实例、内存阈值重启、独立日志、优雅关闭超时给 SSE 留时间）

最小化部署流程：

```bash
# 1) 构建二进制
CGO_ENABLED=0 GOEXPERIMENT=greenteagc go build -o crush .

# 2) PM2 守护
pm2 start ecosystem.config.js
pm2 save && pm2 startup    # 按提示执行 sudo 那一行

# 3) Nginx 反代（先准备好域名 A 记录与证书）
sudo cp nginx.conf /etc/nginx/conf.d/<your-domain>.conf
# 编辑 ssl_certificate / ssl_certificate_key 两行为实际证书路径
sudo nginx -t && sudo systemctl reload nginx
```

如需要 Let's Encrypt 证书：

```bash
sudo mkdir -p /var/www/letsencrypt
sudo certbot certonly --webroot -w /var/www/letsencrypt -d <your-domain>
# 证书路径：/etc/letsencrypt/live/<your-domain>/{fullchain,privkey}.pem
```

### 远程部署建议

- **务必**通过 HTTPS（如 Caddy/Nginx 反代）暴露端口，不要把 `:8080` 直接挂公网
- 使用强密码并避免在命令行明文传入，优先用 `CRUSH_WEB_PASSWORD`
- 不要使用 `--no-auth`
- `--workspace-dir` 会被 LLM 自由读写，建议放在专用目录而不是核心源码仓库；`--pages-dir` 通常作为该工作区下的子目录使用

## Configuration

> [!TIP]
> Crush ships with a builtin `crush-config` skill for configuring itself. In
> many cases you can simply ask Crush to configure itself.

Crush runs great with no configuration. That said, if you do need or want to
customize Crush, configuration can be added either local to the project itself,
or globally, with the following priority:

1. `.crush.json`
2. `crush.json`
3. `$HOME/.config/crush/crush.json`

Configuration itself is stored as a JSON object:

```json
{
  "this-setting": { "this": "that" },
  "that-setting": ["ceci", "cela"]
}
```

As an additional note, Crush also stores ephemeral data, such as application
state, in one additional location:

```bash
# Unix
$HOME/.local/share/crush/crush.json

# Windows
%LOCALAPPDATA%\crush\crush.json
```

> [!TIP]
> You can override the user and data config locations by setting:
>
> - `CRUSH_GLOBAL_CONFIG`
> - `CRUSH_GLOBAL_DATA`

### LSPs

Crush can use LSPs for additional context to help inform its decisions, just
like you would. LSPs can be added manually like so:

```json
{
  "$schema": "https://charm.land/crush.json",
  "lsp": {
    "go": {
      "command": "gopls",
      "env": {
        "GOTOOLCHAIN": "go1.24.5"
      }
    },
    "typescript": {
      "command": "typescript-language-server",
      "args": ["--stdio"]
    },
    "nix": {
      "command": "nil"
    }
  }
}
```

### MCPs

Crush also supports Model Context Protocol (MCP) servers through three transport
types: `stdio` for command-line servers, `http` for HTTP endpoints, and `sse`
for Server-Sent Events.

Shell-style value expansion (`$VAR`, `${VAR:-default}`, `$(command)`, quoting,
nesting) works in `command`, `args`, `env`, `headers`, and `url`, so
file-based secrets work out of the box. You can use values like `"$TOKEN"`
or `"$(cat /path/to/secret/token)"`. Expansion runs through Crush's embedded
shell, so the same syntax works on every supported system, Windows included.

Unset variables expand to the empty string by default, matching bash. For
required credentials, use `${VAR:?message}` so an unset variable fails loudly
at load time with `message` instead of silently resolving to empty:

```json
{ "api_key": "${CODEBERG_TOKEN:?set CODEBERG_TOKEN}" }
```

Headers (both MCP `headers` and provider `extra_headers`) whose value
resolves to the empty string are dropped from the outgoing request rather
than sent as `Header:`. That keeps optional env-gated headers like
`"OpenAI-Organization": "$OPENAI_ORG_ID"` clean when the variable is unset.

Provider `extra_body` is a non-expanding JSON passthrough; put env-driven
values in `extra_headers` or the provider's `api_key` / `base_url`, all of
which do expand.

> **Security note:** `crush.json` is trusted code. Any `$(...)` in it runs at
> load time with your shell's privileges, before the UI appears. Don't launch
> Crush in a directory whose `crush.json` you haven't reviewed.

```json
{
  "$schema": "https://charm.land/crush.json",
  "mcp": {
    "filesystem": {
      "type": "stdio",
      "command": "node",
      "args": ["/path/to/mcp-server.js"],
      "timeout": 120,
      "disabled": false,
      "disabled_tools": ["some-tool-name"],
      "env": {
        "NODE_ENV": "production"
      }
    },
    "github": {
      "type": "http",
      "url": "https://api.githubcopilot.com/mcp/",
      "timeout": 120,
      "disabled": false,
      "disabled_tools": ["create_issue", "create_pull_request"],
      "headers": {
        "Authorization": "Bearer $GH_PAT"
      }
    },
    "streaming-service": {
      "type": "sse",
      "url": "https://example.com/mcp/sse",
      "timeout": 120,
      "disabled": false,
      "headers": {
        "API-Key": "$(echo $API_KEY)"
      }
    }
  }
}
```

### Hooks

Crush has preliminary support for hooks. For details, see
[the hook guide](./docs/hooks/).

### Ignoring Files

Crush respects `.gitignore` files by default, but you can also create a
`.crushignore` file to specify additional files and directories that Crush
should ignore. This is useful for excluding files that you want in version
control but don't want Crush to consider when providing context.

The `.crushignore` file uses the same syntax as `.gitignore` and can be placed
in the root of your project or in subdirectories.

### Allowing Tools

By default, Crush will ask you for permission before running tool calls. If
you'd like, you can allow tools to be executed without prompting you for
permissions. Use this with care.

```json
{
  "$schema": "https://charm.land/crush.json",
  "permissions": {
    "allowed_tools": [
      "view",
      "ls",
      "grep",
      "edit",
      "mcp_context7_get-library-doc"
    ]
  }
}
```

You can also skip all permission prompts entirely by running Crush with the
`--yolo` flag. Be very, very careful with this feature.

### Disabling Built-In Tools

If you'd like to prevent Crush from using certain built-in tools entirely, you
can disable them via the `options.disabled_tools` list. Disabled tools are
completely hidden from the agent.

```json
{
  "$schema": "https://charm.land/crush.json",
  "options": {
    "disabled_tools": ["bash", "sourcegraph"]
  }
}
```

To disable tools from MCP servers, see the [MCP config section](#mcps).

### Disabling Skills

If you'd like to prevent Crush from using certain skills entirely, you can
disable them via the `options.disabled_skills` list. Disabled skills are hidden
from the agent, including builtin skills and skills discovered from disk.

```json
{
  "$schema": "https://charm.land/crush.json",
  "options": {
    "disabled_skills": ["crush-config"]
  }
}
```

### Agent Skills

Crush supports the [Agent Skills](https://agentskills.io) open standard for
extending agent capabilities with reusable skill packages. Skills are folders
containing a `SKILL.md` file with instructions that Crush can discover and
activate on demand.

The global paths we looks for skills are:

* `$CRUSH_SKILLS_DIR`
* `$XDG_CONFIG_HOME/agents/skills` or `~/.config/agents/skills/`
* `$XDG_CONFIG_HOME/crush/skills` or `~/.config/crush/skills/`
* `~/.agents/skills/`
* `~/.claude/skills/`
* On Windows, we _also_ look at
  * `%LOCALAPPDATA%\agents\skills\` or `%USERPROFILE%\AppData\Local\agents\skills\`
  * `%LOCALAPPDATA%\crush\skills\` or `%USERPROFILE%\AppData\Local\crush\skills\`
* Additional paths configured via `options.skills_paths`

On top of that, we _also_ load skills in your project from the following
relative paths:

* `.agents/skills`
* `.crush/skills`
* `.claude/skills`
* `.cursor/skills`

```jsonc
{
  "$schema": "https://charm.land/crush.json",
  "options": {
    "skills_paths": [
      "~/.config/crush/skills", // Windows: "%LOCALAPPDATA%\\crush\\skills",
      "./project-skills",
    ],
  },
}
```

You can get started with example skills from [anthropics/skills](https://github.com/anthropics/skills):

```bash
# Unix
mkdir -p ~/.config/crush/skills
cd ~/.config/crush/skills
git clone https://github.com/anthropics/skills.git _temp
mv _temp/skills/* . && rm -rf _temp
```

```powershell
# Windows (PowerShell)
mkdir -Force "$env:LOCALAPPDATA\crush\skills"
cd "$env:LOCALAPPDATA\crush\skills"
git clone https://github.com/anthropics/skills.git _temp
mv _temp/skills/* . ; rm -r -force _temp
```

#### User-Invocable Skills

Skills can be made invocable as commands from the commands palette (Ctrl+P). Add `user-invocable: true` to the skill's YAML frontmatter:

```yaml
---
name: my-skill
description: A skill that can be invoked as a command.
user-invocable: true
---
```

User-invocable skills appear in the commands palette with a `user:` or `project:` prefix:
- Skills from global directories show as `user:skill-name`
- Skills from project directories show as `project:skill-name`

When invoked, the skill's instructions are loaded into the conversation context.

To prevent the model from auto-triggering a skill (while still allowing user invocation), add `disable-model-invocation: true`:

```yaml
---
name: my-skill
description: Only invocable by users, not the model.
user-invocable: true
disable-model-invocation: true
---
```

Skills with `disable-model-invocation` won't appear in the model's available skills list but can still be invoked manually by users.

### Desktop notifications

Crush sends desktop notifications when a tool call requires permission and when
the agent finishes its turn. They're only sent when the terminal window isn't
focused _and_ your terminal supports reporting the focus state.

```jsonc
{
  "$schema": "https://charm.land/crush.json",
  "options": {
    "disable_notifications": false, // default
  },
}
```

To disable desktop notifications, set `disable_notifications` to `true` in your
configuration. On macOS, notifications currently lack icons due to platform
limitations.

### Initialization

When you initialize a project, Crush analyzes your codebase and creates
a context file that helps it work more effectively in future sessions.
By default, this file is named `AGENTS.md`, but you can customize the
name and location with the `initialize_as` option:

```json
{
  "$schema": "https://charm.land/crush.json",
  "options": {
    "initialize_as": "AGENTS.md"
  }
}
```

This is useful if you prefer a different naming convention or want to
place the file in a specific directory (e.g., `CRUSH.md` or
`docs/LLMs.md`). Crush will fill the file with project-specific context
like build commands, code patterns, and conventions it discovered during
initialization.

### Attribution Settings

By default, Crush adds attribution information to Git commits and pull requests
it creates. You can customize this behavior with the `attribution` option:

```json
{
  "$schema": "https://charm.land/crush.json",
  "options": {
    "attribution": {
      "trailer_style": "co-authored-by",
      "generated_with": true
    }
  }
}
```

- `trailer_style`: Controls the attribution trailer added to commit messages
  (default: `assisted-by`)
  - `assisted-by`: Adds `Assisted-by: Crush:[ModelID]` as specified in [the convention](https://docs.kernel.org/process/coding-assistants.html#attribution)
  - `co-authored-by`: Adds `Co-Authored-By: Crush <crush@charm.land>`
  - `none`: No attribution trailer
- `generated_with`: When true (default), adds `💘 Generated with Crush` line to
  commit messages and PR descriptions

### Custom Providers

Crush supports custom provider configurations for both OpenAI-compatible and
Anthropic-compatible APIs.

> [!NOTE]
> Note that we support two "types" for OpenAI. Make sure to choose the right one
> to ensure the best experience!
>
> - `openai` should be used when proxying or routing requests through OpenAI.
> - `openai-compat` should be used when using non-OpenAI providers that have OpenAI-compatible APIs.

#### OpenAI-Compatible APIs

Here’s an example configuration for Deepseek, which uses an OpenAI-compatible
API. Don't forget to set `DEEPSEEK_API_KEY` in your environment.

```json
{
  "$schema": "https://charm.land/crush.json",
  "providers": {
    "deepseek": {
      "type": "openai-compat",
      "base_url": "https://api.deepseek.com/v1",
      "api_key": "$DEEPSEEK_API_KEY",
      "models": [
        {
          "id": "deepseek-chat",
          "name": "Deepseek V3",
          "cost_per_1m_in": 0.27,
          "cost_per_1m_out": 1.1,
          "cost_per_1m_in_cached": 0.07,
          "cost_per_1m_out_cached": 1.1,
          "context_window": 64000,
          "default_max_tokens": 5000
        }
      ]
    }
  }
}
```

#### Anthropic-Compatible APIs

Custom Anthropic-compatible providers follow this format:

```json
{
  "$schema": "https://charm.land/crush.json",
  "providers": {
    "custom-anthropic": {
      "type": "anthropic",
      "base_url": "https://api.anthropic.com/v1",
      "api_key": "$ANTHROPIC_API_KEY",
      "extra_headers": {
        "anthropic-version": "2023-06-01"
      },
      "models": [
        {
          "id": "claude-sonnet-4-20250514",
          "name": "Claude Sonnet 4",
          "cost_per_1m_in": 3,
          "cost_per_1m_out": 15,
          "cost_per_1m_in_cached": 3.75,
          "cost_per_1m_out_cached": 0.3,
          "context_window": 200000,
          "default_max_tokens": 50000,
          "can_reason": true,
          "supports_attachments": true
        }
      ]
    }
  }
}
```

### Amazon Bedrock

Crush currently supports running Anthropic models through Bedrock, with caching disabled.

- A Bedrock provider will appear once you have AWS configured, i.e. `aws configure`
- Crush also expects the `AWS_REGION` or `AWS_DEFAULT_REGION` to be set
- To use a specific AWS profile set `AWS_PROFILE` in your environment, i.e. `AWS_PROFILE=myprofile crush`
- Alternatively to `aws configure`, you can also just set `AWS_BEARER_TOKEN_BEDROCK`

### Vertex AI Platform

Vertex AI will appear in the list of available providers when `VERTEXAI_PROJECT` and `VERTEXAI_LOCATION` are set. You will also need to be authenticated:

```bash
gcloud auth application-default login
```

To add specific models to the configuration, configure as such:

```json
{
  "$schema": "https://charm.land/crush.json",
  "providers": {
    "vertexai": {
      "models": [
        {
          "id": "claude-sonnet-4@20250514",
          "name": "VertexAI Sonnet 4",
          "cost_per_1m_in": 3,
          "cost_per_1m_out": 15,
          "cost_per_1m_in_cached": 3.75,
          "cost_per_1m_out_cached": 0.3,
          "context_window": 200000,
          "default_max_tokens": 50000,
          "can_reason": true,
          "supports_attachments": true
        }
      ]
    }
  }
}
```

### Local Models

Local models can also be configured via OpenAI-compatible API. Here are two common examples:

#### Ollama

```json
{
  "providers": {
    "ollama": {
      "name": "Ollama",
      "base_url": "http://localhost:11434/v1/",
      "type": "openai-compat",
      "models": [
        {
          "name": "Qwen 3 30B",
          "id": "qwen3:30b",
          "context_window": 256000,
          "default_max_tokens": 20000
        }
      ]
    }
  }
}
```

#### LM Studio

```json
{
  "providers": {
    "lmstudio": {
      "name": "LM Studio",
      "base_url": "http://localhost:1234/v1/",
      "type": "openai-compat",
      "models": [
        {
          "name": "Qwen 3 30B",
          "id": "qwen/qwen3-30b-a3b-2507",
          "context_window": 256000,
          "default_max_tokens": 20000
        }
      ]
    }
  }
}
```

## Logging

Sometimes you need to look at logs. Luckily, Crush logs all sorts of
stuff. Logs are stored in `./.crush/logs/crush.log` relative to the project.

The CLI also contains some helper commands to make perusing recent logs easier:

```bash
# Print the last 1000 lines
crush logs

# Print the last 500 lines
crush logs --tail 500

# Follow logs in real time
crush logs --follow
```

Want more logging? Run `crush` with the `--debug` flag, or enable it in the
config:

```json
{
  "$schema": "https://charm.land/crush.json",
  "options": {
    "debug": true,
    "debug_lsp": true
  }
}
```

## Provider Auto-Updates

By default, Crush automatically checks for the latest and greatest list of
providers and models from [Catwalk](https://github.com/charmbracelet/catwalk),
the open source Crush provider database. This means that when new providers and
models are available, or when model metadata changes, Crush automatically
updates your local configuration.

### Disabling automatic provider updates

For those with restricted internet access, or those who prefer to work in
air-gapped environments, this might not be want you want, and this feature can
be disabled.

To disable automatic provider updates, set `disable_provider_auto_update` into
your `crush.json` config:

```json
{
  "$schema": "https://charm.land/crush.json",
  "options": {
    "disable_provider_auto_update": true
  }
}
```

Or set the `CRUSH_DISABLE_PROVIDER_AUTO_UPDATE` environment variable:

```bash
export CRUSH_DISABLE_PROVIDER_AUTO_UPDATE=1
```

### Manually updating providers

Manually updating providers is possible with the `crush update-providers`
command:

```bash
# Update providers remotely from Catwalk.
crush update-providers

# Update providers from a custom Catwalk base URL.
crush update-providers https://example.com/

# Update providers from a local file.
crush update-providers /path/to/local-providers.json

# Reset providers to the embedded version, embedded at crush at build time.
crush update-providers embedded

# For more info:
crush update-providers --help
```

## Metrics

Crush records pseudonymous usage metrics (tied to a device-specific hash),
which maintainers rely on to inform development and support priorities. The
metrics include solely usage metadata; prompts and responses are NEVER
collected.

Details on exactly what’s collected are in the source code ([here](https://github.com/charmbracelet/crush/tree/main/internal/event)
and [here](https://github.com/charmbracelet/crush/blob/main/internal/llm/agent/event.go)).

You can opt out of metrics collection at any time by setting the environment
variable by setting the following in your environment:

```bash
export CRUSH_DISABLE_METRICS=1
```

Or by setting the following in your config:

```json
{
  "options": {
    "disable_metrics": true
  }
}
```

Crush also respects the [`DO_NOT_TRACK`](https://donottrack.sh/) convention
which can be enabled via `export DO_NOT_TRACK=1`.

## Q&A

### Why is clipboard copy and paste not working?

Installing an extra tool might be needed on Unix-like environments.

| Environment         | Tool                     |
| ------------------- | ------------------------ |
| Windows             | Native support           |
| macOS               | Native support           |
| Linux/BSD + Wayland | `wl-copy` and `wl-paste` |
| Linux/BSD + X11     | `xclip` or `xsel`        |

## Contributing

See the [contributing guide](https://github.com/charmbracelet/crush?tab=contributing-ov-file#contributing).

## Whatcha think?

We’d love to hear your thoughts on this project. Need help? We gotchu. You can find us on:

- [Twitter](https://twitter.com/charmcli)
- [Slack][slack]
- [Discord][discord]
- [The Fediverse](https://mastodon.social/@charmcli)
- [Bluesky](https://bsky.app/profile/charm.land)

[slack]: https://charm.land/slack
[discord]: https://charm.land/discord

## License

[FSL-1.1-MIT](https://github.com/charmbracelet/crush/raw/main/LICENSE.md)

---

Part of [Charm](https://charm.land).

<a href="https://charm.land/"><img alt="The Charm logo" width="400" src="https://stuff.charm.sh/charm-banner-softy.jpg" /></a>

<!--prettier-ignore-->
Charm热爱开源 • Charm loves open source
