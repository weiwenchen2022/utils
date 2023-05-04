package stringx_test

import (
	"strings"
	"testing"

	. "github.com/weiwenchen2022/utils/stringx"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func TestRandString(t *testing.T) {
	t.Parallel()

	const n = 16

	s1 := RandString(n)
	if got := len(s1); got != n {
		t.Errorf("len(s1) = %d, want %d", got, n)
	}
	for _, r := range s1 {
		if !strings.ContainsRune(letters, r) {
			t.Errorf("s Contains not the uppercase or lowercase letters: %q", r)
		}
	}

	s2 := RandString(n)
	if got := len(s2); got != n {
		t.Errorf("len(s2) = %d, want %d", got, n)
	}
	for _, r := range s2 {
		if !strings.ContainsRune(letters, r) {
			t.Errorf("s Contains not the uppercase or lowercase letters: %q", r)
		}
	}

	if s1 == s2 {
		t.Error("s1 == s2")
	}
}
