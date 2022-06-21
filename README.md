# Slack Impersonation

Spoof visual user identity in Slack to send message

## Usage

***Requirements:***
* Slack app with a bot token with scope `users:read` included
* Incoming webhook (do not use app incoming webhook)

*If the different flags are not provided they will be asked as input*

### » to channel
Send a message in ***`#general`*** channel spoofing `Toto RINA` visual identity:
```shell
slack-spoofer -u "Toto RINA" -c "#general" -m 'Hi <!channel|channel>!\nToday I'm the one paying for the meal!' -t $(cat .credentials.json | jq -r .bot_token) -w $(cat .credentials.json | jq -r .webhook)
```

### » to user (direct message)
Send a direct message to ***`Elon MUSK`*** spoofing `Jeff BEZOS` visual identity:
```shell
slack-spoofer dm -u "Jeff BEZOS" -r "Elon MUSK" -m 'Please find all my secrets <https://maliciouscvs|here>' -t $(cat .credentials.json | jq -r .bot_token) -w $(cat .credentials.json | jq -r .webhook)
```

## Install
```shell
git clone https://github.com/ariary/SlackSpoofing
make before.build && make build.slack-spoofer
```
## Notes

The tricks is not revolutionary:
1. Use slack app to map username -> user avatar url
2. Send message with incoming webhook with Bot username = username and Bot avatar url = user avatar url

* ***Pentester idea:*** If you obtain an incoming webhook (leak, compromised etc) you can use it. Keep in mind that the scope is linked with the webhook creator (ie. webhook can publish in private channel where the creator is whitout needing an invitation)
* Step 1 is automated here but can be manual if you do not have the permission to create slack app (view user profile, right-click on avatar, "Copy avatar url" )
* Get user id (useful to mention them in message):  `slack-spoofer getid -u "[USER]" -t $(cat .credentials.json | jq -r .bot_token) -w $(cat .credentials.json | jq -r .webhook)`
* `<!channel|channel>`send a notification to all channel users and `<!here|here>` to all channel online users
