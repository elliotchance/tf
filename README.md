# tf

tf is a microframework for parametrized testing of functions in Go.

* [Function](#functions)
* [ServeHTTP](#servehttp)

# Functions

It offers a simple and intuitive syntax for tests by wrapping the function:

```go
// Remainder returns the quotient and remainder from dividing two integers.
func Remainder(a, b int) (int, int) {
    return a / b, a % b
}

func TestRemainder(t *testing.T) {
    Remainder := tf.Function(t, Remainder)

    Remainder(10, 3).Returns(3, 1)
    Remainder(10, 2).Returns(5, 0)
    Remainder(17, 7).Returns(2, 3)
}
```

Assertions are performed with [testify](https://github.com/stretchr/testify). If
an assertion fails it will point to the correct line so you do not need to
explicitly label tests.

The above test will output (in verbose mode):

```
=== RUN   TestRemainder
--- PASS: TestRemainder (0.00s)
=== RUN   TestRemainder/Remainder#1
--- PASS: TestRemainder/Remainder#1 (0.00s)
=== RUN   TestRemainder/Remainder#2
--- PASS: TestRemainder/Remainder#2 (0.00s)
=== RUN   TestRemainder/Remainder#3
--- PASS: TestRemainder/Remainder#3 (0.00s)
PASS
```

## Struct Functions

You can test struct functions by providing the struct value as the first
parameter followed by any function arguments, if any.

```go
type Item struct {
	a, b float64
}

func (i Item) Add(c float64) float64 {
	return i.a + i.b + c
}

func TestItem_Add(t *testing.T) {
	Sum := tf.Function(t, Item.Add)

	Sum(Item{1.3, 4.5}, 3.4).Returns(9.2)
}
```

# ServeHTTP

Super easy HTTP testing by using the ServeHTTP function. This means that you do
not have to run the server and it is compatible with all HTTP libraries and
frameworks but has all the functionality of the server itself.

The simplest example is to use the default muxer in the `http` package:

```go
http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello, World!")
})
```

And now we can write some tests:

```go
func TestHTTPRouter(t *testing.T) {
	run := tf.ServeHTTP(t, http.DefaultServeMux.ServeHTTP)

	run(&tf.HTTPTest{
		Path:         "/hello",
		Status:       http.StatusOK,
		ResponseBody: strings.NewReader("Hello, World!"),
	})

	run(&tf.HTTPTest{
		Path:   "/world",
		Status: http.StatusNotFound,
	})
}
```

It is compatible with all HTTP frameworks because they must all expose a
ServeHTTP which is the entry point for the request router/handler.

There are many more options for HTTPTest. Some HTTP tests require multiple
operations, you can use `MultiHTTPTest` for this:

```go
run(&tf.MultiHTTPTest{
	Steps: []*tf.HTTPTest{
		{
			Path:        "/save",
			Method:      http.MethodPut,
			RequestBody: strings.NewReader(`{"foo":"bar"}`),
			Status:      http.StatusCreated,
		},
		{
			Path:         "/fetch",
			Method:       http.MethodGet,
			Status:       http.StatusOK,
			ResponseBody: strings.NewReader(`{"foo":"bar"}`),
		},
	},
})
```

Each step will only proceed if the previous step was successful.
