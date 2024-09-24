package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Paranoia8972/PixelBot/internal/app/config"
	"github.com/Paranoia8972/PixelBot/internal/db"
	"github.com/fatih/color"
	"github.com/russross/blackfriday/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var cfg *config.Config

func init() {
	cfg = config.LoadConfig()
}

//TODO:
//TODO: Emojis, Mentions, Attachments, Formatting in Embeds, Interactions (commands, buttons, etc)
//TODO:

type Attachment struct {
	Filename string `json:"filename"`
	Type     string `json:"type"`
	URL      string `json:"url"`
}

type Embed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Fields      []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"fields"`
	URL   string `json:"url"`
	Color int    `json:"color"`
}

type TranscriptMessage struct {
	Attachments    []Attachment `json:"attachments"`
	MessageContent string       `json:"message_content"`
	Pfp            string       `json:"pfp"`
	Timestamp      string       `json:"timestamp"`
	Username       string       `json:"username"`
	Embeds         []Embed      `json:"embeds"`
}

type TranscriptData struct {
	Transcript []TranscriptMessage `json:"transcript"`
}

func main() {
	cfg := config.LoadConfig()
	db.InitMongoDB(cfg.MongoURI)

	http.HandleFunc("/ticket", TranscriptServer)

	color.Green("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func TranscriptServer(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	transcript, err := getTranscript(id)
	if err != nil {
		http.Error(w, "Error fetching transcript: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var data TranscriptData

	err = json.Unmarshal(transcript, &data)
	if err != nil {
		http.Error(w, "Error parsing transcript: "+err.Error(), http.StatusInternalServerError)
		return
	}

	htmlTemplate, err := os.ReadFile("template.html")
	if err != nil {
		http.Error(w, "Error reading HTML template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	messagesHTML := ""
	var lastUsername string
	for i := len(data.Transcript) - 1; i >= 0; i-- {
		msg := data.Transcript[i]
		formattedTime := formatTimestamp(msg.Timestamp)

		// Render markdown to HTML
		renderedContent := string(blackfriday.Run([]byte(msg.MessageContent)))

		if msg.Username != lastUsername {
			messagesHTML += `<div class="message">
                <img src="` + msg.Pfp + `" class="pfp" />
                <div class="content">
                    <div>
                        <span class="username">` + msg.Username + `</span>
                        <span class="timestamp">` + formattedTime + `</span>
                    </div>
                    <div class="message_content">` + renderedContent + `</div>`
			if len(msg.Embeds) > 0 {
				for _, embed := range msg.Embeds {
					messagesHTML += `<div class="embed" style="border-left: 4px solid #` + fmt.Sprintf("%06x", embed.Color) + `;">
                        <div class="embed_title">` + embed.Title + `</div>
                        <div class="embed_description">` + embed.Description + `</div>`
					if len(embed.Fields) > 0 {
						for _, field := range embed.Fields {
							messagesHTML += `<div class="embed_field">
                                <div class="embed_field_name">` + field.Name + `</div>
                                <div class="embed_field_value">` + field.Value + `</div>
                            </div>`
						}
					}
					messagesHTML += `</div>`
				}
			}
			if len(msg.Attachments) > 0 {
				for _, attachment := range msg.Attachments {
					if attachment.Type == "image/jpeg" || attachment.Type == "image/png" || attachment.Type == "image/gif" {
						messagesHTML += `<div class="attachment">
                            <img src="` + attachment.URL + `" alt="` + attachment.Filename + `" />
                        </div>`
					} else {
						messagesHTML += `<div class="attachment">
                            <a href="` + attachment.URL + `" download>` + attachment.Filename + `</a>
                        </div>`
					}
				}
			}
			messagesHTML += `</div>
            </div>`
			lastUsername = msg.Username
		} else {
			messagesHTML += `<div class="message no-padding">
                <div class="content">
                    <div class="message_content">` + renderedContent + `</div>`
			if len(msg.Embeds) > 0 {
				for _, embed := range msg.Embeds {
					messagesHTML += `<div class="embed" style="border-left: 4px solid #` + fmt.Sprintf("%06x", embed.Color) + `;">
                        <h3 class="embed_title">` + embed.Title + `</h3>
                        <div class="embed_description">` + embed.Description + `</div>`
					if len(embed.Fields) > 0 {
						for _, field := range embed.Fields {
							messagesHTML += `<div class="embed_field">
                                <div class="embed_field_name">` + field.Name + `</div>
                                <div class="embed_field_value">` + field.Value + `</div>
                            </div>`
						}
					}
					messagesHTML += `</div>`
				}
			}
			if len(msg.Attachments) > 0 {
				for _, attachment := range msg.Attachments {
					if attachment.Type == "image/jpeg" || attachment.Type == "image/png" || attachment.Type == "image/gif" {
						messagesHTML += `<div class="attachment">
                            <img src="` + attachment.URL[:strings.Index(attachment.URL, "?")] + `" alt="` + attachment.Filename + `" />
                        </div>`
					} else {
						messagesHTML += `<div class="attachment">
                            <a href="` + attachment.URL + `" download>` + attachment.Filename + `</a>
                        </div>`
					}
				}
			}
			messagesHTML += `</div>
            </div>`
		}
	}

	tmpl, err := template.New("transcript").Parse(string(htmlTemplate))
	if err != nil {
		http.Error(w, "Error parsing HTML template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, struct{ Messages template.HTML }{Messages: template.HTML(messagesHTML)})
	if err != nil {
		http.Error(w, "Error executing HTML template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func formatTimestamp(timestamp string) string {
	parsedTime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return timestamp
	}

	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)

	if parsedTime.Year() == now.Year() && parsedTime.YearDay() == now.YearDay() {
		return "Today at " + parsedTime.Format("3:04 PM")
	} else if parsedTime.Year() == yesterday.Year() && parsedTime.YearDay() == yesterday.YearDay() {
		return "Yesterday at " + parsedTime.Format("3:04 PM")
	} else {
		return parsedTime.Format("01/02/2006 3:04 PM")
	}
}

func getTranscript(id string) ([]byte, error) {
	collection := db.GetCollection(cfg.DBName, "tickets")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objID}

	var result struct {
		Transcript []byte `bson:"transcript"`
	}

	err = collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.Transcript, nil
}
