package module

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/damage"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-studio/practice-revamp/moyai/entities"
	"github.com/moyai-studio/practice-revamp/moyai/game"
	"github.com/moyai-studio/practice-revamp/moyai/user"
	"strings"
	"time"
)

// Game is a module that monitors games and makes sure that players are properly removed from them on death or quit.
type Game struct {
	// Host is the host of all games on the server.
	Host game.Host

	NopModule
}

// Priority ...
func (*Game) Priority() Priority {
	return HighPriority()
}
func maxMin(n, n2 float64) (max float64, min float64) {
	if n > n2 {
		return n, n2
	}
	return n2, n
}

// HandleHurt ...
func (g *Game) HandleHurt(u *user.User, ctx *event.Context, d *float64, _ *time.Duration, s damage.Source) {
	src, ok := s.(damage.SourceEntityAttack)
	if !ok || ctx.Cancelled() || u.AttackImmune() {
		ctx.Cancel()
		return
	}

	if _, ok := src.Attacker.(*entities.EnderPearl); ok {
		return
	}

	hurtingPlayer, ok := src.Attacker.(*user.User)
	if !ok {
		ctx.Cancel()
		return
	}
	u.SetCombatWith(hurtingPlayer, 15*time.Second)
	hurtingPlayer.SetCombatWith(u, 15*time.Second)

	if u.WouldDie(u.FinalDamageFrom(*d, s)) {
		ctx.Cancel()
		prov, ok := g.Host.SearchUser(u)
		if ok {
			_, _ = chat.Global.WriteString(fmt.Sprintf("§c%s§8[§f%v§8] §7was slain by §a%s§8[§f%v§8]", u.DisguisedName(), u.Potions(), hurtingPlayer.DisguisedName(), hurtingPlayer.Potions()))
			hurtingPlayer.ReKit(prov.Game().Kit())
			prov.RemoveUser(u)
		}
		u.Kill(hurtingPlayer)
		return
	}
	var height, force = 0.38, 0.38

	if !u.OnGround() {
		max, min := maxMin(hurtingPlayer.Position().Y(), u.Position().Y())
		if max-min >= 2.5 {
			height = 0.38 / 1.25
		}
	}
	u.KnockBack(hurtingPlayer.Position(), force, height)
}

// HandleAttackEntity ...
func (g *Game) HandleAttackEntity(u *user.User, ctx *event.Context, e world.Entity, _, _ *float64, critical *bool) {
	u.SwingArm()
	if ctx.Cancelled() {
		return
	}
	ctx.Cancel()
	// Anti interference
	//if combatWith, ok := u.CombatWith(); !u.Cooldown("combat").Expired() && ok && combatWith.Player != e {
	//	return
	//}
	otherPlayer, ok := e.(*player.Player)
	if !ok || otherPlayer.AttackImmune() || !otherPlayer.GameMode().AllowsTakingDamage() {
		return
	}
	held, _ := u.HeldItems()
	dmg := held.AttackDamage()
	src := damage.SourceEntityAttack{Attacker: u}
	if *critical {
		for _, v := range u.World().Viewers(e.Position()) {
			v.ViewEntityAction(e, entity.CriticalHitAction{})
		}
		dmg *= 1.5
	}
	otherPlayer.Hurt(dmg, src)
	otherPlayer.SetAttackImmunity(465 * time.Millisecond)
}

// HandleItemDrop ...
func (*Game) HandleItemDrop(_ *user.User, ctx *event.Context, e *entity.Item) {
	if (e.Item().Item() != item.GlassBottle{}) {
		ctx.Cancel()
		return
	}
}

// HandleQuit ...
func (g *Game) HandleQuit(u *user.User) {
	prov, ok := g.Host.SearchUser(u)
	if !ok {
		return
	}
	prov.RemoveUser(u)

	if killer, ok := u.CombatWith(); ok && !u.Cooldown("combat").Expired() {
		_, _ = chat.Global.WriteString(fmt.Sprintf("§c%s§8[§f%v§8] §7was slain by §a%s§8[§f%v§8]", u.Name(), u.Potions(), killer.Name(), killer.Potions()))
		u.Kill(killer)
	}

}

// HandleCommandExecution ...
func (g *Game) HandleCommandExecution(u *user.User, ctx *event.Context, command cmd.Command, args []string) {
	ctx.Cancel()
	command.Execute(strings.Join(args, " "), u)
}
