# Assert & Require

[English](README.md) | [中文](README_CN.md)

Go-Spring Testing provides two packages for making assertions in tests: `assert` and `require`. Both packages offer a
fluent API for writing clear and expressive test assertions.

### assert

The `assert` package provides assertion functions that allow tests to continue running even when an assertion fails.

When an assertion in the `assert` package fails, the test function continues executing subsequent assertions. This is
useful when you want to report multiple failures in a single test run.

### require

The `require` package provides assertion functions that stop test execution immediately when an assertion fails.

When an assertion in the `require` package fails, the test function stops executing and no further assertions are
checked. This is useful when subsequent assertions would panic or cause other issues if a critical condition is not met.

### Basic Example

```go
package main

import (
	"testing"

	"github.com/go-spring/spring-base/testing/assert"
	"github.com/go-spring/spring-base/testing/require"
)

func TestSomething(t *testing.T) {
	// Using assert - test continues on failure
	assert.That(t, "hello").Equal("hello")
	assert.ThatNumber(t, 42).GreaterThan(40)

	// Using require - test stops on failure
	require.That(t, someValue).NotNil()

	// Type-specific assertions
	assert.ThatString(t, "user@example.com").IsEmail()
	assert.ThatNumber(t, 100).InRange(0, 200)
	assert.ThatSlice(t, []int{1, 2, 3}).Contains(2)
}
```

## License

Apache License 2.0
