package events

import (
	"github.com/Paranoia8972/PixelBot/internal/db"
	"github.com/Paranoia8972/PixelBot/internal/pkg/status"
	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
)

var StatusManager *status.Manager

func Ready(s *discordgo.Session, r *discordgo.Ready) {
	color.Blue("Logged in as %s", r.User.Username)
	collection := db.GetCollection(cfg.DBName, "statuses")
	StatusManager = status.NewManager(s, collection)
}
