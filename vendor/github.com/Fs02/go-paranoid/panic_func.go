package paranoid

import "fmt"

type fn func()

func PanicFunc(err error, f fn, message string, args ...interface{}) {
	if err != nil {
		f()
		panic(fmt.Errorf(message+": %v", append(args, err)...))
	}
}
