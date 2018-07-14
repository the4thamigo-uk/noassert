package main

import (
	"fmt"
	"go/build"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
	"os"
	"strings"
)

func loadProgram(path string) (*loader.Program, *ssa.Program, error) {
	cfg := loader.Config{
		Build:       &build.Default,
		AllowErrors: true,
	}
	cfg.ImportWithTests(path)
	prog, err := cfg.Load()
	if err != nil {
		return nil, nil, err
	}
	ssaProg := ssautil.CreateProgram(prog, 0)
	ssaProg.Build()
	return prog, ssaProg, nil
}

func findAssertionFreeTests(prog *ssa.Program) ([]*ssa.Function, error) {
	cfg := pointer.Config{
		Mains:          ssautil.MainPackages(prog.AllPackages()),
		BuildCallGraph: true,
	}
	res, err := pointer.Analyze(&cfg)
	if err != nil {
		return nil, err
	}
	var tests []*callgraph.Node
	for fn, node := range res.CallGraph.Nodes {
		if strings.HasPrefix(fn.Name(), "Test") && fn.Signature.String() == "func(t *testing.T)" {
			tests = append(tests, node)
		}
	}
	var found []*ssa.Function
	for _, test := range tests {
		path := callgraph.PathSearch(test, func(node *callgraph.Node) bool {
			return node.Func.String() == "(*testing.common).Fail"
		})
		if len(path) == 0 {
			found = append(found, test.Func)
		}
	}
	return found, nil
}

func main() {
	path := os.Args[1]
	prog, ssaProg, err := loadProgram(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	pkgInfo := prog.Package(path)
	if !pkgInfo.TransitivelyErrorFree {
		for _, err := range pkgInfo.Errors {
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}
	ssaPkg := ssaProg.Package(pkgInfo.Pkg)
	ssaProg.CreateTestMainPackage(ssaPkg)
	tests, err := findAssertionFreeTests(ssaProg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(tests) > 0 {
		for _, test := range tests {
			pos := ssaProg.Fset.PositionFor(test.Pos(), false)
			fmt.Printf("%s: Assertion free test '%s' found \n", pos, test.Name())
		}
		os.Exit(1)
	}
}
