package moyai

import "github.com/moyai-studio/practice-revamp/moyai/rank"

func (m *Moyai) registerAllRanks() {
	m.registerRanks(
		&rank.Owner{},
		&rank.Admin{},
		&rank.Manager{},
		&rank.Moderator{},
		&rank.Famous{},
		&rank.Media{},
		&rank.Star{},
		&rank.Default{},
	)
}

func (m *Moyai) registerRanks(ranks ...rank.Rank) {
	for _, r := range ranks {
		rank.Register(r)
	}
}
