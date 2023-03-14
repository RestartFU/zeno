package rank

import (
	"strings"
	"sync"
)

type Rank interface {
	Name() string
	Flags() uint64
	Level() int
	Color() string
	DisplayName() bool
	Staff() bool
}

var ranks = make(map[string]Rank)
var ranksMu sync.RWMutex

func Register(rank Rank) {
	ranksMu.Lock()
	defer ranksMu.Unlock()
	ranks[strings.ToLower(rank.Name())] = rank
}

func ByName(name string) (Rank, bool) {
	name = strings.ToLower(name)

	ranksMu.RLock()
	defer ranksMu.RUnlock()
	rank, ok := ranks[name]
	return rank, ok
}

func Ranks() (rankList []Rank) {
	ranksMu.RLock()
	for _, r := range ranks {
		rankList = append(rankList, r)
	}
	defer ranksMu.RUnlock()
	return
}
