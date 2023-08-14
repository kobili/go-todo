package models

import "time"

type Todo struct {
	Text     string
	Metadata *TodoMetadata `json:",omitempty" bson:",omitempty"`
}

type TodoMetadata struct {
	Timestamp time.Time
	User      string
}
