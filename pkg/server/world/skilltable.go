/*
 * Copyright (c) 2019 Zachariah Knight <aeros.storkpk@gmail.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

package world

import "sync"

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
	Lock       sync.RWMutex
}

//Current Returns the current level of the skill indicated by idx.
func (s *SkillTable) Current(idx int) int {
	s.Lock.RLock()
	defer s.Lock.RUnlock()
	return s.current[idx]
}

func (s *SkillTable) DecreaseCur(idx, delta int) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.current[idx] -= delta
}

func (s *SkillTable) IncreaseCur(idx, delta int) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.current[idx] += delta
}

func (s *SkillTable) SetCur(idx, val int) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.current[idx] = val
}

func (s *SkillTable) DecreaseMax(idx, delta int) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.maximum[idx] -= delta
}

func (s *SkillTable) IncreaseMax(idx, delta int) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.maximum[idx] += delta
}

func (s *SkillTable) SetMax(idx, val int) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.maximum[idx] = val
}

func (s *SkillTable) SetExp(idx, val int) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.experience[idx] = val
}

func (s *SkillTable) IncExp(idx, val int) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.experience[idx] += val
}

//Maximum Returns the maximum level of the skill indicated by idx.
func (s *SkillTable) Maximum(idx int) int {
	s.Lock.RLock()
	defer s.Lock.RUnlock()
	return s.maximum[idx]
}

//Experience Returns the current level of the skill indicated by idx.
func (s *SkillTable) Experience(idx int) int {
	s.Lock.RLock()
	defer s.Lock.RUnlock()
	return s.experience[idx]
}

//CombatLevel Calculates and returns the combat level for this skill table.
func (s *SkillTable) CombatLevel() int {
	s.Lock.RLock()
	defer s.Lock.RUnlock()
	aggressiveTotal := float32(s.maximum[0] + s.maximum[2])
	defensiveTotal := float32(s.maximum[1] + s.maximum[3])
	spiritualTotal := float32((s.maximum[5] + s.maximum[6]) / 8)
	ranged := float32(s.maximum[4])
	if aggressiveTotal < ranged*1.5 {
		return int((defensiveTotal / 4) + (ranged * 0.375) + spiritualTotal)
	}
	return int((aggressiveTotal / 4) + (defensiveTotal / 4) + spiritualTotal)
}
