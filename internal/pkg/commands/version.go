package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
)

const (
	githubAPIURL = "https://api.github.com/repos/Paranoia8972/PixelBot/actions/runs?status=success&per_page=1"
)

type WorkflowRunsResponse struct {
	WorkflowRuns []struct {
		HeadSHA string `json:"head_sha"`
	} `json:"workflow_runs"`
}

func getLatestGitHash() (string, error) {
	req, err := http.NewRequest("GET", githubAPIURL, nil)
	if err != nil {
		return "", err
	}

	// Set the GitHub token if available
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch workflow runs: %s", resp.Status)
	}

	var result WorkflowRunsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.WorkflowRuns) == 0 {
		return "unknown", nil
	}

	return result.WorkflowRuns[0].HeadSHA, nil
}

func VersionCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	gitHash, err := getLatestGitHash()
	if err != nil {
		gitHash = "unknown"
	}

	embed := &discordgo.MessageEmbed{
		Title:       "PixelBot Version",
		Description: "PixelBot v1.0.0",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Git Hash",
				Value: gitHash,
			},
		},
		Color: 0x00ff00, // Green color
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}
