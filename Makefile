before.build:
	go mod tidy && go mod download

build.slack-spoofer:
	@echo "build in ${PWD}";go build slack-spoofer.go