package testsetup

import "testing"

func expectFail(t *testing.T, value any, expected any, err any) {
	if err == "" {
		err = "not found"
	}

	t.Errorf("🟡 Value : %v | Expect : %v | Error : %v", value, expected, err)
}

func expectError(t *testing.T, err any) {
	if err == "" {
		err = "not found"
	}

	t.Fatalf("🔴 Error : %v", err)
}

func expectSuccess(t *testing.T, value any, expected any) {
	t.Logf("🟢 Value : %v | Expect : %v", value, expected)
}

func LogNameTest(t *testing.T, name string) {
	t.Logf("🔵 Test: %s", name)
}

func Equal(t *testing.T, value any, expected any) {
	if value != expected {
		expectFail(t, value, expected, "not equal")
	} else {
		expectSuccess(t, value, expected)
	}
}

func IsNull(t *testing.T, value any) {
	if value != nil {
		expectFail(t, value, nil, "not <nil>")
	} else {
		expectSuccess(t, value, nil)
	}
}

func IsNotNull(t *testing.T, value any) {
	if value == nil {
		expectFail(t, value, "not <nil>", "is <nil>")
	} else {
		expectSuccess(t, value, "not <nil>")
	}
}

func IsEmptySlice[T any](t *testing.T, value []T) {
	if len(value) != 0 {
		expectFail(t, value, value, "not empty slice")
	} else {
		expectSuccess(t, value, "empty slice")
	}
}

func IsNotEmptySlice[T any](t *testing.T, value []T) {
	if len(value) == 0 {
		expectFail(t, value, value, "empty slice")
	} else {
		expectSuccess(t, value, "not empty slice")
	}
}
