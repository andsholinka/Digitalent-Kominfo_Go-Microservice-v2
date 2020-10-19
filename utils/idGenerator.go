package utils

import "github.com/google/uuid"

func IdGenerator() string{
	id := uuid.New()
	return id.String()
}