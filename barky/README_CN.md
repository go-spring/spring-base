# barky

[English](README.md) | [中文](README_CN.md)

`barky` 是一个用于处理分层键值对数据结构的 Go 语言工具包，主要用于处理 JSON、YAML 或 TOML 等配置格式中的嵌套数据。

## 功能特性

### 1. 数据扁平化 (Flattening)

- 将嵌套的 map、slice 和数组转换为扁平的 `map[string]string`
- 使用点号表示法处理 map（如 `db.hosts`）
- 使用索引表示法处理数组/切片（如 `hosts[0]`）
- 示例：`{"db": {"hosts": ["a", "b"]}}` 转换为 `{"db.hosts[0]": "a", "db.hosts[1]": "b"}`

### 2. 路径处理 (Path handling)

- 定义了 Path 抽象，将分层键表示为类型化段的序列（map 键或数组索引）
- 支持将字符串路径（如 "foo.bar[0]"）解析为 Path 对象
- 支持将 Path 对象重新组合为字符串路径

### 3. 存储管理 (Storage)

- Storage 类型管理扁平化的键值对集合
- 内部构建和维护分层树结构，防止属性冲突
- 关联值与其来源文件，支持多文件合并和来源跟踪

### 4. 查询功能 (Querying)

- 提供辅助方法检索值
- 检查键是否存在
- 枚举子键
- 按确定顺序迭代

## 典型场景

1. 将不同来源的配置文件标准化为扁平的键值对映射，便于比较、合并或差异分析
2. 使用简单字符串路径查询嵌套数据，而无需直接处理反射或嵌套映射结构
3. 构建需要统一来自多个文件的结构化数据的工具，同时保留来源信息并防止冲突

## 使用示例

```go
package main

import (
	"fmt"
	"github.com/go-spring/spring-base/barky"
)

func main() {
	// 创建嵌套数据结构
	data := map[string]interface{}{
		"database": map[string]interface{}{
			"host": "localhost",
			"port": 5432,
			"credentials": map[string]interface{}{
				"username": "admin",
				"password": "secret",
			},
		},
		"features": []interface{}{
			"feature1",
			"feature2",
			map[string]interface{}{
				"name":    "feature3",
				"enabled": true,
			},
		},
	}

	// 扁平化数据
	flat := barky.FlattenMap(data)

	// 输出扁平化结果
	for key, value := range flat {
		fmt.Printf("%s: %s\n", key, value)
	}

	// 使用 Storage 管理数据
	storage := barky.NewStorage()
	fileID := storage.AddFile("config.yaml")

	// 设置值
	storage.Set("database.host", "localhost", fileID)
	storage.Set("database.port", "5432", fileID)

	// 查询值
	host := storage.Get("database.host")
	fmt.Printf("Database host: %s\n", host)

	// 检查键是否存在
	if storage.Has("database.host") {
		fmt.Println("Database host exists")
	}
}
```

## 许可证

Apache License 2.0