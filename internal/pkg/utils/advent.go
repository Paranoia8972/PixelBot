package utils

import (
	"context"
	"sort"
	"strconv"

	"github.com/Paranoia8972/PixelBot/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AdventClick struct {
	UserID      string   `bson:"user_id"`
	Username    string   `bson:"username"`
	Buttons     []string `bson:"buttons"`
	LevelUpDays []string `bson:"level_up_days"`
}

func StoreAdventClick(userID, username, buttonID string, levelUpDay string) error {
	collection := db.GetCollection(cfg.DBName, "advent_clicks")
	filter := bson.M{"user_id": userID}

	var adventClick AdventClick
	err := collection.FindOne(context.Background(), filter).Decode(&adventClick)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	buttonExists := false
	for _, b := range adventClick.Buttons {
		if b == buttonID {
			buttonExists = true
			break
		}
	}
	if !buttonExists {
		adventClick.Buttons = append(adventClick.Buttons, buttonID)
	}

	levelUpExists := false
	for _, day := range adventClick.LevelUpDays {
		if day == levelUpDay {
			levelUpExists = true
			break
		}
	}
	if !levelUpExists && levelUpDay != "" {
		adventClick.LevelUpDays = append(adventClick.LevelUpDays, levelUpDay)
	}

	sort.Slice(adventClick.Buttons, func(i, j int) bool {
		dayI, _ := strconv.Atoi(adventClick.Buttons[i][len("advent_"):])
		dayJ, _ := strconv.Atoi(adventClick.Buttons[j][len("advent_"):])
		return dayI < dayJ
	})

	update := bson.M{
		"$set": bson.M{
			"username":      username,
			"buttons":       adventClick.Buttons,
			"level_up_days": adventClick.LevelUpDays,
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err = collection.UpdateOne(context.Background(), filter, update, opts)
	return err
}

func GetAdventClick(userID string) (*AdventClick, error) {
	collection := db.GetCollection(cfg.DBName, "advent_clicks")
	filter := bson.M{"user_id": userID}
	var adventClick AdventClick
	err := collection.FindOne(context.Background(), filter).Decode(&adventClick)
	if err != nil {
		return nil, err
	}
	return &adventClick, nil
}

func HasButtonBeenClicked(adventClick *AdventClick, buttonID string) bool {
	for _, b := range adventClick.Buttons {
		if b == buttonID {
			return true
		}
	}
	return false
}
