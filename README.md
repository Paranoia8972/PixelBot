# PixelBot

PixelBot is a Discord bot written in Go (Golang) designed to enhance your Discord server with various features and utilities.

## Features

- **Moderation Tools**: Kick, ban, mute, and manage users with ease.
- **Fun Commands**: Engage your community with fun and interactive commands.
- **Utility Commands**: Useful commands to help manage your server.
- **Customizable**: Easily configurable to fit your server's needs.

## Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/Paranoia8972/PixelBot.git
   ```
2. Navigate to the project directory:
   ```sh
   cd PixelBot
   ```
3. Build the project:
   ```sh
   go build -o PixelBot ./cmd/main.go
   ```

## Configuration and Building

1. Copy `.env.example` to `.env` and fill in the variables:

   ```sh
   cp .env.example .env
   ```

2. Build the bot:

   ```sh
   go build -o PixelBot ./cmd/bot/main.go
   ```

3. Run the bot:
   ```sh
   ./PixelBot
   ```

## Usage

## Contributing

Contributions are welcome! Please fork the repository and create a pull request with your changes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
