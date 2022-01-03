package main

import (
	"log"
	"os"

	"github.com/relipocere/gotune/internal/config"
	"github.com/relipocere/gotune/internal/discord/bot"
	l "github.com/relipocere/gotune/internal/logger"
	"github.com/relipocere/gotune/internal/yt"
)

func main() {

	cfg, err := config.New("config", "yml", ".")
	if err != nil {
		log.Fatal(err)
	}

	logger, err := l.NewLogger(l.SetLevel(cfg.LogLevel()))
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	e, err := yt.New(cfg.YtToken(), cfg.FileDirectory())
	if err != nil {
		logger.Fatal(err)
	}

	b := bot.New(cfg, logger, e)

	if len(os.Args) > 1 {
		args := os.Args[1:]
		for _, arg := range args {
			switch arg {
			case "--register-commands":
				b.RegisterSlashCommands()
			default:
				logger.Fatal("invalid argument")
			}
		}
	}
	b.Serve()
}
