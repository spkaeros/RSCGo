package world

type ChatMessage struct {
	Owner  MobileEntity
	Target MobileEntity
	string
}

func NewTargetedMessage(owner MobileEntity, target MobileEntity, content string) ChatMessage {
	return ChatMessage{owner, target, content}
}

func NewChatMessage(owner MobileEntity, content string) ChatMessage {
	return ChatMessage{owner, nil, content}
}
