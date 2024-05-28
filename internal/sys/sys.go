// System functions.
package sys

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

const (
	EX_OK       = 0
	EX_ERROR    = 1
	EX_USAGE    = 2
	EX_DATAERR  = 3
	EX_NOINPUT  = 4
	EX_NOTEXEC  = 126
	EX_NOTFOUND = 127
)

func Print(a ...any) {
	fmt.Fprintln(os.Stdout, a...)
	os.Exit(EX_OK)
}

func Error(err ...any) {
	fmt.Fprintln(os.Stderr, err...)
}

func Fatal(err ...any) {
	fmt.Fprintln(os.Stderr, err...)
	os.Exit(EX_ERROR)
}

func Usage(u string) {
	fmt.Fprintln(os.Stdout, "usage:", u)
	os.Exit(EX_USAGE)
}

func Arg() (p string) {
	l := Args()

	if len(l) > 0 {
		p = l[0]
	}

	return
}

func Args() []string {
	if flag.NArg() > 0 {
		return flag.Args()
	}

	stdin, err, code := Stdin()

	if err != nil {
		Error(err)
		os.Exit(code)
	}

	return strings.Split(stdin, "\n")
}

func Stdin() (in string, err error, code int) {
	fi, err := os.Stdin.Stat()

	if err != nil {
		code = EX_NOINPUT
		return
	}

	if (fi.Mode() & os.ModeCharDevice) == 0 {
		var b []byte

		b, err = io.ReadAll(os.Stdin)

		if err != nil {
			code = EX_DATAERR
			return
		}

		in = strings.TrimSpace(string(b))
	}

	return
}

func StdCall(name string, args ...string) (string, error) {
	stdout := new(strings.Builder)
	stderr := new(strings.Builder)

	if call(stdout, stderr, name, args...) != 0 {
		return stdout.String(), errors.New(stderr.String())
	}

	return stdout.String(), nil
}

func ExitCall(name string, args ...string) {
	os.Exit(call(os.Stdout, os.Stderr, name, args...))
}

func call(stdout, stderr io.Writer, name string, args ...string) int {
	bin, err := exec.LookPath(name)

	if err != nil {
		Error(err)
		return EX_NOTFOUND
	}

	cmd := exec.Command(bin, args...)

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err = cmd.Run()

	if ee, ok := err.(*exec.ExitError); ok {
		return ee.ExitCode()
	} else if err != nil {
		return EX_NOTEXEC
	} else {
		return EX_OK
	}
}
