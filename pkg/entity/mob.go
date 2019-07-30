package entity

type MobState uint8

const (
	Idle MobState = iota
	Walking
	Banking
	Chatting
	MenuChoosing
	Trading
	Dueling
	Fighting
	Batching
	Sleeping
)

// Any data structure that represents a mobile Entity within the game world should be able to implement this
// interface.
type Mob interface {
	// Returns the current MobState of this Mob
	State() MobState
	Entity
}