/*
 * Copyright (c) 2019 Zachariah Knight <aeros.storkpk@gmail.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

package world

import (
	"math"
	"sync"
)

const (
	StatAttack int = iota
	StatDefense
	StatStrength
	StatHits
	StatRanged
	StatPrayer
	StatMagic
	StatCooking
	StatWoodcutting
	StatFletching
	StatFishing
	StatFiremaking
	StatCrafting
	StatSmithing
	StatMining
	StatHerblaw
	StatAgility
	StatThieving
)

//SkillTable Represents a skill table for a mob.
type SkillTable struct {
	current    [18]int
	maximum    [18]int
	experience [18]int
	lock       sync.RWMutex
}

//Current Returns the current level of the skill indicated by idx.
func (s *SkillTable) Current(idx int) int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.current[idx]
}

func (s *SkillTable) DecreaseCur(idx, delta int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.current[idx] -= delta
}

func (s *SkillTable) IncreaseCur(idx, delta int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.current[idx] += delta
}

func (s *SkillTable) SetCur(idx, val int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.current[idx] = val
}

func (s *SkillTable) DecreaseMax(idx, delta int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.maximum[idx] -= delta
}

func (s *SkillTable) IncreaseMax(idx, delta int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.maximum[idx] += delta
}

func (s *SkillTable) SetMax(idx, val int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.maximum[idx] = val
}

func (s *SkillTable) SetExp(idx, val int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.experience[idx] = val
}

func (s *SkillTable) IncExp(idx, val int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.experience[idx] += val
}

//Maximum Returns the maximum level of the skill indicated by idx.
func (s *SkillTable) Maximum(idx int) int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.maximum[idx]
}

//Experience Returns the current level of the skill indicated by idx.
func (s *SkillTable) Experience(idx int) int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.experience[idx]
}

//CombatLevel Calculates and returns the combat level for this skill table.
func (s *SkillTable) CombatLevel() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	aggressiveTotal := float32(s.maximum[0] + s.maximum[2])
	defensiveTotal := float32(s.maximum[1] + s.maximum[3])
	spiritualTotal := float32((s.maximum[5] + s.maximum[6]) / 8)
	ranged := float32(s.maximum[4])
	if aggressiveTotal < ranged*1.5 {
		return int((defensiveTotal / 4) + (ranged * 0.375) + spiritualTotal)
	}
	return int((aggressiveTotal / 4) + (defensiveTotal / 4) + spiritualTotal)
}

//SkillNames Maps skill names to their indexes.
var SkillNames = map[string]int{
	"attack": StatAttack, "defense": StatDefense, "strength": StatStrength, "hits": StatHits, "hitpoints": StatHits, "hp": StatHits,
	"ranged": StatRanged, "prayer": StatPrayer, "magic": StatMagic, "cooking": StatCooking, "woodcutting": StatWoodcutting,
	"fletching": StatFletching, "fishing": StatFishing, "firemaking": StatFiremaking, "crafting": StatCrafting, "smithing": StatSmithing,
	"mining": StatMining, "herblaw": StatHerblaw, "agility": StatAgility, "thieving": StatThieving,
}

//SkillName Returns the skill name for the provided skill index, if any.
// Otherwise returns string("nil")
func SkillName(id int) string {
	for name, idx := range SkillNames {
		if idx == id {
			return name
		}
	}
	return "nil"
}

//SkillIndex Tries to parse the skill indicated in s.  If it is out of skill bounds, returns -1.
func SkillIndex(s string) int {
	if skill, ok := SkillNames[s]; ok {
		return skill
	}
	return -1
}

var experienceLevels [104]int

func init() {
	i := 0
	for lvl := 0; lvl < 104; lvl++ {
		k := float64(lvl + 1)
		i1 := int(k + 300*math.Pow(2, k/7))
		i += i1
		experienceLevels[lvl] = (i & 0xfffffffc) / 4
	}
}

func LevelToExperience(lvl int) int {
	index := lvl - 2
	if index < 0 || index > 104 {
		return 0
	}
	return experienceLevels[index]
}

func ExperienceToLevel(exp int) int {
	for lvl := 0; lvl < 104; lvl++ {
		if exp < experienceLevels[lvl] {
			return lvl + 1
		}
	}

	return 99
}
