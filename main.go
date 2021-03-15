package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/tools/go/packages"
)

func main() {
	var readme bool
	flag.BoolVar(&readme, "readme", false, "output to existing README.md between ``` and ```")
	flag.Parse()

	var output io.Writer = os.Stdout

	if readme {
		output = &strings.Builder{}
	}

	gov, err := exec.Command("go", "version").CombinedOutput()
	if err != nil {
		panic("go version execution returned " + err.Error())
	}
	govs := string(gov)
	if !strings.HasPrefix(govs, "go version go") {
		panic("wrong output from go version: " + govs)
	}
	split := strings.Split(govs, " ")
	if len(split) != 4 {
		panic("wrong output has wrong amount of splits: " + govs)
	}
	fmt.Fprintf(output, split[2])
	fmt.Fprintf(output, "\n")
	// go list -e std
	//
	pkgs, err := packages.Load(nil, "std")
	if err != nil {
		panic(err)
	}
	// builtin is not included, force it to printPadded
	printPadded(output, []string{"builtin"})
	var before string
	for _, pkg := range pkgs {
		split := strings.Split(pkg.ID, "/")
		if among(split, "internal") || split[0] == "vendor" {
			continue
		}
		for s := 1; s < len(split); s++ {
			joined := strings.Join(split[:s], "/")
			if !commonRoot(before, joined) {
				printPadded(output, split[:s])
				before = joined
			}
		}
		printPadded(output, split)
		before = pkg.ID
	}
	if readme {
		f, err := os.ReadFile("README.md")
		if err != nil {
			panic(err)
		}
		from := bytes.Index(f, []byte("```"))
		to := bytes.LastIndex(f, []byte("```"))
		if to > from {
			out, err := os.Create("README.md")
			if err != nil {
				panic(err)
			}
			defer func() {
				out.Sync()
				out.Close()
			}()
			fmt.Fprintf(out, "%s", f[:from+3])
			fmt.Fprintf(out, `
$ go get github.com/iuthere/gostd
$ gostd
`)
			fmt.Fprintf(out, "%s", output)
			fmt.Fprintf(out, "%s", f[to:])
		} else {
			panic("unable to locate ```block``` in README.md")
		}
	}
}

func printPadded(w io.Writer, s []string) {
	fmt.Fprintf(w, "%-20v %v\n", strings.Repeat("  ", len(s)-1)+s[len(s)-1], "https://golang.org/pkg/"+strings.Join(s, "/")+"/")
}

func commonRoot(prev, new string) bool {
	p := strings.Split(prev, "/")
	n := strings.Split(new, "/")
	cnt := 0
	for i := 0; i < len(p) && i < len(n); i++ {
		if p[i] == n[i] {
			cnt++
		} else {
			break
		}
	}
	return cnt >= len(n)
}

func among(l []string, v string) bool {
	for _, s := range l {
		if s == v {
			return true
		}
	}
	return false
}
