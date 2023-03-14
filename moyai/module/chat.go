package module

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/moyai-studio/practice-revamp/moyai/permission"
	"github.com/moyai-studio/practice-revamp/moyai/user"
	"time"
)

type Chat struct {
	NopModule
	Perm interface {
		HasPermission(src cmd.Source, flag uint64) bool
	}
}

func (c *Chat) HandleChat(u *user.User, ctx *event.Context, message *string) {
	ctx.Cancel()

	if cd := u.Cooldown("chat"); cd.Expired() {
		name := u.Name()
		r, ok := u.Rank()
		if !ok {
			u.Disconnect("rank error, contact staff.")
			return
		}

		if u.Disguised() {
			name = u.DisguisedName()
			r, _ = u.DisguisedRank()
		}

		if !c.Perm.HasPermission(u, permission.FlagBypassChatCD) {
			cd.SetCooldown(3 * time.Second)
		}
		var format string
		if !r.DisplayName() {
			format = fmt.Sprintf("§f%s§7: §f%s", name, *message)
		} else {
			format = fmt.Sprintf("§8[%s%s§r§8] §f%s§7: §f%s", r.Color(), r.Name(), name, *message)
		}

		_, err := chat.Global.WriteString(format)
		if err != nil {
			u.Messagef("§can error occurred while trying to send a message: %s", err)
		}
	} else {
		u.Messagef("§cYou're on §lChat§r§c cooldown for §l%.2f Seconds§r", cd.UntilExpiration().Seconds())
	}
}
