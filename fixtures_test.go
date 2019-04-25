package fixtures

import (
	"testing"
)

func TestFixture(t *testing.T) {
	fixtures, closeall := Fixtures([]string{"executable"}, nil)

	expectedLen := 3
	if len(fixtures) != expectedLen {
		t.Errorf("got fixture length %v; want %v", len(fixtures), expectedLen)
	}
	
	if err := closeall.Close(); err != nil {
		t.Fatal(err)
	}
}