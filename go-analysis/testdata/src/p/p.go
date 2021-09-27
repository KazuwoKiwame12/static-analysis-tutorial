package p

func notPrintfFunctionAll() {}

func funcWithEllipsis(args ...interface{}) {}

func printfLikeButWithStrings(format string, args ...string) {}

func printfLikeButWithBadFormat(format int, args ...string) {}

func printfLikeFuncf(format string, args ...interface{}) {}

func prinfLikeFuncWithReturnValuef(format string, args ...interface{}) string {
	return ""
}
