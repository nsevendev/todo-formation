package testsetup

import "testing"

func Except(t *testing.T, value string, expect string) {
	t.Errorf("got label %s, expect label %s", value, expect)
}

func Error(t *testing.T, err string) {
	t.Errorf("got error : %s", err)
}

func ErrorSuccess(t *testing.T) {
	t.Errorf("expected error, but got none")
}
