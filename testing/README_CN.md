# Assert & Require

[English](README.md) | [中文](README_CN.md)

Go-Spring Testing 提供了两个用于测试断言的包：`assert` 和 `require`。这两个包都提供了**流式 API**，便于编写清晰、可读性强的测试断言。

### assert

`assert` 包提供的断言函数在断言失败时不会终止测试函数的执行。

当 `assert` 包中的断言失败时，测试函数会继续执行后续的断言。这在希望在一次测试运行中报告多个失败的情况下非常有用。

### require

`require` 包提供的断言函数在断言失败时会立即停止测试函数的执行。

当 `require` 包中的断言失败时，测试函数会立即停止执行，后续断言将不再被检查。这在关键条件不满足时，后续断言可能会导致 panic
或其他问题的情况下非常有用。

### 基本示例

```go
package main

import (
	"testing"

	"github.com/go-spring/spring-base/testing/assert"
	"github.com/go-spring/spring-base/testing/require"
)

func TestSomething(t *testing.T) {
	// 使用 assert - 断言失败时测试会继续执行
	assert.That(t, "hello").Equal("hello")
	assert.Number(t, 42).GreaterThan(40)

	// 使用 require - 断言失败时测试会立即停止
	require.That(t, someValue).NotNil()

	// 类型专用断言
	assert.String(t, "user@example.com").IsEmail()
	assert.Number(t, 100).InRange(0, 200)
	assert.Slice(t, []int{1, 2, 3}).Contains(2)
}
```

## 许可证

Apache License 2.0
