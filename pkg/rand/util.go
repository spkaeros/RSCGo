package rand

import (
	rand2 "math/rand"
	"time"

	"github.com/spkaeros/rscgo/pkg/isaac"
)

var fastLockedSrc *isaac.ISAAC

var Rng *rand2.Rand

func init() {
	fastLockedSrc = isaac.New(uint64(time.Now().UnixNano()))
	Rng = rand2.New(fastLockedSrc)
	rand2.Seed(Rng.Int63())
}

func Source() *isaac.ISAAC {
	return fastLockedSrc
}
