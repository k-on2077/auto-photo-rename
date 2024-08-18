package common

type RenameRecord struct {
	OldName string `json:"OldName"`
	NewName string `json:"NewName"`
	Error   string `json:"Error"`
}
