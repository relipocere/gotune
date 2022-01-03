package types

import "github.com/bwmarrin/discordgo"

//Song represents a playable track.
type Song struct {
	//Title of the song including author
	Title string
	//Path to the file
	Path string
	//Link to YouTube
	Link string
	//Requester is the user who requested the song
	Requester *discordgo.User
}
