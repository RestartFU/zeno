package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-studio/practice-revamp/moyai/user"
)

func Ping() cmd.Command {
	return cmd.New("ping", "", nil, ping{}, pingTarget{})
}

type ping struct{}

func (ping) Run(src cmd.Source, _ *cmd.Output) {
	if p, ok := src.(*user.User); ok {
		p.Messagef("§bPing§7:§f %vms", p.Ping())
	}
}

type pingTarget struct {
	Target []cmd.Target
}

func (p pingTarget) Run(src cmd.Source, output *cmd.Output) {
	if target, ok := p.Target[0].(*player.Player); ok {
		output.Printf("§b%s's Ping§7:§f %vms", target.Name(), target.Latency().Milliseconds())
	}
}
