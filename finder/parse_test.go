package finder

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseAssembly(t *testing.T) {
	// https://github.com/jeremywohl/nacl/blob/515c9bf085b5154d553ecf04288936850819c788/crypto_stream/salsa20/amd64_xmm6/stream.s#L334-L357
	var asm = `# qhasm: stack64 bytes_backup

# qhasm: enter crypto_stream_salsa20_amd64_xmm6
.text
.p2align 5
.globl _crypto_stream_salsa20_amd64_xmm6
.globl crypto_stream_salsa20_amd64_xmm6
_crypto_stream_salsa20_amd64_xmm6:
crypto_stream_salsa20_amd64_xmm6:
mov %rsp,%r11
and $31,%r11
add $480,%r11
sub %r11,%rsp

# qhasm: r11_stack = r11_caller
# asm 1: movq <r11_caller=int64#9,>r11_stack=stack64#1
# asm 2: movq <r11_caller=%r11,>r11_stack=352(%rsp)
movq %r11,352(%rsp)

# qhasm: r12_stack = r12_caller
# asm 1: movq <r12_caller=int64#10,>r12_stack=stack64#2
# asm 2: movq <r12_caller=%r12,>r12_stack=360(%rsp)
movq %r12,360(%rsp)
`
	expect := []Instruction{
		{Opcode: "mov", Args: []string{"%rsp", "%r11"}, Line: 10},
		{Opcode: "and", Args: []string{"$31", "%r11"}, Line: 11},
		{Opcode: "add", Args: []string{"$480", "%r11"}, Line: 12},
		{Opcode: "sub", Args: []string{"%r11", "%rsp"}, Line: 13},
		{Opcode: "movq", Args: []string{"%r11", "352(%rsp)"}, Line: 18},
		{Opcode: "movq", Args: []string{"%r12", "360(%rsp)"}, Line: 23},
	}
	got, err := ParseAssembly(strings.NewReader(asm))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, expect) {
		t.Errorf("ParseAssembly() = %v; expect %v", got, expect)
	}
}
