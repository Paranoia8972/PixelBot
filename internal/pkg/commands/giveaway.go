package commands

import (
	"context"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/Paranoia8972/PixelBot/internal/db"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
)

type Giveaway struct {
	MessageID    string    `bson:"message_id"`
	ChannelID    string    `bson:"channel_id"`
	EndTime      time.Time `bson:"end_time"`
	WinnersCount int       `bson:"winners_count"`
	Prize        string    `bson:"prize"`
	Participants []string  `bson:"participants"`
	Winners      []string  `bson:"winners"`
}

func GiveawayCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	if len(options) == 0 {
		respondWithMessage(s, i, "Please provide a subcommand.")
		return
	}

	subCommand := options[0].Name
	switch subCommand {
	case "start":
		startGiveaway(s, i, options[0].Options)
	case "end":
		endGiveaway(s, i, options[0].Options)
	case "reroll":
		rerollGiveaway(s, i, options[0].Options)
	}
}

func startGiveaway(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	durationStr := options[0].StringValue()
	winnersCount := int(options[1].IntValue())
	prize := options[2].StringValue()

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		respondWithMessage(s, i, "Invalid duration format.")
		return
	}

	msg, err := s.ChannelMessageSend(i.ChannelID, "**Giveaway Starting...**")
	if err != nil {
		respondWithMessage(s, i, "Failed to send giveaway message.")
		return
	}

	giveawayMessage := &discordgo.MessageEdit{
		ID:      msg.ID,
		Channel: msg.ChannelID,
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Title:       "🎉 Giveaway Started! 🎉",
				Description: "Prize: " + prize + "\nEnds: <t:" + strconv.FormatInt(time.Now().Add(duration).Unix(), 10) + ":R>\nClick the button below to enter!",
				Color:       0x00FF00,
				Footer: &discordgo.MessageEmbedFooter{
					Text: "Good luck!",
				},
			},
		},
		Components: &[]discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Enter Giveaway",
						CustomID: "giveaway_enter_" + msg.ID,
						Style:    discordgo.PrimaryButton,
					},
				},
			},
		},
	}

	_, err = s.ChannelMessageEditComplex(giveawayMessage)
	if err != nil {
		respondWithMessage(s, i, "Failed to edit giveaway message.")
		return
	}

	giveaway := Giveaway{
		MessageID:    msg.ID,
		ChannelID:    i.ChannelID,
		EndTime:      time.Now().Add(duration),
		WinnersCount: winnersCount,
		Prize:        prize,
		Participants: []string{},
		Winners:      []string{},
	}

	collection := db.GetCollection(cfg.DBName, "giveaways")
	_, err = collection.InsertOne(context.TODO(), giveaway)
	if err != nil {
		respondWithMessage(s, i, "Failed to save giveaway to database.")
		return
	}

	go func() {
		time.Sleep(duration)
		endGiveawayLogic(s, giveaway)
	}()

	respondWithMessage(s, i, "Giveaway started!")
}

func endGiveaway(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	messageID := options[0].StringValue()

	collection := db.GetCollection(cfg.DBName, "giveaways")
	var giveaway Giveaway
	err := collection.FindOne(context.TODO(), bson.M{"message_id": messageID}).Decode(&giveaway)
	if err != nil {
		respondWithMessage(s, i, "Giveaway not found.")
		return
	}

	endGiveawayLogic(s, giveaway)

	respondWithMessage(s, i, "Giveaway ended!")
}

func StartBackgroundWorker(s *discordgo.Session) {
	go func() {
		for {
			checkEndedGiveaways(s)
			time.Sleep(10 * time.Second)
		}
	}()
}

func checkEndedGiveaways(s *discordgo.Session) {
	collection := db.GetCollection(cfg.DBName, "giveaways")
	now := time.Now()

	cursor, err := collection.Find(context.TODO(), bson.M{"end_time": bson.M{"$lte": now}, "winners": bson.M{"$size": 0}})
	if err != nil {
		log.Println("Failed to fetch ended giveaways:", err)
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var giveaway Giveaway
		if err := cursor.Decode(&giveaway); err != nil {
			log.Println("Failed to decode giveaway:", err)
			continue
		}

		endGiveawayLogic(s, giveaway)
	}
}

