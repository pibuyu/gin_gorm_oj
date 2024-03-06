package test

import "testing"
import (
	"fmt"
	"github.com/satori/go.uuid"
)

func TestGenerateUUID(t *testing.T) {
	// Creating UUID Version 4
	u1 := uuid.NewV4()
	fmt.Printf("UUIDv4: %s\n", u1)

	// Parsing UUID from string input
	u2, err := uuid.FromString("e37edbe9-7e2f-4a12-b24f-07984d1d49bf")
	if err != nil {
		fmt.Printf("Something gone wrong: %s", err)
	}
	fmt.Printf("Successfully parsed: %s", u2)
}
