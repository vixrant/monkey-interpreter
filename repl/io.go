package repl

import (
	"bufio"
	"io"
)

type ReplIO struct {
	in		io.Reader
	out 	io.Writer
	scanner *bufio.Scanner
}

func SetupIO(in io.Reader, out io.Writer) *ReplIO {
	rio := &ReplIO{in: in, out: out}
	rio.scanner = bufio.NewScanner(in)
	return rio
}

func (rio *ReplIO) Write(msg string) {
	_, err := io.WriteString(rio.out, msg)
	if err != nil {
		panic(err)
	}
}

func (rio *ReplIO) Read() string {
	scanned := rio.scanner.Scan()
	if !scanned {
		return ""
	}

	line := rio.scanner.Text()
	return line
}
