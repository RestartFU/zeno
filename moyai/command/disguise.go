package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/moyai-studio/practice-revamp/moyai"
	"github.com/moyai-studio/practice-revamp/moyai/forms"
	"github.com/moyai-studio/practice-revamp/moyai/permission"
	"github.com/moyai-studio/practice-revamp/moyai/user"
)

func Disguise(m *moyai.Moyai) cmd.Command {
	return cmd.New("disguise", "", nil, disguise{moyai: m}, reveal{moyai: m}, selfReset{}, reset{moyai: m})
}

type disguise struct {
	moyai *moyai.Moyai
}

func (d disguise) Run(src cmd.Source, _ *cmd.Output) {
	u := src.(*user.User)
	u.SendForm(forms.Disguise(u, d.moyai))
}

func (d disguise) Allow(src cmd.Source) bool {
	return d.moyai.HasPermission(src, permission.FlagDisguise)
}
func Nick(m *moyai.Moyai) cmd.Command {
	return cmd.New("nick", "", nil, nick{moyai: m}, reveal{moyai: m}, selfReset{}, reset{moyai: m})
}

type nick struct {
	moyai *moyai.Moyai
}

func (n nick) Run(src cmd.Source, _ *cmd.Output) {
	u := src.(*user.User)
	u.SendForm(forms.Nick(u, n.moyai))
}

func (n nick) Allow(src cmd.Source) bool {
	return n.moyai.HasPermission(src, permission.FlagNick)
}

type revealSub string

// SubName ...
func (revealSub) SubName() string { return "reveal" }

type reveal struct {
	Sub   revealSub
	moyai *moyai.Moyai
	Nick  string
}

func (r reveal) Run(_ cmd.Source, output *cmd.Output) {
	if target, ok := r.moyai.DisguisedUser(r.Nick); ok {
		r, ok := target.Rank()
		if !ok {
			return
		}
		output.Printf("%s§7's real username is %s§r.", target.Format(), r.Color()+target.Name())
	} else {
		output.Errorf("§c%s is not disguised", r.Nick)
	}
}

// Allow ...
func (r reveal) Allow(src cmd.Source) bool {
	return r.moyai.HasPermission(src, permission.FlagNickReveal)
}

type resetSub string

// SubName ...
func (resetSub) SubName() string { return "reset" }

type selfReset struct {
	Sub resetSub
}

func (r selfReset) Run(src cmd.Source, output *cmd.Output) {
	u := src.(*user.User)
	if u.Disguised() {
		u.ResetDisguise()
	} else {
		output.Errorf("You are not disguised")
	}
}

type reset struct {
	Sub   resetSub
	moyai *moyai.Moyai
	Nick  string
}

func (r reset) Run(_ cmd.Source, output *cmd.Output) {
	if target, ok := r.moyai.DisguisedUser(r.Nick); ok {
		target.ResetDisguise()
	} else {
		output.Errorf("§c%s is not disguised", r.Nick)
	}
}

// Allow ...
func (r reset) Allow(src cmd.Source) bool {
	return r.moyai.HasPermission(src, permission.FlagNickReset)
}
