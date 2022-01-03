package player

import (
	"fmt"
	"io"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"github.com/relipocere/gotune/internal/discord/types"
	"go.uber.org/zap"
)

const (
	msgNoPlayer            = "Bot is not playing"
	msgUnexpectedError     = "Unexpected error"
	msgDone                = "👍"
	errPlayerNotResponding = "player is not reading command"
)

//Dispatcher is the player manager that routes songs to the correct guild player.
type Dispatcher struct {
	s       *discordgo.Session
	log     *zap.SugaredLogger
	players *playerMap
}

//command is the message type for player.command channel.
type command struct {
	//Valid actions are: skip, seek, stop, pause, resume
	Action string

	//Seek time in seconds
	SeekTime int
}

//NewDispatcher creates new player dispatcher.
func NewDispatcher(s *discordgo.Session, log *zap.SugaredLogger) *Dispatcher {
	return &Dispatcher{
		s:       s,
		log:     log,
		players: newPlayerMap(),
	}
}

//Play adds songs to the queue of the guild player.
//If player doesn't exist, new one is created and launched in goroutine.
func (d *Dispatcher) Play(gID, vID, cmdID string, songs []types.Song) {
	p, exists := d.players.Load(gID)
	if !exists {
		p = player{
			Command: make(chan command),
			Queue:   newQueue(),
		}
		d.players.Store(gID, p)
	}

	p.Queue.Push(songs)
	if !exists {
		go d.dispatchPlayer(gID, vID, cmdID)
	}
}

//Queue returns titles of queued songs.
func (d *Dispatcher) Queue(gID string) []types.Song {
	p, ok := d.players.Load(gID)
	if !ok {
		return nil
	}

	return p.Queue.ListSongs()
}

//Seek skips song playback to the desired time.
//seekTime is measured in seconds.
func (d *Dispatcher) Seek(gID string, seekTime int) (string, error) {
	p, ok := d.players.Load(gID)
	if !ok {
		return msgNoPlayer, nil
	}

	select {
	case p.Command <- command{Action: "seek", SeekTime: seekTime}:
	case <-time.After(5 * time.Second):
		return msgUnexpectedError, fmt.Errorf(errPlayerNotResponding)
	}

	return msgDone, nil
}

//SkipTo skips to the specified queue position.
func (d *Dispatcher) SkipTo(gID string, pos int) (string, error) {
	p, ok := d.players.Load(gID)
	if !ok {
		return msgNoPlayer, nil
	}

	qLen := p.Queue.Len()
	if qLen < pos {
		return fmt.Sprintf("There are only %d songs in the queue", qLen), nil
	}

	for i := 0; i < pos-1; i++ {
		p.Queue.Pop()
	}

	return d.Skip(gID)
}

//Stop stops music stream and discards the queue.
func (d *Dispatcher) Stop(gID string) (string, error) {
	p, ok := d.players.Load(gID)
	if !ok {
		return msgNoPlayer, nil
	}
	select {
	case p.Command <- command{Action: "stop"}:
	case <-time.After(5 * time.Second):
		return msgUnexpectedError, fmt.Errorf(errPlayerNotResponding)
	}

	return msgDone, nil
}

//Skip skips currently playing track.
func (d *Dispatcher) Skip(gID string) (string, error) {
	p, ok := d.players.Load(gID)
	if !ok {
		return msgNoPlayer, nil
	}
	select {
	case p.Command <- command{Action: "skip"}:
	case <-time.After(5 * time.Second):
		return msgUnexpectedError, fmt.Errorf(errPlayerNotResponding)
	}

	return msgDone, nil
}

//Pause pauses currently playing track.
func (d *Dispatcher) Pause(gID string) (string, error) {
	p, ok := d.players.Load(gID)
	if !ok {
		return msgNoPlayer, nil
	}
	select {
	case p.Command <- command{Action: "pause"}:
	case <-time.After(5 * time.Second):
		return msgUnexpectedError, fmt.Errorf(errPlayerNotResponding)
	}

	return msgDone, nil
}

//Resume resumes track that was playing.
func (d *Dispatcher) Resume(gID string) (string, error) {
	p, ok := d.players.Load(gID)
	if !ok {
		return msgNoPlayer, nil
	}
	select {
	case p.Command <- command{Action: "resume"}:
	case <-time.After(5 * time.Second):
		return msgUnexpectedError, fmt.Errorf(errPlayerNotResponding)
	}

	return msgDone, nil
}

//dispatchPlayer creates new player for the guild.
func (d *Dispatcher) dispatchPlayer(gID, vID, cmdID string) {
	p, _ := d.players.Load(gID)
	defer d.players.Delete(gID)

	vc, err := d.s.ChannelVoiceJoin(gID, vID, false, true)
	if err != nil {
		d.s.ChannelMessageSendEmbed(cmdID, types.ErrorEmbed("Can't join the channel"))
		d.log.Errorw(err.Error(), "guildID", gID, "voiceID", vID)
		return
	}
	defer vc.Disconnect()

	for {
		if p.Queue.Len() < 1 {
			d.log.Debugw("done playing", "guildID", gID)
			return
		}

		song := p.Queue.Pop()
		d.s.ChannelMessageSendEmbed(cmdID, types.TrackEmbed("Now playing", song))
		d.log.Debugw("playing", "guildID", gID, "song", song)

		stop, err := encodeAndPlay(vc, song.Path, p.Command)
		if err != nil {
			d.s.ChannelMessageSendEmbed(cmdID, types.ErrorEmbed("Unable to play the song"))
			d.log.Errorw(fmt.Sprintf("encodeAndPlay: %s", err.Error()), "path", song.Path)
		}
		if stop {
			return
		}
	}
	d.log.Debugw("Deleting guild player", "guildID", gID)
}

//encodeAndPlay encodes the file into a dca session and plays it.
//If stop signal was sent on command channel stop flag will be true.
func encodeAndPlay(vc *discordgo.VoiceConnection, path string, command <-chan command) (stop bool, rErr error) {
	opts := dca.StdEncodeOptions
	opts.StartTime = 0

	encodeSession, err := dca.EncodeFile(path, opts)
	if err != nil {
		rErr = err
		return
	}
	defer encodeSession.Cleanup()

	vc.Speaking(true)
	defer vc.Speaking(false)
	for {
		frame, err := encodeSession.OpusFrame()
		if err != nil {
			if err != io.EOF {
				rErr = err
			}
			return
		}

		select {
		case vc.OpusSend <- frame:
		case <-time.After(5 * time.Second):
			rErr = fmt.Errorf("connection is broken, unable to send a frame for more than 1 second")
			return

		case cmd := <-command:
		pauseLoop:
			for {
				switch cmd.Action {
				case "stop":
					stop = true
					return
				case "skip":
					return
				case "resume":
					//Continue playing song
					break pauseLoop
				case "pause":
					//Stay in the loop, wait for the next command
					cmd = <-command
				case "seek":
					//Re-encode to recover sent frames
					encodeSession.Cleanup()
					opts.StartTime = cmd.SeekTime
					encodeSession, err = dca.EncodeFile(path, opts)
					if err != nil {
						rErr = err
						return
					}
					break pauseLoop
				}
			}
		}
	}
}
