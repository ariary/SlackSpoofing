# Slack Impersonation

Spoof visual user identity in Slack to send message

## Usage

***Requirements:***
* Slack app with a bot token with scope `users:read` included
* Incoming webhook (do not use app incoming webhook)

Send a message in ***`#general`*** channel spoofing `Toto RINA` visual identity:
```shell
./slackspoof -u "Toto RINA" -c "#general" -m "Hi all!\nToday I'm the one paying for the meal!" -t $(cat .credentials.json | jq .bot_token) -w $(cat .credentials.json | jq .webhook)
```
If the different flags are not provided they will be asked as input

## Notes

The tricks is not revolutionary:
1. Use slack app to map username -> user avatar url
2. Send message with incoming webhook with Bot username = username and Bot avatar url = user avatar url

* If you obtain a incoming webhook (leak,compromised etc) you can use it. keep in mind that the scope is linked with the webhook creator (ie. webhook can publish in private channel where the creator is whitout needing an invitation)
* Step 1 is automated here but can be manual if you do not have the permission to create slack app (view user profile, right-click on avatar, "Copy avatar url" )