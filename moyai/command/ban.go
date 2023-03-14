package command

import (
	"fmt"
	"github.com/RestartFU/list"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/moyai-studio/practice-revamp/moyai"
	"github.com/moyai-studio/practice-revamp/moyai/permission"
	"github.com/moyai-studio/practice-revamp/moyai/user"
)

func Ban(l *list.List, m *moyai.Moyai) cmd.Command {
	return cmd.New("ban", "", nil, BANOnline{l: l, moyai: m}, BANOffline{moyai: m, l: l})
}

type banArgument string

func (banArgument) Type() string { return "Ban Argument" }
func (banArgument) Options(_ cmd.Source) []string {
	return []string{
		"-s",
	}
}

type BANOffline struct {
	Player   string
	l        *list.List
	moyai    *moyai.Moyai
	Argument banArgument `optional:""`
}

func (b BANOffline) Run(src cmd.Source, output *cmd.Output) {
	b.l.Add(b.Player)

	u := src.(*user.User)

	switch b.Argument {
	case "-s":
		b.moyai.Staffs().Messagef("§3[S] §r%s §r§bbanned §r%s§b.", u.Format(), b.Player)
	default:
		chat.Global.WriteString(fmt.Sprintf("%s §7has permanently banned §r%s.", u.Format(), b.Player))
	}
}
func (b BANOffline) Allow(src cmd.Source) bool {
	return b.moyai.HasPermission(src, permission.FlagBan)
}

type BANOnline struct {
	Player   []cmd.Target
	l        *list.List
	moyai    *moyai.Moyai
	Argument banArgument `optional:""`
}

func (b BANOnline) Run(src cmd.Source, output *cmd.Output) {
	if src.Name() == b.Player[0].Name() {
		output.Errorf("don't ban yourself retard")
		return
	}
	u := src.(*user.User)
	target, _ := b.moyai.User(b.Player[0].Name())

	switch b.Argument {
	case "-s":
		b.moyai.Staffs().Messagef("§3[S] §r%s §r§bbanned §r%s§b.", u.Format(), target.Format())
	default:
		chat.Global.WriteString(fmt.Sprintf("%s §7has permanently banned §r%s.", u.Format(), target.Format()))
	}

	target.Disconnect("You are now banned")
	b.l.Add(b.Player[0].Name())

}
func (b BANOnline) Allow(src cmd.Source) bool {
	return b.moyai.HasPermission(src, permission.FlagBan)
}

type UNBAN struct {
	Player string
	l      *list.List
	moyai  *moyai.Moyai
}

func Unban(l *list.List, m *moyai.Moyai) cmd.Command {
	return cmd.New("unban", "", nil, UNBAN{l: l, moyai: m})
}

func (b UNBAN) Run(src cmd.Source, output *cmd.Output) {
	if b.l.Listed(b.Player) {
		b.l.Remove(b.Player)
		u := src.(*user.User)
		b.moyai.Staffs().Messagef("§3[S] §r%s §r§bunbanned §r%s", u.Format(), b.Player)
	} else {
		output.Errorf("%s is not banned", b.Player)
	}
}

func (b UNBAN) Allow(src cmd.Source) bool {
	return b.moyai.HasPermission(src, permission.FlagBan)
}
