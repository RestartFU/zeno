package items

import (
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-studio/practice-revamp/moyai/forms"
	"github.com/moyai-studio/practice-revamp/moyai/user"
)

type Settings struct{}

func (Settings) Name() string { return "§r§bSettings" }

func (s Settings) Use(_ *event.Context, _ item.Stack, u *user.User) {
	u.SendForm(forms.Settings(u))
}

func (Settings) Item() world.Item { return item.Clock{} }
