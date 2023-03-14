package kits

import (
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-plus/kit"
	"github.com/moyai-studio/practice-revamp/moyai/usable"
	"time"
)

// NoDebuff kit.
func NoDebuff() kit.Kit {
	return nodebuff{}
}

type nodebuff struct{}

func (nodebuff) Name() string { return "§8NoDebuff" }
func (nodebuff) Items() kit.Items {

	return kit.Items{
		kit.Set{Slot: 0}: item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithEnchantments(item.NewEnchantment(enchantment.Unbreaking{}, 3)),
		kit.Set{Slot: 1}: item.NewStack(usable.EnderPearl{}, 16),
		kit.Add{}:        item.NewStack(usable.SplashPotion{Type: potion.StrongHealing()}, 36),
	}
}
func (nodebuff) Armour() [4]item.Stack {
	return [4]item.Stack{
		item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond}, 1).WithEnchantments(item.NewEnchantment(enchantment.Unbreaking{}, 3)),
		item.NewStack(item.Chestplate{Tier: item.ArmourTierDiamond}, 1).WithEnchantments(item.NewEnchantment(enchantment.Unbreaking{}, 3)),
		item.NewStack(item.Leggings{Tier: item.ArmourTierDiamond}, 1).WithEnchantments(item.NewEnchantment(enchantment.Unbreaking{}, 3)),
		item.NewStack(item.Boots{Tier: item.ArmourTierDiamond}, 1).WithEnchantments(item.NewEnchantment(enchantment.Unbreaking{}, 3)),
	}
}

func (nodebuff) Effects() []effect.Effect {
	return []effect.Effect{effect.New(effect.Speed{}, 1, 30*time.Hour).WithoutParticles()}
}

// Lobby kit.
func Lobby() kit.Kit { return lobby{} }

type lobby struct{}

func (lobby) Name() string { return "Lobby" }
func (lobby) Items() kit.Items {
	return kit.Items{
		kit.Set{Slot: 0}: item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithCustomName("§r§bFree For All"),
		kit.Set{Slot: 1}: item.NewStack(item.Sword{Tier: item.ToolTierIron}, 1).WithCustomName("§r§aDuels"),
		kit.Set{Slot: 8}: item.NewStack(item.Clock{}, 1).WithCustomName("§r§bSettings"),
	}
}
func (lobby) Armour() [4]item.Stack {
	return [4]item.Stack{{}, {}, {}, {}}
}
