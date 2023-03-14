package module

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/damage"
	"github.com/df-mc/dragonfly/server/entity/healing"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player/skin"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-studio/practice-revamp/moyai/user"
	"net"
	"time"
)

// Module represents a player handler along with its priority, for use with multi-handlers.
type Module interface {
	// Priority returns the priority of the module.
	Priority() Priority

	// HandleJoin handles the join of a player.
	HandleJoin(u *user.User)
	// HandleMove handles the movement of a player. ctx.Cancel() may be called to cancel the movement event.
	// The new position, yaw and pitch are passed.
	HandleMove(u *user.User, ctx *event.Context, newPos mgl64.Vec3, newYaw, newPitch float64)
	// HandleJump handles the player jumping.
	HandleJump()
	// HandleTeleport handles the teleportation of a player. ctx.Cancel() may be called to cancel it.
	HandleTeleport(u *user.User, ctx *event.Context, pos mgl64.Vec3)
	// HandleChangeWorld handles when the player is added to a new world. before may be nil.
	HandleChangeWorld(u *user.User, before, after *world.World)
	// HandleToggleSprint handles when the player starts or stops sprinting.
	// After is true if the player is sprinting after toggling (changing their sprinting state).
	HandleToggleSprint(u *user.User, ctx *event.Context, after bool)
	// HandleToggleSneak handles when the player starts or stops sneaking.
	// After is true if the player is sneaking after toggling (changing their sneaking state).
	HandleToggleSneak(u *user.User, ctx *event.Context, after bool)
	// HandleChat handles a message sent in the chat by a player. ctx.Cancel() may be called to cancel the
	// message being sent in chat.
	// The message may be changed by assigning to *message.
	HandleChat(u *user.User, ctx *event.Context, message *string)
	// HandleFoodLoss handles the food bar of a player depleting naturally, for example because the player was
	// sprinting and jumping. ctx.Cancel() may be called to cancel the food points being lost.
	HandleFoodLoss(u *user.User, ctx *event.Context, from, to int)
	// HandleHeal handles the player being healed by a healing source. ctx.Cancel() may be called to cancel
	// the healing.
	// The health added may be changed by assigning to *health.
	HandleHeal(u *user.User, ctx *event.Context, health *float64, src healing.Source)
	// HandleHurt handles the player being hurt by any damage source. ctx.Cancel() may be called to cancel the
	// damage being dealt to the player.
	// The damage dealt to the player may be changed by assigning to *damage.
	HandleHurt(u *user.User, ctx *event.Context, damage *float64, attackImmunity *time.Duration, src damage.Source)
	// HandleDeath handles the player dying to a particular damage cause.
	HandleDeath(u *user.User, src damage.Source)
	// HandleRespawn handles the respawning of the player in the world. The spawn position passed may be
	// changed by assigning to *pos.
	HandleRespawn(u *user.User, pos *mgl64.Vec3)
	// HandleSkinChange handles the player changing their skin. ctx.Cancel() may be called to cancel the skin
	// change.
	HandleSkinChange(u *user.User, ctx *event.Context, skin *skin.Skin)
	// HandleStartBreak handles the player starting to break a block at the position passed. ctx.Cancel() may
	// be called to stop the player from breaking the block completely.
	HandleStartBreak(u *user.User, ctx *event.Context, pos cube.Pos)
	// HandleBlockBreak handles a block that is being broken by a player. ctx.Cancel() may be called to cancel
	// the block being broken. A pointer to a slice of the block's drops is passed, and may be altered
	// to change what items will actually be dropped.
	HandleBlockBreak(u *user.User, ctx *event.Context, pos cube.Pos, drops *[]item.Stack)
	// HandleBlockPlace handles the player placing a specific block at a position in its world. ctx.Cancel()
	// may be called to cancel the block being placed.
	HandleBlockPlace(u *user.User, ctx *event.Context, pos cube.Pos, b world.Block)
	// HandleBlockPick handles the player picking a specific block at a position in its world. ctx.Cancel()
	// may be called to cancel the block being picked.
	HandleBlockPick(u *user.User, ctx *event.Context, pos cube.Pos, b world.Block)
	// HandleItemUse handles the player using an item in the air. It is called for each item, although most
	// will not actually do anything. Items such as snowballs may be thrown if HandleItemUse does not cancel
	// the context using ctx.Cancel(). It is not called if the player is holding no item.
	HandleItemUse(u *user.User, ctx *event.Context)
	// HandleItemUseOnBlock handles the player using the item held in its main hand on a block at the block
	// position passed. The face of the block clicked is also passed, along with the relative click position.
	// The click position has X, Y and Z values which are all in the range 0.0-1.0. It is also called if the
	// player is holding no item.
	HandleItemUseOnBlock(u *user.User, ctx *event.Context, pos cube.Pos, face cube.Face, clickPos mgl64.Vec3)
	// HandleItemUseOnEntity handles the player using the item held in its main hand on an entity passed to
	// the method.
	// HandleItemUseOnEntity is always called when a player uses an item on an entity, regardless of whether
	// the item actually does anything when used on an entity. It is also called if the player is holding no
	// item.
	HandleItemUseOnEntity(u *user.User, ctx *event.Context, e world.Entity)
	// HandleAttackEntity handles the player attacking an entity using the item held in its hand. ctx.Cancel()
	// may be called to cancel the attack, which will cancel damage dealt to the target and will stop the
	// entity from being knocked back.
	// The entity attacked may not be alive (implements entity.Living), in which case no damage will be dealt
	// and the target won't be knocked back.
	// The entity attacked may also be immune when this method is called, in which case no damage and knock-
	// back will be dealt.
	// The knock back force and height is also provided which can be modified.
	// The attack can be a critical attack, which would increase damage by a factor of 1.5 and
	// spawn critical hit particles around the target entity. These particles will not be displayed
	// if no damage is dealt.
	HandleItemConsume(*event.Context, item.Stack)
	HandleAttackEntity(u *user.User, ctx *event.Context, e world.Entity, force, height *float64, critical *bool)
	// HandlePunchAir handles the player punching air.
	HandlePunchAir(u *user.User, ctx *event.Context)
	// HandleSignEdit handles the player editing a sign. It is called for every keystroke while editing a sign and
	// has both the old text passed and the text after the edit. This typically only has a change of one character.
	HandleSignEdit(u *user.User, ctx *event.Context, oldText, newText string)
	// HandleItemDamage handles the event wherein the item either held by the player or as armour takes
	// damage through usage.
	// The type of the item may be checked to determine whether it was armour or a tool used. The damage to
	// the item is passed.
	HandleItemDamage(u *user.User, ctx *event.Context, i item.Stack, damage int)
	// HandleItemPickup handles the player picking up an item from the ground. The item stack laying on the
	// ground is passed. ctx.Cancel() may be called to prevent the player from picking up the item.
	HandleItemPickup(u *user.User, ctx *event.Context, i item.Stack)
	// HandleItemDrop handles the player dropping an item on the ground. The dropped item entity is passed.
	// ctx.Cancel() may be called to prevent the player from dropping the entity.Item passed on the ground.
	// e.Item() may be called to obtain the item stack dropped.
	HandleItemDrop(u *user.User, ctx *event.Context, e *entity.Item)
	// HandleTransfer handles a player being transferred to another server. ctx.Cancel() may be called to
	// cancel the transfer.
	HandleTransfer(u *user.User, ctx *event.Context, addr *net.UDPAddr)
	// HandleCommandExecution handles the command execution of a player, who wrote a command in the chat.
	// ctx.Cancel() may be called to cancel the command execution.
	HandleCommandExecution(u *user.User, ctx *event.Context, command cmd.Command, args []string)
	// HandleQuit handles the closing of a player. It is always called when the player is disconnected,
	// regardless of the reason.
	HandleQuit(u *user.User)
}

