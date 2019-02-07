package paranoid

type fn func()

// PanicFunc panics if err is not nil after executing function f.
func PanicFunc(err error, f fn, message string, args ...interface{}) {
	if err != nil {
		f()
		logger.Printf(message+"\n", args...)
		panic(err)
	}
}
