package items

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-studio/practice-revamp/moyai/usable"
	"github.com/moyai-studio/practice-revamp/moyai/user"
	"time"
)

type EnderPearl struct {
}

func (e EnderPearl) Use(ctx *event.Context, _ item.Stack, u *user.User) {
	if cd := u.Cooldown("ender_pearl"); !cd.Expired() {
		ctx.Cancel()
		u.SendTip(fmt.Sprintf("§cYou are on §lEnder Pearl§r§c cooldown for §l%.2f Seconds§r", cd.UntilExpiration().Seconds()))
	} else {
		cd.SetCooldown(10 * time.Second)
		u.SendTip("§cCooldown Started")
		time.AfterFunc(9900*time.Millisecond, func() {
			if !cd.Expired() {
				u.SendTip("§aCooldown Ended")
			}
		})
	}
}

func (EnderPearl) Item() world.Item { return usable.EnderPearl{} }
