package items

import (
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-studio/practice-revamp/moyai/forms"
	"github.com/moyai-studio/practice-revamp/moyai/game"
	"github.com/moyai-studio/practice-revamp/moyai/user"
)

func FFASword(host game.Host) *ffaSword {
	return &ffaSword{host: host}
}

func DuelSword(host game.Host) *duelSword {
	return &duelSword{host: host}
}

type ffaSword struct {
	host game.Host
}

func (*ffaSword) Name() string { return "§r§bFree For All" }

func (f *ffaSword) Use(_ *event.Context, _ item.Stack, u *user.User) {
	if !u.Immobile() {
		u.SendForm(forms.NewFFAMenu(u, f.host))
	}
}

func (*ffaSword) Item() world.Item { return item.Sword{Tier: item.ToolTierDiamond} }

type duelSword struct {
	host game.Host
}

func (*duelSword) Name() string { return "§r§aDuels" }

func (f *duelSword) Use(_ *event.Context, _ item.Stack, u *user.User) {
	u.Messagef("§cThis feature is not released yet")
	return
	if !u.Immobile() {
		u.SendForm(forms.NewDuelsMenu(u, f.host))
	}
}

func (*duelSword) Item() world.Item { return item.Sword{Tier: item.ToolTierIron} }
