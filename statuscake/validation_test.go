package statuscake_test

import (
	"reflect"
	"testing"

	"github.com/StatusCakeDev/terraform-provider-statuscake/statuscake"
)

func TestIsEmailAddress(t *testing.T) {
	t.Run("returns no errors when the given value is a valid email address", func(t *testing.T) {
		_, errs := statuscake.IsEmailAddress("picard@starfleet.com", "email")
		if errs != nil {
			t.Error("expected no errors but errors were returned")
		}
	})

	t.Run("returns an error when the value is not of type string", func(t *testing.T) {
		expected := []string{"expected type of email to be string"}

		_, errs := statuscake.IsEmailAddress(1701, "email")
		if errs == nil {
			t.Error("expected errors but no errors were returned")
		}

		if !reflect.DeepEqual(collect(errs), expected) {
			t.Error("unexpected error message")
		}
	})

	t.Run("returns an error when the value is an empty string", func(t *testing.T) {
		expected := []string{"expected email email address to not be empty, got "}

		_, errs := statuscake.IsEmailAddress("", "email")
		if errs == nil {
			t.Error("expected errors but no errors were returned")
		}

		if !reflect.DeepEqual(collect(errs), expected) {
			t.Error("unexpected error message")
		}
	})

	t.Run("returns an error when the value is not a valid email address", func(t *testing.T) {
		expected := "expected email to be a valid email address, got enterprise"

		_, errs := statuscake.IsEmailAddress("enterprise", "email")
		if errs == nil {
			t.Error("expected errors but no errors were returned")
		}

		if !reflect.DeepEqual(collect(errs), expected) {
			t.Error("unexpected error message")
		}
	})
}

func TestStringIsNumerical(t *testing.T) {
	t.Run("returns no errors when the given value is a string representation of a numerical value", func(t *testing.T) {
		_, errs := statuscake.StringIsNumerical("123", "number")
		if errs != nil {
			t.Error("expected no errors but errors were returned")
		}
	})

	t.Run("returns an error when the value is not of type string", func(t *testing.T) {
		expected := "expected type of number to be string"

		_, errs := statuscake.StringIsNumerical(1701, "number")
		if errs == nil {
			t.Error("expected errors but no errors were returned")
		}

		if !reflect.DeepEqual(collect(errs), expected) {
			t.Error("unexpected error message")
		}
	})

	t.Run("returns an error when the value is an empty string", func(t *testing.T) {
		expected := "expected number number to not be empty, got "

		_, errs := statuscake.StringIsNumerical("", "number")
		if errs == nil {
			t.Error("expected errors but no errors were returned")
		}

		if !reflect.DeepEqual(collect(errs), expected) {
			t.Error("unexpected error message")
		}
	})

	t.Run("returns an error when the value is not a strign representation of a numerical value", func(t *testing.T) {
		expected := "expected number to be a valid number, got enterprise"

		_, errs := statuscake.StringIsNumerical("enterprise", "email")
		if errs == nil {
			t.Error("expected errors but no errors were returned")
		}

		if !reflect.DeepEqual(collect(errs), expected) {
			t.Error("unexpected error message")
		}
	})
}

func TestInt32InSlice(t *testing.T) {
	t.Run("returns no errors when the given value is contained within the validation slice", func(t *testing.T) {
		_, errs := statuscake.Int32InSlice([]int32{1, 2, 3})(2, "number")
		if errs != nil {
			t.Error("expected no errors but errors were returned")
		}
	})

	t.Run("returns an error when the value is not of type int", func(t *testing.T) {
		expected := "expected type of number to be integer"

		_, errs := statuscake.Int32InSlice([]int32{1, 2, 3})("1701", "number")
		if errs == nil {
			t.Error("expected errors but no errors were returned")
		}

		if !reflect.DeepEqual(collect(errs), expected) {
			t.Error("unexpected error message")
		}
	})

	t.Run("returns an error when the value is not contained within the validation slice", func(t *testing.T) {
		expected := "expected number to be one of {1, 2, 3}, got 4"

		_, errs := statuscake.Int32InSlice([]int32{1, 2, 3})(4, "number")
		if errs == nil {
			t.Error("expected errors but no errors were returned")
		}

		if !reflect.DeepEqual(collect(errs), expected) {
			t.Error("unexpected error message")
		}
	})
}

func collect(errs []error) []string {
	strs := make([]string, len(errs))
	for i, err := range errs {
		strs[i] = err.Error()
	}
	return strs
}
