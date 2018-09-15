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
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// F wrapper around a func which handles testing instance, agrs and reveals function name
type F struct {
	t         *testing.T
	fn        interface{}
	args      []interface{}
	fnArgsIn  []reflect.Type
	fnArgsOut []reflect.Type
	fnName    string
}

var (
	funcMap = map[string]int{}
)

// Returns matches if expected result matches actual
//
//	Remainder := tf.Function(t, func(a,b int) int { return a + b })
//	Remainder(1, 2).Returns(3)
//
func (f *F) Returns(expected ...interface{}) {
	if _, ok := funcMap[f.fnName]; !ok {
		funcMap[f.fnName] = 0
	}

	funcMap[f.fnName]++

	f.t.Run(fmt.Sprintf("%s#%d", f.fnName, funcMap[f.fnName]), func(t *testing.T) {
		// Casting calling arguments
		argsIn := make([]reflect.Value, len(f.args))
		for idx, arg := range f.args {
			if arg == nil {
				argsIn[idx] = reflect.Zero(f.fnArgsIn[idx])
			} else {
				argsIn[idx] = reflect.ValueOf(arg).Convert(f.fnArgsIn[idx])
			}
		}

		returns := make([]interface{}, len(f.fnArgsOut))
		for idx, r := range reflect.ValueOf(f.fn).Call(argsIn) {
			returns[idx] = r.Interface()
		}

		for idx, e := range expected {
			if e == nil {
				expected[idx] = reflect.Zero(f.fnArgsOut[idx]).Interface()
			} else {
				expected[idx] = reflect.ValueOf(e).Convert(f.fnArgsOut[idx]).Interface()
			}
		}

		assert.Equal(t, expected, returns)
	})
}

// True matches is function returns true as a result
//
//   func Switch() bool {
//       return true
//   }
//
//   func TestSwitch(t *testing.T) {
//       Switch := tf.Function(t, Switch)
//
//       Switch().True()
//   }
//
func (f *F) True() {
	f.Returns(true)
}

// False matches is function returns false as a result
//
//   func Switch() bool {
//       return false
//   }
//
//   func TestSwitch(t *testing.T) {
//       Switch := tf.Function(t, Switch)
//
//       Switch().False()
//   }
//
func (f *F) False() {
	f.Returns(false)
}

func getFunctionName(fn interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	parts := strings.Split(name, ".")

	return parts[len(parts)-1]
}

func getFunctionArgs(fn interface{}) []reflect.Type {
	ref := reflect.ValueOf(fn).Type()
	argsCount := ref.NumIn()
	args := make([]reflect.Type, argsCount)
	for i := 0; i < argsCount; i++ {
		args[i] = ref.In(i)
	}

	return args
}

func getFunctionReturns(fn interface{}) []reflect.Type {
	ref := reflect.ValueOf(fn).Type()
	argsCount := ref.NumOut()
	args := make([]reflect.Type, argsCount)
	for i := 0; i < argsCount; i++ {
		args[i] = ref.Out(i)
	}

	return args
}

// Function wraps fn into F testing type and returns back function to which you can use
// as a regular function in e.g:
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
func Function(t *testing.T, fn interface{}) func(args ...interface{}) *F {
	return func(args ...interface{}) *F {
		return &F{
			t:         t,
			fn:        fn,
			args:      args,
			fnArgsIn:  getFunctionArgs(fn),
			fnArgsOut: getFunctionReturns(fn),
			fnName:    getFunctionName(fn),
		}
	}
}
