package world

import (
	"github.com/spkaeros/rscgo/pkg/game/entity"
)

type ChatMessage struct {
	Owner  entity.MobileEntity
	Target entity.MobileEntity
	string
}

func NewTargetedMessage(owner entity.MobileEntity, target entity.MobileEntity, content string) ChatMessage {
	return ChatMessage{owner, target, content}
}

func NewChatMessage(owner entity.MobileEntity, content string) ChatMessage {
	return ChatMessage{owner, nil, content}
}
