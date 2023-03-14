package usable

import (
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/moyai-studio/practice-revamp/moyai/entities"
	"time"
)

// EnderPearl is a smooth, greenish-blue item used to teleport and to make an eye of ender.
type EnderPearl struct{}

// Use ...
func (EnderPearl) Use(w *world.World, user item.User, ctx *item.UseContext) bool {
	yaw, pitch := user.Rotation()
	e := entities.NewEnderPearl(entity.EyePosition(user), entity.DirectionVector(user).Mul(2.35), yaw, pitch, user)

	w.AddEntity(e)
	w.PlaySound(user.Position(), sound.ItemThrow{})
	ctx.SubtractFromCount(1)

	return true
}

// Cooldown ...
func (EnderPearl) Cooldown() time.Duration {
	return time.Second
}

// MaxCount ...
func (EnderPearl) MaxCount() int {
	return 16
}

// EncodeItem ...
func (EnderPearl) EncodeItem() (name string, meta int16) {
	return "minecraft:ender_pearl", 0
}
