package command

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/moyai-studio/practice-revamp/moyai"
	"github.com/moyai-studio/practice-revamp/moyai/permission"
	"github.com/moyai-studio/practice-revamp/moyai/user"
)

func SetWorldSpawn(m *moyai.Moyai) cmd.Command {
	return cmd.New("setworldspawn", "", nil, setWorldSpawn{moyai: m})
}

type setWorldSpawn struct {
	moyai *moyai.Moyai
}

func (s setWorldSpawn) Run(src cmd.Source, output *cmd.Output) {
	if u, ok := src.(*user.User); ok {
		pos := u.Position()
		u.World().SetSpawn(cube.PosFromVec3(pos))
		output.Printf("set new world spawn to %.0f %.0f %.0f", pos.X(), pos.Y(), pos.Z())
	}
}

// Allow ...
func (s setWorldSpawn) Allow(src cmd.Source) bool {
	return s.moyai.HasPermission(src, permission.FlagSetWorldSpawn)
}
