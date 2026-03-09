package entity

type Comment struct {
	CommentID      int
	SenderID       int
	MessageUserID  int
	MessageGroupID int
	Reaction       string
}
