// Package tf is a microframework for parametrized testing of functions.
//
// I wrote this because I was tired of creating a []struct{} fixture for most of
// my tests. I knew there had to be an easier and more reliable way.
//
// It offers a simple and intuitive syntax for tests by wrapping the function:
//
//   // Remainder returns the quotient and remainder from dividing two integers.
//   func Remainder(a, b int) (int, int) {
//       return a / b, a % b
//   }
//
//   func TestRemainder(t *testing.T) {
//       Remainder := tf.Function(t, Remainder)
//
//       Remainder(10, 3).Returns(3, 1)
//       Remainder(10, 2).Returns(5, 0)
//       Remainder(17, 7).Returns(2, 3)
//   }
//
// Assertions are performed with github.com/stretchr/testify/assert. If an
// assertion fails it will point to the correct line so you do not need to
// explicitly label tests.
//
// The above test will output (in verbose mode):
//
//   === RUN   TestRemainder
//   --- PASS: TestRemainder (0.00s)
//   === RUN   TestRemainder/Remainder#1
//   --- PASS: TestRemainder/Remainder#1 (0.00s)
//   === RUN   TestRemainder/Remainder#2
//   --- PASS: TestRemainder/Remainder#2 (0.00s)
//   === RUN   TestRemainder/Remainder#3
//   --- PASS: TestRemainder/Remainder#3 (0.00s)
//   PASS
//
package tf

import (
	"testing"
	"fmt"
	"reflect"
	"github.com/stretchr/testify/assert"
	"runtime"
	"strings"
)

type F struct {
	t      *testing.T
	fn     interface{}
	args   []interface{}
	fnName string
}

var funcMap = map[string]int{}

func (f *F) getTestName() string {
	if _, ok := funcMap[f.fnName]; !ok {
		funcMap[f.fnName] = 0
	}

	funcMap[f.fnName] += 1

	return fmt.Sprintf("%s#%d", f.fnName, funcMap[f.fnName])
}

func stringValue(x interface{}) string {
	switch a := x.(type) {
	case string:
		return a

	case fmt.Stringer:
		return a.String()

	default:
		return fmt.Sprintf("%#+v", a)
	}
}

// Panic tests if the function will panic. You can provide a single optional
// expectedValue that is tested against the panic value. If the expectedValue is
// a string but the actual panic is not a string it will try to convert the
// panic value into a string with:
//
//   fmt.Sprintf("%#+v", r)
//
// For example each of the following are more verbose versions of the previous:
//
//
func (f *F) Panics(expectedValue ...interface{}) {
	if len(expectedValue) > 1 {
		f.t.Fatal("more than one value provided for Panics()")
	}

	f.t.Run(f.getTestName(), func(t *testing.T) {
		defer func() {
			r := recover()
			if assert.NotNil(t, r, "did not panic") && len(expectedValue) > 0 {
				if expected, ok := expectedValue[0].(string); ok {
					assert.Equal(t, expected, stringValue(r))
				} else {
					assert.Equal(t, expected, r)
				}
			}
		}()

		f.invoke(t)
	})
}

func (f *F) Returns(expected ...interface{}) {
	f.t.Run(f.getTestName(), func(t *testing.T) {
		assert.Equal(t, expected, f.invoke(t))
	})
}

func (f *F) invoke(t *testing.T) []interface{} {
	args := []reflect.Value{}
	for argIndex, arg := range f.args {
		if arg == nil {
			t.Fatalf("ambiguous nil for argument %d", argIndex+1)
		}

		args = append(args, reflect.ValueOf(arg))
	}

	returns := []interface{}{}
	for _, r := range reflect.ValueOf(f.fn).Call(args) {
		returns = append(returns, r.Interface())
	}

	return returns
}

func getFunctionName(fn interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	parts := strings.Split(name, ".")

	return parts[len(parts)-1]
}

func Function(t *testing.T, fn interface{}) func(args ...interface{}) *F {
	return func(args ...interface{}) *F {
		return &F{
			t:      t,
			fn:     fn,
			args:   args,
			fnName: getFunctionName(fn),
		}
	}
}
