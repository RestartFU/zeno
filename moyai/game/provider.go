package game

import "github.com/moyai-studio/practice-revamp/moyai/user"

// Host represents a structure storing game data. This allows us to check if players are in a game, or what
// game they're in.
type Host interface {
	// SearchUser searches for a provider from a user. If no provider is found, the second return value is false.
	SearchUser(*user.User) (Provider, bool)
	// RequestDuelProvider requests a game provider for the given game. If no provider was found, the second return
	// value will be false. It is essentially the same as RequestFFAProvider, but only supports duels.
	RequestDuelProvider(Game) (*DuelsProvider, bool)
	// RequestFFAProvider requests a game provider for the given game. If no provider was found, the second return
	// value will be false. It is essentially the same as RequestFFAProvider, but only supports FFA games.
	RequestFFAProvider(Game) (*FFAProvider, bool)

	// Playing returns the count of all players playing a game under this host.
	Playing() int
}

// Provider provides important information and functions for games.
type Provider interface {
	// Game returns the game being provided.
	Game() Game
	// Users returns a list of all users in the game.
	Users() []*user.User
	// HasUser returns true if the given user is in the game.
	HasUser(*user.User) bool
	// RemoveUser removes a user from the game.
	RemoveUser(*user.User)
}
