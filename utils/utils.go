package utils

import (
	"github.com/Paranoia8972/PixelBot/config"
)

var cfg *config.Config

func init() {
	cfg = config.LoadConfig()
}
