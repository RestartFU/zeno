package forms

import (
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/moyai-studio/practice-revamp/moyai/rank"
	"github.com/moyai-studio/practice-revamp/moyai/user"
	"strings"
)

var validCharacters = "abcdefghijklmnopqrstuvwxyz1234567890 "

func validRune(l rune) bool {
	for _, i := range validCharacters {
		if strings.EqualFold(string(i), string(l)) {
			return true
		}
	}
	return false
}

func validName(name string) bool {
	for _, l := range name {
		if !validRune(l) {
			return false
		}
	}
	return true
}

var forbiddenDisguises = []string{
	"Chayn", "F5", "Dasalia", "xWqter", "Ozan", "Ozan12", "Czeoh", "Hulcuh", "egirlresort", "xraidse", "Aluzay",
	"zyn", "SharpnessXII", "ItsThreatz", "ImThreatz", "D3", "D3MoNiqK", "Khxqs", "qiexa", "IDontWantToDiie", "OhSad",
	"Patinal", "xWilliqm7w7", "Aweigs", "cameronisgey123", "harmfuldillon1", "ChxmpVI", "VelqtePvP",
	"ShElD0nMC", "Shappel", "Ninox", "Ninohx", "SmhCoreyy", "asitkiller", "ReachIsToggled", "Rememory", "Fipzi",
	"DocPenguinPlays", "ChaqsXP", "FoulDrip", "Lyreoz", "CurtPlayzMC", "yunghesu", "Frite", "FriteQc", "Myma",
	"MymaQc", "Dreacho", "DatPigmaster", "JavaJar", "xernah", "LeafyQuan", "Patar", "PatarHD", "PatarHD123",
	"Evident", "Swimfan", "Swimfan72", "Qwimston", "QwimstonYT", "Cranexe", "ItzDiecies", "Diecies", "ImPizzas",
	"xWqki", "EGIRLSLAYER2", "rooted", "RootedZ", "Restart", "RestartFU", "D3",
}

func Nick(u *user.User, m moyai) form.Form {
	return form.New(nick{
		u:     u,
		moyai: m,
		Disguise: form.NewInput(
			"Nickname",
			"",
			"ex: RestartIsCute",
		),
		Ranks: form.NewDropdown(
			"Rank",
			ranks,
			0,
		),
	})
}

type nick struct {
	u        *user.User
	Disguise form.Input
	Ranks    form.Dropdown
	moyai    moyai
}

func (n nick) Submit(_ form.Submitter) {
	name := n.Disguise.Value()

	if len(name) < 3 || len(name) > 15 || !validName(name) {
		n.u.Message("§cYour nickname must be between 3 - 15 characters and only contain characters a - z or numbers 0 - 9")
		return
	}

	for _, i := range forbiddenDisguises {
		if strings.EqualFold(i, name) {
			n.u.Message("§cYou can't use this nickname.")
			return
		}
	}

	_, ok := n.moyai.User(name)
	if ok {
		n.u.Message("§cSomeone is already online with that username.")
		return
	}
	_, ok = n.moyai.DisguisedUser(name)
	if ok {
		n.u.Message("§cSomeone is already online with that username.")
		return
	}
	r, _ := rank.ByName(ranks[n.Ranks.Value()])
	n.u.Disguise(name, r)
}
