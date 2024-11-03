package commands

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
)

type MinecraftServerStatus struct {
	IP    string `json:"ip"`
	Port  int    `json:"port"`
	Debug struct {
		Ping          bool `json:"ping"`
		Query         bool `json:"query"`
		Srv           bool `json:"srv"`
		QueryMismatch bool `json:"querymismatch"`
		IpInSrv       bool `json:"ipinsrv"`
		CnameInSrv    bool `json:"cnameinsrv"`
		AnimatedMotd  bool `json:"animatedmotd"`
		CacheHit      bool `json:"cachehit"`
		CacheTime     int  `json:"cachetime"`
		CacheExpire   int  `json:"cacheexpire"`
		ApiVersion    int  `json:"apiversion"`
		Error         struct {
			Query string `json:"query"`
		} `json:"error"`
	} `json:"debug"`
	Motd struct {
		Raw   []string `json:"raw"`
		Clean []string `json:"clean"`
		Html  []string `json:"html"`
	} `json:"motd"`
	Players struct {
		Online int `json:"online"`
		Max    int `json:"max"`
	} `json:"players"`
	Version      string `json:"version"`
	Online       bool   `json:"online"`
	Protocol     int    `json:"protocol"`
	ProtocolName string `json:"protocol_name"`
	Hostname     string `json:"hostname"`
	Icon         string `json:"icon"`
	Software     string `json:"software"`
	EulaBlocked  bool   `json:"eula_blocked"`
}

func MinecraftStatusCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	serverIP := i.ApplicationCommandData().Options[0].StringValue()
	apiURL := fmt.Sprintf("https://api.mcsrvstat.us/2/%s", serverIP)

	resp, err := http.Get(apiURL)
	if err != nil {
		RespondWithMessage(s, i, "Failed to fetch server status.")
		return
	}
	defer resp.Body.Close()

	var status MinecraftServerStatus
	err = json.NewDecoder(resp.Body).Decode(&status)
	if err != nil {
		log.Printf("Error decoding server status: %v", err)
		RespondWithMessage(s, i, "Failed to parse server status.")
		return
	}

	if status.Icon != "" {
		iconData, err := base64.StdEncoding.DecodeString(status.Icon)
		if err != nil {
			log.Printf("Error decoding server icon: %v", err)
		} else {
			err = os.WriteFile("/assets/servericon.png", iconData, 0644)
			if err != nil {
				log.Printf("Error saving server icon: %v", err)
			}
		}
	}

	motd := ""
	if len(status.Motd.Clean) > 0 {
		motd = status.Motd.Clean[0]
	}

	message := fmt.Sprintf("**Server Status for %s**\n", serverIP)
	message += fmt.Sprintf("Online: %t\n", status.Online)
	message += fmt.Sprintf("Players Online: %d/%d\n", status.Players.Online, status.Players.Max)
	message += fmt.Sprintf("Version: %s\n", status.Version)
	message += fmt.Sprintf("Protocol: %s\n", status.Software)
	message += fmt.Sprintf("MOTD:\n```\n%s\n```", motd)

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("Minecraft Server Status for %s", serverIP),
		Description: message,
	}

	if status.Icon != "" {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: "attachment://assets/servericon.png",
		}
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
	}
}
