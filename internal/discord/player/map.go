package player

import (
	"sync"
)

//player represent guild player.
type player struct {
	Command chan command
	Queue   *queue
}

//playerMap is map safe for concurrent use.
type playerMap struct {
	mux      *sync.RWMutex
	internal map[string]player
}

func newPlayerMap() *playerMap {
	return &playerMap{
		mux:      &sync.RWMutex{},
		internal: make(map[string]player),
	}
}

//Load gets player from the map.
func (pm *playerMap) Load(key string) (value player, ok bool) {
	pm.mux.RLock()
	p, ok := pm.internal[key]
	pm.mux.RUnlock()
	return p, ok
}

//Delete removes player from the map.
func (pm *playerMap) Delete(key string) {
	pm.mux.Lock()
	delete(pm.internal, key)
	pm.mux.Unlock()
}

//Store adds player to the map.
func (pm *playerMap) Store(key string, value player) {
	pm.mux.Lock()
	pm.internal[key] = value
	pm.mux.Unlock()
}
