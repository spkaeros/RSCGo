package world

import (
	"github.com/spkaeros/rscgo/pkg/game/entity"
)

type Projectile struct {
	Kind   int
	Owner  entity.MobileEntity
	Target entity.MobileEntity
}

func NewProjectile(owner, target entity.MobileEntity, kind int) Projectile {
	return Projectile{kind, owner, target}
}
