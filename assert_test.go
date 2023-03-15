package assert_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/unkiwii/assert"
)

func TestFailOnError_failsWithAnError(t *testing.T) {
	var mock tMock
	assert.FailOnError(&mock, errors.New("abc"))

	errWhenDiff(t, "FailOnError", mock.fatal, "unexpected error: abc")
}

func TestFailOnError_successWithNil(t *testing.T) {
	var mock tMock
	assert.FailOnError(&mock, nil)

	errWhenDiff(t, "FailOnError", mock.fatal, "")
}

func TestNil_failsWithNotNil(t *testing.T) {
	var mock tMock
	assert.Nil(&mock, "abc")

	errWhenDiff(t, "Nil", mock.err, `assert.Nil failed
want: nil
 got: "abc"`)
}

func TestNil_successWithNil(t *testing.T) {
	var mock tMock
	assert.Nil(&mock, nil)

	errWhenDiff(t, "Nil", mock.err, "")
}

func TestNil_successWithNilValueTypedInterface(t *testing.T) {
	var buf *bytes.Buffer      // concrete nil value
	var reader io.Reader = buf // interface with nil value

	if reader == nil {
		t.Fatal("interface with nil value should not be nil")
	}

	var mock tMock
	assert.Nil(&mock, reader)

	errWhenDiff(t, "Nil", mock.err, "")
}

func TestEquals_failsWhenValuesAreDifferent(t *testing.T) {
	a := 42
	b := 24

	var mock tMock
	assert.Equals(&mock, a, b)

	errWhenDiff(t, "Equals", mock.err, `assert.Equals failed
want: int(24)
 got: int(42)`)
}

func TestEquals_failsWhenValuesHaveDifferentTypes(t *testing.T) {
	var a int32 = 42
	var b int64 = 42

	var mock tMock
	assert.Equals(&mock, a, b)

	errWhenDiff(t, "Equals", mock.err, `assert.Equals failed
want: int64(42)
 got: int32(42)`)
}

func TestEquals_successWhenValueAreTheSame(t *testing.T) {
	a, b := 3, 3

	var mock tMock
	assert.Equals(&mock, a, b)

	errWhenDiff(t, "Equals", mock.err, "")
}

func TestIsError_failsWhenErrorsAreDifferent(t *testing.T) {
	want := errors.New("some error")
	got := errors.New("some error")

	var mock tMock
	assert.IsError(&mock, got, want)

	errWhenDiff(t, "IsError", mock.err, `assert.IsError failed
want: *errors.errorString(some error)
 got: *errors.errorString(some error)`)
}

func TestIsError_successWhenErrorsAreDifferent(t *testing.T) {
	want := errors.New("some error")
	got := func() error {
		return want
	}()

	var mock tMock
	assert.IsError(&mock, got, want)

	errWhenDiff(t, "IsError", mock.err, "")
}

func TestAsError_failWhenErrorsAreDifferent(t *testing.T) {
	var target *os.PathError
	err := io.EOF

	var mock tMock
	assert.AsError(&mock, err, &target)

	errWhenDiff(t, "AsError", mock.err, `assert.AsError failed
want: *fs.PathError
 got: *errors.errorString`)
}

func TestAsError_failWhenTargetIsNotAPointer(t *testing.T) {
	defer func() {
		if e, ok := recover().(string); ok {
			errWhenDiff(t, "AsError", e, "errors: target must be a non-nil pointer")
		} else {
			t.Fatal("this should have panicked")
		}
	}()

	var target *os.PathError
	err := io.EOF

	var mock tMock
	assert.AsError(&mock, err, target)
}

func TestAsError_failWhenTargetIsIncorretType(t *testing.T) {
	defer func() {
		if e, ok := recover().(string); ok {
			errWhenDiff(t, "AsError", e, "errors: *target must be interface or implement error")
		} else {
			t.Fatal("this should have panicked")
		}
	}()

	var target os.PathError
	err := io.EOF

	var mock tMock
	assert.AsError(&mock, err, &target)
}

func TestAsError_success(t *testing.T) {
	var target *os.PathError
	err := &os.PathError{
		Op:   "read",
		Path: "/some/path",
		Err:  errors.New("some error"),
	}

	var mock tMock
	assert.AsError(&mock, err, &target)

	errWhenDiff(t, "AsError", mock.err, "")
	errWhenDiff(t, "AsError", err.Error(), target.Error())
}

func errWhenDiff(t *testing.T, name, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("%s\nwant: %q\n got: %q", name, want, got)
	}
}

type tMock struct {
	err   string
	fatal string
}

func (t *tMock) Errorf(format string, args ...interface{}) {
	t.err = fmt.Sprintf(format, args...)
}

func (t *tMock) Fatalf(format string, args ...interface{}) {
	t.fatal = fmt.Sprintf(format, args...)
}

func (t tMock) Helper() {}
