// ecosystem.config.js — PM2 进程管理配置
//
// 用法：
//   构建：     CGO_ENABLED=0 GOEXPERIMENT=greenteagc go build -o crush .
//   启动：     pm2 start ecosystem.config.js
//   查看日志： pm2 logs crush-web
//   重启：     pm2 restart crush-web
//   停止：     pm2 stop crush-web
//   开机自启： pm2 save && pm2 startup   (按提示执行 sudo 命令)
//
// 日志默认写到 ./logs/，目录会在首次启动时由 PM2 自动创建。

module.exports = {
  apps: [
    {
      name: 'crush-web',
      script: './crush',
      args: 'web --addr :8085 --pages-dir ./pages --no-auth',
      cwd: __dirname,
      interpreter: 'none', // 直接执行二进制，不走 node

      // 进程模式
      exec_mode: 'fork',
      instances: 1,
      autorestart: true,
      watch: false,
      max_restarts: 10,
      restart_delay: 3000,
      min_uptime: '10s',
      kill_timeout: 8000, // 给 SSE 连接和 DB 留出优雅关闭时间
      wait_ready: false,

      // 资源限制（按需调整）
      max_memory_restart: '1G',

      // 日志
      out_file: './logs/crush-web.out.log',
      error_file: './logs/crush-web.err.log',
      merge_logs: true,
      log_date_format: 'YYYY-MM-DD HH:mm:ss',

      // 环境变量
      env: {
        NODE_ENV: 'production',
        // 如需指定终端类型避免 TUI 库异常输出：
        TERM: 'xterm-256color',
        // 在这里追加 LLM provider 的 API key，例如：
        // ANTHROPIC_API_KEY: '...',
        // OPENAI_API_KEY: '...',
      },
    },
  ],
};
