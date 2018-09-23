package tf_test

import (
	"errors"
	"testing"

	"github.com/elliotchance/tf"
)

type Optional struct {
	Name string
}

func NewNil(o *Optional) string {
	if o == nil {
		o = &Optional{
			Name: "default",
		}
	}

	return o.Name
}

func TestNil(t *testing.T) {
	NewNil := tf.Function(t, NewNil)
	NewNil(nil).Returns("default")
}

type CustomError struct {
	err error
}

func (c CustomError) Error() string {
	return c.err.Error()
}

func TestError(t *testing.T) {
	t.Run("Catches basic error", func(t *testing.T) {
		NewError := tf.Function(t, func() error {
			return errors.New("some error")
		})

		NewError().Errors()
	})

	t.Run("Catches named error", func(t *testing.T) {
		NewError := tf.Function(t, func() error {
			return errors.New("some error")
		})

		NewError().Errors("some error")
	})

	t.Run("Catches particular class error", func(t *testing.T) {
		custom := CustomError{errors.New("custom error")}

		NewError := tf.Function(t, func() error {
			return custom
		})

		NewError().Errors(custom)
	})
}

type Item struct {
	a, b float64
}

func (i Item) Sum() float64 {
	return i.a + i.b
}

func (i Item) Add(c float64) float64 {
	return i.a + i.b + c
}

func TestCasting(t *testing.T) {
	t.Run("int8/int64 casting", func(t *testing.T) {
		Run := tf.Function(t, func(a, b int8) int64 {
			return int64(a + b)
		})

		Run(1, 2).Returns(3)
	})

	t.Run("nil in return casting", func(t *testing.T) {
		type MyStruct struct{}

		Run := tf.Function(t, func(a, b int8) (int64, *MyStruct) {
			return int64(a + b), nil
		})

		Run(1, 2).Returns(3, nil)
	})
}

func TestItem_Average(t *testing.T) {
	Sum := tf.Function(t, Item.Sum)

	Sum(Item{4.2, 5.1}).Returns(9.3)
}

func TestItem_Add(t *testing.T) {
	Sum := tf.Function(t, Item.Add)

	Sum(Item{1.3, 4.5}, 3.4).Returns(9.2)
}

func Remainder(a, b int) (int, int) {
	return a / b, a % b
}

func TestRemainder(t *testing.T) {
	Remainder := tf.Function(t, Remainder)

	Remainder(10, 3).Returns(3, 1)
	Remainder(10, 2).Returns(5, 0)
	Remainder(17, 7).Returns(2, 3)
}

func Booler(b bool) bool {
	return b
}

func TestTrueFalse(t *testing.T) {
	Booler := tf.Function(t, Booler)

	Booler(true).True()
	Booler(false).False()
}
