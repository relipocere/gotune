# builder image
FROM golang:1.17.5-alpine as builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/music_bot cmd/bot/main.go


# generate clean, final image
FROM ubuntu:20.04
RUN apt-get -y update && \
apt-get -y install python3 && \
apt-get -y install ffmpeg && \
apt-get -y install wget && \
wget https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -O /usr/local/bin/yt-dlp && \
chmod a+rx /usr/local/bin/yt-dlp

WORKDIR /app
RUN mkdir /audio
COPY --from=builder /build/bin/music_bot .
COPY --from=builder /build/config.yml .
ENTRYPOINT ["./music_bot"]
CMD []