// NopModule implements the Module interface but does not execute any code when an event is called.
// Users may embed NopModule to avoid having to implement each method.
type NopModule struct{}

// Compile time check to make sure NopModule implements Handler.
var _ Module = (*NopModule)(nil)

// Priority ...
func (NopModule) Priority() Priority { return LowPriority() }

func (NopModule) HandleItemConsume(*event.Context, item.Stack) {}

// HandleJoin ...
func (NopModule) HandleJoin(*user.User) {}

// HandleItemDrop ...
func (NopModule) HandleItemDrop(*user.User, *event.Context, *entity.Item) {}

// HandleMove ...
func (NopModule) HandleMove(*user.User, *event.Context, mgl64.Vec3, float64, float64) {}

// HandleJump ...
func (NopModule) HandleJump() {}

// HandleTeleport ...
func (NopModule) HandleTeleport(*user.User, *event.Context, mgl64.Vec3) {}

// HandleChangeWorld ...
func (NopModule) HandleChangeWorld(*user.User, *world.World, *world.World) {}

// HandleToggleSprint ...
func (NopModule) HandleToggleSprint(*user.User, *event.Context, bool) {}

// HandleToggleSneak ...
func (NopModule) HandleToggleSneak(*user.User, *event.Context, bool) {}

// HandleCommandExecution ...
func (NopModule) HandleCommandExecution(*user.User, *event.Context, cmd.Command, []string) {}

