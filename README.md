# 🎬 ClipHarborBot

> A Telegram bot that downloads videos from **TikTok**, **YouTube**, **Instagram** and more — right into your chat.

---

## 📸 Screenshots

| Start | Download |
|-------|----------|
| ![Start screen](docs/screenshots/start.png) | ![Download screen](docs/screenshots/download.png) |

---

## ✨ Features

- 📥 Download videos from **TikTok**, **YouTube**, and **Instagram** _(more platforms coming soon)_
- 🌐 Multi-language support
- ⚡ Fast and reliable via [yt-dlp](https://github.com/yt-dlp/yt-dlp)
- 🐳 Ready to deploy with **Docker Compose**
- 🔄 Supports both **Polling** and **Webhook** modes
- 🗄️ PostgreSQL for user data persistence

---

## 🚀 Getting Started

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) & [Docker Compose](https://docs.docker.com/compose/)
- A Telegram Bot token (get one from [@BotFather](https://t.me/BotFather))

### 1. Clone the repository

```bash
git clone https://github.com/FortiBrine/ClipHarborBot.git
cd ClipHarborBot
```

### 2. Configure environment variables

Create a `.env` file in the project root.

> **Note:** `BOT_MODE` can be set to `polling` or `webhook`. Defaults to `polling` if not set or invalid.

#### Polling mode _(simpler, no public URL required)_

```env
BOT_TOKEN=your_bot_token_here
BOT_MODE=polling

POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=clipharborbot_db
```

#### Webhook mode _(requires a public HTTPS URL)_

```env
BOT_TOKEN=your_bot_token_here
BOT_MODE=webhook
WEBHOOK_URL=https://your-domain.com/webhook
WEBHOOK_SECRET=your_webhook_secret_here

# Optional: only needed if you use Cloudflare Tunnel to expose the webhook
CLOUDFLARE_TUNNEL_TOKEN=your_cloudflare_tunnel_token_here

POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=clipharborbot_db
```

### 3. Start the bot

#### Polling mode

```bash
docker compose up --build -d
```

#### Webhook mode — with your own reverse proxy / public IP

```bash
docker compose up --build -d
```

#### Webhook mode — with Cloudflare Tunnel

```bash
docker compose --profile tunnel up --build -d
```

The bot will automatically:
- Run database migrations
- Register the webhook with Telegram _(webhook mode only)_
- Start listening for messages

### 4. Use the bot

1. Open Telegram and find your bot.
2. Send `/start` to begin.
3. Paste a **TikTok**, **YouTube**, or **Instagram** link — the bot will download and send the video back to you.

---

## 🛠️ Development

### Run locally (without Docker)

Make sure you have **Go 1.25+** and **yt-dlp** installed.

```bash
# Install yt-dlp
pip install yt-dlp

# Install Go dependencies
go mod download

# Run the bot (polling mode by default)
go run ./cmd/bot/main.go
```

### Project structure

```
ClipHarborBot/
├── cmd/bot/              # Application entry point
├── internal/
│   ├── bot/              # Bot initialization & routing
│   ├── config/           # Configuration loading (BOT_MODE, tokens, DB)
│   ├── database/         # Database connection & migrations
│   ├── handler/          # Telegram update handlers
│   │   ├── default.go    # Fallback/unknown message handler
│   │   ├── language.go   # Language selection handler
│   │   ├── start.go      # /start command handler
│   │   └── video.go      # Video download handler
│   ├── messages/         # Localized messages
│   ├── model/            # Data models
│   ├── platform/         # Platform detection logic
│   ├── repository/       # Database repositories
│   └── service/          # Business logic
│       ├── downloader.go     # yt-dlp wrapper
│       └── message_service.go
├── compose.yml           # Docker Compose configuration
├── Dockerfile
└── .env.example          # Example environment variables
```

---

## 🤝 Contributing

Contributions are welcome! Here's how to get started:

1. **Fork** the repository
2. **Create** a feature branch: `git checkout -b feat/your-feature`
3. **Commit** your changes: `git commit -m "feat: add your feature"`
4. **Push** to the branch: `git push origin feat/your-feature`
5. **Open a Pull Request** — describe what you've changed and why

### Guidelines

- Follow the existing code style and project structure.
- Write clear, descriptive commit messages (preferably [Conventional Commits](https://www.conventionalcommits.org/)).
- If you're adding a new platform, create a dedicated service file under `internal/service/`.
- Open an issue first for large changes so we can discuss the approach.

---

## 💖 Support the Developer

If you find this project useful, consider supporting its development:

| Method | Link                                  |
|--------|---------------------------------------|
| ☕ Buy Me a Coffee | Coming soon                           |
| ⭐ GitHub Star | Just star this repo — it helps a lot! |

---

## 📄 License

This project is licensed under the **Apache License 2.0** — see the [LICENSE](LICENSE) file for details.
