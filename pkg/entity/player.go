package entity

import "strconv"

//TODO: Probably remove the Attribute type-alias.
//Attribute Type-alias for attribute names.  Might not need this, was just so to provide methods for them, don't think I'm doing it anymore
type Attribute string

//AttributeList A type alias for a map of strings to empty interfaces, to hold generic player information for easy serialization and to provide dynamic insertion/deletion of new player properties easily
type AttributeList map[Attribute]interface{}

//Player Represents a single player.
type Player struct {
	state         MobState
	Username      string
	UserBase37    uint64
	Password      string
	Path          *Pathway
	FriendList    map[uint64]bool
	IgnoreList    []uint64
	LocalPlayers  *EntityList
	LocalObjects  *EntityList
	Connected     bool
	Updating      bool
	Appearances   []int
	Skillset      *SkillTable
	DatabaseIndex int
	Rank          int
	Attributes    AttributeList
	TransAttrs    AttributeList
	Appearance    *AppearanceTable
	Entity
}

//SetVar Sets the attribute mapped at name to value in the attribute map.
func (attributes AttributeList) SetVar(name Attribute, value interface{}) {
	attributes[name] = value
}

//VarInt If there is an attribute assigned to the specified name, returns it.  Otherwise, returns zero
func (attributes AttributeList) VarInt(name Attribute, zero int) int {
	if _, ok := attributes[name].(int); !ok {
		attributes[name] = zero
	}

	return attributes[name].(int)
}

//VarLong If there is an attribute assigned to the specified name, returns it.  Otherwise, returns zero
func (attributes AttributeList) VarLong(name Attribute, zero uint64) uint64 {
	if _, ok := attributes[name].(uint64); !ok {
		attributes[name] = zero
	}

	return attributes[name].(uint64)
}

//VarBool If there is an attribute assigned to the specified name, returns it.  Otherwise, returns zero
func (attributes AttributeList) VarBool(name Attribute, zero bool) bool {
	if _, ok := attributes[name].(bool); !ok {
		attributes[name] = zero
	}

	return attributes[name].(bool)
}

//AppearanceTable Represents a mobs appearance.
type AppearanceTable struct {
	Head   int
	Body   int
	Male   bool
	Hair   int
	Top    int
	Bottom int
	Skin   int
}

//NewAppearanceTable Returns a reference to a new appearance table with specified parameters
func NewAppearanceTable(head, body int, male bool, hair, top, bottom, skin int) *AppearanceTable {
	return &AppearanceTable{head, body, male, hair, top, bottom, skin}
}

//SkillTable Represents a skill table for a mob.
type SkillTable struct {
	Current    [18]int
	Maximum    [18]int
	Experience [18]int
}

//SetIndex Sets the server index to idx
func (p *Player) SetIndex(idx int) {
	p.Index = idx
}

//VarEquipment Returns the attribute mapped to by name, and if it doesn't exist, returns 1.
func (p *Player) VarEquipment(name Attribute) int {
	return p.Attributes.VarInt(name, 1)
}

//TransVarInt Returns the transient attribute mapped to by name as an int, and if it doesn't exist, returns -1.
func (p *Player) TransVarInt(name Attribute) int {
	return p.TransAttrs.VarInt(name, -1)
}

//FriendsWith Returns true if specified username is in our friend list.
func (p *Player) FriendsWith(other uint64) bool {
	for hash := range p.FriendList {
		if hash == other {
			return true
		}
	}
	return false
}

