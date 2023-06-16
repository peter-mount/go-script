package fmt

import "fmt"

func (_ FMT) Appendf(b []byte, f string, a ...interface{}) []byte {
	return fmt.Appendf(b, f, a...)
}

func (_ FMT) Append(b []byte, a ...interface{}) []byte {
	return fmt.Append(b, a...)
}

func (_ FMT) Appendln(b []byte, a ...interface{}) []byte {
	return fmt.Appendln(b, a...)
}
