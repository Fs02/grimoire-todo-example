package paranoid

import (
	"log"
	"os"
)

var logger = log.New(os.Stderr, "warning: ", 0)

// Panic if err is not nil.
func Panic(err error, message string, args ...interface{}) {
	if err != nil {
		logger.Printf(message+"\n", args...)
		panic(err)
	}
}
