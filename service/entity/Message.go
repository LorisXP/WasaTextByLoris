package entity

type Message struct {
	MessageID      int
	ConversationID int
	Between_users  bool
	Sender         int
	Receiver       int //groupID if between_users is false
	Sent_at        string
	Status         string
	Type           string
	Content        string
}
