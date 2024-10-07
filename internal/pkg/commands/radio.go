package commands

import (
	"log"
	"time"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

var voiceConnection *discordgo.VoiceConnection

func RadioCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if len(i.ApplicationCommandData().Options) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No subcommand provided.",
				Flags:   64,
			},
		})
		return
	}

	switch i.ApplicationCommandData().Options[0].Name {
	case "start":
		startRadio(s, i)
	case "stop":
		stopRadio(s, i)
	}
}

func startRadio(s *discordgo.Session, i *discordgo.InteractionCreate) {
	guildID := cfg.GuildID
	userID := i.Member.User.ID

	guild, err := s.State.Guild(guildID)
	if err != nil {
		log.Fatalf("Failed to get guild: %v", err)
	}

	var channelID string
	for _, vs := range guild.VoiceStates {
		if vs.UserID == userID {
			channelID = vs.ChannelID
			break
		}
	}

	if channelID == "" {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You need to be in a voice channel to start the radio.",
				Flags:   64,
			},
		})
		if err != nil {
			log.Printf("Failed to respond to interaction: %v", err)
		}
		return
	}

	var errJoin error
	voiceConnection, errJoin = s.ChannelVoiceJoin(guildID, channelID, false, true)
	if errJoin != nil {
		log.Fatalf("Failed to join voice channel: %v", errJoin)
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: "Radio started!",
					Color: 0x00ff00,
				},
			},
		},
	})
	if err != nil {
		log.Printf("Failed to respond to interaction: %v", err)
	}

	go monitorVoiceChannel(s, guildID, channelID)

	dgvoice.PlayAudioFile(voiceConnection, "https://onthepixel.stream.laut.fm/onthepixel", make(chan bool))
}

func stopRadio(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if voiceConnection == nil {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Radio is currently not playing.",
				Flags:   64,
			},
		})
		if err != nil {
			log.Printf("Failed to respond to interaction: %v", err)
		}
		return
	}

	voiceConnection.Disconnect()
	voiceConnection = nil

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Radio stopped.",
		},
	})
	if err != nil {
		log.Printf("Failed to respond to interaction: %v", err)
	}
}

func monitorVoiceChannel(s *discordgo.Session, guildID, channelID string) {
	for {
		time.Sleep(10 * time.Second)

		guild, err := s.State.Guild(guildID)
		if err != nil {
			log.Printf("Failed to get guild: %v", err)
			continue
		}

		var userCount int
		for _, vs := range guild.VoiceStates {
			if vs.ChannelID == channelID {
				userCount++
			}
		}

		if userCount <= 1 {
			if voiceConnection != nil {
				err := voiceConnection.Disconnect()
				if err != nil {
					log.Printf("Failed to disconnect: %v", err)
				} else {
					voiceConnection = nil
				}
			}
			break
		}
	}
}
