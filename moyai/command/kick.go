package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/moyai-studio/practice-revamp/moyai"
	"github.com/moyai-studio/practice-revamp/moyai/permission"
)

func Kick(m *moyai.Moyai) cmd.Command {
	return cmd.New("kick", "", nil, kick{moyai: m})
}

type kick struct {
	moyai  *moyai.Moyai
	Target []cmd.Target
	Reason cmd.Varargs
}

func (k kick) Run(src cmd.Source, out *cmd.Output) {
	if src.Name() == k.Target[0].Name() {
		out.Errorf("You can't kick yourself")
		return
	}
	if p, ok := k.moyai.User(k.Target[0].Name()); ok {
		p.Disconnect("You've been kicked:\n" + k.Reason)
	}
}

// Allow ...
func (k kick) Allow(src cmd.Source) bool {
	return k.moyai.HasPermission(src, permission.FlagKick)
}
