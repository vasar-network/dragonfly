package entity

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/cube/trace"
	"github.com/df-mc/dragonfly/server/entity/damage"
	"github.com/df-mc/dragonfly/server/internal/nbtconv"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/go-gl/mathgl/mgl64"
)

// Snowball is a throwable projectile which damages entities on impact.
type Snowball struct {
	transform
	age   int
	close bool

	owner world.Entity

	c *ProjectileComputer
}

// NewSnowball ...
func NewSnowball(pos mgl64.Vec3, owner world.Entity) *Snowball {
	s := &Snowball{
		c: &ProjectileComputer{&MovementComputer{
			Gravity:           0.03,
			Drag:              0.01,
			DragBeforeGravity: true,
		}},
		owner: owner,
	}
	s.transform = newTransform(s, pos)

	return s
}

// Name ...
func (s *Snowball) Name() string {
	return "Snowball"
}

// EncodeEntity ...
func (s *Snowball) EncodeEntity() string {
	return "minecraft:snowball"
}

// BBox ...
func (s *Snowball) BBox() cube.BBox {
	return cube.Box(-0.125, 0, -0.125, 0.125, 0.25, 0.125)
}

// Tick ...
func (s *Snowball) Tick(w *world.World, current int64) {
	if s.close {
		_ = s.Close()
		return
	}
	s.mu.Lock()
	m, result := s.c.TickMovement(s, s.pos, s.vel, 0, 0, s.ignores)
	s.pos, s.vel = m.pos, m.vel
	s.mu.Unlock()

	s.age++
	m.Send()

	if m.pos[1] < float64(w.Range()[0]) && current%10 == 0 {
		s.close = true
		return
	}

	if result != nil {
		for i := 0; i < 6; i++ {
			w.AddParticle(result.Position(), particle.SnowballPoof{})
		}

		if r, ok := result.(trace.EntityResult); ok {
			if l, ok := r.Entity().(Living); ok {
				if _, vulnerable := l.Hurt(0.0, damage.SourceProjectile{Projectile: s, Owner: s.Owner()}); vulnerable {
					l.KnockBack(m.pos, 0.45, 0.3608)
				}
			}
		}

		s.close = true
	}
}

// ignores returns whether the snowball should ignore collision with the entity passed.
func (s *Snowball) ignores(entity world.Entity) bool {
	_, ok := entity.(Living)
	return !ok || entity == s || (s.age < 5 && entity == s.owner)
}

// New creates a snowball with the position, velocity, yaw, and pitch provided. It doesn't spawn the snowball,
// only returns it.
func (s *Snowball) New(pos, vel mgl64.Vec3, owner world.Entity) world.Entity {
	snow := NewSnowball(pos, owner)
	snow.vel = vel
	return snow
}

// Explode ...
func (s *Snowball) Explode(explosionPos mgl64.Vec3, impact float64, _ block.ExplosionConfig) {
	s.mu.Lock()
	s.vel = s.vel.Add(s.pos.Sub(explosionPos).Normalize().Mul(impact))
	s.mu.Unlock()
}

// Owner ...
func (s *Snowball) Owner() world.Entity {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.owner
}

// DecodeNBT decodes the properties in a map to a Snowball and returns a new Snowball entity.
func (s *Snowball) DecodeNBT(data map[string]any) any {
	return s.New(
		nbtconv.MapVec3(data, "Pos"),
		nbtconv.MapVec3(data, "Motion"),
		nil,
	)
}

// EncodeNBT encodes the Snowball entity's properties as a map and returns it.
func (s *Snowball) EncodeNBT() map[string]any {
	return map[string]any{
		"Pos":    nbtconv.Vec3ToFloat32Slice(s.Position()),
		"Yaw":    0.0,
		"Pitch":  0.0,
		"Motion": nbtconv.Vec3ToFloat32Slice(s.Velocity()),
		"Damage": 0.0,
	}
}
