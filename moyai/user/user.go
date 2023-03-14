package user

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/healing"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/scoreboard"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-plus/cooldown"
	"github.com/df-plus/kit"
	"github.com/jmoiron/sqlx"
	"github.com/moyai-studio/practice-revamp/moyai/game/kits"
	"github.com/moyai-studio/practice-revamp/moyai/permission"
	"github.com/moyai-studio/practice-revamp/moyai/rank"
	"github.com/moyai-studio/practice-revamp/moyai/specs"
	"go.uber.org/atomic"
	"strconv"
	"strings"
	"sync"
	"time"
	_ "unsafe"
)

type staffMap interface {
	AddStaff(u *User)
	RemoveStaff(u *User)
}

type moyai interface {
	HasPermission(src cmd.Source, flag uint64) bool
}

type User struct {
	*player.Player
	session *session.Session

	disguisedName string
	disguisedRank rank.Rank

	dead atomic.Bool

	clickMu sync.Mutex
	clicks  []time.Time

	*cooldown.Manager

	data *data

	rank rank.Rank

	combatWith  *User
	lastWhisper *User

	moyai    moyai
	staffMap staffMap

	logs bool
}

// New returns a new *User
func New(p *player.Player, db *sqlx.DB, m moyai, staffMap staffMap) *User {
	u := &User{
		Player:        p,
		Manager:       cooldown.NewManager(),
		session:       player_session(p),
		moyai:         m,
		staffMap:      staffMap,
		disguisedName: p.Name(),
	}
	u.data = u.SelectData(db)
	if r, ok := rank.ByName(u.data.Role); ok {
		u.SetRank(r)
	} else {
		u.SetRank(&rank.Default{})
	}
	u.UpdateData(db)
	return u
}

func (u *User) ResetCoolDown(coolDown string) {
	u.Cooldown(coolDown).SetCooldown(0)
}
func (u *User) SetCoolDown(coolDown string, duration time.Duration) {
	u.Cooldown(coolDown).SetCooldown(duration)
}

func (u *User) TeleportToUser(target *User) {
	if u.World() != target.World() {
		target.World().AddEntity(u.Player)
	}
	u.Teleport(target.Position())
}

// Disguised ...
func (u *User) Disguised() bool { return u.Name() != u.disguisedName }

// DisguisedName ...
func (u *User) DisguisedName() string { return u.disguisedName }

// SetDisguisedName ...
func (u *User) SetDisguisedName(name string) { u.disguisedName = name }

// ResetDisguisedName ...
func (u *User) ResetDisguisedName() { u.disguisedName = u.Name() }

// DisguisedRank ...
func (u *User) DisguisedRank() (rank.Rank, bool) { return u.disguisedRank, u.disguisedRank != nil }

// SetDisguisedRank ...
func (u *User) SetDisguisedRank(r rank.Rank) {
	u.disguisedRank = r
}

// ResetDisguisedRank ...
func (u *User) ResetDisguisedRank() {
	u.disguisedRank = nil
}

// ResetDisguise ...
func (u *User) ResetDisguise() {
	u.ResetDisguisedName()
	u.ResetDisguisedRank()
	u.SetNameTag(u.Format())
}

// Disguise ...
func (u *User) Disguise(name string, rank rank.Rank) {
	u.SetLastWhisper(nil)
	u.SetDisguisedName(name)
	u.SetDisguisedRank(rank)
	u.SetNameTag(u.Format())
}

// SetLogs ...
func (u *User) SetLogs(b bool) { u.logs = b }

// Logs ...
func (u *User) Logs() bool { return u.logs }

// SetLastWhisper ...
func (u *User) SetLastWhisper(usr *User) {
	u.lastWhisper = usr
}

// LastWhisper ...
func (u *User) LastWhisper() (*User, bool) {
	return u.lastWhisper, u.lastWhisper != nil
}

// Spawn teleports the user to spawn.
// It also cleans up the user and reset their cooldowns.
func (u *User) Spawn(w *world.World) {
	u.CleanUp()
	w.AddEntity(u.Player)
	u.Teleport(w.Spawn().Vec3())
	u.Cooldown("combat").SetCooldown(0)
	u.Cooldown("ender_pearl").SetCooldown(0)
	kit.GiveKit(u.Player, kits.Lobby())
	u.SetGameMode(world.GameModeAdventure)
	u.MarkAlive()
}

// ReKit rekits the user with the given kit.
func (u *User) ReKit(k kit.Kit) {
	u.CleanUp()
	kit.GiveKit(u.Player, k)
}

// RemoveAllEffects removes all effect the user has.
func (u *User) RemoveAllEffects() {
	for _, eff := range u.Effects() {
		u.RemoveEffect(eff.Type())
	}
}

// Settings returns the settings of the player.
func (u *User) Settings() *data { return u.data }

// DeathAnimation starts a death animation.
func (u *User) DeathAnimation(killer world.Entity) {
	w := u.World()
	c := player.New(u.Name(), u.Skin(), u.Position())
	c.SetScale(u.Scale())
	c.SetNameTag(u.NameTag())
	w.AddEntity(c)
	for _, viewer := range w.Viewers(c.Position()) {
		viewer.ViewEntityAction(c, entity.DeathAction{})
	}
	c.KnockBack(killer.Position(), 0.5, 0.2)
	u.KnockBack(killer.Position(), 0.6, 0.2)
	time.AfterFunc(time.Millisecond*1200, func() {
		_ = c.Close()
	})
}

