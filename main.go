package main

import (
	_ "embed"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Paranoia8972/PixelBot/internal/app/config"
	"github.com/Paranoia8972/PixelBot/internal/db"
	"github.com/Paranoia8972/PixelBot/internal/events"
	"github.com/Paranoia8972/PixelBot/internal/pkg/commands"
	"github.com/Paranoia8972/PixelBot/internal/pkg/commands/moderation"
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
	dg.AddHandler(events.Welcome)
	dg.AddHandler(commands.GiveawayInteractionHandler)
	dg.AddHandler(commands.HandleVoiceStateUpdate)
	dg.AddHandler(events.MessageCreate)

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
			case "giveaway":
				commands.GiveawayCommand(s, i)
			case "edit":
				commands.EditCommand(s, i)
			case "mcstatus":
				commands.MinecraftStatusCommand(s, i)
			case "ban":
				moderation.BanCommand(s, i)
			case "unban":
				moderation.UnbanCommand(s, i)
			case "coinflip":
				commands.CoinFlipCommand(s, i)
			case "randomnumber":
				commands.RandomNumberCommand(s, i)
			case "chooser":
				commands.ChooserCommand(s, i)
			case "version":
				commands.VersionCommand(s, i)
			case "advent":
				commands.AdventCommand(s, i)
			case "level":
				commands.LevelCommand(s, i)
			case "leveling":
				commands.LevelingCommand(s, i)
			}
		}
	})

	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionMessageComponent:
			switch i.MessageComponentData().CustomID {
			case "ticket_menu":
				commands.TicketSelectHandler(s, i)
			case "close_ticket":
				commands.TicketCloseHandler(s, i)
			case "stop_radio":
				commands.StopRadio(s, i)
			default:
				if strings.HasPrefix(i.MessageComponentData().CustomID, "advent_") {
					commands.HandleAdventButton(s, i)
				}
			}
		case discordgo.InteractionModalSubmit:
			switch {
			case i.ModalSubmitData().CustomID == "say_modal":
				commands.SayCommand(s, i)
			case strings.HasPrefix(i.ModalSubmitData().CustomID, "edit_modal_"):
				commands.EditCommand(s, i)
			default:
				commands.TicketModalSubmitHandler(s, i)
			}
		case discordgo.InteractionApplicationCommandAutocomplete:
			switch i.ApplicationCommandData().Name {
			case "unban":
				moderation.UnbanAutocomplete(s, i)
			}
		}
	})

	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommand {
			switch i.ApplicationCommandData().Name {
			case "Say Hello":
				commands.SayHello(s, i)
			}
		}
	})

	// start transcript server in goroutine
	go transcript.StartTranscriptServer()
	go commands.StartBackgroundWorker(dg)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop
	color.Yellow("\nShutting down gracefully...")

	err = dg.Close()
	if err != nil {
		color.Red("Error closing Discord session: %v", err)
	}
}
