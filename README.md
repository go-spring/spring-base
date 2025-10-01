# spring-base

<div>
   <img src="https://img.shields.io/github/license/go-spring/spring-base" alt="license"/>
   <img src="https://img.shields.io/github/go-mod/go-version/go-spring/spring-base" alt="go-version"/>
   <img src="https://img.shields.io/github/v/release/go-spring/spring-base?include_prereleases" alt="release"/>
</div>

[English](README.md) | [中文](README_CN.md)

> The project has been officially released, welcome to use!

A collection of foundational libraries that provide core support for the `go-spring` framework.

## `barky` - Hierarchical Key-Value Data Processing

The `barky` package offers tools for working with hierarchical key-value data structures, mainly for handling nested
data in configuration formats such as `JSON`, `YAML`, or `TOML`.

- `flatten` - Unfold nested data structures into a single-layer structure.

## `testing` - Testing Utilities

The `testing` directory contains a complete set of testing utilities that provide assertions and validations.

* `assert` – Non-blocking Assertions

Provides helper utilities for test assertions with a functional and fluent style. When an assertion fails, the test
continues executing.

* `require` – Blocking Assertions

Provides helper utilities for test assertions that stop test execution immediately upon failure.

## License

Apache License Version 2.0
