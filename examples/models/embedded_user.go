package models

type EmbeddedUser struct {
	EmbeddedAudit
	Company Company
	Name    string
}

type TaggedEmbeddedUser struct {
	Filter EmbeddedAudit `gorm:"embedded;embeddedPrefix:filter_"`
	Name   string
}
