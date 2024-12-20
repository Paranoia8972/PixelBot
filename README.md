# 🤖 PixelBot

A feature-rich Discord bot written in Go for managing gaming communities and Minecraft servers. Built using the [DiscordGo](https://github.com/bwmarrin/discordgo) library.

## ✨ Features

### 🎫 Ticket System

- Custom ticket categories with modals
- Ticket transcripts
- Automated ticket management

### 🎮 Server Management

- Welcome messages with custom images
- Auto-role assignment
- Level system with rewards
- Message cleanup
- Role management

### 🎁 Fun & Engagement

- Giveaways system
- Random number generator
- Coin flip
- Random chooser from reactions/lists

### 🎵 Radio

- Music streaming capabilities
- Basic playback controls

### 🎯 Minecraft Integration

- Server status checking
- Player statistics

## 🚧 Work in Progress

- [ ] Advanced Ticket system
- [x] Counting game
- [ ] Advamced Leveling system
- [ ] Setup command

## 📝 Prerequisites

- Go 1.23 or higher
- MongoDB Database
- Discord Bot Token

## ⚙️ Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/PixelBot.git
cd PixelBot
```

1. Copy the example environment file:

```bash
cp .env.example .env
```

1. Configure your `.env`

file:

```env
DISCORD_BOT_TOKEN="abcDEFghiJKLmnoPQRstuvwxyz.123456.789-_ABC"
GUILD_ID="123456789012345678"
MONGO_URI="mongodb://root:password123@127.0.0.1:27017"
DB_NAME="example"
SERVER_PORT=8080
TRANSCRIPT_URL="transcripts.example.com"
```

## 🚀 Building and Running

Build the project:

```bash
go build -o bin/PixelBot
```

Run the bot:

```bash
./bin/PixelBot
```

## 📚 Usage

1. Invite the bot to your server using the OAuth2 URL
2. Use `/help` to see available commands
3. Set up basic configurations:
   - Welcome channel: `/welcome set`
   - Auto roles: `/autorole add`
   - Level rewards: `/leveling set_reward`

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- [DiscordGo](https://github.com/bwmarrin/discordgo)
- [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver)
