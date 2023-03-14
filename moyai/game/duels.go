package game

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-studio/practice-revamp/moyai/user"
	"sync"
	"time"
)

// DuelsProvider is a simple duels provider for any duels game.
type DuelsProvider struct {
	arenas []string
	lobby  *world.World

	game Game

	userMu sync.Mutex
	users  []*user.User

	queueMu sync.Mutex
	queue   []*user.User
}

// NewDuelsProvider ...
func NewDuelsProvider(game Game, arena []string, lobby *world.World) *DuelsProvider {
	return &DuelsProvider{
		arenas: arena,
		lobby:  lobby,
		game:   game,
	}
}

// Game ...
func (d *DuelsProvider) Game() Game {
	return d.game
}

// Users ...
func (d *DuelsProvider) Users() []*user.User {
	d.userMu.Lock()
	defer d.userMu.Unlock()
	return d.users
}

// QueueUser ...
func (d *DuelsProvider) QueueUser(u *user.User) {
	d.queueMu.Lock()
	defer d.queueMu.Unlock()

	u.Inventory().Clear()
	u.Inventory().SetItem(8, item.NewStack(item.Dye{Colour: item.ColourRed()}, 1).WithCustomName("§r§cLeave Queue"))

	notify := func() {
		u.Message("§9You have been queued to duel. Match making settings:")
		u.Messagef("§bDevice Group: §9%v", u.DeviceGroup().Name())
		u.Messagef("§bPing Range: §9%v", u.PingRange().Name())
	}

	go func() {
		ticker := time.NewTicker(time.Second * 10)
		defer ticker.Stop()
		for range ticker.C {
			if !d.UserQueued(u) {
				return
			}
			notify()
		}
	}()

	notify()

	d.queue = append(d.queue, u)
	if match, ok := d.searchForMatch(u); ok {
		d.RemoveQueuedUser(u)
		d.RemoveQueuedUser(match)

		u.Message("Matched with another user, generating arena..")
		match.Message("Matched with another user, generating arena...")
	}
}

// UserQueued ...
func (d *DuelsProvider) UserQueued(u *user.User) bool {
	d.queueMu.Lock()
	defer d.queueMu.Unlock()

	for _, v := range d.queue {
		if v == u {
			return true
		}
	}
	return false
}

// RemoveQueuedUser ...
func (d *DuelsProvider) RemoveQueuedUser(u *user.User) {
	d.queueMu.Lock()
	defer d.queueMu.Unlock()

	for i, v := range d.queue {
		if v == u {
			d.queue = append(d.queue[:i], d.queue[i+1:]...)
			return
		}
	}
}

// HasUser ...
func (d *DuelsProvider) HasUser(user *user.User) bool {
	d.userMu.Lock()
	defer d.userMu.Unlock()
	for _, u := range d.users {
		if u == user {
			return true
		}
	}
	return false
}

// RemoveUser ...
func (d *DuelsProvider) RemoveUser(user *user.User) {
	d.userMu.Lock()
	defer d.userMu.Unlock()
	for i, u := range d.users {
		if u == user {
			d.users = append(d.users[:i], d.users[i+1:]...)
			u.Spawn(d.lobby)
			return
		}
	}
}

// searchForMatch searches for another user under the passed user's ping range and device group.
func (d *DuelsProvider) searchForMatch(user *user.User) (*user.User, bool) {
	ourPingRange, ourDeviceGroup := user.PingRange(), user.DeviceGroup()
	for _, u := range d.queue {
		if u == user {
			// Can't match with ourselves, unfortunately.
			continue
		}

		theirPingRange, theirDeviceGroup := u.PingRange(), u.DeviceGroup()

		satisfiesBothPingRanges := theirPingRange.Compare(user.Ping(), ourPingRange.Unrestricted()) && ourPingRange.Compare(u.Ping(), theirPingRange.Unrestricted())
		satisfiesBothDeviceGroups := theirDeviceGroup.Compare(user.InputMode(), ourDeviceGroup.Unrestricted()) && ourDeviceGroup.Compare(u.InputMode(), theirDeviceGroup.Unrestricted())

		if satisfiesBothPingRanges && satisfiesBothDeviceGroups {
			return u, true
		}
	}
	return nil, false
}
