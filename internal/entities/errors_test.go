package entities

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrors(t *testing.T) {
	err := fmt.Errorf("test")
	stackTrace := NewStackTrace("test-id")
	stackTrace.Add(err)
	stackTrace.Add(err)
	stackTrace.Add(err)
	stackTrace.Add(err)
	stackTrace.Add(ErrorInitializeFailed)
	fmt.Println(stackTrace.Error())
	fmt.Print(stackTrace.StackTrace())

	t.Log(errors.Is(stackTrace, ErrorInitializeFailed))
	t.Log(errors.Is(stackTrace, err))
}
