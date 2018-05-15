package go_buff_alloc

import (
	"testing"
)

func TestNewBuffAllocator(t *testing.T) {
	shouldPanic(t, func() {
		NewBuffAllocator(0, 2)
	})
	shouldPanic(t, func() {
		NewBuffAllocator(128, 0)
	})
	shouldPanic(t, func() {
		NewBuffAllocator(128, 1)
	})
	NewBuffAllocator(128, 2)
}

func shouldPanic(t *testing.T, f func()) {
	defer func() {
		err := recover()
		if err != nil {
			// good, panic as expected
		} else {
			t.Errorf("should panic")
		}
	}()
	f()
}

func TestBuffAllocator_Allocate(t *testing.T) {
	ba := NewBuffAllocator(128, 2)

	buf := ba.Allocate(100)
	if cap(buf) != 128 {
		t.Errorf("buf cap is %d, should be %d", cap(buf), 128)
	}

	if len(buf) != 0 {
		t.Errorf("buff len should be 0")
	}

	buf = ba.Allocate(256)
	if cap(buf) != 256 {
		t.Errorf("buf cap is %d, should be %d", cap(buf), 256)
	}

	buf = ba.Allocate(150)
	if cap(buf) != 256 {
		t.Errorf("buf cap is %d, should be %d", cap(buf), 256)
	}
}

func TestBuffAllocator_Release(t *testing.T) {
	ba := NewBuffAllocator(128, 2)
	goodBuff := ba.Allocate(1024)
	ba.Release(goodBuff)
	shouldPanic(t, func() {
		badBuff := make([]byte, 1000)
		ba.Release(badBuff)
	})

	buff := ba.Allocate(256)
	buff = buff[:12]
	ba.Release(buff)
	buff = ba.Allocate(256)
	if len(buff) != 0 {
		t.Errorf("buff len should be 0")
	}
}

func BenchmarkBuffAllocator_Allocate_1024(b *testing.B) {
	ba := NewBuffAllocator(128, 2)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buff := ba.Allocate(1024)
		ba.Release(buff)
	}
}

func BenchmarkWithoutBuffAllocator_Allocate_1024(b *testing.B) {

	for i := 0; i < b.N; i++ {
		_ = make([]byte, 1024)
	}
}

func BenchmarkBuffAllocator_Allocate_10240(b *testing.B) {
	ba := NewBuffAllocator(128, 2)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buff := ba.Allocate(10240)
		ba.Release(buff)
	}
}

func BenchmarkWithoutBuffAllocator_Allocate_10240(b *testing.B) {

	for i := 0; i < b.N; i++ {
		_ = make([]byte, 10240)
	}
}

func BenchmarkBuffAllocator_Allocate_102400(b *testing.B) {
	ba := NewBuffAllocator(128, 2)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buff := ba.Allocate(102400)
		ba.Release(buff)
	}
}

func BenchmarkWithoutBuffAllocator_Allocate_102400(b *testing.B) {

	for i := 0; i < b.N; i++ {
		_ = make([]byte, 102400)
	}
}
