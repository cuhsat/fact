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

func Print(a ...any) {
	fmt.Fprintln(os.Stdout, a...)
	os.Exit(0)
}

func Error(err ...any) {
	fmt.Fprintln(os.Stderr, err...)
}

func Fatal(err ...any) {
	fmt.Fprintln(os.Stderr, err...)
	os.Exit(1)
}

func Usage(u string) {
	fmt.Fprintln(os.Stdout, "usage:", u)
	os.Exit(2)
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

	stdin, err := Stdin()

	if err != nil {
		Fatal(err)
	}

	return strings.Split(stdin, "\n")
}

func Stdin() (in string, err error) {
	fi, err := os.Stdin.Stat()

	if err != nil {
		return
	}

	if (fi.Mode() & os.ModeCharDevice) == 0 {
		b, _ := io.ReadAll(os.Stdin)

		in = strings.TrimSpace(string(b))
	}

	return
}

func StdCall(name string, arg ...string) (string, error) {
	stdout := new(strings.Builder)
	stderr := new(strings.Builder)

	if call(stdout, stderr, name, arg...) != 0 {
		return stdout.String(), errors.New(stderr.String())
	}

	return stdout.String(), nil
}

func ExitCall(name string, arg ...string) {
	os.Exit(call(os.Stdout, os.Stderr, name, arg...))
}

func call(stdout, stderr io.Writer, name string, arg ...string) int {
	bin, err := exec.LookPath(name)

	if err != nil {
		Error(err)
		return 1
	}

	cmd := exec.Command(bin, arg...)

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err = cmd.Run()

	if ee, ok := err.(*exec.ExitError); ok {
		return ee.ExitCode()
	} else if err != nil {
		return 1
	} else {
		return 0
	}
}
