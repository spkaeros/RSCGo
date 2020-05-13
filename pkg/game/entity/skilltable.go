/*
 * Copyright (c) 2020 Zachariah Knight <aeros.storkpk@gmail.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

package entity

import (
	"math"
	"strconv"
	"sync"

	"github.com/lithammer/fuzzysearch/fuzzy"
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
	sync.RWMutex
}

//Current returns the current level of the skill indicated by idx.
func (s *SkillTable) Current(idx int) int {
	s.RLock()
	defer s.RUnlock()
	return s.current[idx]
}

//DeltaMax returns the delta between maximum and current for the skill at idx.
func (s *SkillTable) DeltaMax(idx int) int {
	return s.Maximum(idx) - s.Current(idx)
}

//DecreaseCur decreases the current level of the skill at idx by delta
func (s *SkillTable) DecreaseCur(idx, delta int) {
	s.Lock()
	defer s.Unlock()
	s.current[idx] -= delta
}

//IncreaseCur increases the current level of the skill at idx by delta
func (s *SkillTable) IncreaseCur(idx, delta int) {
	s.Lock()
	defer s.Unlock()
	s.current[idx] += delta
}

//SetCur sets the current level of the skill at idx to val
func (s *SkillTable) SetCur(idx, val int) {
	s.Lock()
	defer s.Unlock()
	s.current[idx] = val
}

//DecreaseMax decreases the maximum level of the skill at idx by delta
func (s *SkillTable) DecreaseMax(idx, delta int) {
	s.Lock()
	defer s.Unlock()
	s.maximum[idx] -= delta
}

//IncreaseMax increases the maximum level of the skill at idx by delta
func (s *SkillTable) IncreaseMax(idx, delta int) {
	s.Lock()
	defer s.Unlock()
	s.maximum[idx] += delta
}

//SetMax sets the maximum level of the skill at idx to val
func (s *SkillTable) SetMax(idx, val int) {
	s.Lock()
	defer s.Unlock()
	s.maximum[idx] = val
}

//SetExp Sets the experience of the skill at idx to val
func (s *SkillTable) SetExp(idx, val int) {
	s.Lock()
	defer s.Unlock()
	s.experience[idx] = val
}

//IncExp Increases the experience of the skill at idx by val
func (s *SkillTable) IncExp(idx, val int) {
	s.Lock()
	defer s.Unlock()
	s.experience[idx] += val
}

//Maximum Returns the maximum level of the skill indicated by idx.
func (s *SkillTable) Maximum(idx int) int {
	s.RLock()
	defer s.RUnlock()
	return s.maximum[idx]
}

//Experience Returns the current level of the skill indicated by idx.
func (s *SkillTable) Experience(idx int) int {
	s.RLock()
	defer s.RUnlock()
	return s.experience[idx]
}

func (s *SkillTable) String() (s1 string) {
	s1 += "{"
	for i := 0; i < 17; i++ {
		s1 += strconv.Itoa(s.Current(i))
		s1 += "/" + strconv.Itoa(s.Maximum(i))
		s1 += " (" + strconv.Itoa(s.Experience(i)) + "),"
	}
	return s1 + "}"
}

//CombatLevel Calculates and returns the combat level for this skill table.
// Melee stats are .25 combat levels per skill level
// Pray and magic are .125 combat levels per skill level
// Ranged is .375 combat levels per skill level
// Ranged has more impact here because its a single skill, where the other skills had partners
// they tacked on a half-melee-level step per ranged level to compensate (.25+.125=.375)
func (s *SkillTable) CombatLevel() int {
	s.RLock()
	defer s.RUnlock()
	strengthAvg := (float64(s.maximum[StatAttack]+s.maximum[StatStrength]) * .25)
	defenseAvg := (float64(s.maximum[StatDefense]+s.maximum[StatHits]) * .25)
	magicAvg := (float64(s.maximum[StatPrayer]+s.maximum[StatMagic]) * .125)
	rangedAvg := (float64(s.maximum[4]) * .375)
	return int(math.Floor(defenseAvg + magicAvg + math.Max(strengthAvg, rangedAvg)))
}

var skillNames = []string{"attack", "defense", "strength", "hits", "ranged", "prayer", "magic", "cooking",
	"woodcutting", "fletching", "fishing", "firemaking", "crafting", "smithing", "mining", "herblaw", "agility",
	"thieving"}

//SkillName Returns the skill name for the provided skill index, if any.
// Otherwise returns string("nil")
func SkillName(id int) string {
	for idx, name := range skillNames {
		if idx == id {
			return name
		}
	}

	return "nil"
}

//SkillIndex Returns the skill index for the skill with the closest matching name, or -1 if there was
// no close match found.
func SkillIndex(s string) int {
	matches := fuzzy.RankFindNormalizedFold(s, skillNames)
	if len(matches) == 0 {
		return -1
	}
	return matches[0].OriginalIndex
}

// TODO: Configuration value for max level
var experienceLevels [105]uint64

func init() {
	accumulativeExp := uint64(0)
	for lvl := 0; lvl < len(experienceLevels); lvl++ {
		curLvl := float64(lvl + 1)
		exp := curLvl + 300*math.Pow(2, curLvl/7)
		accumulativeExp += uint64(exp)

		// mask off final 2 bits
		experienceLevels[lvl] = accumulativeExp & 0xFFFFFFFFFFFFFFFC
	}
}

//LevelToExperience Finds the experience required for the specified level
func LevelToExperience(lvl int) int {
	if lvl < 2 {
		return 0
	}
	if lvl > len(experienceLevels)-2 {
		return int(experienceLevels[len(experienceLevels)-1] / 4)
	}

	return int(experienceLevels[lvl-2] / 4)
}

//ExperienceToLevel Finds the maximum level for the provided experience amount.
func ExperienceToLevel(exp int) int {
	for lvl := 1; lvl < len(experienceLevels)-1; lvl++ {
		if exp < int(experienceLevels[lvl-1]/4) {
			return lvl
		}
	}

	return len(experienceLevels) - 1
}

// obase=10; ibase=16; 7FFFFFFFFFFFFFFC 9223372036854775804
// obase=10; ibase=16; FFFFFFFFFFFFFFFC 18446744073709551612
