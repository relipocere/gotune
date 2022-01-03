package player

import (
	"sync"

	"github.com/relipocere/gotune/internal/discord/types"
)

//queue is queue safe for concurrent use/
type queue struct {
	mux   *sync.RWMutex
	songs []types.Song
}

//newQueue returns an initialized empty song queue.
func newQueue() *queue {
	return &queue{
		mux:   &sync.RWMutex{},
		songs: make([]types.Song, 0)}
}

//Len returns length of the queue.
func (q *queue) Len() int {
	q.mux.RLock()
	defer q.mux.RUnlock()
	return len(q.songs)
}

//ListSongs lists song that are still in the queue.
func (q *queue) ListSongs() []types.Song {
	q.mux.RLock()
	defer q.mux.RUnlock()
	return q.songs
}

//Push adds songs to the queue.
func (q *queue) Push(s []types.Song) {
	q.mux.Lock()
	defer q.mux.Unlock()
	q.songs = append(q.songs, s...)
}

//Pop removes first element from the queue and returns it.
func (q *queue) Pop() types.Song {
	q.mux.Lock()
	defer q.mux.Unlock()
	s := q.songs[0]
	q.songs = q.songs[1:]
	return s
}
