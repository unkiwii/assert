# unkiwii/assert

[![Go Report Card](https://goreportcard.com/badge/github.com/unkiwii/assert)](https://goreportcard.com/report/github.com/unkiwii/assert)

A very simple package to compare test results in [Go](https://go.dev/)

## Very simple API

This package has only 5 functions:

  * FailOnError: fails the test if an error occurred
  * Nil: fails the test if the value isn't nil
  * Equals: fails the test if the values are no equal
  * IsError: fails the test if 2 errors are not equivalent
  * AsError: fails the test if 2 error types are not equivalent

Complete reference [here](https://pkg.go.dev/github.com/unkiwii/assert)

## Usage

With go modules just run:

```
$ go get github.com/unkiwii/assert
```

Import and use it:

```go
package something_test

import "github.com/unkiwii/assert"

func TestSomething(t *testing.T) {
	// setup your test case
	want := "the result you want"

	// run your test case
	got, err := ...

	// check the results
	assert.Nil(t, err)
	assert.Equals(t, got, want)
}
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md)

## License

This project is licensed under the terms of the [MIT license](https://opensource.org/license/mit/)

See [LICENSE](LICENSE)
