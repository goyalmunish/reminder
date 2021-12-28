package models

/*
BaseStruct represents set of common fields
*/
type BaseStruct struct {
	UpdatedAt int64 `json:"updated_at,omitempty"`
	CreatedAt int64 `json:"created_at,omitempty"`
}
