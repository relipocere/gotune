# Go Tune - Discordgo music bot

## Non-GO Dependencies
[yt-dlp](https://github.com/yt-dlp/yt-dlp) is a YouTube downloader.

[ffmpeg](https://www.ffmpeg.org/) for format conversion and audio extraction.

## Features
* Multi-server support
* Song search
* Pause, resume, skip, skip to, stop, queue, and seek
* Auto-disconnect when done playing
* Cleans-up and leaves if kicked or forcefully moved to another channel
* Support for playlists
* Cache which is cleaned after the container restart

## Limits
* YouTube is the only supported platform
* Doesn't support streams
* Downloads only first 10 songs of a playlist (can be manually changed in yt-dlp options)

## Usage
Clone the repository:
```sh
https://github.com/relipocere/gotune.git
```
Rename **example_config.yml** to **config.yml** and fill it.

### Running using Docker
Build Docker image:
```sh
sudo docker build -t gotune .
```

If you're running bot for the first time you need to register commands,
so pass a --register-commands argument:
```sh
sudo docker run -d -l bot gotune --register-commands
```

Run the Docker image:
```sh
sudo docker run -d -l bot gotune
```

### Running locally on Ubuntu-based Linux
Install dependencies manually or use:
```sh
make get-deps-ubuntu 
```

If you're running bot for the first time you need to register slash commands:
```sh
make run-register
```

Run the bot:
```sh
make run
```
