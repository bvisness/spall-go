package spall

// #include <stdlib.h>
//
// #include "spall.h"
// #include "wrap.h"
import "C"

import (
	"io"
	"runtime/cgo"
	"time"
	"unsafe"
)

type TimestampUnit float64

const (
	UnitNanoseconds  TimestampUnit = 0.001
	UnitMicroseconds TimestampUnit = 1
	UnitMilliseconds TimestampUnit = 1000
)

type Profile struct {
	w io.Writer

	selfHandle cgo.Handle
	sp         *C.SpallProfile
}

func NewProfile(w io.Writer, unit TimestampUnit) *Profile {
	p := &Profile{
		w: w,
	}

	p.selfHandle = cgo.NewHandle(p)
	p.sp = C.NewSpallProfile(C.uintptr_t(p.selfHandle), C.double(unit))

	return p
}

func (p *Profile) Now() float64 {
	switch p.sp.timestamp_unit {
	case C.double(UnitNanoseconds):
		return float64(time.Now().UnixNano())
	case C.double(UnitMilliseconds):
		return float64(time.Now().UnixMilli())
	default:
		return float64(time.Now().UnixMicro())
	}
}

func (p *Profile) Close() {
	C.SpallQuit(p.sp)
	C.FreeSpallProfile(p.sp)
	p.selfHandle.Delete()
}

type Eventer struct {
	p  *Profile
	sb *C.SpallBuffer
}

const bufferSize = 100 * 1024 * 1024

func (p *Profile) NewEventer() Eventer {
	e := Eventer{
		p:  p,
		sb: C.NewSpallBuffer(bufferSize),
	}
	C.SpallBufferInit(e.p.sp, e.sb)

	return e
}

func (e *Eventer) Begin(name string, when float64) {
	nameC := C.CString(name)
	C.SpallTraceBeginLen(e.p.sp, e.sb, nameC, C.long(len(name)), C.double(when))
	C.free(unsafe.Pointer(nameC))
}

func (e *Eventer) End(when float64) {
	C.SpallTraceEnd(e.p.sp, e.sb, C.double(when))
}

func (e *Eventer) BeginNow(name string) {
	e.Begin(name, e.p.Now())
}

func (e *Eventer) EndNow() {
	e.End(e.p.Now())
}

func (e *Eventer) Close() {
	C.SpallBufferQuit(e.p.sp, e.sb)
	C.FreeSpallBuffer(e.sb)
}

type Flusher interface {
	Flush() error
}

//export gowrite
func gowrite(pHandle C.uintptr_t, data unsafe.Pointer, dataLen uint64) bool {
	p := cgo.Handle(pHandle).Value().(*Profile)
	_, err := p.w.Write(unsafe.Slice((*byte)(data), dataLen))
	if err != nil {
		panic(err)
	}
	return true
}

//export goflush
func goflush(pHandle C.uintptr_t) bool {
	p := cgo.Handle(pHandle).Value().(*Profile)
	if flusher, ok := p.w.(Flusher); ok {
		err := flusher.Flush()
		if err != nil {
			panic(err)
		}
	}
	return true
}

//export goclose
func goclose(pHandle C.uintptr_t) {
	p := cgo.Handle(pHandle).Value().(*Profile)
	if closer, ok := p.w.(io.Closer); ok {
		err := closer.Close()
		if err != nil {
			panic(err)
		}
	}
}
