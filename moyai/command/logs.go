package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/moyai-studio/practice-revamp/moyai"
	"github.com/moyai-studio/practice-revamp/moyai/permission"
	"github.com/moyai-studio/practice-revamp/moyai/user"
)

func Logs(m *moyai.Moyai) cmd.Command { return cmd.New("logs", "", nil, logs{moyai: m}) }

type logStatus string

func (logStatus) Type() string                  { return "log status" }
func (logStatus) Options(_ cmd.Source) []string { return []string{"on", "off"} }

type logs struct {
	moyai  *moyai.Moyai
	Status logStatus
}

func (l logs) Run(src cmd.Source, output *cmd.Output) {
	if u, ok := src.(*user.User); ok {
		switch l.Status {
		case "on":
			output.Errorf("Logs are now enabled.")
			u.SetLogs(true)
		case "off":
			output.Errorf("Logs are now disabled.")
			u.SetLogs(false)
		}
	}
}
func (l logs) Allow(src cmd.Source) bool {
	return l.moyai.HasPermission(src, permission.FlagLogs)
}
