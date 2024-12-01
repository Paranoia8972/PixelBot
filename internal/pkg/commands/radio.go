package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

var voiceConnection *discordgo.VoiceConnection

type Song struct {
	ID          int    `json:"id"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	Album       string `json:"album"`
	Length      int    `json:"length"`
	Genre       string `json:"genre"`
	ReleaseYear int    `json:"releaseyear"`
	Artist      struct {
		Name    string `json:"name"`
		LautID  int    `json:"laut_id"`
		URL     string `json:"url"`
		LautURL string `json:"laut_url"`
		Image   string `json:"image"`
		Thumb   string `json:"thumb"`
	} `json:"artist"`
	StartedAt string `json:"started_at"`
	EndsAt    string `json:"ends_at"`
}

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
		StopRadio(s, i)
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

	song, err := getCurrentSong()
	if err != nil {
		log.Printf("Failed to get current song: %v", err)
		return
	}

	embedMessage := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("**%s** by **%s**", song.Title, song.Artist.Name),
		Description: fmt.Sprintf("Album: %s\nGenre: %s\nRelease Year: %d\nLength: %ss", song.Album, song.Genre, song.ReleaseYear, strconv.Itoa(song.Length)),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: song.Artist.Thumb,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Started at: %s", song.StartedAt),
		},
		Color: 0x248045,
	}

	stopButton := discordgo.Button{
		Label:    "Stop Radio",
		Style:    discordgo.DangerButton,
		CustomID: "stop_radio",
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Radio started.",
			Embeds:  []*discordgo.MessageEmbed{embedMessage},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{stopButton},
				},
			},
		},
	})
	if err != nil {
		log.Printf("Failed to respond to interaction: %v", err)
		return
	}

	go monitorVoiceChannel(s, guildID, channelID)
	go updateNowPlaying(s, i.ChannelID, i.ID)

	dgvoice.PlayAudioFile(voiceConnection, "https://onthepixel.stream.laut.fm/onthepixel", make(chan bool))
}

func StopRadio(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
			Flags:   64,
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

func getCurrentSong() (*Song, error) {
	resp, err := http.Get("https://api.laut.fm/station/onthepixel/current_song")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var song Song
	err = json.Unmarshal(body, &song)
	if err != nil {
		return nil, err
	}
	return &song, nil
}

func HandleVoiceStateUpdate(s *discordgo.Session, vsu *discordgo.VoiceStateUpdate, channelID, messageID string) {
	if vsu.UserID == s.State.User.ID && vsu.ChannelID == "" {
		if voiceConnection != nil {
			voiceConnection = nil
			_, err := s.ChannelMessageEdit(channelID, messageID, "Radio stopped")
			if err != nil {
				log.Printf("Error editing message: %v", err)
			}
		}
	}
}

func updateNowPlaying(s *discordgo.Session, channelID, messageID string) {
	for {
		if voiceConnection == nil {
			_, err := s.ChannelMessageEdit(channelID, messageID, "Radio stopped")
			if err != nil {
				log.Printf("Error editing message: %v", err)
			}
			return
		}

		song, err := getCurrentSong()
		if err != nil {
			log.Printf("Error fetching current song: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}

		embedMessage := &discordgo.MessageEmbed{
			Title:       fmt.Sprintf("**%s** by **%s**", song.Title, song.Artist.Name),
			Description: fmt.Sprintf("Album: %s\nGenre: %s\nRelease Year: %d\nLength: %ss", song.Album, song.Genre, song.ReleaseYear, strconv.Itoa(song.Length)),
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: song.Artist.Thumb,
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Started at: %s", song.StartedAt),
			},
			Color: 0x248045,
		}

		_, err = s.ChannelMessageEditEmbed(channelID, messageID, embedMessage)
		if err != nil {
			log.Printf("Error updating message: %v", err)
		}

		endsAt, err := time.Parse("2006-01-02 15:04:05 -0700", song.EndsAt)
		if err != nil {
			log.Printf("Error parsing end time: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}
		time.Sleep(time.Until(endsAt))
	}
}
