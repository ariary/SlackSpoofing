before.build:
	go mod tidy && go mod download

build.slackoff:
	@echo "build in ${PWD}";go build slack-spoofer.go