package main

import (
	"fmt"
	"unsafe"
)

// Assume createInt will not be inlined.
func createInt() *int {
	return new(int)
}

func foo() {
	p0, y, z := createInt(), createInt(), createInt()
	var p1 = unsafe.Pointer(y)
	var p2 = uintptr(unsafe.Pointer(z))

	// At the time, even if the address of the int
	// value referenced by z is still stored in p2,
	// the int value has already become unused, so
	// garbage collector can collect the memory
	// allocated for it now. On the other hand, the
	// int values referenced by p0 and p1 are still
	// in using.

	// uintptr can participate arithmetic operations.
	p2 += 2
	p2--
	p2--

	*p0 = 1                         // okay
	*(*int)(p1) = 2                 // okay
	*(*int)(unsafe.Pointer(p2)) = 3 // dangerous!
}

type T struct {
	x bool
	y [3]int16
}

const N = unsafe.Offsetof(T{}.y)
const M = unsafe.Sizeof([3]int16{}[0])

func foo2() {
	t := T{y: [3]int16{123, 456, 789}}
	p := unsafe.Pointer(&t)
	// ty2 := (*int16)(unsafe.Pointer(uintptr(p)+N+M+M))
	addr := uintptr(p) + N + M + M
	// Now the t value becomes unused, its memory may be
	// garbage collected at this time. So the following
	// use of the address of t.y[2] may become invalid
	// and dangerous!
	// Another potential danger is, if some operations
	// make the stack grow or shrink here, then the
	// address of t might change, so that the address
	// saved in addr will become invalid (fact 3),
	// though this danger doesn't exist for this
	// specified example.
	ty2 := (*int16)(unsafe.Pointer(addr))
	fmt.Println(*ty2)
}
