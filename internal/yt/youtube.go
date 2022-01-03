package yt

import (
	"context"
	"fmt"

	"os/exec"
	"regexp"
	"strings"

	"github.com/relipocere/gotune/internal/discord/types"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Extractor struct {
	folder string
	api    *youtube.Service
}

//New creates YouTube API client and is wrapper for ytdl caller.
func New(token, folder string) (Extractor, error) {
	ext := Extractor{
		folder: folder,
	}

	service, err := youtube.NewService(context.Background(), option.WithAPIKey(token))
	if err != nil {
		return ext, err
	}

	ext.api = service
	return ext, nil
}

//search returns ID of the best matching video.
func (e Extractor) search(query string) (link, title string, err error) {
	call := e.api.Search.List([]string{"id,snippet"}).Type("video").Q(query).MaxResults(1)
	response, err := call.Do()
	if err != nil {
		return "", "", err
	}

	if len(response.Items) < 1 || response.Items[0] == nil {
		return "", "", fmt.Errorf("yt response is empty")
	}

	if response.Items[0].Id == nil {
		return "", "", fmt.Errorf("yt response doesn't cointain video ID")
	}

	link = fmt.Sprintf("https://www.youtube.com/watch?v=%s", response.Items[0].Id.VideoId)
	if response.Items[0].Snippet != nil {
		title = response.Items[0].Snippet.Title
	}
	return
}

//Get downloads the file into the specified folder then returns the filename.
func (e Extractor) Get(query string) ([]types.Song, error) {
	var songs []types.Song

	link, title, err := e.search(query)
	if err != nil {
		if !isLink(query) {
			return songs, err
		}
		link = query
	}

	args := []string{"--no-colors", "--no-simulate",
		"--max-downloads", "10",
		"-P", e.folder,
		"--format", "ba",
		"--restrict-filenames",
		"-o", "%(title)s.%(ext)s", link}

	cmd := exec.Command("yt-dlp", args...)
	binOut, _ := cmd.CombinedOutput()
	out := string(binOut)

	paths := extractPaths(out, e.folder)
	if len(paths) < 1 {
		return songs, fmt.Errorf("no song files were found: %s", out)
	}

	unknownTitle := title == ""
	for _, path := range paths {
		if unknownTitle {
			title = extractTitleFromPath(path, e.folder)
		}

		songs = append(songs, types.Song{
			Title: title,
			Path:  path,
			Link:  link,
		})
	}
	return songs, nil
}

func isLink(s string) bool {
	return strings.HasPrefix(s, "https://www.youtube.com/watch")
}

//extractTitleFromPath extracts readable title from the file name.
func extractTitleFromPath(path, folder string) string {
	fLen := len(folder)
	//Accounting for closing /
	if !strings.HasSuffix(folder, "/") {
		fLen += 1
	}
	extPos := strings.LastIndex(path, ".")
	return strings.Replace(path[fLen:extPos], "_", " ", -1)
}

//extractPaths gets the names of the downloaded songs.
func extractPaths(output, folder string) []string {
	expr := fmt.Sprintf(`%s/.+\.(3gp|aac|flv|m4a|mp3|mp4|ogg|wav|webm)`, folder)
	re := regexp.MustCompile(expr)
	found := re.FindAllString(output, -1)
	return found
}
