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
	"sync"
	"time"
)

//AttributeList A concurrency-safe collection data type for storing misc. variables by a descriptive name.
type AttributeList struct {
	set  map[string]interface{}
	lock sync.RWMutex
}

func NewAttributeList() *AttributeList {
	return &AttributeList{set: make(map[string]interface{})}
}

//Range runs fn(key, value) for every entry in the attributes collection.  If fn returns true, returns to caller.
func (a *AttributeList) Range(fn func(string, interface{}) bool) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	for k, v := range a.set {
		if fn(k, v) {
			return
		}
	}
}

//Contains checks if there is an attribute in the collection set with the provided name, and returns true if so.
// Otherwise, returns false.
func (a *AttributeList) Contains(name string) bool {
	a.lock.RLock()
	defer a.lock.RUnlock()
	_, ok := a.set[name]
	return ok
}

//SetVar Sets the attribute with the provided name to value.
// NOTE: Even if there is already an attribute with this name, it will be overridden by calling this.  Maybe check
// first if any attributes exist with this name using Contains(name)
func (a *AttributeList) SetVar(name string, value interface{}) {
	a.lock.Lock()
	a.set[name] = value
	a.lock.Unlock()
}

//DecVar If there is an integer attribute with the provided name, it will decrease it by delta.
func (a *AttributeList) DecVar(name string, delta int) {
	a.lock.Lock()
	if val, ok := a.set[name].(int); ok {
		a.set[name] = val - delta
	}
	a.lock.Unlock()
}

//IncVar If there is an integer attribute with the provided name, it will increase it by delta.
func (a *AttributeList) IncVar(name string, delta int) {
	a.lock.Lock()
	if val, ok := a.set[name].(int); ok {
		a.set[name] = val + delta
	} else {
		a.set[name] = delta
	}
	a.lock.Unlock()
}

//Var Returns the attribute associated with name as a blank interface.  Needs to be cast to be useful, typically.
func (a *AttributeList) Var(name string) (interface{}, bool) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	val, ok := a.set[name]
	return val, ok
}

//UnsetVar Removes the attribute with the provided name from the collection set, if any exist.
func (a *AttributeList) UnsetVar(name string) {
	if a.Contains(name) {
		a.lock.Lock()
		delete(a.set, name)
		a.lock.Unlock()
	}
}

//VarString checks if there is a string attribute assigned to the specified name, and returns it.
// Otherwise, returns zero.
func (a *AttributeList) VarString(name string, zero string) string {
	a.lock.RLock()
	defer a.lock.RUnlock()
	if _, ok := a.set[name].(string); !ok {
		return zero
	}

	return a.set[name].(string)
}

//VarInt checks if there is an int attribute assigned to the specified name, and returns it.
// Otherwise, returns zero.
func (a *AttributeList) VarInt(name string, zero int) int {
	a.lock.RLock()
	defer a.lock.RUnlock()
	if _, ok := a.set[name].(int); !ok {
		return zero
	}

	return a.set[name].(int)
}

//RemoveMask checks for an int attribute with the specified name, and if it exists, tries to apply mask to it.
// If it doesn't exist, or is another type, it will set it to whatever the mask value is.
// NOTE: mask parameter should be the index of the bit from the right most bit that you want to activate.
func (a *AttributeList) StoreMask(name string, mask int) {
	a.lock.Lock()
	defer a.lock.Unlock()
	if val, ok := a.set[name].(int); ok {
		a.set[name] = val | 1<<mask
		return
	}
	a.set[name] = 1 << mask
}

//HasMasks checks if there is an int attribute assigned to the specified name, and if there is, checks if each mask in
// masks is set on it, individually.
// If there is no such attribute, or the type of the attribute is not an int, it will return false.
// NOTE: masks parameter should be the indexes of the bits from the right most bit that you want to check.
func (a *AttributeList) HasMasks(name string, masks ...int) bool {
	a.lock.RLock()
	defer a.lock.RUnlock()
	for _, mask := range masks {
		if a.VarInt(name, 0)&(1<<mask) != 0 {
			return true
		}
	}
	return false
}

//RemoveMask checks for an int attribute with the specified name, and if it exists, tries to unset mask from it.
// If it doesn't exist, or is another type, it will set it to 0.
// NOTE: mask parameter should be the index of the bit from the right most bit that you want to deactivate.
func (a *AttributeList) RemoveMask(name string, mask int) {
	a.lock.Lock()
	defer a.lock.Unlock()
	if val, ok := a.set[name].(int); ok {
		a.set[name] = val & ^(1 << mask)
		return
	}
	a.set[name] = 0
}

//CheckMask checks if there is an int attribute assigned to the specified name, and returns true if mask is set on it.
// Otherwise, returns false.
// NOTE: mask parameter should be the index of the bit from the right most bit that you want to check.
func (a *AttributeList) CheckMask(name string, mask int) bool {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return a.VarInt(name, 0)&mask != 0
}

//VarMob checks if there is a entity.MobileEntity attribute assigned to the specified name, and returns it.
// Otherwise, returns nil.
func (a *AttributeList) VarMob(name string) MobileEntity {
	a.lock.RLock()
	defer a.lock.RUnlock()
	if _, ok := a.set[name].(MobileEntity); !ok {
		return nil
	}

	return a.set[name].(MobileEntity)
}

//VarSkills checks if there is a *SkillTable attribute assigned to the specified name, and returns it.
// Otherwise, returns nil.
func (a *AttributeList) VarSkills(name string) *SkillTable {
	a.lock.RLock()
	defer a.lock.RUnlock()
	if _, ok := a.set[name].(*SkillTable); !ok {
		return nil
	}

	return a.set[name].(*SkillTable)
}

//VarLong checks if there is a uint64 attribute assigned to the specified name, and returns it.
// Otherwise, returns zero.
func (a *AttributeList) VarLong(name string, zero uint64) uint64 {
	a.lock.RLock()
	defer a.lock.RUnlock()
	if _, ok := a.set[name].(uint64); !ok {
		return zero
	}

	return a.set[name].(uint64)
}

//VarBool checks if there is a bool attribute assigned to the specified name, and returns it.
// Otherwise, returns zero.
func (a *AttributeList) VarBool(name string, zero bool) bool {
	a.lock.RLock()
	defer a.lock.RUnlock()
	if _, ok := a.set[name].(bool); !ok {
		return zero
	}

	return a.set[name].(bool)
}

//VarTime checks if there is a time.Time attribute assigned to the specified name, and returns it.
// Otherwise, returns time.Time{}
func (a *AttributeList) VarTime(name string) time.Time {
	a.lock.RLock()
	defer a.lock.RUnlock()
	if _, ok := a.set[name].(time.Time); !ok {
		return time.Time{}
	}

	return a.set[name].(time.Time)
}
