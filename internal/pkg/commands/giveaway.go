package commands

import (
	"context"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/Paranoia8972/PixelBot/internal/db"
	"github.com/Paranoia8972/PixelBot/internal/pkg/utils"

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
	Ended        bool      `bson:"ended"`
}

func GiveawayCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	if len(options) == 0 {
		RespondWithMessage(s, i, "Please provide a subcommand.")
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

	duration, err := utils.ParseDuration(durationStr)
	if err != nil {
		RespondWithMessage(s, i, "Invalid duration format. Use format like: 2d, 4h, 30m")
		return
	}

	msgSend := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "🎉 Giveaway Started! 🎉",
				Description: "Prize: " + prize + "\nEnds: <t:" + strconv.FormatInt(time.Now().Add(duration).Unix(), 10) + ":R>\nClick the button below to enter!",
				Color:       0x248045,
				Footer: &discordgo.MessageEmbedFooter{
					Text: "Good luck!",
				},
			},
		},
	}

	msg, err := s.ChannelMessageSendComplex(i.ChannelID, msgSend)
	if err != nil {
		RespondWithMessage(s, i, "Failed to send giveaway message.")
		return
	}

	msgEdit := &discordgo.MessageEdit{
		ID:      msg.ID,
		Channel: i.ChannelID,
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

	_, err = s.ChannelMessageEditComplex(msgEdit)
	if err != nil {
		RespondWithMessage(s, i, "Failed to add button to giveaway message.")
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
		RespondWithMessage(s, i, "Failed to save giveaway to database.")
		return
	}

	go func() {
		time.Sleep(duration)
		endGiveawayLogic(s, giveaway)
	}()

	RespondWithMessage(s, i, "Giveaway started!")
}

func endGiveaway(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	messageID := options[0].StringValue()

	collection := db.GetCollection(cfg.DBName, "giveaways")
	var giveaway Giveaway
	err := collection.FindOne(context.TODO(), bson.M{"message_id": messageID}).Decode(&giveaway)
	if err != nil {
		RespondWithMessage(s, i, "Giveaway not found.")
		return
	}

	endGiveawayLogic(s, giveaway)
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
		RespondWithMessage(s, i, "Giveaway not found.")
		return
	}

	selectWinners(&giveaway)

	s.ChannelMessageSend(giveaway.ChannelID, "Giveaway has been rerolled!")

	RespondWithMessage(s, i, "Giveaway rerolled!")
}

func endGiveawayLogic(s *discordgo.Session, giveaway Giveaway) {
	collection := db.GetCollection(cfg.DBName, "giveaways")
	err := collection.FindOne(context.TODO(), bson.M{"message_id": giveaway.MessageID}).Decode(&giveaway)
	if err != nil {
		s.ChannelMessageSend(giveaway.ChannelID, "Error fetching giveaway data.")
		return
	}

	if giveaway.Ended {
		return
	}

	if len(giveaway.Participants) == 0 {
		s.ChannelMessageSend(giveaway.ChannelID, "No participants entered the giveaway.")
	} else {
		winners := selectWinners(&giveaway)
		giveaway.Winners = winners
		message := "**Giveaway Ended!**\nPrize: " + giveaway.Prize + "\nWinners: " + formatWinners(winners) + "https://discord.com/channels/" + s.State.Guilds[0].ID + "/" + giveaway.ChannelID + "/" + giveaway.MessageID
		s.ChannelMessageSend(giveaway.ChannelID, message)
	}

	giveaway.Ended = true

	_, err = collection.UpdateOne(
		context.TODO(),
		bson.M{"message_id": giveaway.MessageID},
		bson.M{"$set": bson.M{"ended": true, "winners": giveaway.Winners}},
	)
	if err != nil {
		s.ChannelMessageSend(giveaway.ChannelID, "Error updating giveaway status.")
		return
	}

	msg, err := s.ChannelMessage(giveaway.ChannelID, giveaway.MessageID)
	if err != nil {
		s.ChannelMessageSend(giveaway.ChannelID, "Error fetching original giveaway message.")
		return
	}

	if len(msg.Embeds) > 0 {
		msg.Embeds[0].Title = "Giveaway Ended"
	}

	for _, component := range msg.Components {
		for _, actionRow := range component.(*discordgo.ActionsRow).Components {
			if button, ok := actionRow.(*discordgo.Button); ok {
				button.Disabled = true
			}
		}
	}

	_, err = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID:         msg.ID,
		Channel:    msg.ChannelID,
		Embeds:     &msg.Embeds,
		Components: &msg.Components,
	})
	if err != nil {
		s.ChannelMessageSend(giveaway.ChannelID, "Error editing giveaway message.")
		return
	}
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
				RespondWithMessage(s, i, "Giveaway not found.")

				return
			}

			if time.Now().After(giveaway.EndTime) {
				RespondWithMessage(s, i, "Giveaway has ended.")
				return
			}

			var userID string
			if i.User != nil {
				userID = i.User.ID
			} else if i.Member != nil && i.Member.User != nil {
				userID = i.Member.User.ID
			} else {
				RespondWithMessage(s, i, "Unable to retrieve user information.")

				return
			}

			if giveaway.Participants == nil {
				giveaway.Participants = []string{}
			}

			for _, participant := range giveaway.Participants {
				if participant == userID {
					RespondWithMessage(s, i, "You have already entered this giveaway.")
					return
				}
			}

			giveaway.Participants = append(giveaway.Participants, userID)

			_, err = collection.UpdateOne(context.TODO(), bson.M{"message_id": messageID}, bson.M{"$set": bson.M{"participants": giveaway.Participants}})
			if err != nil {
				log.Println("Failed to update participants:", err)
				return
			}

			RespondWithMessage(s, i, "You have entered the giveaway!")
		}
	}
}
