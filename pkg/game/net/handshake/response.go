/*
 * Copyright (c) 2020 Zachariah Knight <aeros.storkpk@gmail.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

package handshake

import (
	"github.com/spkaeros/rscgo/pkg/throttle"
)

var LoginThrottle = ipThrottle.NewThrottle()
var RegisterThrottle = ipThrottle.NewThrottle()

type (
	//ResponseType A networking handshake response identifier code.
	ResponseType int
	//ResponseCode A valid networking handshake response code.  In practice as long as this type is being used,
	// it's a safe bet that it is a valid handshake response.
	ResponseCode int
)

//IsValid is used to determine whether the ResponseCode is for a successful handshake or not.
// Returns true if the handshake was a success and the client is now logged in, otherwise returns false.
func (r ResponseCode) IsValid() bool {
	for _, i := range [...]ResponseCode{ResponseLoginSuccess, ResponseReconnected, ResponseModerator,
		ResponseAdministrator, ResponseRegisterSuccess} {
		if i == r {
			return true
		}
	}
	return int(r)&64 == 64
}

const (
	//ResponseLoginSuccess is sent when we have successfully validated and started to load the player profile
	// that was requested.  Generally, once a handshake receives this response, the handshake is over and the
	// client will want more data about the characters identity and the world where it is currently at.
	// Network protocols <= 204 rely on this.  See ResponseLoginAcceptBit for protocols > 204
	ResponseLoginSuccess ResponseCode = iota
	//ResponseReconnected is sent after a brief connectivity problem; indicates that the client should skip
	// certain character-specific cleanup routines since the current variables it's got are good still.
	ResponseReconnected
	//ResponseRegisterSuccess Sent when a new player was created and saved successfully.
	ResponseRegisterSuccess
	//ResponseBadPassword is sent when the username or password supplied weren't valid.
	// This is a safe and very generic response that informs the client the username and password
	// they had provided in the handshake could not be validated with our records.
	// This response could indicate either a bad username OR a bad password.  Vague by design.
	ResponseBadPassword
	//ResponseLoggedIn is sent when the server realizes that the account requested is already logged in.
	ResponseLoggedIn
	//ResponseUpdated is sent when the client version number being used during the handshake was too old for this
	// server, and this indicates to the client that it needs to update to the latest client version before
	// being able to handshake with us.
	//TODO: Consider to put some less serious update ``suggestion'' responses in the placeholders in case
	// there is no protocol breaking incompatibility involved and the only problems from not updating would
	// be inconsequential to everyone involved?
	ResponseUpdated
	//ResponseSingleIp is sent when the server realizes that handshaking client is already online with another
	// character.  The message states IP limits of 1 character per IP, but traditionally this is triggered by
	// using 2 clients at once with one set of cache files, which happen to store an ID that was made to assist
	// in tracking our user bases client<->server interactions without relying on IP addresses and that.
	ResponseSingleIp
	//ResponseSpamTimeout is sent when the server has received a lot of invalid login attempts from one IP address,
	// consecutively, within a relatively short time frame.  For the sake of security, the server had decided to
	// temporarily deny login attempts from the handshaking client's IP address on the previous failed handshake.
	// Breaching this invalid login spamming threshold can only mean the client is trying to steal a character
	// from somebody using the slowest bruteforce attempt on earth, or have forgotten their own login credentials
	// and are trying to remember them.
	// In either case, there are better ways to reach their goal, and hopefully getting this reply for the next
	// 5 minutes will make them realize this.
	ResponseSpamTimeout
	// TODO: With the recent use of the TLS stack over RSA/ISAAC, is codes 8/9 legacy now or what?
	//ResponseServerRejection is sent when the game server rejects the session being used.
	ResponseServerRejection
	//ResponseLoginServerRejection is sent when the login server rejects the session being used.
	ResponseLoginServerRejection
	// TODO: Distinct from LoggedIn?  LoggedIn maybe for logged in on another world, and this same world??
	//ResponseInUse is sent when the server realizes that the username provided is already in use.
	ResponseInUse
	//ResponseTempBan is sent when the player profile requested was found, but is temporarily banned from
	// logging in.
	ResponseTempBan
	//ResponsePernBan is sent when the player profile requested was found, but is banned from logging in,
	// and apparently will never be unbanned.
	ResponsePermBan
	//ResponseLoggedIn2 Registration response check reads this as the same thing that LoggedIn says.  Unsure why its here
	ResponseLoggedIn2
	//ResponseWorldFull is sent when the world is completely out of player slots.
	// This requires a lot of players to happen.
	ResponseWorldFull
	//ResponseMembersWorld was for segregating P2P and F2P players and the exclusive P2P content from free
	// non-paying players of the game.  I never liked that.
	ResponseMembersWorld
	//ResponseNoReply is sent, I think, when the underlying data source for player saves is not reachable.
	ResponseNoReply
	//ResponseDecodeFailure is sent when we could not unmarshal the character from the player save data.
	ResponseDecodeFailure
	//ResponseSuspectedStolenLocked is sent when there's strong reason to believe the account has been stolen
	// and as such was disabled pending resolution via help from human customer support.
	ResponseSuspectedStolenLocked
	//ResponseBadInputLength Either username or password length was out of bounds.
	ResponseBadInputLength
	//ResponseMismatchedLogin is sent probably when the player data is accessed via a separate data service,
	// and that service has some sort of conflicting data to what this server has.
	ResponseMismatchedLogin
	// FIXME: VVVVVVV Legacy and thus a placeholder for a new reply if needed VVVVVVV
	//ResponseNeedClassicAccount was sent originally to any post-2005ish newly-registered characters,
	// because they were deprecating this version of the game for RS2.  Shame, shame.
	ResponseNeedClassicAccount
	//ResponseSuspectedStolen is sent when there's strong reason to believe the account is being targeted
	// to be stolen and as such was disabled pending a change of password.  This is a less worrisome message.
	ResponseSuspectedStolen
	//ResponsePlaceholder Placeholder for new responses.
	ResponsePlaceholder
	//ResponseModerator Used to indicate to the client that the character it wants to load is a player moderator
	ResponseModerator
	//ResponseAdministrator Used to indicate to the client that the character it wants to load is a stadd member/admin/owner.
	ResponseAdministrator
	//ResponseLoginAcceptBit bitmask to use for valid login response.  the network protocol of revisions > 204 rely on this.
	ResponseLoginAcceptBit = 1 << 6
	ResponseUsernameTaken  = ResponseBadPassword
)

const (
	LoginCode ResponseType = iota
	RegisterCode
)