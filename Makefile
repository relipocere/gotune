.PHONY: run
run: build make-audio-folder
	./bin/bot

.PHONY: run-register
run-register: build make-audio-folder
	./bin/bot --register-commands

.PHONY:build
build:
	go build -o ./bin/bot ./cmd/bot/main.go

.PHONY:make-audio-folder
make-audio-folder:
	mkdir -p ./audio

.PHONY: get-deps-ubuntu
get-deps-ubuntu:
	sudo apt update
	sudo apt install ffmpeg
	sudo wget https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -O /usr/local/bin/yt-dlp
	sudo chmod a+rx /usr/local/bin/yt-dlp