// HandleTransfer ...
func (NopModule) HandleTransfer(*user.User, *event.Context, *net.UDPAddr) {}

// HandleChat ...
func (NopModule) HandleChat(*user.User, *event.Context, *string) {}

// HandleSkinChange ...
func (NopModule) HandleSkinChange(*user.User, *event.Context, *skin.Skin) {}

// HandleStartBreak ...
func (NopModule) HandleStartBreak(*user.User, *event.Context, cube.Pos) {}

// HandleBlockBreak ...
func (NopModule) HandleBlockBreak(*user.User, *event.Context, cube.Pos, *[]item.Stack) {}

// HandleBlockPlace ...
func (NopModule) HandleBlockPlace(*user.User, *event.Context, cube.Pos, world.Block) {}

// HandleBlockPick ...
func (NopModule) HandleBlockPick(*user.User, *event.Context, cube.Pos, world.Block) {}

// HandleSignEdit ...
func (NopModule) HandleSignEdit(*user.User, *event.Context, string, string) {}

// HandleItemPickup ...
func (NopModule) HandleItemPickup(*user.User, *event.Context, item.Stack) {}

// HandleItemUse ...
func (NopModule) HandleItemUse(*user.User, *event.Context) {}

// HandleItemUseOnBlock ...
func (NopModule) HandleItemUseOnBlock(*user.User, *event.Context, cube.Pos, cube.Face, mgl64.Vec3) {}

// HandleItemUseOnEntity ...
func (NopModule) HandleItemUseOnEntity(*user.User, *event.Context, world.Entity) {}

// HandleItemDamage ...
func (NopModule) HandleItemDamage(*user.User, *event.Context, item.Stack, int) {
}

// HandleAttackEntity ...
func (NopModule) HandleAttackEntity(*user.User, *event.Context, world.Entity, *float64, *float64, *bool) {
}

// HandlePunchAir ...
func (NopModule) HandlePunchAir(*user.User, *event.Context) {}

// HandleHurt ...
func (NopModule) HandleHurt(*user.User, *event.Context, *float64, *time.Duration, damage.Source) {}

// HandleHeal ...
func (NopModule) HandleHeal(*user.User, *event.Context, *float64, healing.Source) {}

// HandleFoodLoss ...
func (NopModule) HandleFoodLoss(*user.User, *event.Context, int, int) {}

// HandleDeath ...
func (NopModule) HandleDeath(*user.User, damage.Source) {}

// HandleRespawn ...
func (NopModule) HandleRespawn(*user.User, *mgl64.Vec3) {}

// HandleQuit ...
func (NopModule) HandleQuit(*user.User) {}
