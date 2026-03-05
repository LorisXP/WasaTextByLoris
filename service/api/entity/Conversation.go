package entity

type Conversation struct {
	ConversationID int
	Between_users  bool
	Sender         int
	Receiver       int // groupID if between_users is false
	LastMessageID  int
}
