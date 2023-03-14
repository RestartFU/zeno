package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/moyai-studio/practice-revamp/moyai"
	"github.com/moyai-studio/practice-revamp/moyai/permission"
	rank2 "github.com/moyai-studio/practice-revamp/moyai/rank"
	"github.com/moyai-studio/practice-revamp/moyai/user"
	"strings"
)

func Rank(m *moyai.Moyai) cmd.Command {
	return cmd.New("rank", "", nil, setRankOnline{moyai: m})
}

// set ...
type set string

// SubName ...
func (set) SubName() string { return "set" }

type rank string

// Type ...
func (rank) Type() string { return "Rank" }

// Options ...
func (rank) Options(_ cmd.Source) (s []string) {
	for _, r := range rank2.Ranks() {
		s = append(s, strings.ToLower(r.Name()))
	}
	return
}

type setRankOnline struct {
	moyai  *moyai.Moyai
	Sub    set
	Target []cmd.Target
	Rank   rank
}

func (s setRankOnline) Run(src cmd.Source, out *cmd.Output) {
	m := s.moyai

	if p, ok := m.User(s.Target[0].Name()); ok {

		if r, ok := rank2.ByName(string(s.Rank)); ok {

			if p2, ok := src.(*user.User); ok {
				if r2, ok := p2.Rank(); ok && r2.Level() < r.Level() && !m.Operator(src.Name()) {
					out.Errorf("you cannot give this rank, it is higher than yours")
					return
				}
			}
			p.SetRank(r)

			out.Printf("§a%s's rank is now %s", p.Name(), r.Name())
		}
	}
}

// Allow ...
func (s setRankOnline) Allow(src cmd.Source) bool {
	return s.moyai.HasPermission(src, permission.FlagSetRank)
}
