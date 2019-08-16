package entity

//MobState Mob state.
type MobState uint8

const (
	//Idle The default MobState, means doing nothing.
	Idle MobState = iota
	//Walking The mob is walking.
	Walking
	//Banking The mob is banking.
	Banking
	//Chatting The mob is chatting with a NPC
	Chatting
	//MenuChoosing The mob is in a query menu
	MenuChoosing
	//Trading The mob is negotiating a trade.
	Trading
	//Dueling The mob is negotiating a duel.
	Dueling
	//Fighting The mob is fighting.
	Fighting
	//Batching The mob is performing a skill that repeats itself an arbitrary number of times.
	Batching
	//Sleeping The mob is using a bed or sleeping bag, and trying to solve a CAPTCHA
	Sleeping
)

//Mob Any data structure that represents a mobile Entity within the game world should be able to implement this
type Mob interface {
	// Returns the current MobState of this Mob
	State() MobState
	Entity
}
