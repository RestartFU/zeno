package usable

import (
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/moyai-studio/practice-revamp/moyai/entities"
)

// SplashPotion is an item that grants effects when thrown.
type SplashPotion struct {
	// Type is the type of splash potion.
	Type potion.Potion
}

// MaxCount ...
func (s SplashPotion) MaxCount() int {
	return 1
}

// Use ...
func (s SplashPotion) Use(w *world.World, user item.User, ctx *item.UseContext) bool {

	splash := entities.SplashPotion{}

	yaw, pitch := user.Rotation()
	e := splash.New(entity.EyePosition(user), entity.DirectionVector(user).Mul(0.5), yaw, pitch, s.Type)
	if o, ok := e.(entities.Owned); ok {
		o.Own(user)
	}

	ctx.SubtractFromCount(1)

	w.PlaySound(user.Position(), sound.ItemThrow{})

	w.AddEntity(e)

	return true
}

// EncodeItem ...
func (s SplashPotion) EncodeItem() (name string, meta int16) {
	return "minecraft:splash_potion", int16(s.Type.Uint8())
}
