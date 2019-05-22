# tf

tf is a microframework for parametrized testing of functions in Go.

- [Function](#function)
  * [Grouping](#grouping)
  * [Struct Functions](#struct-functions)
- [ServeHTTP](#servehttp)
- [Environment Variables](#environment-variables)

# Function

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

## Grouping

Use `NamedFunction` to specify a custom name for the function/group:

```go
func TestNamedSum(t *testing.T) {
	Sum := tf.NamedFunction(t, "Sum1", Item.Add)

	Sum(Item{1.3, 4.5}, 3.4).Returns(9.2)
	Sum(Item{1.3, 4.6}, 3.5).Returns(9.4)

	Sum = tf.NamedFunction(t, "Sum2", Item.Add)

	Sum(Item{1.3, 14.5}, 3.4).Returns(19.2)
	Sum(Item{21.3, 4.6}, 3.5).Returns(29.4)
}
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

# Environment Variables

`SetEnv` sets an environment variable and returns a reset function to ensure
the environment is always returned to it's previous state:

```go
resetEnv := tf.SetEnv(t, "HOME", "/somewhere/else")
defer resetEnv()
```

If you would like to set multiple environment variables, you can use `SetEnvs`
in the same way:

```go
resetEnv := tf.SetEnvs(t, map[string]string{
    "HOME":  "/somewhere/else",
    "DEBUG": "on",
})
defer resetEnv()
```
