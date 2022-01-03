package bot

import (
	"fmt"

	"github.com/relipocere/gotune/internal/discord/types"

	"github.com/bwmarrin/discordgo"
)

const (
	msgInternalErr = "Internal error"
)

//play is the handler for play command.
func (b *Bot) play(s *discordgo.Session, i *discordgo.InteractionCreate) {
	vID := findVoiceChannel(s, i.Member.User.ID)
	if vID == "" {
		s.InteractionRespond(i.Interaction, types.TextInteractionResp("You must be in a voice channel"))
		return
	}
	s.InteractionRespond(i.Interaction, types.TextInteractionResp("Retrieving the tunes 🎶 (this may take a while)"))

	req := i.ApplicationCommandData().Options[0].StringValue()
	songs, err := b.extractor.Get(req)
	if err != nil {
		s.ChannelMessageSendEmbed(i.ChannelID, types.ErrorEmbed("Unable to download song(s)"))
		b.log.Errorw(fmt.Sprintf("extractor: %s", err.Error()), "query", req)
		return
	}

	addedMsg := "Song was added to the queue"
	if len(songs) > 1 {
		addedMsg = fmt.Sprintf("%d songs were added to the queue", len(songs))
	}
	s.ChannelMessageSend(i.ChannelID, addedMsg)

	for ind := range songs {
		songs[ind].Requester = i.Member.User
	}

	b.dispatcher.Play(i.GuildID, vID, i.ChannelID, songs)
}

//findVoiceChannel attempts to find in what voice channel user is in.
func findVoiceChannel(s *discordgo.Session, uID string) string {
	var voiceID string
	for _, guild := range s.State.Guilds {
		for _, vs := range guild.VoiceStates {
			if vs.UserID == uID {
				voiceID = vs.ChannelID
			}
		}
	}
	return voiceID
}

//seek is the handler for seek command.
func (b *Bot) seek(s *discordgo.Session, i *discordgo.InteractionCreate) {
	min := int(i.ApplicationCommandData().Options[0].IntValue())
	sec := int(i.ApplicationCommandData().Options[1].IntValue())

	if !validTime(min, sec) {
		s.InteractionRespond(i.Interaction, types.TextInteractionResp("Invalid seek time"))
		return
	}

	timeSec := (min * 60) + sec
	msg, err := b.dispatcher.Seek(i.GuildID, timeSec)
	if err != nil {
		s.InteractionRespond(i.Interaction, types.EmbedInteractionResp(types.ErrorEmbed(msgInternalErr)))
		b.log.Errorw(err.Error(), "guildID", i.GuildID)
		return
	}
	s.InteractionRespond(i.Interaction, types.TextInteractionResp(msg))
}

//validTime checks wheter seek time is valid.
func validTime(min, sec int) bool {
	if min < 0 {
		return false
	}

	if sec < 0 || sec > 60 {
		return false
	}

	return true
}

//queue is the handler for queue command.
func (b *Bot) queue(s *discordgo.Session, i *discordgo.InteractionCreate) {
	titles := b.dispatcher.Queue(i.GuildID)
	err := s.InteractionRespond(i.Interaction, types.EmbedInteractionResp(types.QueueEmbed(titles)))
	if err != nil {
		b.log.Errorw(err.Error(), "guildID", i.GuildID, "msgLength", len(titles))
	}
}

//pause is the handler for pause command.
func (b *Bot) pause(s *discordgo.Session, i *discordgo.InteractionCreate) {
	msg, err := b.dispatcher.Pause(i.GuildID)
	if err != nil {
		s.InteractionRespond(i.Interaction, types.EmbedInteractionResp(types.ErrorEmbed(msgInternalErr)))
		b.log.Errorw(err.Error(), "guildID", i.GuildID)
		return
	}
	s.InteractionRespond(i.Interaction, types.TextInteractionResp(msg))
}

//resume is the handler for resume command.
func (b *Bot) resume(s *discordgo.Session, i *discordgo.InteractionCreate) {
	msg, err := b.dispatcher.Resume(i.GuildID)
	if err != nil {
		s.InteractionRespond(i.Interaction, types.EmbedInteractionResp(types.ErrorEmbed(msgInternalErr)))
		b.log.Errorw(err.Error(), "guildID", i.GuildID)
		return
	}
	s.InteractionRespond(i.Interaction, types.TextInteractionResp(msg))
}

//skip is the handler for skip command.
func (b *Bot) skip(s *discordgo.Session, i *discordgo.InteractionCreate) {
	msg, err := b.dispatcher.Skip(i.GuildID)
	if err != nil {
		s.InteractionRespond(i.Interaction, types.EmbedInteractionResp(types.ErrorEmbed(msgInternalErr)))
		b.log.Errorw(err.Error(), "guildID", i.GuildID)
		return
	}
	s.InteractionRespond(i.Interaction, types.TextInteractionResp(msg))
}

//skipTo is the handler for skipto command.
func (b *Bot) skipTo(s *discordgo.Session, i *discordgo.InteractionCreate) {
	b.log.Debugw("skipTo is called", "guildID", i.GuildID, "userID", i.Member.User.ID)

	pos := int(i.ApplicationCommandData().Options[0].IntValue())
	if !validQueuePosition(pos) {
		s.InteractionRespond(i.Interaction, types.TextInteractionResp("Invalid queue position"))
		return
	}

	msg, err := b.dispatcher.SkipTo(i.GuildID, pos)
	if err != nil {
		s.InteractionRespond(i.Interaction, types.EmbedInteractionResp(types.ErrorEmbed(msgInternalErr)))
		b.log.Errorw(err.Error(), "guildID", i.GuildID)
		return
	}
	s.InteractionRespond(i.Interaction, types.TextInteractionResp(msg))
}

//validQueuePosition checks whether skip to position is valid.
func validQueuePosition(pos int) bool {
	if pos < 1 {
		return false
	}
	return true
}

//stop is the handler for stop command.
func (b *Bot) stop(s *discordgo.Session, i *discordgo.InteractionCreate) {
	msg, err := b.dispatcher.Stop(i.GuildID)
	if err != nil {
		s.InteractionRespond(i.Interaction, types.EmbedInteractionResp(types.ErrorEmbed(msgInternalErr)))
		b.log.Errorw(err.Error(), "guildID", i.GuildID)
		return
	}
	s.InteractionRespond(i.Interaction, types.TextInteractionResp(msg))
}

//preventVoiceStateChange removes guild player, if bot state is forcefully changed.
func (b *Bot) preventVoiceStateChange(_ *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	if v.VoiceState == nil || v.BeforeUpdate == nil {
		return
	}

	if v.UserID != b.cfg.AppID() {
		return
	}

	//If bot was moved from channel
	if v.ChannelID != v.BeforeUpdate.ChannelID {
		_, err := b.dispatcher.Stop(v.GuildID)
		if err != nil {
			b.log.Errorw(err.Error(), "guildID", v.GuildID)
		}
		b.log.Debug("preventVoiceStateChange trigerred")
	}
}
