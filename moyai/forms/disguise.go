package forms

import (
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/moyai-studio/practice-revamp/moyai/rank"
	"github.com/moyai-studio/practice-revamp/moyai/user"
)

type moyai interface {
	DisguisedUser(name string) (*user.User, bool)
	User(name string) (*user.User, bool)
}

var ranks = []string{
	"Default", "Star", "Media", "Famous",
}

var disguises = []string{
	"Nakoso", "Staind", "Lastro", "Jewdah", "Jerseys", "M0DIFIER", "PuffedUp",
	"Darthrai", "ZIBLACKINGGG", "Verzide", "Adviser", "Marcel", "bcz", "Dream",
	"PoissonChat42", "God", "Devil", "Demon", "Hate", "Love", "Sad", "Cry", "tragic",
	"Hurt", "curse", "Killua", "Sasuke", "Kyavn", "Hayonn", "Vhimzi", "Zertiify",
	"AciDicMix", "Luhso", "Humanoides", "Tyzuko", "JZRA", "Aezuhfy", "Ghostsphere",
	"Hetoku", "UkSx CobrA", "Nebulq", "Sledds", "Mazurah", "TrickyOffer72", "FrillyGland9823",
	"WeeklyAtom1434", "GrubbySack83", "PlainFrock0731", "LuckyExile5625",
}

func Disguise(u *user.User, m moyai) form.Form {
	return form.New(disguise{
		u:     u,
		moyai: m,
		Disguises: form.NewDropdown(
			"Disguise",
			disguises,
			0,
		),
		Ranks: form.NewDropdown(
			"Rank",
			ranks,
			0,
		),
	})
}

type disguise struct {
	u         *user.User
	Disguises form.Dropdown
	Ranks     form.Dropdown
	moyai     moyai
}

func (d disguise) Submit(_ form.Submitter) {
	name := disguises[d.Disguises.Value()]
	_, ok := d.moyai.User(name)
	if ok {
		d.u.Message("§cSomeone is already online with that username.")
		return
	}
	_, ok = d.moyai.DisguisedUser(name)
	if ok {
		d.u.Message("§cSomeone is already online with that username.")
		return
	}
	r, _ := rank.ByName(ranks[d.Ranks.Value()])
	d.u.Disguise(name, r)
}
