package typrest

import (
	"io"
	"os"
)

// Stdout standard output
var Stdout io.Writer = os.Stdout
