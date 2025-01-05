package games

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type TicTacToeGame struct {
	Board     [3][3]string
	Player1ID string
	Player2ID string
	CurrentID string
	GameOver  bool
}

var activeTicTacToeGames = make(map[string]*TicTacToeGame)

func CreateTicTacToeCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	opponent := i.ApplicationCommandData().Options[0].UserValue(s)

	if opponent.ID == i.Member.User.ID || opponent.Bot {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You cannot play against yourself or a bot!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	gameID := i.ChannelID
	if _, exists := activeTicTacToeGames[gameID]; exists {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "A game is already in progress in this channel!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	game := &TicTacToeGame{
		Player1ID: i.Member.User.ID,
		Player2ID: opponent.ID,
		CurrentID: i.Member.User.ID,
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("<@%s> vs <@%s>\nIt's <@%s>'s turn!",
				game.Player1ID, game.Player2ID, game.CurrentID),
			Components: createTicTacToeButtons(game),
		},
	})

	if err != nil {
		log.Printf("Failed to create game: %v", err)
		return
	}

	activeTicTacToeGames[gameID] = game
}

func HandleTicTacToeButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	gameID := i.ChannelID
	game, exists := activeTicTacToeGames[gameID]
	if !exists {
		log.Printf("Button press for non-existent game: channel=%s user=%s",
			gameID, i.Member.User.ID)
		return
	}

	if i.Member.User.ID != game.CurrentID {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "It's not your turn!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	var row, col int
	fmt.Sscanf(i.MessageComponentData().CustomID, "tictactoe_%d_%d", &row, &col)

	if game.Board[row][col] != "" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "That space is already taken!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	symbol := "X"
	if game.CurrentID == game.Player2ID {
		symbol = "O"
	}
	game.Board[row][col] = symbol

	winner := checkWin(game.Board)
	isDraw := checkDraw(game.Board)

	content := ""
	if winner != "" {
		game.GameOver = true
		content = fmt.Sprintf("Game Over! <@%s> wins!", game.CurrentID)
	} else if isDraw {
		game.GameOver = true
		content = "Game Over! It's a draw!"
	} else {
		if game.CurrentID == game.Player1ID {
			game.CurrentID = game.Player2ID
		} else {
			game.CurrentID = game.Player1ID
		}
		content = fmt.Sprintf("It's <@%s>'s turn!", game.CurrentID)
	}

	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content:    content,
			Components: createTicTacToeButtons(game),
		},
	}

	if err := s.InteractionRespond(i.Interaction, response); err != nil {
		log.Printf("Failed to update game state: channel=%s user=%s error=%v",
			gameID, i.Member.User.ID, err)
	}

	if game.GameOver {
		log.Printf("Game completed: channel=%s winner=%s",
			gameID, game.CurrentID)
		delete(activeTicTacToeGames, gameID)
	}
}

func createTicTacToeButtons(game *TicTacToeGame) []discordgo.MessageComponent {
	var components []discordgo.MessageComponent

	for row := 0; row < 3; row++ {
		var buttons []discordgo.MessageComponent

		for col := 0; col < 3; col++ {
			label := "⬜"

			if game.Board[row][col] == "X" {
				label = "❌"
			} else if game.Board[row][col] == "O" {
				label = "⭕"
			}

			buttons = append(buttons, discordgo.Button{
				Label:    label,
				CustomID: fmt.Sprintf("tictactoe_%d_%d", row, col),
				Style:    discordgo.PrimaryButton,
				Disabled: game.GameOver || game.Board[row][col] != "",
			})
		}

		components = append(components, discordgo.ActionsRow{
			Components: buttons,
		})
	}

	return components
}

func checkWin(board [3][3]string) string {
	for i := 0; i < 3; i++ {
		if board[i][0] != "" && board[i][0] == board[i][1] && board[i][1] == board[i][2] {
			return board[i][0]
		}
	}
	for j := 0; j < 3; j++ {
		if board[0][j] != "" && board[0][j] == board[1][j] && board[1][j] == board[2][j] {
			return board[0][j]
		}
	}
	if board[0][0] != "" && board[0][0] == board[1][1] && board[1][1] == board[2][2] {
		return board[0][0]
	}
	if board[0][2] != "" && board[0][2] == board[1][1] && board[1][1] == board[2][0] {
		return board[0][2]
	}
	return ""
}

func checkDraw(board [3][3]string) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == "" {
				return false
			}
		}
	}
	return true
}
