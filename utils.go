package redirector

import (
	"errors"
	"unsafe"
)

// struct string {
// 		uint8 *str;
// 		int len;
// }
func string2byte(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// struct []uint8 {
//     uint8 *array;
//     int len;
//     int cap;
// }
func byte2string(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// int2byte return byte array in big endian
func uint322byte(i uint32) []byte {
	x := (*[4]byte)(unsafe.Pointer(&i))
	return x[:]
}

//
func byte2uint32(b []byte) (uint32, error) {
	if len(b) != 4 {
		return 0, errors.New("invalid")
	}
	x := (**[4]uint8)(unsafe.Pointer(&b))
	inv := [4]byte{(*x)[3], (*x)[2], (*x)[1], (*x)[0]}
	t := 0x1234
	p := (*byte)(unsafe.Pointer(&t))
	if *p == 0 {
		// big endian
		return *(*uint32)(unsafe.Pointer(*x)), nil
	}
	// little endian
	return *(*uint32)(unsafe.Pointer(&inv)), nil
}