func rerollGiveaway(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	messageID := options[0].StringValue()

	collection := db.GetCollection(cfg.DBName, "giveaways")
	var giveaway Giveaway
	err := collection.FindOne(context.TODO(), bson.M{"message_id": messageID}).Decode(&giveaway)
	if err != nil {
		respondWithMessage(s, i, "Giveaway not found.")
		return
	}

	selectWinners(&giveaway)

	s.ChannelMessageSend(giveaway.ChannelID, "Giveaway has been rerolled!")

	respondWithMessage(s, i, "Giveaway rerolled!")
}

func endGiveawayLogic(s *discordgo.Session, giveaway Giveaway) {
	collection := db.GetCollection(cfg.DBName, "giveaways")
	err := collection.FindOne(context.TODO(), bson.M{"message_id": giveaway.MessageID}).Decode(&giveaway)
	if err != nil {
		s.ChannelMessageSend(giveaway.ChannelID, "Error fetching giveaway data.")
		return
	}

	if len(giveaway.Participants) == 0 {
		s.ChannelMessageSend(giveaway.ChannelID, "No participants entered the giveaway.")
		return
	}

	winners := selectWinners(&giveaway)

	giveaway.Winners = winners
	_, err = collection.UpdateOne(context.TODO(), bson.M{"message_id": giveaway.MessageID}, bson.M{"$set": bson.M{"winners": winners}})
	if err != nil {
		s.ChannelMessageSend(giveaway.ChannelID, "Error updating giveaway winners.")
		return
	}

	message := "**Giveaway Ended!**\nPrize: " + giveaway.Prize + "\nWinners: " + formatWinners(winners) + "https://discord.com/channels/" + s.State.Guilds[0].ID + "/" + giveaway.ChannelID + "/" + giveaway.MessageID
	s.ChannelMessageSend(giveaway.ChannelID, message)
}

func selectWinners(giveaway *Giveaway) []string {
	winners := pickWinners(giveaway.Participants, giveaway.WinnersCount)

	winnerMentions := ""
	for _, winnerID := range winners {
		winnerMentions += "<@" + winnerID + "> "
	}

	return winners
}

func pickWinners(participants []string, count int) []string {
	if count > len(participants) {
		count = len(participants)
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(participants), func(i, j int) { participants[i], participants[j] = participants[j], participants[i] })
	return participants[:count]
}

func formatWinners(winners []string) string {
	winnerMentions := ""
	for _, winnerID := range winners {
		winnerMentions += "<@" + winnerID + "> "
	}
	return winnerMentions
}

func GiveawayInteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionMessageComponent:
		data := i.MessageComponentData()
		if len(data.CustomID) >= len("giveaway_enter_") && data.CustomID[:len("giveaway_enter_")] == "giveaway_enter_" {
			messageID := data.CustomID[len("giveaway_enter_"):]

			collection := db.GetCollection(cfg.DBName, "giveaways")
			if collection == nil {
				log.Println("Failed to get collection")
				return
			}
			var giveaway Giveaway
			err := collection.FindOne(context.TODO(), bson.M{"message_id": messageID}).Decode(&giveaway)
			if err != nil {
				respondWithMessage(s, i, "Giveaway not found.")

				return
			}

			if time.Now().After(giveaway.EndTime) {
				respondWithMessage(s, i, "Giveaway has ended.")
				return
			}

			var userID string
			if i.User != nil {
				userID = i.User.ID
			} else if i.Member != nil && i.Member.User != nil {
				userID = i.Member.User.ID
			} else {
				respondWithMessage(s, i, "Unable to retrieve user information.")

				return
			}

			if giveaway.Participants == nil {
				giveaway.Participants = []string{}
			}

			for _, participant := range giveaway.Participants {
				if participant == userID {
					respondWithMessage(s, i, "You have already entered this giveaway.")
					return
				}
			}

			giveaway.Participants = append(giveaway.Participants, userID)

			_, err = collection.UpdateOne(context.TODO(), bson.M{"message_id": messageID}, bson.M{"$set": bson.M{"participants": giveaway.Participants}})
			if err != nil {
				log.Println("Failed to update participants:", err)
				return
			}

			respondWithMessage(s, i, "You have entered the giveaway!")
		}
	}
}
