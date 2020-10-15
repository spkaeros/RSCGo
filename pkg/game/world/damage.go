package world

import "github.com/spkaeros/rscgo/pkg/game/entity"

type HitSplat struct {
	Owner  entity.MobileEntity
	Damage int
}

func NewHitsplat(target entity.MobileEntity, damage int) HitSplat {
	return HitSplat{target, damage}
}
