package module

import (
	"github.com/df-mc/dragonfly/server/event"
	"github.com/moyai-studio/practice-revamp/moyai/items"
	"github.com/moyai-studio/practice-revamp/moyai/user"
)

// ItemInteraction ...
type ItemInteraction struct {
	NopModule
}

func (*ItemInteraction) Priority() Priority { return MediumPriority() }

// HandleItemUse will handle when an item has been used in the air.
// It makes sure that the item held is compatible (registered) and that it is usable.
func (*ItemInteraction) HandleItemUse(u *user.User, ctx *event.Context) {
	p := u.Player                // Handled player
	heldItem, _ := p.HeldItems() // Main hand

	i, compatible := items.Compatible(heldItem)
	if usable, ok := i.(items.UsableItem); compatible && ok {
		usable.Use(ctx, heldItem, u)
	}
}
