//go:build !darwin

package sessionwatch

import "context"

func watch(context.Context) (<-chan Event, error) { return nil, ErrUnsupported }