// Potions returns the amount of potions the user has in their inventory.
func (u *User) Potions() (n int) {
	inv := u.Inventory()
	pot, _ := item.SplashPotion{Type: potion.StrongHealing()}.EncodeItem()
	for _, i := range inv.Items() {
		uPot, _ := i.Item().EncodeItem()
		if uPot == pot {
			n++
		}
	}
	return
}

// Kill kills the user.
func (u *User) Kill(killer *User) {
	u.SetCombatWith(nil, 0)
	u.data.Deaths += 1

	killer.data.Kills += 1
	killer.ResetCoolDown("ender_pearl")

}

// CombatWith returns the last user that hurt our user and a bool of if the user is nil or not.
func (u *User) CombatWith() (*User, bool) { return u.combatWith, u.combatWith != nil }

// SetCombatWith sets the given user as the last user that hurt our user.
func (u *User) SetCombatWith(user *User, duration time.Duration) {
	u.SetCoolDown("combat", duration)
	u.combatWith = user
}

func (u *User) WouldDie(dmg float64) bool { return u.Health()-dmg <= 0 }

// SendScoreBoard sends the scoreboard to the player. (TO IMPROVE)
func (u *User) SendScoreBoard(playing int) {
	if u.data.ScoreBoard {
		combat := u.Cooldown("combat").UntilExpiration().Seconds()
		if combat < 0 {
			combat = 0
		}
		s := scoreboard.New("§b§lZeno")

		lines := []string{
			"-------------------------",
			fmt.Sprintf("§b§lPlaying§r§7: §f%v", playing),
			fmt.Sprintf("§c§lCombat§r§7: §f%v", fmt.Sprintf("%.0f", combat)),
			"§r-------------------------",
		}
		_, err := s.WriteString(strings.Join(lines, "\n"))
		if err != nil {
			fmt.Println("scoreboard error: ", err)
			return
		}
		u.SendScoreboard(s)
	} else {
		u.RemoveScoreboard()
	}
}

// CleanUp cleans up the user.
// It clears the inventory and armour.
// Removes all effects.
// Heals the player to max health.
func (u *User) CleanUp() {
	u.Inventory().Clear()
	u.Armour().Clear()
	u.RemoveAllEffects()
	u.Heal(u.MaxHealth(), healing.SourceCustom{})
}

// Kills returns the amount of kills the user has.
func (u *User) Kills() int { return u.data.Kills }

// Deaths returns the amount of deaths the user has.
func (u *User) Deaths() int { return u.data.Deaths }

// SetRank sets the given rank to the user.
func (u *User) SetRank(r rank.Rank) {
	if r == nil {
		return
	}
	if r.Staff() {
		u.staffMap.AddStaff(u)
	} else {
		u.staffMap.RemoveStaff(u)
	}
	if !u.moyai.HasPermission(u.Player, permission.FlagDisguise) && u.Disguised() {
		u.ResetDisguisedName()
	}
	u.rank = r
	u.data.Role = r.Name()
	u.SetNameTag(u.Format())
}

// Rank returns the rank of the user and a bool of if the rank is nil or not.
func (u *User) Rank() (rank.Rank, bool) { return u.rank, u.rank != nil }

// AddClick adds a click to the user's click history.
func (u *User) AddClick() {
	u.clickMu.Lock()
	u.clicks = append(u.clicks, time.Now())
	if len(u.clicks) >= 100 {
		u.clicks = u.clicks[1:]
	}
	u.clickMu.Unlock()
	if u.data.CPS {
		u.SendTip("§bCPS" + "§f " + strconv.Itoa(u.CPS()))
	}
}

// CPS returns the user's current click per second.
func (u *User) CPS() int {
	u.clickMu.Lock()
	defer u.clickMu.Unlock()

	var clicks int
	for _, past := range u.clicks {
		if time.Since(past) <= time.Second {
			clicks++
		}
	}
	return clicks
}

// Ping returns the ping of the user.
func (u *User) Ping() int {
	return int(u.Latency().Milliseconds())
}

// MarkKilled marks the user as killed.
func (u *User) MarkKilled() {
	u.dead.Store(true)
}

// MarkedKilled returns whether the user has been marked as killed.
func (u *User) MarkedKilled() bool {
	return u.dead.Load()
}

// MarkAlive marks the user as alive.
func (u *User) MarkAlive() {
	u.dead.Store(false)
}

// DeviceGroup returns the preferred device group of the user.
func (u *User) DeviceGroup() specs.DeviceGroup {
	return specs.DeviceGroupUnrestricted()
}

// PingRange returns the preferred ping range of the user.
func (u *User) PingRange() specs.PingRange {
	return specs.PingRangeUnrestricted()
}

// InputMode returns the current input mode of the user.
func (u *User) InputMode() int {
	return u.session.ClientData().CurrentInputMode
}

// Connected returns true if the user has a world.
func (u *User) Connected() bool {
	return player_session(u.Player) != session.Nop
}

// Format returns the format of the player with their rank.
func (u *User) Format() string {
	if u.Disguised() {
		if r, ok := u.DisguisedRank(); ok {
			return r.Color() + u.disguisedName + "§r"
		}
		return u.disguisedName
	}

	if r, ok := u.Rank(); ok {
		return r.Color() + u.Name() + "§r"
	}
	return u.Name()
}

//go:linkname player_session github.com/df-mc/dragonfly/server/player.(*Player).session
//noinspection ALL
func player_session(*player.Player) *session.Session
