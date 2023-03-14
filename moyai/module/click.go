package module

import (
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-studio/practice-revamp/moyai/user"
)

// Click is a module that is used to track player clicks for CPS or for launching.
type Click struct {
	// Lobby is the lobby world that players are restricted to.
	Lobby *world.World

	NopModule
}

// HandlePunchAir ...
func (c *Click) HandlePunchAir(u *user.User, _ *event.Context) {
	u.AddClick()
}

// HandleAttackEntity ...
func (*Click) HandleAttackEntity(u *user.User, _ *event.Context, _ world.Entity, _ *float64, _ *float64, _ *bool) {
	u.AddClick()
}
