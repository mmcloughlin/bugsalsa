package finder

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Instruction in an assembly file.
type Instruction struct {
	Error  error
	Opcode string
	Args   []string
	Line   int
}

func (i Instruction) Arity() int {
	return len(i.Args)
}

func (i Instruction) Arg(n int) string {
	if n >= i.Arity() {
		return ""
	}
	return i.Args[n]
}

func (i Instruction) String() string {
	s := i.Opcode
	if len(i.Args) > 0 {
		s += "\t" + strings.Join(i.Args, ",")
	}
	return s
}

// ParseAssembly parses out instructions from a qhasm-generated assembly file.
func ParseAssembly(r io.Reader) ([]Instruction, error) {
	is := []Instruction{}
	scanner := bufio.NewScanner(r)
	linenum := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		linenum++

		// Skip: empty lines, comments, metadata, labels
		n := len(line)
		if n == 0 || line[0] == '#' || line[0] == '.' || line[n-1] == ':' {
			continue
		}

		inst := Instruction{
			Line: linenum,
		}
		parts := strings.Fields(line)
		inst.Opcode = parts[0]
		if len(parts) > 2 {
			inst.Error = fmt.Errorf("expected at most two fields: %q", line)
		}
		if len(parts) > 1 {
			inst.Args = strings.Split(parts[1], ",")
		}
		is = append(is, inst)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return is, nil
}
