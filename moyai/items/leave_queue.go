package items

import (
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-plus/kit"
	"github.com/moyai-studio/practice-revamp/moyai/game"
	"github.com/moyai-studio/practice-revamp/moyai/game/kits"
	"github.com/moyai-studio/practice-revamp/moyai/user"
)

func LeaveQueue(duelsProvider map[string]*game.DuelsProvider) *leaveQueue {
	return &leaveQueue{duelsProvider: duelsProvider}
}

type leaveQueue struct {
	duelsProvider map[string]*game.DuelsProvider
}

func (*leaveQueue) Name() string { return "§r§cLeave Queue" }

func (l *leaveQueue) Use(_ *event.Context, _ item.Stack, u *user.User) {
	for _, d := range l.duelsProvider {
		if d.UserQueued(u) {
			d.RemoveQueuedUser(u)
			kit.GiveKit(u.Player, kits.Lobby())
		}
	}
}

func (*leaveQueue) Item() world.Item { return item.Dye{Colour: item.ColourRed()} }
