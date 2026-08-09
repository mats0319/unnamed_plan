package utils

import (
	"testing"

	"github.com/google/uuid"
)

func TestGenerateRandomStr(t *testing.T) {
	for range 5 {
		t.Log(GenerateRandomBytes[string](32))
	}
}

func TestGenerateUUID(t *testing.T) {
	for range 5 {
		t.Log(uuid.New())
	}
}
