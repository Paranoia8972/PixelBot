package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Giveaway struct {
	MessageID    string    `bson:"message_id"`
	ChannelID    string    `bson:"channel_id"`
	EndTime      time.Time `bson:"end_time"`
	WinnersCount int       `bson:"winners_count"`
	Prize        string    `bson:"prize"`
	Participants []string  `bson:"participants"`
	Winners      []string  `bson:"winners"`
	Ended        bool      `bson:"ended"`
}

func ParseDuration(durationStr string) (time.Duration, error) {
	value := durationStr[:len(durationStr)-1]
	unit := durationStr[len(durationStr)-1:]

	numeric, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid duration format")
	}

	switch strings.ToLower(unit) {
	case "d":
		return time.Duration(numeric) * 24 * time.Hour, nil
	case "h":
		return time.Duration(numeric) * time.Hour, nil
	case "m":
		return time.Duration(numeric) * time.Minute, nil
	default:
		return 0, fmt.Errorf("invalid duration unit")
	}
}
