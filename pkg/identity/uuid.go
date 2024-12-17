package identity

import "github.com/google/uuid"

// NewUUID creates a UUIDV7 or panics
func NewUUID() string {
	newUUID, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}

	return newUUID.String()
}
