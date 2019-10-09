package typictx

import (
	"fmt"
)

func invalidContextError(msg string) error {
	return fmt.Errorf("Invalid Context: %s", msg)
}
