package entity

type Group struct {
	GroupID int
	Name    string
	Photo   string
	AdminID int
}

// GroupInfoResponse è la struct di risposta per getGroupInfo
type GroupInfoResponse struct {
	GroupID       int              `json:"groupID"`
	Name          string           `json:"name"`
	Photo         string           `json:"photo"`
	Administrator string           `json:"administrator"`
	Members       []MemberResponse `json:"members"`
}

// MemberResponse è la struct per ogni membro
type MemberResponse struct {
	UserName string `json:"userName"`
	Photo    string `json:"photo"`
}
