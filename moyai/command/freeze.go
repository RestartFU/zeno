package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/moyai-studio/practice-revamp/moyai"
	"github.com/moyai-studio/practice-revamp/moyai/permission"
	"github.com/moyai-studio/practice-revamp/moyai/user"
)

func Freeze(m *moyai.Moyai) cmd.Command {
	return cmd.New("freeze", "", nil, freeze{moyai: m})
}

type freeze struct {
	moyai  *moyai.Moyai
	Target []cmd.Target
}

func (f freeze) Run(src cmd.Source, out *cmd.Output) {
	m := f.moyai
	u := src.(*user.User)
	target, _ := m.User(f.Target[0].Name())

	if target.Immobile() {
		target.SetMobile()
		f.moyai.Staffs().Messagef("§3[S] §r%s §r§bunfroze %s§b.", u.Format(), target.Format())
		out.Printf("§cYou have unfrozen %s", target.Format())
		target.Message("§cYou have been unfrozen.")
	} else {
		target.SetImmobile()
		f.moyai.Staffs().Messagef("§3[S] §r%s §r§bfroze %s§b.", u.Format(), target.Format())
		out.Printf("§cYou have frozen %s", target.Format())
		target.Message("§cYou have been frozen.")
	}

}

// Allow ...
func (f freeze) Allow(src cmd.Source) bool {
	return f.moyai.HasPermission(src, permission.FlagFreeze)
}