//Ignored Returns true if specified username is in our ignore list.
func (p *Player) Ignored(hash uint64) bool {
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

//ResetPrivacySettings Resets privacy settings to specified values.
func (p *Player) ResetPrivacySettings(chatBlocked, friendBlocked, tradeBlocked, duelBlocked bool) {
	p.Attributes.SetVar(Attribute("chat_block"), chatBlocked)     // General?
	p.Attributes.SetVar(Attribute("friend_block"), friendBlocked) // Privacy?
	p.Attributes.SetVar(Attribute("trade_block"), tradeBlocked)   // Trade?
	p.Attributes.SetVar(Attribute("duel_block"), duelBlocked)     // Duel?
}

//SetClientSetting Sets the specified client setting to flag.
func (p *Player) SetClientSetting(id int, flag bool) {
	p.Attributes.SetVar(Attribute("client_setting_"+strconv.Itoa(id)), flag)
}

//GetClientSetting Looks up the client setting with the specified ID, and returns it.  If it can't be found, returns false.
func (p *Player) GetClientSetting(id int) bool {
	return p.Attributes.VarBool(Attribute("client_setting_"+strconv.Itoa(id)), false)
}

//SetPath Sets the player's current pathway to path.  If path is nil, effectively clears the players path.
func (p *Player) SetPath(path *Pathway) {
	p.Path = path
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
	return p.TransVarInt("followrad")
}

//FollowIndex Returns the index of the mob we are following, or -1 if we aren't following anyone.
func (p *Player) FollowIndex() int {
	return p.TransVarInt("plrfollowing")
}

//ResetFollowing Resets the transient attribute for storing the server index of the player we want to follow.
func (p *Player) ResetFollowing() {
	p.SetFollowing(-1)
	p.ResetPath()
}

//FinishedPath Returns true if the players path is nil or if we are already on the path's next tile.
func (p *Player) FinishedPath() bool {
	if p.Path == nil {
		return true
	}
	next := p.Path.NextTile(p.X, p.Y)
	return p.AtLocation(&next)
}

//ArmourPoints Returns the players armour points.
func (p *Player) ArmourPoints() int {
	return p.VarEquipment("armour_points")
}

//SetArmourPoints Sets the players armour points to i.
func (p *Player) SetArmourPoints(i int) {
	p.Attributes["armour_points"] = i
}

//PowerPoints Returns the players power points.
func (p *Player) PowerPoints() int {
	return p.VarEquipment("power_points")
}

//SetPowerPoints Sets the players power points to i
func (p *Player) SetPowerPoints(i int) {
	p.Attributes["power_points"] = i
}

//AimPoints Returns the players aim points
func (p *Player) AimPoints() int {
	return p.VarEquipment("aim_points")
}

//SetAimPoints Sets the players aim points to i.
func (p *Player) SetAimPoints(i int) {
	p.Attributes["aim_points"] = i
}

//MagicPoints Returns the players magic points
func (p *Player) MagicPoints() int {
	return p.VarEquipment("magic_points")
}

//SetMagicPoints Sets the players magic points to i
func (p *Player) SetMagicPoints(i int) {
	p.Attributes["magic_points"] = i
}

//PrayerPoints Returns the players prayer points
func (p *Player) PrayerPoints() int {
	return p.VarEquipment("prayer_points")
}

//SetPrayerPoints Sets the players prayer points to i
func (p *Player) SetPrayerPoints(i int) {
	p.Attributes["prayer_points"] = i
}

//RangedPoints Returns the players ranged points.
func (p *Player) RangedPoints() int {
	return p.VarEquipment("ranged_points")
}

//SetRangedPoints Sets the players ranged points tp i.
func (p *Player) SetRangedPoints(i int) {
	p.Attributes["ranged_points"] = i
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

//TraversePath If the player has a path, calling this method will change the players location to the next location
//  described by said Path data structure.  This should be called no more than once per game tick.
func (p *Player) TraversePath() {
	if p.Path == nil {
		return
	}
	path := p.Path
	if p.AtLocation(path.Waypoint(path.CurrentWaypoint)) {
		path.CurrentWaypoint++
	}
	newLocation := path.NextTile(p.X, p.Y)
	if path.CurrentWaypoint >= len(path.WaypointsX) || newLocation.X == -1 || newLocation.Y == -1 {
		p.ResetPath()
		return
	}
	p.TransAttrs["plrmoved"] = true
	p.SetLocation(newLocation)
}

//NearbyPlayers Returns the nearby players from the current and nearest adjacent regions in a slice.
func (p *Player) NearbyPlayers() (players []*Player) {
	for _, r := range SurroundingRegions(p.X, p.Y) {
		players = append(players, r.Players.NearbyPlayers(p)...)
	}

	return
}

//ResetPath Sets the players path to nil, to stop the traversal of the path instantly
func (p *Player) ResetPath() {
	p.SetPath(nil)
}

//UpdateDirection Updates the direction the player is facing based on where the player is trying to move, and
// where the player is currently at.
func (p *Player) UpdateDirection(destX, destY int) {
	sprites := [3][3]int{{3, 2, 1}, {4, -1, 0}, {5, 6, 7}}
	xIndex := p.X - destX + 1
	yIndex := p.Y - destY + 1
	if xIndex >= 0 && yIndex >= 0 && xIndex < 3 && yIndex < 3 {
		p.SetDirection(sprites[xIndex][yIndex])
	} else {
		p.SetDirection(int(North))
	}
}

//SetLocation Sets the players location.
func (p *Player) SetLocation(location Location) {
	p.SetCoords(location.X, location.Y)
}

//SetCoords Sets the players locations coordinates.
func (p *Player) SetCoords(x, y int) {
	curArea := GetRegion(p.X, p.Y)
	newArea := GetRegion(x, y)
	if newArea != curArea {
		curArea.Players.RemovePlayer(p)
		newArea.Players.AddPlayer(p)
	}
	p.UpdateDirection(x, y)
	p.X = x
	p.Y = y
}

//AtLocation Returns true if the player is at the specified location, otherwise returns false
func (p *Player) AtLocation(location *Location) bool {
	return p.AtCoords(location.X, location.Y)
}

//AtCoords Returns true if the player is at the specified coordinates, otherwise returns false
func (p *Player) AtCoords(x, y int) bool {
	return p.X == x && p.Y == y
}

//State Returns the players state.
func (p *Player) State() MobState {
	return p.state
}

//SetState Sets the players state.
func (p *Player) SetState(state MobState) {
	p.state = state
}

//Direction Returns the players direction.
func (p *Player) Direction() int {
	return p.Attributes.VarInt("direction", 0)
}

//SetDirection Sets the players direction.
func (p *Player) SetDirection(direction int) {
	p.Attributes["direction"] = direction
}

//NewPlayer Returns a reference to a new player.
func NewPlayer() *Player {
	return &Player{Entity: Entity{Index: -1}, state: Idle, Attributes: make(AttributeList), TransAttrs: make(AttributeList), LocalPlayers: &EntityList{}, LocalObjects: &EntityList{}, Skillset: &SkillTable{}, Appearance: NewAppearanceTable(1, 2, true, 2, 8, 14, 0), Connected: false, FriendList: make(map[uint64]bool)}
}
