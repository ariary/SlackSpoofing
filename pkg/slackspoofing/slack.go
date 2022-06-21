package slackspoofing

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const listEndpoint = "https://slack.com/api/users.list"

type Profile struct {
	Image string `json:"image_192"`
}

type Member struct {
	Id       string  `json:"id"`
	RealName string  `json:"real_name"`
	Profile  Profile `json:"profile"`
}

//Result of the list user slack api
type Result struct {
	Ok      bool     `json:"ok"`
	Members []Member `json:"members"`
}

//Struct representing a message
type Message struct {
	Username string `json:"username"`
	Channel  string `json:"channel"`
	IconUrl  string `json:"icon_url"`
	Text     string `json:"text"`
}

//PostMessage: send a message using slack API and an incoming webhook (url parameter)
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

//GetUserList: retieve worskpace users using slack API
func GetUsersList(cfg Config) (usersList Result, err error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", listEndpoint, nil)
	req.Header.Set("Authorization", "Bearer "+cfg.BotToken)
	resp, _ := client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return usersList, err
	}

	var userListsResult Result
	err = json.Unmarshal(body, &userListsResult)
	if err != nil {
		return usersList, err
	}

	return userListsResult, nil
}
