package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/moyai-studio/practice-revamp/moyai/forms"
	"github.com/moyai-studio/practice-revamp/moyai/user"
)

func Settings() cmd.Command {
	return cmd.New("settings", "", nil, settings{})
}

type settings struct {
}

func (s settings) Run(src cmd.Source, _ *cmd.Output) {
	if u, ok := src.(*user.User); ok {
		u.SendForm(forms.Settings(u))
	}
}
