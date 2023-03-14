package moyai

import (
	"github.com/moyai-studio/practice-revamp/moyai/game"
)

func (m *Moyai) registerAll() {
	ndbf, _ := m.WorldManager.World("NoDebuff")
	m.registerDuels(
		game.NewDuelsProvider(game.NoDebuff(), []string{}, m.DefaultWorld()),
	)
	m.registerFFA(
		game.NewFFAProvider(game.NoDebuff(), ndbf.World, m.DefaultWorld()),
	)
}
func (m *Moyai) registerDuels(providers ...*game.DuelsProvider) {
	for _, p := range providers {
		m.duelProviders[p.Game().Name()] = p
	}
}
func (m *Moyai) registerFFA(providers ...*game.FFAProvider) {
	for _, p := range providers {
		m.ffaProviders[p.Game().Name()] = p
	}
}
