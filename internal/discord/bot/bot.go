package bot

import (
	"os"
	"os/signal"

	"github.com/relipocere/gotune/internal/discord/player"

	"github.com/relipocere/gotune/internal/config"

	"github.com/relipocere/gotune/internal/discord/types"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

//Bot ...
type Bot struct {
	s          *discordgo.Session
	log        *zap.SugaredLogger
	cfg        *config.Config
	extractor  types.Extractor
	dispatcher types.Dispatcher
}

//New creates new Bot.
func New(cfg *config.Config, l *zap.SugaredLogger, e types.Extractor) *Bot {
	s, err := discordgo.New("Bot " + cfg.Token())
	if err != nil {
		l.Fatal(err)
	}

	s.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildVoiceStates
	b := &Bot{
		s:          s,
		log:        l,
		cfg:        cfg,
		extractor:  e,
		dispatcher: player.NewDispatcher(s, l),
	}

	s.AddHandler(b.routeCommand)
	s.AddHandler(b.preventVoiceStateChange)
	return b
}

//RegisterSlashCommands registers bot commands.
func (b *Bot) RegisterSlashCommands() {
	var commands = []*discordgo.ApplicationCommand{
		{
			Name:        "play",
			Description: "Play a song or an album",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "song",
					Description: "name of the song or youtube link",
					Required:    true,
				},
			},
		},
		{
			Name:        "queue",
			Description: "List the song queue",
		},
		{
			Name:        "pause",
			Description: "Pause the song",
		},
		{
			Name:        "resume",
			Description: "Resume the song",
		},
		{
			Name:        "skip",
			Description: "Skip currently playing song",
		},
		{
			Name:        "skipto",
			Description: "Skip to a certain queued song",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "position",
					Description: "position of the song to skip to",
					Required:    true,
				},
			},
		},
		{
			Name:        "stop",
			Description: "Stop playing and leave",
		},
		{
			Name:        "seek",
			Description: "Play track at the specified minute and second",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "minute",
					Description: "0-...",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "second",
					Description: "0-60",
					Required:    true,
				},
			},
		},
	}

	for _, v := range commands {
		_, err := b.s.ApplicationCommandCreate(b.cfg.AppID(), "", v)
		if err != nil {
			b.log.Fatalw("Cannot create command", "err", err)
		}
	}
}

//routeCommand is the interaction command router.
func (b *Bot) routeCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	commandHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"play":   b.play,
		"queue":  b.queue,
		"pause":  b.pause,
		"resume": b.resume,
		"skip":   b.skip,
		"skipto": b.skipTo,
		"stop":   b.stop,
		"seek":   b.seek,
	}
	if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
		h(s, i)
		return
	}
}

//Serve starts the bot and blocks until termination signal is received.
func (b *Bot) Serve() {
	b.log.Warn("Bot is online")
	err := b.s.Open()
	if err != nil {
		b.log.Fatal(err)
	}
	defer b.s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	b.log.Warn("Bot is offline")
}
