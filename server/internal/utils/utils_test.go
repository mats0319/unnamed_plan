package utils

import (
	"testing"

	"github.com/google/uuid"
)

func TestGenerateRandomBytes(t *testing.T) {
	for range 5 {
		t.Log(string(GenerateRandomBytes(32)))
	}
}

func TestGenerateUUID(t *testing.T) {
	for range 5 {
		t.Log(uuid.New())
	}
}
