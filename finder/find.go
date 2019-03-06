package finder

import (
	"reflect"
)

// Result contains information about a discovered instance of the bug.
type Result struct {
	Instructions []Instruction
	LowRegister  string
	HighRegister string
}

func (r Result) StartLine() int {
	return r.Instructions[0].Line
}

// Find looks for the pattern of instructions present in the salsa20 bug.
func Find(is []Instruction) *Result {
	for ; len(is) > 0; is = is[1:] {
		if r := match(is); r != nil {
			return r
		}
	}
	return nil
}

// The offending sequence has the form:
//	add  $1,%rdx
//	shl  $32,%rcx
//	add  %rcx,%rdx
//	mov  %rdx,%rcx
//	shr  $32,%rcx
//	movl %edx,4+288(%rsp)
//	movl %ecx,4+304(%rsp)

const seqlength = 7

// match looks for the salsa20 pattern starting at the beginning of the given instruction stream.
func match(is []Instruction) *Result {
	if len(is) < 2*seqlength {
		return nil
	}

	// Look for the first two instructions, from which we can identify the low and
	// high registers.
	add := is[0].Opcode == "add" && is[0].Arg(0) == "$1"
	shl := is[1].Opcode == "shl" && is[1].Arg(0) == "$32"
	if !add || !shl {
		return nil
	}
	lo := is[0].Arg(1)
	hi := is[1].Arg(1)

	// Now expect this form to be repeated twice.
	matches := hasform(is, lo, hi) && hasform(is[seqlength:], lo, hi)
	if !matches {
		return nil
	}

	return &Result{
		Instructions: is,
		LowRegister:  lo,
		HighRegister: hi,
	}
}

// hasform checks for the pattern with the given lo and hi registers.
func hasform(is []Instruction, lo, hi string) bool {
	expect := []Instruction{
		{Opcode: "add", Args: []string{"$1", lo}},  //	add  $1,%rdx
		{Opcode: "shl", Args: []string{"$32", hi}}, //	shl  $32,%rcx
		{Opcode: "add", Args: []string{hi, lo}},    //	add  %rcx,%rdx
		{Opcode: "mov", Args: []string{lo, hi}},    //	mov  %rdx,%rcx
		{Opcode: "shr", Args: []string{"$32", hi}}, //	shr  $32,%rcx
		{Opcode: "movl", Args: nil},                //	movl %edx,4+288(%rsp)
		{Opcode: "movl", Args: nil},                //	movl %ecx,4+304(%rsp)
	}
	if len(is) < len(expect) {
		return false
	}
	for i, required := range expect {
		if is[i].Error != nil {
			return false
		}
		if is[i].Opcode != required.Opcode {
			return false
		}
		if required.Args != nil && !reflect.DeepEqual(is[i].Args, required.Args) {
			return false
		}
	}
	return true
}
