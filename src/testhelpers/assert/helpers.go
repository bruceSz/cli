package assert

import (
	"testing"
	"fmt"
	"runtime"
	"strings"
)

// Fail reports a failure through
func Fail(t *testing.T, failureMessage string, msgAndArgs ...interface{}) bool {

	message := messageFromMsgAndArgs(msgAndArgs...)

	if len(message) > 0 {
		t.Errorf("\r%s\r\tLocation:\t%s\n\r\tError:\t\t%s\n\r\tMessages:\t%s\n\r", getWhitespaceString(), CallerInfo(), failureMessage, message)
	} else {
		t.Errorf("\r%s\r\tLocation:\t%s\n\r\tError:\t\t%s\n\r", getWhitespaceString(), CallerInfo(), failureMessage)
	}

	return false
}

/* CallerInfo is necessary because the assert functions use the testing object
internally, causing it to print the file:line of the assert method, rather than where
the problem actually occured in calling code.*/

// CallerInfo returns a string containing the file and line number of the assert call
// that failed.
func CallerInfo() string {

	file := ""
	line := 0
	ok := false

	for i := 0; ; i++ {
		_, file, line, ok = runtime.Caller(i)
		if !ok {
			return ""
		}
		parts := strings.Split(file, "/")
		dir := parts[len(parts)-2]
		file = parts[len(parts)-1]
		if (dir != "assert" && dir != "mock") || file == "mock_test.go" {
			break
		}
	}

	return fmt.Sprintf("%s:%d", file, line)
}

// getWhitespaceString returns a string that is long enough to overwrite the default
// output from the go testing framework.
func getWhitespaceString() string {

	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	parts := strings.Split(file, "/")
	file = parts[len(parts)-1]

	return strings.Repeat(" ", len(fmt.Sprintf("%s:%d:      ", file, line)))

}


func messageFromMsgAndArgs(msgAndArgs ...interface{}) string {
	if len(msgAndArgs) == 0 || msgAndArgs == nil {
		return ""
	}
	if len(msgAndArgs) == 1 {
		return msgAndArgs[0].(string)
	}
	if len(msgAndArgs) > 1 {
		return fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}
	return ""
}
