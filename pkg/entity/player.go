package entity

import (
	"strconv"
	"time"
)

//Player Represents a single player.
type Player struct {
	Username      string
	UserBase37    uint64
	Password      string
	FriendList    map[uint64]bool
	IgnoreList    []uint64
	LocalPlayers  *List
	LocalObjects  *List
	Connected     bool
	Updating      bool
	Appearances   []int
	DatabaseIndex int
	Rank          int
	Appearance    *AppearanceTable
	Mob
}

//RunDistancedAction Creates a distanced action belonging to this player, that runs action once the player arrives at dest, or cancels if we become busy, or we become unreasonably far from dest.
func (p *Player) RunDistancedAction(dest Location, action func()) {
	go func() {
		// TODO: Is checking 10 times per second good enough?  Probably.
		for range time.Tick(100 * time.Millisecond) {
			if p.State != Idle {
				// We became busy somehow, so cancel
				return
			}
			if !p.WithinRange(dest, 16) {
				// Somehow our destination is no longer in our view area, so cancel.
				return
			}
			if p.WithinRange(dest, 1) {
				action()
				return
			}
		}
	}()
}

//Friends Returns true if specified username is in our friend list.
func (p *Player) Friends(other uint64) bool {
	for hash := range p.FriendList {
		if hash == other {
			return true
		}
	}
	return false
}

//Ignoring Returns true if specified username is in our ignore list.
func (p *Player) Ignoring(hash uint64) bool {
	for _, v := range p.IgnoreList {
		if v == hash {
			return true
		}
	}
	return false
}

//ChatBlocked Returns true if public chat is blocked for this player.
func (p *Player) ChatBlocked() bool {
	return p.Attributes.VarBool("chat_block", false)
}

//FriendBlocked Returns true if private chat is blocked for this player.
func (p *Player) FriendBlocked() bool {
	return p.Attributes.VarBool("friend_block", false)
}

//TradeBlocked Returns true if trade requests are blocked for this player.
func (p *Player) TradeBlocked() bool {
	return p.Attributes.VarBool("trade_block", false)
}

//DuelBlocked Returns true if duel requests are blocked for this player.
func (p *Player) DuelBlocked() bool {
	return p.Attributes.VarBool("duel_block", false)
}

//SetPrivacySettings Sets privacy settings to specified values.
func (p *Player) SetPrivacySettings(chatBlocked, friendBlocked, tradeBlocked, duelBlocked bool) {
	p.Attributes.SetVar(Attribute("chat_block"), chatBlocked)
	p.Attributes.SetVar(Attribute("friend_block"), friendBlocked)
	p.Attributes.SetVar(Attribute("trade_block"), tradeBlocked)
	p.Attributes.SetVar(Attribute("duel_block"), duelBlocked)
}

//SetClientSetting Sets the specified client setting to flag.
func (p *Player) SetClientSetting(id int, flag bool) {
	// TODO: Meaningful names mapped to IDs
	p.Attributes.SetVar(Attribute("client_setting_"+strconv.Itoa(id)), flag)
}

//GetClientSetting Looks up the client setting with the specified ID, and returns it.  If it can't be found, returns false.
func (p *Player) GetClientSetting(id int) bool {
	// TODO: Meaningful names mapped to IDs
	return p.Attributes.VarBool(Attribute("client_setting_"+strconv.Itoa(id)), false)
}

//IsFollowing Returns true if the player is following another mob, otherwise false.
func (p *Player) IsFollowing() bool {
	return p.FollowIndex() != -1
}

//ServerSeed Returns the seed for the ISAAC cipher provided by the server for this player, if set, otherwise returns 0
func (p *Player) ServerSeed() uint64 {
	return p.TransAttrs.VarLong("server_seed", 0)
}

//SetServerSeed Sets the player's stored server seed to seed for later comparison to ensure we decrypted the login block properly and the player received the proper seed.
func (p *Player) SetServerSeed(seed uint64) {
	p.TransAttrs.SetVar("server_seed", seed)
}

//Reconnecting Returns true if the player is reconnecting, false otherwise.
func (p *Player) Reconnecting() bool {
	return p.TransAttrs.VarBool("reconnecting", false)
}

//SetReconnecting Sets the player's reconnection status to flag.
func (p *Player) SetReconnecting(flag bool) {
	p.TransAttrs.SetVar("reconnecting", flag)
}

