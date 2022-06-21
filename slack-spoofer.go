package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/ariary/go-utils/pkg/color"
	"github.com/spf13/cobra"
)

type Config struct {
	Username   string
	Channel    string
	Message    string
	WebhookUrl string
	BotToken   string
}

type Profile struct {
	Image string `json:"image_192"`
}

type Member struct {
	Id       string  `json:"id"`
	RealName string  `json:"real_name"`
	Profile  Profile `json:"profile"`
}

type Result struct {
	Ok      bool     `json:"ok"`
	Members []Member `json:"members"`
}

type Message struct {
	Username string `json:"username"`
	Channel  string `json:"channel"`
	IconUrl  string `json:"icon_url"`
	Text     string `json:"text"`
}

func PostMessage(url string, message Message) error {

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	client := &http.Client{}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	err = response.Body.Close()
	if err != nil {
		return err
	}

	return nil
}

func checkConfig(cfg *Config) {

	if cfg.BotToken == "" {
		cfg.BotToken = waitInput("bot token")
	}

	if cfg.WebhookUrl == "" {
		cfg.WebhookUrl = waitInput("incoming webhook url")
	}

	if cfg.Channel == "" {
		cfg.Channel = waitInput("channel")
	}

	if cfg.Username == "" {
		cfg.Username = waitInput("username")
	}

	if cfg.Message == "" {
		cfg.Message = waitInput("message")
	}
}

//waitInput: wait user input and return it. If no value is typed, it retries asking for input
func waitInput(name string) (input string) {
	msg := color.Blue("»") + " " + name + ":"
	msg += " "
	fmt.Printf(msg)
	inputReader := bufio.NewReader(os.Stdin)
	input, err := inputReader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return waitInput(name)
	}
	if input == "" {
		return waitInput(name)
	} else {
		return input[:len(input)-1]
	}
}

const listEndpoint = "https://slack.com/api/users.list"

func main() {
	//CMD ROOT
	var cfg Config
	var rootCmd = &cobra.Command{Use: "slackctl",
		Short: "impersonificate user to send messages",
		Run: func(cmd *cobra.Command, args []string) {
			// Init
			checkConfig(&cfg)

			//retrieve avatar
			client := &http.Client{}
			req, _ := http.NewRequest("GET", listEndpoint, nil)
			req.Header.Set("Authorization", "Bearer "+cfg.BotToken)
			resp, _ := client.Do(req)
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}

			var userListsResult Result
			err = json.Unmarshal(body, &userListsResult)
			if err != nil {
				fmt.Println(err)
				os.Exit(0)
			}

			var avatarUrl string
			for i := 0; i < len(userListsResult.Members); i++ {
				if userListsResult.Members[i].RealName == cfg.Username {
					avatarUrl = userListsResult.Members[i].Profile.Image
					message := Message{Username: cfg.Username, Channel: cfg.Channel, IconUrl: avatarUrl, Text: cfg.Message}
					if err := PostMessage(cfg.WebhookUrl, message); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					break
				}
			}
			if avatarUrl == "" {
				fmt.Println(color.RedForeground("Failed to retrieve user: " + cfg.Username))
			}

		},
	}
	rootCmd.Flags().StringVarP(&cfg.Channel, "channel", "c", "", "specify channel")
	rootCmd.Flags().StringVarP(&cfg.Username, "username", "u", "", "specify username to impersonate")
	rootCmd.Flags().StringVarP(&cfg.BotToken, "token", "t", "", "specify bot token with users.identities scope")
	rootCmd.Flags().StringVarP(&cfg.WebhookUrl, "webhook", "w", "", "specify incoming webhook used to send message")
	rootCmd.Flags().StringVarP(&cfg.Message, "message", "m", "", "specify the message to send")

	rootCmd.Execute()

}
