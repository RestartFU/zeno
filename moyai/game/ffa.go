package game

import (
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-plus/kit"
	"github.com/moyai-studio/practice-revamp/moyai/user"
	"sync"
)

// FFAProvider is a simple FFA provider for any FFA game.
type FFAProvider struct {
	arena, lobby *world.World
	game         Game

	userMu sync.Mutex
	users  []*user.User
}

// NewFFAProvider ...
func NewFFAProvider(game Game, arena, lobby *world.World) *FFAProvider {
	return &FFAProvider{
		arena: arena,
		lobby: lobby,
		game:  game,
	}
}

// Game ...
func (s *FFAProvider) Game() Game {
	return s.game
}

// Users ...
func (s *FFAProvider) Users() []*user.User {
	s.userMu.Lock()
	defer s.userMu.Unlock()
	return s.users
}

// AddUser ...
func (s *FFAProvider) AddUser(user *user.User) {
	s.userMu.Lock()
	defer s.userMu.Unlock()

	s.users = append(s.users, user)
	s.arena.AddEntity(user.Player)
	user.Teleport(s.arena.Spawn().Vec3())
	user.Inventory().Clear()
	user.Armour().Clear()
	kit.GiveKit(user.Player, s.game.Kit())
}

// HasUser ...
func (s *FFAProvider) HasUser(user *user.User) bool {
	s.userMu.Lock()
	defer s.userMu.Unlock()
	for _, u := range s.users {
		if u == user {
			return true
		}
	}
	return false
}

// RemoveUser ...
func (s *FFAProvider) RemoveUser(user *user.User) {
	s.userMu.Lock()
	defer s.userMu.Unlock()
	for i, u := range s.users {
		if u == user {
			s.users = append(s.users[:i], s.users[i+1:]...)
			u.Spawn(s.lobby)
			return
		}
	}
}
