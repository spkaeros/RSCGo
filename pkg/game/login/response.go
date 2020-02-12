/*
 * Copyright (c) 2020 Zachariah Knight <aeros.storkpk@gmail.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

package login

import (
	"time"

	"github.com/spkaeros/rscgo/pkg/engine/tasks"
	"github.com/spkaeros/rscgo/pkg/game/world"
	"github.com/spkaeros/rscgo/pkg/ipthrottle"
	"github.com/spkaeros/rscgo/pkg/log"
)

var LoginThrottler = ipthrottle.NewThrottle()
var RegisterThrottler = ipthrottle.NewThrottle()

type (
	ResponseType int
	ResponseCode int
	Response     struct {
		kind     ResponseType
		listener chan ResponseCode
		player   *world.Player
	}
)

func NewRegistrationListener(p *world.Player) *Response {
	return &Response{player: p, listener: make(chan ResponseCode), kind: RegisterCode}
}

func NewLoginListener(p *world.Player) *Response {
	return &Response{player: p, listener: make(chan ResponseCode), kind: LoginCode}
}

func (r ResponseCode) IsValid() bool {
	valid := [...]ResponseCode{ResponseLoginSuccess, ResponseReconnected, ResponseModerator, ResponseAdministrator}
	for _, i := range valid {
		if i == r {
			return true
		}
	}
	return false
}

const (
	ResponseLoginSuccess ResponseCode = iota
	ResponseReconnected
	ResponsePlaceholder1
	ResponseBadPassword
	ResponseLoggedIn
	ResponseUpdated
	ResponseSingleIp
	ResponseSpamTimeout
	ResponseServerRejection
	ResponseLoginServerRejection
	ResponseInUse // TODO: Distinction between logged in?
	ResponseTempBan
	ResponsePermBan
	ResponsePlaceholder2
	ResponseWorldFull
	ResponseMembersWorld
	ResponseNoReply
	ResponseDecodeFailure
	ResponseSuspectedStolenLocked
	ResponsePlaceholder3
	ResponseMismatchedLogin
	ResponseNeedClassicAccount
	ResponseSuspectedStolen
	ResponsePlaceholder4
	ResponseModerator
	ResponseAdministrator
)
const (
	LoginCode ResponseType = iota
	RegisterCode
)

const (
	ResponseRegisterSuccess ResponseCode = 2 + iota
	ResponseUsernameTaken
)

//ResponseListener This method will block until a byte is sent down the reply channel with the login Response to send to the client, or if this doesn't occur, it will timeout after 10 seconds.
func (r *Response) ResponseListener() chan ResponseCode {
	// schedules the channel listener on the game engines thread
	tasks.TickerList.Add("playerCreating", func() bool {
		defer close(r.listener)
		select {
		case code := <-r.listener:
			switch r.kind {
			case LoginCode:
				r.player.SendPacket(world.SessionResponse(int(code)))
				if code.IsValid() {
					r.player.Initialize()
					log.Info.Printf("Registered: %v\n", r.player)
					return true
				}
				r.player.Destroy()
				if code == ResponseBadPassword {
					LoginThrottler.Add(r.player.CurrentIP())
				}
				log.Info.Printf("Denied: %v (ResponseCode='%v')\n", r.player.String(), r)
			case RegisterCode:
				r.player.SendPacket(world.SessionResponse(int(code)))
				r.player.Destroy()
				RegisterThrottler.Add(r.player.CurrentIP())
				return true
			}
			return true
		case <-time.After(time.Second * 10):
			r.player.SendPacket(world.SessionResponse(-1))
			r.player.Destroy()
			return true
		}
	})
	return r.listener
}
