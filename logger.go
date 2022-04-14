package backend

func printLog(loggerFunc func(string), msg string) {
	if loggerFunc != nil {
		loggerFunc(msg)
	}
}
