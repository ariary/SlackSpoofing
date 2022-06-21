package slackspoofing

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ariary/go-utils/pkg/color"
)

//Config: configuration structure
type Config struct {
	Username   string
	Channel    string
	Message    string
	WebhookUrl string
	BotToken   string
	Recipient  string
}

//CheckConfig: check if all the necesary information has been provided
func CheckConfig(cfg *Config, dm bool) {

	if cfg.BotToken == "" {
		cfg.BotToken = WaitInput("bot token")
	}

	if cfg.WebhookUrl == "" {
		cfg.WebhookUrl = WaitInput("incoming webhook url")
	}
	if dm {
		if cfg.Recipient == "" {
			cfg.Recipient = WaitInput("recipient username")
		}
	} else {
		if cfg.Channel == "" {
			cfg.Channel = WaitInput("channel")
		}
	}

	if cfg.Username == "" {
		cfg.Username = WaitInput("username")
	}

	if cfg.Message == "" {
		cfg.Message = WaitInput("message")
	}
	cfg.Message = strings.ReplaceAll(cfg.Message, "\\n", "\n") //correct new line in message
}

//WaitInput: wait user input and return it. If no value is typed, it retries asking for input
func WaitInput(name string) (input string) {
	msg := color.Blue("»") + " " + name + ":"
	msg += " "
	fmt.Printf(msg)
	inputReader := bufio.NewReader(os.Stdin)
	input, err := inputReader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return WaitInput(name)
	}
	if input == "" {
		return WaitInput(name)
	} else {
		return input[:len(input)-1]
	}
}
