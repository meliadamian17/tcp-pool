package utils

import (
	"reflect"
	"testing"
)

func AssertEqual(t *testing.T, expected, actual interface{}, message string) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%s: expected %v, got %v", message, expected, actual)
	}
}

func AssertNotEqual(t *testing.T, unexpected, actual interface{}, message string) {
	if reflect.DeepEqual(unexpected, actual) {
		t.Errorf("%s: did not expect %v", message, unexpected)
	}
}

func AssertNil(t *testing.T, actual interface{}, message string) {
	if !isNil(actual) {
		t.Errorf("%s: expected nil, got %v", message, actual)
	}
}

func AssertNotNil(t *testing.T, actual interface{}, message string) {
	if isNil(actual) {
		t.Errorf("%s: expected not nil, but got nil", message)
	}
}

func AssertTrue(t *testing.T, condition bool, message string) {
	if !condition {
		t.Errorf("%s: expected true, got false", message)
	}
}

func AssertFalse(t *testing.T, condition bool, message string) {
	if condition {
		t.Errorf("%s: expected false, got true", message)
	}
}

func isNil(value interface{}) bool {
	if value == nil {
		return true
	}
	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return val.IsNil()
	}
	return false
}
