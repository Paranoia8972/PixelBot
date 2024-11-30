package transcript

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Paranoia8972/PixelBot/internal/app/config"
	"github.com/Paranoia8972/PixelBot/internal/db"
	"github.com/fatih/color"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Cfg *config.Config

func init() {
	Cfg = config.LoadConfig()
}

type Attachment struct {
	Filename string `json:"filename"`
	Type     string `json:"type"`
	URL      string `json:"url"`
}

type Reactions struct {
	Emoji string `json:"emoji"`
	Count int    `json:"count"`
}

type Embed struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Fields      []EmbedField `json:"fields"`
	URL         string       `json:"url"`
	Color       int          `json:"color"`
	Image       string       `json:"image.url"`
}

type EmbedField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type TranscriptMessage struct {
	Attachments    []Attachment `json:"attachments"`
	MessageContent string       `json:"message_content"`
	Pfp            string       `json:"pfp"`
	Timestamp      string       `json:"timestamp"`
	Username       string       `json:"username"`
	Embeds         []Embed      `json:"embeds"`
	Reactions      []Reactions  `json:"reactions"`
}

type TranscriptData struct {
	Transcript []TranscriptMessage `json:"transcript"`
}

func StartTranscriptServer() {
	http.HandleFunc("/ticket", TranscriptServer)
	http.Handle("/downloads/", http.StripPrefix("/downloads/", http.FileServer(http.Dir("downloads"))))
	color.Green("Transcript server is running on http://localhost:" + Cfg.Port + "/ticket | https://" + Cfg.TranscriptUrl + "/ticket")
	log.Fatal(http.ListenAndServe(":"+Cfg.Port, nil))
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

		if msg.Username != lastUsername {
			messagesHTML += `<div class="message">
				<img src="` + msg.Pfp + `" class="pfp" />
				<div class="content">
					<div>
						<span class="username">` + msg.Username + `</span>
						<span class="timestamp">` + formattedTime + `</span>
					</div>
					<div class="message_content">` + msg.MessageContent + `</div>`
			if len(msg.Embeds) > 0 {
				for _, embed := range msg.Embeds {
					messagesHTML += `<div class="embed" style="border-left: 4px solid #` + fmt.Sprintf("%06x", embed.Color) + `;">
						<div class="embed-title">` + embed.Title + `</div>
						<div class="embed-description">` + embed.Description + `</div>`
					if len(embed.Fields) > 0 {
						messagesHTML += `<div class="embed-fields">`
						for _, field := range embed.Fields {
							messagesHTML += `<div class="embed-field">
								<div class="embed-field-name">` + field.Name + `:</div>
								<div class="embed-field-value">` + field.Value + `</div>
							</div>`
						}
						messagesHTML += `</div>`
					}
					messagesHTML += `</div>`
				}
			}
			if len(msg.Attachments) > 0 {
				for _, attachment := range msg.Attachments {
					if strings.HasPrefix(attachment.Type, "image/") {
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
			if len(msg.Reactions) > 0 {
				messagesHTML += `<div class="reactions">`
				for _, reaction := range msg.Reactions {
					messagesHTML += `<span class="reaction">` + reaction.Emoji + ` ` + strconv.Itoa(reaction.Count) + `</span>`
				}
				messagesHTML += `</div>`
			}
			messagesHTML += `</div>
			</div>`
			lastUsername = msg.Username
		} else {
			messagesHTML += `<div class="message no-padding">
				<div class="content">
					<div class="message_content">` + msg.MessageContent + `</div>`
			if len(msg.Embeds) > 0 {
				for _, embed := range msg.Embeds {
					messagesHTML += `<div class="embed" style="border-left: 4px solid #` + fmt.Sprintf("%06x", embed.Color) + `;">
						<div class="embed-title">` + embed.Title + `</div>
						<div class="embed-description">` + embed.Description + `</div>`
					if len(embed.Fields) > 0 {
						for _, field := range embed.Fields {
							messagesHTML += `<div class="embed-field">
								<div class="embed-field-name">` + field.Name + `</div>
								<div class="embed-field-value">` + field.Value + `</div>
							</div>`
						}
					}
					if embed.Image != "" {
						messagesHTML += `<div class="embed-image">
							<img src="` + embed.Image + `" alt="embed image" />
						</div>`
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
			if len(msg.Reactions) > 0 {
				messagesHTML += `<div class="reactions" style="margin-top: 5px;">`
				for _, reaction := range msg.Reactions {
					messagesHTML += `<span class="reaction">` + reaction.Emoji + ` ` + strconv.Itoa(reaction.Count) + `</span>`
				}
				messagesHTML += `</div>`
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
	collection := db.GetCollection(Cfg.DBName, "tickets")
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
