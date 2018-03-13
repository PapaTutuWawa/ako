package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Polynomdivision/ako/trello"
	"github.com/fatih/color"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

// Comment
func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}

	return false
}

func printCardPreview(card trello.TrelloCard, list_name string, labels map[string]trello.TrelloLabel) {
	used_labels := make([]string, 0)
	for _, l := range card.Labels {
		label, ok := labels[l]
		if !ok {
			used_labels = append(used_labels, "ERROR")
			continue
		}

		used_labels = append(used_labels, label.Format())
	}

	red := color.New(color.FgRed).SprintFunc()
	fmt.Printf("[%s] %s: %s (%s)\n", card.Id, list_name, red(card.Name), strings.Join(used_labels, ","))
}

func printBoardPreview(board trello.TrelloBoard) {
	red := color.New(color.FgRed).SprintFunc()
	fmt.Printf("[%s] %s\n", board.Id, red(board.Name))
}

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.config/ako")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	user := trello.TrelloUser{
		ApiKey:   viper.GetString("key"),
		ApiToken: viper.GetString("token"),
	}

	if user.ApiKey == "" || user.ApiToken == "" {
		fmt.Println("API Key or Token not found")
		return
	}

	// Load in aliases
	aliases := viper.GetStringMap("aliases")

	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:  "boards",
			Usage: "List the boards",
			Action: func(c *cli.Context) error {
				boards, err := user.GetBoards()
				if err != nil {
					return err
				}

				for _, board := range boards {
					printBoardPreview(board)
				}

				return nil
			},
		},
		{
			Name:  "cards",
			Usage: "Get the cards from a board",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "id",
					Usage: "The ID or alias of the board",
				},
			},
			Action: func(c *cli.Context) error {
				id := c.String("id")

				if id == "" {
					fmt.Println("Error: Empty ID or alias")
					// TODO: No
					return nil
				}

				// Check if it is an alias
				aliasId, ok := aliases[id]
				if ok {
					id = aliasId.(string)
				}

				board, err := user.GetBoard(id)
				if err != nil {
					panic(err)
				}

				cards, err := board.GetCards(user)
				if err != nil {
					panic(err)
				}
				lists, err := board.GetLists(user)
				if err != nil {
					panic(err)
				}
				labels, err := board.GetLabels(user)
				if err != nil {
					panic(err)
				}

				for _, card := range cards {
					list_name, ok := lists[card.IdList]
					if !ok {
						list_name = ""
					}

					printCardPreview(card, list_name, labels)
				}

				return nil
			},
		},
		{
			Name:  "self",
			Usage: "Get the cards from a board that are assigned to me",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "id",
					Usage: "The ID or alias of the board",
				},
				cli.BoolFlag{
					Name:  "all",
					Usage: "List all cards from all boards that are assigned to me",
				},
			},
			Action: func(c *cli.Context) error {
				id := c.String("id")

				if id == "" && c.Bool("all") == false {
					fmt.Println("Error: Empty ID or alias")
					// TODO: No
					return nil
				}

				// Check if it is an alias
				aliasId, ok := aliases[id]
				if ok {
					id = aliasId.(string)
				}

				selfId, err := user.GetUserId()
				if err != nil {
					panic(err)
				}

				var boards []trello.TrelloBoard
				if c.Bool("all") {
					boards, err = user.GetBoards()
					if err != nil {
						panic(err)
					}
				} else {
					board, err := user.GetBoard(id)
					if err != nil {
						panic(err)
					}

					boards = []trello.TrelloBoard{board}
				}

				for _, board := range boards {
					cards, err := board.GetCards(user)
					if err != nil {
						panic(err)
					}
					lists, err := board.GetLists(user)
					if err != nil {
						panic(err)
					}
					labels, err := board.GetLabels(user)
					if err != nil {
						panic(err)
					}

					fmt.Printf("%s\n", color.New(color.FgRed).SprintFunc()(board.Name))
					for _, card := range cards {
						if !contains(card.Users, selfId) {
							continue
						}

						list_name, ok := lists[card.IdList]
						if !ok {
							list_name = ""
						}

						printCardPreview(card, list_name, labels)
					}
				}

				return nil
			},
		},
		{
			Name:  "card",
			Usage: "View a card",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "id",
					Usage: "The ID or alias of the card",
				},
			},
			Action: func(c *cli.Context) error {
				id := c.String("id")

				if id == "" {
					fmt.Println("Error: Empty ID or alias")
					// TODO: No
					return nil
				}

				// Check if it is an alias
				aliasId, ok := aliases[id]
				if ok {
					id = aliasId.(string)
				}

				card, err := user.GetCard(id)
				if err != nil {
					panic(err)
				}

				usernames := make([]string, 0)
				for _, userId := range card.Users {
					name, err := user.GetUsernameFromId(userId)
					if err != nil {
						name = "<ERROR>"
					}

					usernames = append(usernames, name)
				}

				red := color.New(color.FgRed).SprintFunc()
				fmt.Printf("[%s]\n", red(card.Name))
				fmt.Printf("Users: %s\n", strings.Join(usernames, ", "))
				fmt.Println("")
				fmt.Println(card.Desc)

				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
