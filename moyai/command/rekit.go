package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/moyai-studio/practice-revamp/moyai"
	"github.com/moyai-studio/practice-revamp/moyai/user"
	"math"
)

func Rekit(m *moyai.Moyai) cmd.Command {
	return cmd.New("rekit", "", nil, rekit{moyai: m})
}

type rekit struct {
	moyai *moyai.Moyai
}

func (r rekit) Run(src cmd.Source, out *cmd.Output) {
	if u, ok := src.(*user.User); ok {
		if cd := u.Cooldown("combat"); !cd.Expired() && !r.moyai.Operator(src.Name()) {
			out.Errorf("You're still in combat for %v seconds.", math.Round(cd.UntilExpiration().Seconds()))
			return
		}
		if prov, ok := r.moyai.SearchUser(u); ok {
			u.ReKit(prov.Game().Kit())
		} else {
			out.Errorf("You are not in any game.")
		}
	}
}
