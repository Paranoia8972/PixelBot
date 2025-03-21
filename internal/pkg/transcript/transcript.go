package transcript

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
	ChannelID      string       `json:"channel_id"`
	MessageID      string       `json:"message_id"`
}

type TranscriptData struct {
	Transcript []TranscriptMessage `json:"transcript"`
}

type Message struct {
	Username       string       `json:"username"`
	Pfp            string       `json:"pfp"`
	MessageContent string       `json:"message_content"`
	Timestamp      int64        `json:"timestamp"`
	Embeds         []Embed      `json:"embeds"`
	Attachments    []Attachment `json:"attachments"`
	Reactions      []Reactions  `json:"reactions"`
	ChannelID      string       `json:"channel_id"`
	MessageID      string       `json:"message_id"`
}

func StartTranscriptServer() {
	http.HandleFunc("/ticket", TranscriptServer)
	http.HandleFunc("/attachments/", ServeAttachment)
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
					savedFilename, err := downloadAttachment(attachment.URL, msg.ChannelID, msg.MessageID, attachment.Filename)
					if err != nil {
						log.Printf("Error downloading attachment: %v", err)
						continue
					}

					if strings.HasPrefix(attachment.Type, "image/") {
						messagesHTML += `<div class="attachment">
							<img src="/attachments/` + savedFilename + `" alt="` + attachment.Filename + `" />
						</div>`
					} else {
						messagesHTML += `<div class="attachment">
							<a href="/attachments/` + savedFilename + `" download>` + attachment.Filename + `</a>
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
					savedFilename, err := downloadAttachment(attachment.URL, msg.ChannelID, msg.MessageID, attachment.Filename)
					if err != nil {
						log.Printf("Error downloading attachment: %v", err)
						continue
					}

					if attachment.Type == "image/jpeg" || attachment.Type == "image/png" || attachment.Type == "image/gif" {
						messagesHTML += `<div class="attachment">
							<img src="/attachments/` + savedFilename + `" alt="` + attachment.Filename + `" />
						</div>`
					} else {
						messagesHTML += `<div class="attachment">
							<a href="/attachments/` + savedFilename + `" download>` + attachment.Filename + `</a>
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

func downloadAttachment(url, channelID, messageID, filename string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Create attachments directory if it doesn't exist
	if err := os.MkdirAll("attachments", 0755); err != nil {
		return "", err
	}

	// Create a unique filename using channel and message IDs
	savedFilename := fmt.Sprintf("%s-%s-%s", channelID, messageID, filename)
	filePath := filepath.Join("attachments", savedFilename)

	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	return savedFilename, nil
}

func ServeAttachment(w http.ResponseWriter, r *http.Request) {
	// Extract filename from URL path
	filename := strings.TrimPrefix(r.URL.Path, "/attachments/")
	if filename == "" {
		http.Error(w, "Missing filename", http.StatusBadRequest)
		return
	}

	// Ensure the path is safe and within attachments directory
	filePath := filepath.Join("attachments", filepath.Clean(filename))
	if !strings.HasPrefix(filePath, filepath.Clean("attachments/")) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	// Serve the file
	http.ServeFile(w, r, filePath)
}
