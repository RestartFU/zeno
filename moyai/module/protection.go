package module

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity/damage"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-studio/practice-revamp/moyai/user"
	"time"
)

// Protection is a module that ensures that players cannot break blocks or attack players in the lobby.
type Protection struct {
	// Lobby is the lobby world that players are restricted to.
	Lobby *world.World

	NopModule
}

// Priority ...
func (*Protection) Priority() Priority {
	return HighPriority()
}

// HandleAttackEntity ...
func (p *Protection) HandleAttackEntity(u *user.User, ctx *event.Context, _ world.Entity, _ *float64, _ *float64, _ *bool) {
	if u.World() == p.Lobby {
		ctx.Cancel()
	}
}

// HandleHurt ...
func (p *Protection) HandleHurt(u *user.User, ctx *event.Context, _ *float64, _ *time.Duration, s damage.Source) {
	if u.World() == p.Lobby || (s == damage.SourceFall{}) || u.MarkedKilled() {
		ctx.Cancel()
	}
}

// HandleFoodLoss ...
func (*Protection) HandleFoodLoss(_ *user.User, ctx *event.Context, _ int, _ int) {
	ctx.Cancel()
}

// HandleBlockBreak ...
func (*Protection) HandleBlockBreak(_ *user.User, ctx *event.Context, _ cube.Pos, _ *[]item.Stack) {
	ctx.Cancel()
}

// HandleItemDamage ...
func (*Protection) HandleItemDamage(_ *user.User, ctx *event.Context, _ item.Stack, _ int) {
	ctx.Cancel()
}
