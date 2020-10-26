package main

import (
	"io"
	"os"
	"os/exec"
	"strings"
)

// ExpectStream is a stream that catch outputted data, to write on input pipe.
type ExpectStream struct {
	output  *os.File
	input   *os.File
	process *exec.Cmd
	expects []*Expect
	execCmd func(index int, stdout *os.File)
}

func (w *ExpectStream) Write(p []byte) (int, error) {
	text := string(p)
	for i, e := range w.expects {
		if e.Runs != 0 && strings.Contains(text, e.Output) {
			if e.Runs > 0 {
				e.Runs--
			}

			if len(e.Send) != 0 {
				w.input.Write([]byte(e.Send))
			}
			if e.Cmd != nil {
				w.execCmd(i, w.input)
			}
		}
	}
	return w.output.Write(p)
}

func expectPipes(process *exec.Cmd, expects []*Expect, execCmd func(index int, stdout *os.File)) (*ExpectStream, *ExpectStream, *os.File, error) {
	// Pipe for stdin
	stdin, inputWriter, err := os.Pipe()
	if err != nil {
		return nil, nil, nil, err
	}

	// Create stdout + stderr that catch expected strings
	stdout := &ExpectStream{os.Stdout, inputWriter, process, expects, execCmd}
	stderr := &ExpectStream{os.Stderr, inputWriter, process, expects, execCmd}

	// Copy real stdin to fake stdin
	go func() {
		defer inputWriter.Close()
		if _, err := io.Copy(inputWriter, os.Stdin); err != nil {
			panic(err)
		}
	}()

	return stdout, stderr, stdin, nil
}
