package world

type Projectile struct {
	Kind   int
	Owner  MobileEntity
	Target MobileEntity
}

func NewProjectile(owner, target MobileEntity, kind int) Projectile {
	return Projectile{kind, owner, target}
}
