# Marla.ONE // Shepard

Guards the security of your data types.

**Work In Progress**

## Description

Library to achieve higher typesafety in Go by implementing the [Result](https://doc.rust-lang.org/std/result/enum.Result.html) and [Option](https://doc.rust-lang.org/std/option/enum.Option.html) enums of Rust. Powered with Go Generics.

## Packages
 
### [num](https://github.com/marlaone/shepard/tree/main/num)

Implements checked implementation for mul/div/add/sub to keep you save from overflows.

### [option](https://github.com/marlaone/shepard/tree/main/option)

Package with helper functions to map or check Options of two different types.

### [result](https://github.com/marlaone/shepard/tree/main/result)

Package with helper functions to map or check Result of two different types.

### [shepard_json](https://github.com/marlaone/shepard/tree/main/shepard_json)

Package implements types to json.(Un-)Marshal Results or Options.

### [iter](https://github.com/marlaone/shepard/tree/main/iter)

Package implements a type safe generic Iterator for any slices.

### [slice](https://github.com/marlaone/shepard/tree/main/slice)

Package implements a type safe generic Slice for any types.