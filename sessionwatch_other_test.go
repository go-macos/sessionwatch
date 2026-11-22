//go:build !darwin

package sessionwatch

import (
	"context"
	"errors"
	"testing"
)

func TestWatchOther(t *testing.T) {
	if _, err := Watch(context.Background()); !errors.Is(err, ErrUnsupported) {
		t.Fatal(err)
	}
}
