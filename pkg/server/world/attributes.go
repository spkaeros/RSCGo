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
	"sync"
	"time"
)

//AttrList A type alias for a map of strings to empty interfaces, to hold generic mob information for easy serialization and to provide dynamic insertion/deletion of new mob properties easily
type AttrList map[string]interface{}

//AttributeList A concurrency-safe collection data type for storing misc. variables by a descriptive name.
type AttributeList struct {
	Set  map[string]interface{}
	Lock sync.RWMutex
}

//Range Runs fn(key, value) for every entry in this attribute list.
func (attributes *AttributeList) Range(fn func(string, interface{})) {
	attributes.Lock.RLock()
	defer attributes.Lock.RUnlock()
	for k, v := range attributes.Set {
		fn(k, v)
	}
}

//SetVar Sets the attribute mapped at name to value in the attribute map.
func (attributes *AttributeList) SetVar(name string, value interface{}) {
	attributes.Lock.Lock()
	attributes.Set[name] = value
	attributes.Lock.Unlock()
}

//UnsetVar Removes the attribute with the key `name` from this attribute set.
func (attributes *AttributeList) UnsetVar(name string) {
	attributes.Lock.Lock()
	delete(attributes.Set, name)
	attributes.Lock.Unlock()
}

//VarInt If there is an attribute assigned to the specified name, returns it.  Otherwise, returns zero
func (attributes *AttributeList) VarString(name string, zero string) string {
	attributes.Lock.RLock()
	defer attributes.Lock.RUnlock()
	if _, ok := attributes.Set[name].(string); !ok {
		return zero
	}

	return attributes.Set[name].(string)
}

//VarInt If there is an attribute assigned to the specified name, returns it.  Otherwise, returns zero
func (attributes *AttributeList) VarInt(name string, zero int) int {
	attributes.Lock.RLock()
	defer attributes.Lock.RUnlock()
	if _, ok := attributes.Set[name].(int); !ok {
		return zero
	}

	return attributes.Set[name].(int)
}

//MaskInt Mask attribute `name` with the specified bitmask.
func (attributes *AttributeList) StoreMask(name string, mask int) {
	attributes.Lock.Lock()
	defer attributes.Lock.Unlock()
	if val, ok := attributes.Set[name].(int); ok {
		attributes.Set[name] = val | 1<<mask
		return
	}
	attributes.Set[name] = 0|1<<mask
}

func (attributes *AttributeList) HasMask(name string, mask int) bool {
	attributes.Lock.RLock()
	defer attributes.Lock.RUnlock()
	return attributes.VarInt(name, 0) & (1 << mask) != 0
}

//UnmaskInt Mask attribute `name` with the specified bitmask.
func (attributes *AttributeList) RemoveMask(name string, mask int) {
	attributes.Lock.Lock()
	defer attributes.Lock.Unlock()
	if val, ok := attributes.Set[name].(int); ok {
		attributes.Set[name] = val & ^(1 << mask)
		return
	}
	attributes.Set[name] = 0 & ^(1<<mask)
}

//CheckMask Check if a bitmask attribute has a mask set.
func (attributes *AttributeList) CheckMask(name string, mask int) bool {
	attributes.Lock.RLock()
	defer attributes.Lock.RUnlock()
	return attributes.VarInt(name, 0) & mask != 0
}

//VarMob If there is a MobileEntity attribute assigned to the specified name, returns it.  Otherwise, returns nil
func (attributes *AttributeList) VarMob(name string) MobileEntity {
	attributes.Lock.RLock()
	defer attributes.Lock.RUnlock()
	if _, ok := attributes.Set[name].(MobileEntity); !ok {
		return nil
	}

	return attributes.Set[name].(MobileEntity)
}

//VarPlayer If there is a *Player attribute assigned to the specified name, returns it.  Otherwise, returns nil
func (attributes *AttributeList) VarPlayer(name string) *Player {
	attributes.Lock.RLock()
	defer attributes.Lock.RUnlock()
	if _, ok := attributes.Set[name].(*Player); !ok {
		return nil
	}

	return attributes.Set[name].(*Player)
}

//VarSkills If there is a *SkillTable attribute assigned to the specified name, returns it.  Otherwise, returns nil
func (attributes *AttributeList) VarSkills(name string) *SkillTable {
	attributes.Lock.RLock()
	defer attributes.Lock.RUnlock()
	if _, ok := attributes.Set[name].(*SkillTable); !ok {
		return nil
	}

	return attributes.Set[name].(*SkillTable)
}

//VarLong If there is an attribute assigned to the specified name, returns it.  Otherwise, returns zero
func (attributes *AttributeList) VarLong(name string, zero uint64) uint64 {
	attributes.Lock.RLock()
	defer attributes.Lock.RUnlock()
	if _, ok := attributes.Set[name].(uint64); !ok {
		return zero
	}

	return attributes.Set[name].(uint64)
}

//VarBool If there is an attribute assigned to the specified name, returns it.  Otherwise, returns zero
func (attributes *AttributeList) VarBool(name string, zero bool) bool {
	attributes.Lock.RLock()
	defer attributes.Lock.RUnlock()
	if _, ok := attributes.Set[name].(bool); !ok {
		return zero
	}

	return attributes.Set[name].(bool)
}

//VarTime If there is a time.Duration attribute assigned to the specified name, returns it.  Otherwise, returns zero
func (attributes *AttributeList) VarTime(name string) time.Time {
	attributes.Lock.RLock()
	defer attributes.Lock.RUnlock()
	if _, ok := attributes.Set[name].(time.Time); !ok {
		return time.Time{}
	}

	return attributes.Set[name].(time.Time)
}

//VarTime If there is a time.Duration attribute assigned to the specified name, returns it.  Otherwise, returns zero
func (attributes *AttributeList) VarPath(name string) *Pathway {
	attributes.Lock.RLock()
	defer attributes.Lock.RUnlock()
	if _, ok := attributes.Set[name].(*Pathway); !ok {
		return nil
	}

	return attributes.Set[name].(*Pathway)
}
