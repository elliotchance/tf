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
