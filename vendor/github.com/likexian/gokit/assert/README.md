# GoKit - assert

Assert kits for Golang development.

## Installation

    go get -u github.com/likexian/gokit

## Importing

    import (
        "github.com/likexian/gokit/assert"
    )

## Documentation

Visit the docs on [GoDoc](https://godoc.org/github.com/likexian/gokit/assert)

## Example

### assert panic

```go
func willItPanic() {
    panic("failed")
}
assert.Panic(t, willItPanic)
```

### assert err is nil

```go
fp, err := os.Open("/data/dev/gokit/LICENSE")
assert.Nil(t, err)
```

### assert equal

```go
x := map[string]int{"a": 1, "b": 2}
y := map[string]int{"a": 1, "b": 2}
assert.Equal(t, x, y, "x shall equal to y")
```

### check string in array

```go
ok := assert.IsContains([]string{"a", "b", "c"}, "b")
if ok {
    fmt.Println("value in array")
} else {
    fmt.Println("value not in array")
}
```

### check string in interface array

```go
ok := assert.IsContains([]interface{}{0, "1", 2}, "1")
if ok {
    fmt.Println("value in array")
} else {
    fmt.Println("value not in array")
}
```

### check object in struct array

```go
ok := assert.IsContains([]A{A{0, 1}, A{1, 2}, A{1, 3}}, A{1, 2})
if ok {
    fmt.Println("value in array")
} else {
    fmt.Println("value not in array")
}
```

### a := c ? x : y

```go
a := 1
// b := a == 1 ? true : false
b := assert.If(a == 1, true, false)
```

## License

Copyright 2012-2020 [Li Kexian](https://www.likexian.com/)

Licensed under the Apache License 2.0

## Donation

If this project is helpful, please share it with friends.

If you want to thank me, you can [give me a cup of coffee](https://www.likexian.com/donate/).
