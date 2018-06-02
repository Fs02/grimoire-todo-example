package paranoid

import "fmt"

func Panic(err error, message string, args ...interface{}) {
	if err != nil {
		panic(fmt.Errorf(message+": %v", append(args, err)...))
	}
}
