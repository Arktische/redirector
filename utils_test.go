package redirector

import (
	"strings"
	"testing"
)

func TestString2byte(t *testing.T) {
	s := strings.Repeat("abc", 3)
	b := string2byte(s)
	s1 := byte2string(b)
	s2 := uint322byte(12)
	s3, _ := byte2uint32([]byte{192, 168, 3, 1})

	t.Log(b, s1, s2, s3)
}