//SetFollowing Sets the transient attribute for storing the server index of the player we want to follow to index.
func (p *Player) SetFollowing(index int) {
	if index != -1 {
		p.TransAttrs.SetVar("plrfollowing", index)
		p.TransAttrs.SetVar("followrad", 2)
	} else {
		delete(p.TransAttrs, "plrfollowing")
		delete(p.TransAttrs, "followrad")
	}
}

//FollowRadius Returns the radius within which we should follow whatever mob we are following, or -1 if we aren't following anyone.
func (p *Player) FollowRadius() int {
	return p.TransAttrs.VarInt("followrad", -1)
}

//FollowIndex Returns the index of the mob we are following, or -1 if we aren't following anyone.
func (p *Player) FollowIndex() int {
	return p.TransAttrs.VarInt("plrfollowing", -1)
}

//ResetFollowing Resets the transient attribute for storing the server index of the player we want to follow.
func (p *Player) ResetFollowing() {
	p.SetFollowing(-1)
	p.ResetPath()
}

//ArmourPoints Returns the players armour points.
func (p *Player) ArmourPoints() int {
	return p.TransAttrs.VarInt("armour_points", 1)
}

//SetArmourPoints Sets the players armour points to i.
func (p *Player) SetArmourPoints(i int) {
	p.TransAttrs["armour_points"] = i
}

//PowerPoints Returns the players power points.
func (p *Player) PowerPoints() int {
	return p.TransAttrs.VarInt("power_points", 1)
}

//SetPowerPoints Sets the players power points to i
func (p *Player) SetPowerPoints(i int) {
	p.TransAttrs["power_points"] = i
}

//AimPoints Returns the players aim points
func (p *Player) AimPoints() int {
	return p.TransAttrs.VarInt("aim_points", 1)
}

//SetAimPoints Sets the players aim points to i.
func (p *Player) SetAimPoints(i int) {
	p.TransAttrs["aim_points"] = i
}

//MagicPoints Returns the players magic points
func (p *Player) MagicPoints() int {
	return p.TransAttrs.VarInt("magic_points", 1)
}

//SetMagicPoints Sets the players magic points to i
func (p *Player) SetMagicPoints(i int) {
	p.TransAttrs["magic_points"] = i
}

//PrayerPoints Returns the players prayer points
func (p *Player) PrayerPoints() int {
	return p.TransAttrs.VarInt("prayer_points", 1)
}

//SetPrayerPoints Sets the players prayer points to i
func (p *Player) SetPrayerPoints(i int) {
	p.TransAttrs["prayer_points"] = i
}

//RangedPoints Returns the players ranged points.
func (p *Player) RangedPoints() int {
	return p.TransAttrs.VarInt("ranged_points", 1)
}

//SetRangedPoints Sets the players ranged points tp i.
func (p *Player) SetRangedPoints(i int) {
	p.TransAttrs["ranged_points"] = i
}

//Fatigue Returns the players current fatigue.
func (p *Player) Fatigue() int {
	return p.Attributes.VarInt("fatigue", 0)
}

//SetFatigue Sets the players current fatigue to i.
func (p *Player) SetFatigue(i int) {
	p.Attributes["fatigue"] = i
}

//FightMode Returns the players current fight mode.
func (p *Player) FightMode() int {
	return p.Attributes.VarInt("fight_mode", 0)
}

//SetFightMode Sets the players fightmode to i.  0=all,1=attack,2=defense,3=strength
func (p *Player) SetFightMode(i int) {
	p.Attributes["fight_mode"] = i
}

//NearbyPlayers Returns the nearby players from the current and nearest adjacent regions in a slice.
func (p *Player) NearbyPlayers() (players []*Player) {
	for _, r := range SurroundingRegions(p.X, p.Y) {
		players = append(players, r.Players.NearbyPlayers(p)...)
	}

	return
}

//NewPlayer Returns a reference to a new player.
func NewPlayer() *Player {
	return &Player{Mob: Mob{Entity: Entity{Index: -1}, Skillset: &SkillTable{}, State: Idle, Attributes: make(AttributeList), TransAttrs: make(AttributeList)}, LocalPlayers: &List{}, LocalObjects: &List{}, Appearance: NewAppearanceTable(1, 2, true, 2, 8, 14, 0), Connected: false, FriendList: make(map[uint64]bool)}
}
