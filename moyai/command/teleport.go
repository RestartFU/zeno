package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/moyai-studio/practice-revamp/moyai"
	"github.com/moyai-studio/practice-revamp/moyai/game"
	"github.com/moyai-studio/practice-revamp/moyai/permission"
	"github.com/moyai-studio/practice-revamp/moyai/user"
)

func Teleport(m *moyai.Moyai) cmd.Command {
	return cmd.New("teleport", "", []string{"tp"},
		teleport{moyai: m},
		teleportTarget{moyai: m},
	)
}

type teleport struct {
	moyai *moyai.Moyai
}

func (teleport) Run(_ cmd.Source, output *cmd.Output) {
	output.Errorf("§7Usage: §b/tp §f<player>")
}
func (t teleport) Allow(src cmd.Source) bool {
	return t.moyai.HasPermission(src, permission.FlagTeleport)
}

type teleportTarget struct {
	moyai   *moyai.Moyai
	Target  []cmd.Target
	Target2 []cmd.Target `optional:""`
}

func (t teleportTarget) Run(src cmd.Source, output *cmd.Output) {
	teleporter := src.(*user.User)
	target, _ := t.moyai.User(t.Target[0].Name())

	if len(t.Target2) >= 1 {
		teleporter = target
		target, _ = t.moyai.User(t.Target2[0].Name())
	}

	if prov, ok := t.moyai.SearchUser(teleporter); ok && !prov.HasUser(target) {
		prov.RemoveUser(teleporter)

		if targetProv, ok := t.moyai.SearchUser(target); ok {
			if ffa, ok := targetProv.(*game.FFAProvider); ok {
				ffa.AddUser(teleporter)
			}
		}
	}

	teleporter.TeleportToUser(target)
	teleporter.Messagef("§7You got §bteleported §7to §r%s.", target.Format())
	if t.moyai.Staffs().Staff(target.Name()) {
		target.Messagef("%s §r§7has been §bteleported §7to you.", teleporter.Format())
	}
}

func (t teleportTarget) Allow(src cmd.Source) bool {
	return t.moyai.HasPermission(src, permission.FlagTeleport)
}
