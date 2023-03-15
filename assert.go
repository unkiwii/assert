package assert

import (
	"errors"
	"reflect"
)

// T is the minimal interface that we need to perform every assertion.
// *testing.T implements this interface so does other testing wrappers use the
// one you want
type T interface {
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Helper()
}

// FailOnError call t.Fatalf when error is not nil. This must be used to check
// that everything need to setup the test works as intended.
func FailOnError(t T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// Nil asserts that a value is nil or it will call t.Errorf when the value is
// not nil.
//
// This isn't equivalent of executing `v == nil`, it uses reflection to check so
// it will return true even for cases when v is an interface with a type but
// with a nil value (which will not be equal to nil).
//
// For more info on this see a video from Francesc Campoy, Understanding nil:
// https://www.youtube.com/watch?v=ynoY2xz-F8s
func Nil(t T, v interface{}) {
	t.Helper()
	if !isNil(v) {
		t.Errorf("assert.Nil failed\nwant: nil\n got: %#v", v)
	}
}

// Equals asserts that both values are equals using reflect.DeepEqual or it will
// call t.Errorf when the values are not equals.
func Equals(t T, got, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("assert.Equals failed\nwant: %T(%v)\n got: %T(%v)", want, want, got, got)
	}
}

// IsError uses errors.Is from the standard library to assert that both errors
// are equivalent.
//
// For information on how this works see: https://pkg.go.dev/errors#Is
func IsError(t T, got, want error) {
	t.Helper()
	if !errors.Is(got, want) {
		t.Errorf("assert.IsError failed\nwant: %T(%v)\n got: %T(%v)", want, want, got, got)
	}
}

// AsError uses errors.As from the standard library to assert that both errors
// are equivalent.
//
// For information on how this works see: https://pkg.go.dev/errors#As
func AsError(t T, got error, want interface{}) {
	t.Helper()
	if !errors.As(got, want) {
		// want will always be a pointer to something, here we get the type of
		// want removing one indirection. i.e: if want is *T here we got T
		// even when T is a pointer itself (**T transforms to *T)
		typ := reflect.Indirect(reflect.ValueOf(want)).Type()
		t.Errorf("assert.AsError failed\nwant: %v\n got: %T", typ, got)
	}
}

func isNil(v interface{}) bool {
	if v == nil {
		return true
	}

	value := reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.Interface, reflect.Chan, reflect.Slice, reflect.Map,
		reflect.Func, reflect.Pointer:
		return value.IsNil()
	}

	return false
}
