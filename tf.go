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

func (f *F) Returns(expected ...interface{}) {
	if _, ok := funcMap[f.fnName]; !ok {
		funcMap[f.fnName] = 0
	}

	funcMap[f.fnName] += 1

	f.t.Run(fmt.Sprintf("%s#%d", f.fnName, funcMap[f.fnName]), func(t *testing.T) {
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

		assert.Equal(t, expected, returns)
	})
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
