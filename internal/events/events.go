package events

import (
	"github.com/Paranoia8972/PixelBot/internal/app/config"
)

var cfg *config.Config

func init() {
	cfg = config.LoadConfig()
}
