package model

// Columns stores information about column
// it does not store relation info
type Column struct {
	Name       string
	Type       string
	IsArray    bool
	Dimensions int
	IsNullable bool
	IsPK       bool
	IsFK       bool
}
