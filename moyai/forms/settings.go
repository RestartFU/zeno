package forms

import (
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/moyai-studio/practice-revamp/moyai/user"
)

func Settings(u *user.User) form.Form {
	return form.New(settings{
		u:          u,
		CPS:        form.NewToggle("CPS", u.Settings().CPS),
		ScoreBoard: form.NewToggle("Scoreboard", u.Settings().ScoreBoard),
	})
}

type settings struct {
	u               *user.User
	CPS, ScoreBoard form.Toggle
}

func (s settings) Submit(_ form.Submitter) {
	s.u.Settings().CPS = s.CPS.Value()
	s.u.Settings().ScoreBoard = s.ScoreBoard.Value()
}
