package forms

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/moyai-studio/practice-revamp/moyai/game"
	"github.com/moyai-studio/practice-revamp/moyai/user"
)

type FreeForAll struct {
	host game.Host
	user *user.User
}

// Submit ...
func (f FreeForAll) Submit(_ form.Submitter, pressed form.Button) {
	// TODO: gay shit
	//g, ok := f.host.RequestFFAProvider(strings.Split(pressed.Text, "\n")[0])
	//if !ok {
	//	return
	//}
	prov, ok := f.host.RequestFFAProvider(game.NoDebuff())
	if !ok {
		f.user.Message("§cThis game isn't implemented yet!")
		return
	}
	prov.AddUser(f.user)
}

// NewFFAMenu creates a new FFA menu.
func NewFFAMenu(u *user.User, h game.Host) form.Menu {
	var buttons []form.Button
	for _, g := range game.FFA() {
		if prov, ok := h.RequestFFAProvider(g); ok {
			buttons = append(buttons, form.NewButton(
				g.Name()+"\n"+fmt.Sprintf("Playing: %v", len(prov.Users())),
				g.Texture()),
			)
		}
	}
	return form.NewMenu(FreeForAll{host: h, user: u}, "Free For All").WithButtons(buttons...)
}
