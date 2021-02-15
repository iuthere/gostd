package main

import (
	"fmt"
	"strings"

	"golang.org/x/tools/go/packages"
)

func main() {
	// go list -e std
	//
	pkgs, err := packages.Load(nil, "std")
	if err != nil {
		panic(err)
	}
	// builtin is not included, force it to printPadded
	printPadded([]string{"builtin"})

	var before string
	for _, pkg := range pkgs {
		split := strings.Split(pkg.ID, "/")
		if among(split, "internal") || split[0] == "vendor" {
			continue
		}
		if len(split) > 1 {
			for s := 1; s < len(split); s++ {
				if !commonRoot(before, strings.Join(split[0:s], "/")) {
					printPadded(split[:s])
					before = strings.Join(split[0:s], "/")
				}
			}
		}
		printPadded(split)
		before = pkg.ID
	}
}

func printPadded(s []string) {
	fmt.Printf("%-20v %v\n", strings.Repeat("  ", len(s)-1)+s[len(s)-1], "https://golang.org/pkg/"+strings.Join(s, "/")+"/")
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
