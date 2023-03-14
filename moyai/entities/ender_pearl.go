package entities

import (
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/damage"
	"github.com/df-mc/dragonfly/server/entity/physics"
	"github.com/df-mc/dragonfly/server/entity/physics/trace"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)

// EnderPearl is a smooth, greenish-blue item used to teleport and to make an eye of ender.
type EnderPearl struct {
	transform
	yaw, pitch float64

	age   int
	close bool

	owner world.Entity

	c *entity.ProjectileComputer
}

// NewEnderPearl ...
func NewEnderPearl(pos, vel mgl64.Vec3, yaw, pitch float64, owner world.Entity) *EnderPearl {
	e := &EnderPearl{
		yaw:   yaw,
		pitch: pitch,
		c: &entity.ProjectileComputer{MovementComputer: &entity.MovementComputer{
			Gravity:           0.05,
			Drag:              0.01,
			DragBeforeGravity: true,
		}},
		owner: owner,
	}
	e.transform = newTransform(e, pos)
	e.vel = vel
	return e
}

// Name ...
func (e *EnderPearl) Name() string {
	return "Ender Pearl"
}

// EncodeEntity ...
func (e *EnderPearl) EncodeEntity() string {
	return "minecraft:ender_pearl"
}

// AABB ...
func (e *EnderPearl) AABB() physics.AABB {
	return physics.NewAABB(mgl64.Vec3{-0.125, 0, -0.125}, mgl64.Vec3{0.125, 0.25, 0.125})
}

// Rotation ...
func (e *EnderPearl) Rotation() (float64, float64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.yaw, e.pitch
}

// teleporter represents a living entity that can teleport.
type teleporter interface {
	// Teleport teleports the entity to the position given.
	Teleport(pos mgl64.Vec3)
	entity.Living
}

// Tick ...
func (e *EnderPearl) Tick(w *world.World, current int64) {
	if e.close {
		_ = e.Close()
		return
	}
	e.mu.Lock()
	m, result := e.c.TickMovement(e, e.pos, e.vel, e.yaw, e.pitch, e.ignores)
	yaw, pitch := m.Rotation()
	e.pos, e.vel, e.yaw, e.pitch = m.Position(), m.Velocity(), yaw, pitch
	e.mu.Unlock()

	e.age++
	m.Send()

	if m.Position()[1] < float64(w.Range()[0]) && current%10 == 0 {
		e.close = true
		return
	}

	if result != nil {
		if r, ok := result.(trace.EntityResult); ok {
			if l, ok := r.Entity().(entity.Living); ok {
				if _, vulnerable := l.Hurt(0.0, damage.SourceEntityAttack{Attacker: e}); vulnerable {
					l.KnockBack(m.Position(), 0.38, 0.540)
				}
			}
		}

		if owner := e.Owner(); owner != nil {
			if user, ok := owner.(teleporter); ok {
				w.PlaySound(user.Position(), sound.EndermanTeleport{})

				user.Teleport(m.Position())

				w.AddParticle(m.Position(), particle.EndermanTeleportParticle{})
				w.PlaySound(m.Position(), sound.EndermanTeleport{})

				user.Hurt(5, damage.SourceFall{})
			}
		}

		e.close = true
	}
}

// ignores returns whether the ender pearl should ignore collision with the entity passed.
func (et *EnderPearl) ignores(e world.Entity) bool {
	_, ok := e.(entity.Living)
	return !ok || e == et || (et.age < 5 && e == et.owner)
}

// New creates an ender pearl with the position, velocity, yaw, and pitch provided. It doesn't spawn the ender pearl,
// only returns it.
func (e *EnderPearl) New(pos, vel mgl64.Vec3, yaw, pitch float64) world.Entity {
	return NewEnderPearl(pos, vel, yaw, pitch, nil)
}

// Owner ...
func (e *EnderPearl) Owner() world.Entity {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.owner
}

// Own ...
func (e *EnderPearl) Own(owner world.Entity) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.owner = owner
}

// DecodeNBT decodes the properties in a map to a EnderPearl and returns a new EnderPearl entity.
func (e *EnderPearl) DecodeNBT(data map[string]interface{}) interface{} {
	return nil
}

// EncodeNBT encodes the EnderPearl entity's properties as a map and returns it.
func (e *EnderPearl) EncodeNBT() map[string]interface{} {
	return map[string]interface{}{}
}
