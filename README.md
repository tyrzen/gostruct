# gostruct

`gostruct` is a Go package that provides functions for working with structs and their tags using reflection.

## Installation

To use `gostruct`, you need to install Go and set your Go workspace first.

1. Download and install it:

```sh
$ go get -u github.com/openai/gostruct
```

2. Import it in your code:

```go
import "github.com/openai/gostruct"
```

## Usage

The `gostruct` package provides the following functions:

### MapStructFieldTags

```go
func MapStructFieldTags[T any](key string) map[string]string
```

MapStructFieldTags takes a generic type T and a string key. It returns a map of field names to their corresponding tag values for all fields in the struct type T that contain a tag key that matches the given key. It uses reflection to get the field tags.

### GetTagValue

```go
func GetTagValue(tag reflect.StructTag, key string) (string, bool)
```

GetTagValue is designed to work with MapStructFieldTags and retrieves the value in the given key from struct fields.

## Examples

```go
type Person struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string
}

func main() {
	tags := gostruct.MapStructFieldTags[Person]("json")
	fmt.Println(tags) // Output: map[Name:name Age:age]
}
```

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

