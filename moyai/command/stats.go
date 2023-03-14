package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/moyai-studio/practice-revamp/moyai/user"
)

func Stats() cmd.Command { return cmd.New("stats", "", nil, stats{}) }

type stats struct{}

func (stats) Run(src cmd.Source, output *cmd.Output) {
	if u, ok := src.(*user.User); ok {
		output.Printf("§b%s's stats§7:\n§f%v §aKills\n§f%v §cDeaths", u.Name(), u.Kills(), u.Deaths())
	}
}
