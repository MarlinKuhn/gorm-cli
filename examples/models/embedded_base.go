package models

type EmbeddedMeta struct {
	Code string
	Pets []Pet
}

type EmbeddedAudit struct {
	EmbeddedMeta
	Company Company
}
