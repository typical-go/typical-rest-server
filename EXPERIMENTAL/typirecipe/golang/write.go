package golang

import (
	"fmt"
	"io"
)

func write(w io.Writer, str string) (n int, err error) {
	return w.Write([]byte(str))
}
func writef(w io.Writer, format string, args ...interface{}) (n int, err error) {
	return write(w, fmt.Sprintf(format, args...))
}

func writelnf(w io.Writer, format string, args ...interface{}) (n int, err error) {
	return write(w, fmt.Sprintf(format, args...)+"\n")
}

func writeln(w io.Writer, s string) (n int, err error) {
	return write(w, s+"\n")
}
