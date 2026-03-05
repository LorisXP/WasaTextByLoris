package entity

type Message struct {
	MessageID          int
	ConversationID     int
	Between_users      bool // destinazione: true=utente-utente, false=utente-gruppo
	SourceBetweenUsers bool // sorgente del messaggio originale da inoltrare
	Sender             int
	Receiver           int // groupID if between_users is false
	Type               string
	Sent_at            string
	Status             string
	ContentType        string
	Content            string
}
