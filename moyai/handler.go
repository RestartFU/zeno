package moyai

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/damage"
	"github.com/df-mc/dragonfly/server/entity/healing"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player/skin"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-studio/practice-revamp/moyai/user"
	"net"
	"time"
)

// handler is a base handler that forwards all events to their respective modules.
type handler struct {
	user   *user.User
	server *Moyai
}

func (h *handler) HandleJump() {
	for _, m := range h.server.modules {
		m.HandleJump()
	}
}

func (h *handler) HandleItemConsume(ctx *event.Context, i item.Stack) {
	for _, m := range h.server.modules {
		m.HandleItemConsume(ctx, i)
	}
}

func (h *handler) HandleChangeWorld(before, after *world.World) {
	for _, m := range h.server.modules {
		m.HandleChangeWorld(h.user, before, after)
	}
}

// HandleMove ...
func (h *handler) HandleMove(ctx *event.Context, newPos mgl64.Vec3, newYaw, newPitch float64) {
	for _, m := range h.server.modules {
		m.HandleMove(h.user, ctx, newPos, newYaw, newPitch)
	}
}

// HandleTeleport ...
func (h *handler) HandleTeleport(ctx *event.Context, pos mgl64.Vec3) {
	for _, m := range h.server.modules {
		m.HandleTeleport(h.user, ctx, pos)
	}
}

// HandleToggleSprint ...
func (h *handler) HandleToggleSprint(ctx *event.Context, after bool) {
	for _, m := range h.server.modules {
		m.HandleToggleSprint(h.user, ctx, after)
	}
}

// HandleToggleSneak ...
func (h *handler) HandleToggleSneak(ctx *event.Context, after bool) {
	for _, m := range h.server.modules {
		m.HandleToggleSneak(h.user, ctx, after)
	}
}

// HandleChat ...
func (h *handler) HandleChat(ctx *event.Context, message *string) {
	for _, m := range h.server.modules {
		m.HandleChat(h.user, ctx, message)
	}
}

// HandleFoodLoss ...
func (h *handler) HandleFoodLoss(ctx *event.Context, from, to int) {
	for _, m := range h.server.modules {
		m.HandleFoodLoss(h.user, ctx, from, to)
	}
}

// HandleHeal ...
func (h *handler) HandleHeal(ctx *event.Context, health *float64, src healing.Source) {
	for _, m := range h.server.modules {
		m.HandleHeal(h.user, ctx, health, src)
	}
}

// HandleHurt ...
func (h *handler) HandleHurt(ctx *event.Context, damage *float64, attackImmunity *time.Duration, src damage.Source) {
	for _, m := range h.server.modules {
		m.HandleHurt(h.user, ctx, damage, attackImmunity, src)
	}
}

// HandleDeath ...
func (h *handler) HandleDeath(src damage.Source) {
	for _, m := range h.server.modules {
		m.HandleDeath(h.user, src)
	}
}

// HandleRespawn ...
func (h *handler) HandleRespawn(pos *mgl64.Vec3, w **world.World) {
	for _, m := range h.server.modules {
		m.HandleRespawn(h.user, pos)
	}
}

// HandleSkinChange ...
func (h *handler) HandleSkinChange(ctx *event.Context, skin *skin.Skin) {
	for _, m := range h.server.modules {
		m.HandleSkinChange(h.user, ctx, skin)
	}
}

// HandleStartBreak ...
func (h *handler) HandleStartBreak(ctx *event.Context, pos cube.Pos) {
	for _, m := range h.server.modules {
		m.HandleStartBreak(h.user, ctx, pos)
	}
}

// HandleBlockBreak ...
func (h *handler) HandleBlockBreak(ctx *event.Context, pos cube.Pos, drops *[]item.Stack) {
	for _, m := range h.server.modules {
		m.HandleBlockBreak(h.user, ctx, pos, drops)
	}
}

// HandleBlockPlace ...
func (h *handler) HandleBlockPlace(ctx *event.Context, pos cube.Pos, b world.Block) {
	for _, m := range h.server.modules {
		m.HandleBlockPlace(h.user, ctx, pos, b)
	}
}

// HandleBlockPick ...
func (h *handler) HandleBlockPick(ctx *event.Context, pos cube.Pos, b world.Block) {
	for _, m := range h.server.modules {
		m.HandleBlockPick(h.user, ctx, pos, b)
	}
}

// HandleItemUse ...
func (h *handler) HandleItemUse(ctx *event.Context) {
	for _, m := range h.server.modules {
		m.HandleItemUse(h.user, ctx)
	}
}

// HandleItemUseOnBlock ...
func (h *handler) HandleItemUseOnBlock(ctx *event.Context, pos cube.Pos, face cube.Face, clickPos mgl64.Vec3) {
	for _, m := range h.server.modules {
		m.HandleItemUseOnBlock(h.user, ctx, pos, face, clickPos)
	}
}

// HandleItemUseOnEntity ...
func (h *handler) HandleItemUseOnEntity(ctx *event.Context, e world.Entity) {
	for _, m := range h.server.modules {
		m.HandleItemUseOnEntity(h.user, ctx, e)
	}
}

// HandleAttackEntity ...
func (h *handler) HandleAttackEntity(ctx *event.Context, e world.Entity, force, height *float64, critical *bool) {
	for _, m := range h.server.modules {
		m.HandleAttackEntity(h.user, ctx, e, force, height, critical)
	}
	if h.user.Sprinting() {
		*critical = false
	}
}

// HandlePunchAir ...
func (h *handler) HandlePunchAir(ctx *event.Context) {
	for _, m := range h.server.modules {
		m.HandlePunchAir(h.user, ctx)
	}
}

// HandleSignEdit ...
func (h *handler) HandleSignEdit(ctx *event.Context, oldText, newText string) {
	for _, m := range h.server.modules {
		m.HandleSignEdit(h.user, ctx, oldText, newText)
	}
}

// HandleItemDamage ...
func (h *handler) HandleItemDamage(ctx *event.Context, i item.Stack, damage int) {
	for _, m := range h.server.modules {
		m.HandleItemDamage(h.user, ctx, i, damage)
	}
}

// HandleItemPickup ...
func (h *handler) HandleItemPickup(ctx *event.Context, i item.Stack) {
	for _, m := range h.server.modules {
		m.HandleItemPickup(h.user, ctx, i)
	}
}

// HandleItemDrop ...
func (h *handler) HandleItemDrop(ctx *event.Context, e *entity.Item) {
	for _, m := range h.server.modules {
		m.HandleItemDrop(h.user, ctx, e)
	}
}

// HandleTransfer ...
func (h *handler) HandleTransfer(ctx *event.Context, addr *net.UDPAddr) {
	for _, m := range h.server.modules {
		m.HandleTransfer(h.user, ctx, addr)
	}
}

// HandleCommandExecution ...
func (h *handler) HandleCommandExecution(ctx *event.Context, command cmd.Command, args []string) {
	for _, m := range h.server.modules {
		m.HandleCommandExecution(h.user, ctx, command, args)
	}
}

// HandleQuit ...
func (h *handler) HandleQuit() {
	for _, m := range h.server.modules {
		m.HandleQuit(h.user)
	}
}
