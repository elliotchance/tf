# tf

tf is a microframework for parametrized testing of functions in Go.

I wrote this because I was tired of creating a []struct{} fixture for most of
my tests. I knew there had to be an easier and more reliable way.

It offers a simple and intuitive syntax for tests by wrapping the function:

```go
import "github.com/elliotchance/tf"

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

# Testing Struct Functions

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

