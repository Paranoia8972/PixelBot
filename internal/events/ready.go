package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
)

func Ready(s *discordgo.Session, r *discordgo.Ready) {
	color.Blue("Logged in as %s", r.User.Username)
	status := discordgo.UpdateStatusData{
		IdleSince: nil,
		Activities: []*discordgo.Activity{
			{
				Name: "on OnThePixel.net",
				Type: discordgo.ActivityTypeGame,
			},
		},
		Status: "online",
		AFK:    false,
	}

	s.UpdateStatusComplex(status)
}
