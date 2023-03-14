package forms

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/moyai-studio/practice-revamp/moyai/game"
	"github.com/moyai-studio/practice-revamp/moyai/user"
)

type Duels struct {
	host game.Host
	user *user.User
}

// Submit ...
func (f Duels) Submit(_ form.Submitter, pressed form.Button) {
	//g, ok := f.host.FromName(strings.Split(pressed.Text, "\n")[0])
	//if !ok {
	//	return
	//}
	prov, ok := f.host.RequestDuelProvider(game.NoDebuff())
	if !ok {
		f.user.Message("§cThis game isn't implemented yet!")
		return
	}
	prov.QueueUser(f.user)
}

// NewDuelsMenu creates a new duels menu.
func NewDuelsMenu(u *user.User, h game.Host) form.Menu {
	var buttons []form.Button
	for _, g := range game.Duels() {
		if prov, ok := h.RequestDuelProvider(g); ok {
			buttons = append(buttons, form.NewButton(
				g.Name()+"\n"+fmt.Sprintf("Playing: %v", len(prov.Users())),
				g.Texture()),
			)
		}
	}
	return form.NewMenu(Duels{host: h, user: u}, "Duels").WithButtons(buttons...)
}
