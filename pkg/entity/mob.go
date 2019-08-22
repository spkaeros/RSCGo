/**
 * @Author: Zachariah Knight <zach>
 * @Date:   08-20-2019
 * @Email:  aeros.storkpk@gmail.com
 * @Project: RSCGo
 * @Last modified by:   zach
 * @Last modified time: 08-22-2019
 * @License: Use of this source code is governed by the MIT license that can be found in the LICENSE file.
 * @Copyright: Copyright (c) 2019 Zachariah Knight <aeros.storkpk@gmail.com>
 */

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
