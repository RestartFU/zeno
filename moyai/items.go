package moyai

import "github.com/moyai-studio/practice-revamp/moyai/items"

func (m *Moyai) registerAllItems() {
	items.RegisterUsable(
		items.FFASword(m),
		items.DuelSword(m),
		items.EnderPearl{},
		items.Settings{},
		items.LeaveQueue(m.duelProviders),
	)
}
