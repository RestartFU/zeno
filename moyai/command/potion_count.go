package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/moyai-studio/practice-revamp/moyai"
	"github.com/moyai-studio/practice-revamp/moyai/permission"
)

func PotionCount(m *moyai.Moyai) cmd.Command {
	return cmd.New("potions", "", []string{"pots"}, potionCount{moyai: m})
}

type potionCount struct {
	Target []cmd.Target
	moyai  *moyai.Moyai
}

func (p potionCount) Run(src cmd.Source, output *cmd.Output) {
	if t, ok := p.moyai.User(p.Target[0].Name()); ok {
		output.Printf("%s §7Has §b%v §7potions", t.Format(), t.Potions())
	}
}

func (p potionCount) Allow(src cmd.Source) bool {
	return p.moyai.HasPermission(src, permission.FlagAdministrator)
}
