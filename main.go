package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Paranoia8972/PixelBot/internal/app/config"
	"github.com/Paranoia8972/PixelBot/internal/db"
	"github.com/Paranoia8972/PixelBot/internal/events"
	"github.com/Paranoia8972/PixelBot/internal/pkg/commands"
	"github.com/Paranoia8972/PixelBot/internal/pkg/transcript"
	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
)

func main() {
	log.SetFlags(0)
	log.SetOutput(io.Discard)

	// Load config and initialize MongoDB
	cfg := config.LoadConfig()
	db.InitMongoDB(cfg.MongoURI)

	dg, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		color.Red("Error creating Discord session: %v", err)
	}

	dg.Identify.Intents = discordgo.IntentsAll

	// Handlers
	dg.AddHandler(events.Ready)
	dg.AddHandler(commands.TicketCloseHandler)
	// dg.AddHandler(commands.ModalSubmitHandler)

	err = dg.Open()
	if err != nil {
		color.Red("Error opening Discord session: %v", err)
	}
	color.Green("Bot is now running. Press CTRL+C to exit.")

	// Register commands
	commands.RegisterCommands(dg, cfg)
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			switch i.ApplicationCommandData().Name {
			case "ping":
				commands.PingCommand(s, i)
			case "clear":
				commands.RemoveMessagesCommand(s, i)
			case "radio":
				commands.RadioCommand(s, i)
			case "say":
				commands.SayCommand(s, i)
			case "welcome":
				commands.WelcomeCommand(s, i)
			case "ticket":
				commands.TicketCommand(s, i)
			case "role":
				commands.RoleCommand(s, i)
			case "autorole":
				commands.AutoRoleCommand(s, i)
			}
		case discordgo.InteractionModalSubmit:
			commands.ModalSubmitHandler(s, i)
		case discordgo.InteractionMessageComponent:
			switch i.MessageComponentData().CustomID {
			case "ticket_menu":
				commands.TicketSelectHandler(s, i)
			}
		}
	})

	// Start the transcript server in a separate goroutine
	go transcript.StartTranscriptServer()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop
	color.Yellow("\nShutting down gracefully...")

	err = dg.Close()
	if err != nil {
		color.Red("Error closing Discord session: %v", err)
	}
}
