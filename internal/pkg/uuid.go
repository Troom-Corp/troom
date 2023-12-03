package pkg

import "github.com/google/uuid"

func GenerateUUID() string {
	uuidCode := uuid.New()
	return uuidCode.String()
}
