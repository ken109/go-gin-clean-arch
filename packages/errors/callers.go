package errors

import (
	"encoding/json"
	"fmt"
	"runtime"
)

type Frame uintptr

func (f Frame) pc() uintptr { return uintptr(f) - 1 }

func (f Frame) file() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown"
	}
	file, _ := fn.FileLine(f.pc())
	return file
}

func (f Frame) line() int {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return 0
	}
	_, line := fn.FileLine(f.pc())
	return line
}

func (f Frame) name() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown"
	}
	return fn.Name()
}

func (f Frame) String() string {
	name := f.name()
	if name == "unknown" {
		return name
	}
	return fmt.Sprintf("%s %s:%d", name, f.file(), f.line())
}

func (f Frame) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name string `json:"name"`
		File string `json:"file"`
		Line int    `json:"line"`
	}{
		Name: f.name(),
		File: f.file(),
		Line: f.line(),
	})
}

func callers() []Frame {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])

	var stack []Frame
	for _, pc := range pcs[0:n] {
		stack = append(stack, Frame(pc))
	}
	return stack
}
