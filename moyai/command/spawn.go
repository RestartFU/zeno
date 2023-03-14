package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-studio/practice-revamp/moyai/game"
	"github.com/moyai-studio/practice-revamp/moyai/user"
	"math"
)

func Spawn(host game.Host, lobby *world.World) cmd.Command {
	return cmd.New("spawn", "teleport to spawn", []string{"hub", "lobby"}, spawn{host: host, lobby: lobby})
}

type spawn struct {
	host  game.Host
	lobby *world.World
}

func (s spawn) Run(src cmd.Source, out *cmd.Output) {
	if u, ok := src.(*user.User); ok {
		if cd := u.Cooldown("combat"); !cd.Expired() {
			out.Errorf("You're still in combat for %v seconds.", math.Round(cd.UntilExpiration().Seconds()))
			return
		}
		if prov, ok := s.host.SearchUser(u); ok {
			prov.RemoveUser(u)
		} else {
			u.Spawn(s.lobby)
		}
	}
}
