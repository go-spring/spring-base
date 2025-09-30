# spring-base

<div>
   <img src="https://img.shields.io/github/license/go-spring/spring-base" alt="license"/>
   <img src="https://img.shields.io/github/go-mod/go-version/go-spring/spring-base" alt="go-version"/>
   <img src="https://img.shields.io/github/v/release/go-spring/spring-base?include_prereleases" alt="release"/>
   <a href="https://codecov.io/gh/go-spring/spring-base" > 
      <img src="https://codecov.io/gh/go-spring/spring-base/graph/badge.svg?token=SX7CV1T0O8" alt="test-coverage"/>
   </a>
</div>

> 该项目已经正式发布，欢迎使用！

为 `go-spring` 框架提供基础功能支持的库集合。

## `barky` - 分层键值数据处理包

`barky` 包提供了处理分层键值数据结构的工具，主要用于处理 `JSON`、`YAML` 或 `TOML` 等配置格式中的嵌套数据。

- `flatten` - 将嵌套数据结构展开为单层结构。

## `testing` - 测试工具包

`testing` 目录包含了一套完整的测试工具，提供断言和验证功能。

- `assert` - 非中断式断言包

提供测试断言帮助工具，采用功能性和流畅的断言风格。当断言失败时，测试将继续运行。

- `require` - 中断式断言包

提供在失败时停止测试执行的断言帮助工具。

## 许可证

Apache License Version 2.0