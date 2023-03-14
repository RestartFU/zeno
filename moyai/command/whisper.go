package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/moyai-studio/practice-revamp/moyai/user"
	"strings"
)

type searchUser interface {
	User(name string) (*user.User, bool)
}

func Whisper(u searchUser) cmd.Command {
	return cmd.New("whisper", "", []string{"msg", "message", "w"}, whisper{u: u})
}

type whisper struct {
	Target  []cmd.Target
	Message cmd.Varargs
	u       searchUser
}

func (w whisper) Run(src cmd.Source, output *cmd.Output) {
	if strings.TrimSpace(string(w.Message)) == "" {
		output.Errorf("Please send a valid message")
		return
	}
	if t, ok := w.u.User(w.Target[0].Name()); ok {
		if u, ok := src.(*user.User); ok {
			targetFormat := t.Format()
			userFormat := u.Format()

			if r, ok := t.Rank(); ok && t.Disguised() {
				targetFormat = r.Color() + t.Name()
			}
			if r, ok := u.Rank(); ok && u.Disguised() {
				userFormat = r.Color() + u.Name()
			}

			output.Printf("§b(To %s§b) %s", targetFormat, w.Message)
			t.Messagef("§b(From %s§b) %s", userFormat, w.Message)

			t.SetLastWhisper(u)
		}
	} else {
		output.Errorf("Target is not online")
	}
}

func Reply() cmd.Command { return cmd.New("reply", "", []string{"r"}, reply{}) }

type reply struct {
	Message cmd.Varargs
}

func (r reply) Run(src cmd.Source, output *cmd.Output) {
	if strings.TrimSpace(string(r.Message)) == "" {
		output.Errorf("Please send a valid message")
		return
	}
	if u, ok := src.(*user.User); ok {
		if t, ok := u.LastWhisper(); ok {
			targetFormat := t.Format()
			userFormat := u.Format()

			if r, ok := t.Rank(); ok && t.Disguised() {
				targetFormat = r.Color() + t.Name()
			}
			if r, ok := u.Rank(); ok && u.Disguised() {
				userFormat = r.Color() + u.Name()
			}

			output.Printf("§b(To %s§b) %s", targetFormat, r.Message)
			t.Messagef("§b(From %s§b) %s", userFormat, r.Message)

			u.SetLastWhisper(t)
			t.SetLastWhisper(u)
		} else {
			output.Errorf("Target is not online")
		}
	}
}
