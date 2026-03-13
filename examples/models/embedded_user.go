package models

type EmbeddedUser struct {
	EmbeddedAudit
	Company Company
	Name    string
}
