package module

import (
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-plus/kit"
	"github.com/moyai-studio/practice-revamp/moyai/game"
	"github.com/moyai-studio/practice-revamp/moyai/game/kits"
	"github.com/moyai-studio/practice-revamp/moyai/user"
	"time"
)

// Welcome is a module that welcomes the user with the season and info in the configuration file.
type Welcome struct {
	Moyai interface {
		DisguisedUser(string) (*user.User, bool)
	}
	StaffMap interface {
		AddStaff(*user.User)
	}
	// Host is the host of all games on the server
	Host game.Host
	NopModule
}

// HandleJoin ...
func (w *Welcome) HandleJoin(u *user.User) {
	go func() {
		for u.Connected() {
			u.SendScoreBoard(w.Host.Playing())
			time.Sleep(1 * time.Second)
		}
	}()
	if r, ok := u.Rank(); ok && r.Staff() {
		w.StaffMap.AddStaff(u)
	}
	if disguisedUser, ok := w.Moyai.DisguisedUser(u.Name()); ok {
		disguisedUser.Message("§cYour nickname has been reset due to someone online with the same nickname.")
		disguisedUser.ResetDisguise()
	}

	kit.GiveKit(u.Player, kits.Lobby())
	u.SetGameMode(world.GameModeAdventure)
}
