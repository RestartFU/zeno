package game

import (
	"github.com/df-plus/kit"
	"github.com/moyai-studio/practice-revamp/moyai/game/kits"
)

type Game interface {
	Kit() kit.Kit
	Name() string
	Texture() string
}

func FFA() []Game {
	return []Game{NoDebuff()}
}

func Duels() []Game {
	return []Game{NoDebuff()}
}

func NoDebuff() Game { return nodebuff{} }

type nodebuff struct{}

func (nodebuff) Kit() kit.Kit    { return kits.NoDebuff() }
func (nodebuff) Name() string    { return "NoDebuff" }
func (nodebuff) Texture() string { return "textures/items/potion_bottle_splash_heal.png" }
