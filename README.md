# ClipHarborBot

A Telegram bot that downloads videos from TikTok, YouTube, and Instagram, right into your chat.

## Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Quick start](#quick-start)
- [Webhook mode](#webhook-mode)
- [Environment variables](#environment-variables)
- [Usage](#usage)
- [Local development](#local-development)
- [Contributing](#contributing)
- [Support the developer](#support-the-developer)
- [License](#license)

## Features

- Download videos from TikTok, YouTube, and Instagram
- Multi-language support (English, Polish, Ukrainian)
- Downloads via yt-dlp, capped to Telegram's 50 MB upload limit
- Long-polling or webhook transport (auto-selected based on config)
- PostgreSQL for per-user language persistence
- Deploys with Docker Compose, optional Cloudflare Tunnel for webhook mode

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/)
- A Telegram bot token from [@BotFather](https://t.me/BotFather)

## Quick start

1. Clone the repository and create a `.env` file (see `.env.example`).
2. Set at least `TELEGRAM_TOKEN` and the `POSTGRES_*` variables.
3. Start the bot:

   ```bash
   docker compose up --build -d
   ```

The bot runs in long-polling mode by default (no public URL needed), runs database migrations automatically, and starts listening for messages.

## Webhook mode

Set `WEBHOOK_URL` in `.env` to a publicly reachable HTTPS URL to switch the bot from polling to webhook mode. There's no separate mode flag: the bot infers the transport from whether `WEBHOOK_URL` is set.

```bash
docker compose up --build -d
```

To expose the bot via Cloudflare Tunnel instead of your own reverse proxy, also set `CLOUDFLARE_TUNNEL_TOKEN` and start with the `tunnel` profile:

```bash
docker compose --profile tunnel up --build -d
```

## Environment variables

| Variable                  | Required | Default | Description |
|----------------------------|----------|---------|--------------|
| `TELEGRAM_TOKEN`           | yes      | -       | Telegram bot token |
| `WEBHOOK_URL`              | no       | -       | Enables webhook mode when set |
| `HTTP_ADDRESS`             | no       | `:2000` | Address for the webhook/health HTTP server |
| `CLOUDFLARE_TUNNEL_TOKEN`  | no       | -       | Only needed with the `tunnel` compose profile |
| `ENVIRONMENT`              | no       | `dev`   | `dev` or `prod`; controls log level |
| `DEFAULT_LANG`             | no       | `en`    | Language for users who haven't run `/lang` |
| `DOWNLOAD_TIMEOUT`         | no       | `5m`    | Per-download yt-dlp timeout |
| `SHUTDOWN_TIMEOUT`         | no       | `10s`   | Grace period for shutdown on SIGINT |
| `POSTGRES_HOST`            | yes      | -       | |
| `POSTGRES_PORT`            | yes      | -       | |
| `POSTGRES_USER`            | yes      | -       | |
| `POSTGRES_PASSWORD`        | yes      | -       | |
| `POSTGRES_DB`              | yes      | -       | |

## Usage

1. Open Telegram and find your bot.
2. Send `/start`.
3. Send `/lang` to pick English, Polish, or Ukrainian.
4. Paste a TikTok, YouTube, or Instagram link, the bot downloads it and sends the video back.

## Local development

Requires Go 1.26+, [yt-dlp](https://github.com/yt-dlp/yt-dlp) on `PATH`, [sqlc](https://docs.sqlc.dev/en/latest/overview/install.html), and a reachable Postgres instance.

```bash
pip install yt-dlp
go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.31.1
go mod download
go run ./cmd/bot
```

Database migrations run automatically on startup. If you change a query in `db/queries/` or a migration in `db/migrations/`, regenerate the sqlc-generated code:

```bash
sqlc generate
```

## Contributing

1. Fork the repository.
2. Create a feature branch: `git checkout -b feat/your-feature`.
3. Commit with a clear message (preferably [Conventional Commits](https://www.conventionalcommits.org/)).
4. Open a pull request describing what changed and why.

Open an issue first for large changes so we can discuss the approach.

## Support the developer

If you find this project useful, consider supporting its development:

| Method | Link |
|--------|------|
| Buy Me a Coffee | Coming soon |
| GitHub Star | Just star this repo, it helps a lot |

## License

Dual-licensed under Apache License 2.0 or MIT, at your option. See [LICENSE-APACHE](LICENSE-APACHE) and [LICENSE-MIT](LICENSE-MIT).
