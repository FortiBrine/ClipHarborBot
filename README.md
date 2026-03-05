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
- 🔒 Webhook-based architecture for security and performance
- 🗄️ PostgreSQL for user data persistence

---

## 🚀 Getting Started

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) & [Docker Compose](https://docs.docker.com/compose/)
- A Telegram Bot token (get one from [@BotFather](https://t.me/BotFather))
- A [Cloudflare](https://www.cloudflare.com/) account with **Cloudflare Tunnel** — required to expose the webhook endpoint over a public HTTPS URL (the Docker mode won't work without it)

### 1. Clone the repository

```bash
git clone https://github.com/FortiBrine/ClipHarborBot.git
cd ClipHarborBot
```

### 2. Configure environment variables

Create a `.env` file in the project root:

```env
BOT_TOKEN=your_bot_token_here
WEBHOOK_URL=https://your-domain.com/webhook
WEBHOOK_SECRET=your_webhook_secret_here
CLOUDFLARE_TUNNEL_TOKEN=your_cloudflare_tunnel_token_here

POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=clipharborbot_db
```

### 3. Start the bot

```bash
docker compose --profile tunnel up --build -d
```

The bot will automatically:
- Run database migrations
- Register the webhook with Telegram
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

# Run the bot
go run ./cmd/bot/main.go
```

### Project structure

```
ClipHarborBot/
├── cmd/bot/          # Application entry point
├── internal/
│   ├── bot/          # Bot initialization & routing
│   ├── config/       # Configuration loading
│   ├── database/     # Database connection
│   ├── handler/      # Telegram update handlers
│   ├── messages/     # Localized messages
│   ├── model/        # Data models
│   ├── repository/   # Database repositories
│   └── service/      # Business logic (downloader, etc.)
├── compose.yml       # Docker Compose configuration
└── Dockerfile
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
