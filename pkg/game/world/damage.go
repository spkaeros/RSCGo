package world

import `github.com/spkaeros/rscgo/pkg/game/entity`
import `math`

type HitSplat struct {
	Owner entity.MobileEntity
	Damage int
}

func NewHitsplat(target entity.MobileEntity, damage int) interface{} {
	return &HitSplat{target, damage}
}

func (p *Player) DamageFrom(m entity.MobileEntity, damage int, kind int) bool {
	damage = int(math.Min(float64(p.Skills().Current(entity.StatHits)), float64(damage)))
	splat := NewHitsplat(p, damage)
	p.Enqueue(playerEvents, splat)
	p.Skills().SetCur(entity.StatHits, p.Skills().Current(entity.StatHits) - damage)
	if p.Skills().Current(entity.StatHits) <= 0 {
		if attacker := AsPlayer(m); attacker != nil {
			attacker.PlaySound("victory")
		}
		p.Killed(m)
		return true
	}
	sound := "combat1"
	if damage == 0 {
		// a is a miss
		sound += "a"
	} else {
		// b is a hit
		sound += "b"
	}
	// we send to both parties here, so long as they happen to be a player
	if attackerp := AsPlayer(m); attackerp != nil {
		attackerp.PlaySound(sound)
	}
	p.PlaySound(sound)
	return false
}

func (n *NPC) DamageFrom(m entity.MobileEntity, damage int, kind int) bool {
	damage = int(math.Min(float64(n.Skills().Current(entity.StatHits)), float64(damage)))
	splat := NewHitsplat(n, damage)
	n.enqueueArea(npcEvents, splat)
	n.Skills().SetCur(entity.StatHits, n.Skills().Current(entity.StatHits) - damage)
	if damage > 0 && m.IsPlayer() {
		if kind == 0 {
			n.meleeRangeDamage.Put(AsPlayer(m).UsernameHash(), damage)
		} else if kind == 1 {
			n.magicDamage.Put(AsPlayer(m).UsernameHash(), damage)
		}
	}
	if n.Skills().Current(entity.StatHits) <= 0 {
		if attacker := AsPlayer(m); attacker != nil {
			attacker.PlaySound("victory")
		}
		n.Killed(m)
		return true
	}
	sound := "combat1"
	if damage == 0 {
		// a is a miss
		sound += "a"
	} else {
		// b is a hit
		sound += "b"
	}
	if attackerp := AsPlayer(m); attackerp != nil {
		attackerp.PlaySound(sound)
	}
	return false
}
