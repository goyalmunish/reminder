package models

/*
BaseStruct represents set of common fields
*/
type BaseStruct struct {
	CreatedAt int64 `json:"created_at,omitempty"`
	UpdatedAt int64 `json:"updated_at,omitempty"`
}
