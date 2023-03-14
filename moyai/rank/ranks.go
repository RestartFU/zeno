package rank

import "github.com/moyai-studio/practice-revamp/moyai/permission"

// Default rank.
type Default struct {
}

func (Default) DisplayName() bool { return false }

// Name returns the name of the rank.
func (*Default) Name() string { return "Default" }

// Flags returns the flag manager of the rank.
func (r *Default) Flags() uint64 {
	return 0
}

// Level ...
func (r *Default) Level() int {
	return 0
}

// Color ...
func (r *Default) Color() string { return "§f" }

// Staff ...
func (r *Default) Staff() bool { return false }

// Media rank.
type Media struct{}

func (Media) DisplayName() bool { return true }

// Name returns the name of the rank.
func (*Media) Name() string { return "Media" }

// Flags returns the flag manager of the rank.
func (r *Media) Flags() uint64 {
	return permission.FlagDisguise
}

// Level ...
func (r *Media) Level() int {
	return 1
}

// Color ...
func (r *Media) Color() string { return "§b§o" }

// Staff ...
func (r *Media) Staff() bool { return false }

// Famous rank.
type Famous struct{}

func (Famous) DisplayName() bool { return true }

// Name returns the name of the rank.
func (*Famous) Name() string { return "Famous" }

// Flags returns the flag manager of the rank.
func (r *Famous) Flags() uint64 {
	return permission.FlagDisguise | permission.FlagNick
}

// Level ...
func (r *Famous) Level() int {
	return 2
}

// Color ...
func (r *Famous) Color() string { return "§d§o" }

// Staff ...
func (r *Famous) Staff() bool { return false }

// Moderator rank.
type Moderator struct{}

func (Moderator) DisplayName() bool { return true }

// Name returns the name of the rank.
func (*Moderator) Name() string { return "Moderator" }

// Flags returns the flag manager of the rank.
func (r *Moderator) Flags() uint64 {
	return permission.FlagKick | permission.FlagBan | permission.FlagDisguise | permission.FlagNick | permission.FlagNickReveal | permission.FlagFreeze
}

// Level ...
func (r *Moderator) Level() int {
	return 4
}

// Color ...
func (r *Moderator) Color() string { return "§3" }

// Staff ...
func (r *Moderator) Staff() bool { return true }

// Admin rank.
type Admin struct{}

func (Admin) DisplayName() bool { return true }

// Name returns the name of the rank.
func (*Admin) Name() string { return "Admin" }

// Flags returns the flag manager of the rank.
func (r *Admin) Flags() uint64 {
	return permission.FlagKick | permission.FlagBan | permission.FlagDisguise | permission.FlagNick | permission.FlagNickReveal | permission.FlagFreeze
}

// Level ...
func (r *Admin) Level() int {
	return 5
}

// Color ...
func (r *Admin) Color() string { return "§c" }

// Staff ...
func (r *Admin) Staff() bool { return true }

// Manager rank.
type Manager struct{}

func (Manager) DisplayName() bool { return true }

// Name returns the name of the rank.
func (*Manager) Name() string { return "Manager" }

// Flags returns the permission of the rank
func (r *Manager) Flags() uint64 {
	return permission.FlagKick | permission.FlagBan | permission.FlagDisguise | permission.FlagNick | permission.FlagNickReveal | permission.FlagFreeze
}

// Level ...
func (r *Manager) Level() int {
	return 6
}

// Color ...
func (r *Manager) Color() string { return "§5" }

// Staff ...
func (r *Manager) Staff() bool { return true }

// Star rank.
type Star struct{}

func (Star) DisplayName() bool { return true }

// Name returns the name of the rank.
func (*Star) Name() string { return "Star" }

// Flags returns the permission of the rank
func (r *Star) Flags() uint64 {
	return permission.FlagDisguise | permission.FlagNick
}

// Level ...
func (r *Star) Level() int {
	return 3
}

// Color ...
func (r *Star) Color() string { return "§9§l" }

// Staff ...
func (r *Star) Staff() bool { return false }

// Owner rank.
type Owner struct{}

func (Owner) DisplayName() bool { return true }

// Name returns the name of the rank.
func (*Owner) Name() string { return "Owner" }

// Flags returns the permission of the rank
func (r *Owner) Flags() uint64 {
	return permission.FlagKick | permission.FlagBan | permission.FlagDisguise | permission.FlagNick | permission.FlagNickReveal | permission.FlagFreeze
}

// Level ...
func (r *Owner) Level() int {
	return 7
}

// Color ...
func (r *Owner) Color() string { return "§4" }

// Staff ...
func (r *Owner) Staff() bool { return true }
