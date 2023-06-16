package fmt

import (
	"fmt"
	"io"
)

func (_ FMT) Print(a ...interface{}) {
	fmt.Print(a...)
}

func (_ FMT) Println(a ...interface{}) {
	fmt.Println(a...)
}

func (_ FMT) Printf(f string, a ...interface{}) {
	fmt.Printf(f, a...)
}

func (_ FMT) Sprintln(a ...interface{}) string {
	return fmt.Sprintln(a...)
}

func (_ FMT) Sprintf(f string, a ...interface{}) string {
	return fmt.Sprintf(f, a...)
}

func (_ FMT) Fprint(w io.Writer, a ...interface{}) (int, error) {
	return fmt.Fprint(w, a...)
}

func (_ FMT) Fprintln(w io.Writer, a ...interface{}) (int, error) {
	return fmt.Fprintln(w, a...)
}

func (_ FMT) Fprintf(w io.Writer, f string, a ...interface{}) (int, error) {
	return fmt.Fprintf(w, f, a...)
}
