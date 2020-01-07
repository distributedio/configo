//Use configo-build to build a conf generator
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"golang.org/x/tools/go/packages"
)

var source = `package main
import(
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/distributedio/configo"

	"%s" // the config package
)
func main() {
    base := ""
	test := ""
    flag.StringVar(&base, "patch", "", "the base conf that to patch")
    flag.StringVar(&test, "test", "", "try to load the configuration file")
    flag.Parse()

    config := %s{} //the struct

    if base != "" {
        baseData, err := ioutil.ReadFile(base)
        if err != nil {
            fmt.Println(err)
            return
        }
        if data, err := configo.Patch(baseData, config); err != nil {
            fmt.Println(err)
            return
        } else {
			if err := ioutil.WriteFile(base, data, os.ModePerm&0644); err != nil {
				fmt.Println(err)
			}
        }
        return
    }
	if test != "" {
		if err := configo.Load(test, &config); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(config)
			fmt.Println("Load successfully !")
		}
		return
	}

    if data, err := configo.Marshal(config); err != nil {
        fmt.Println(err)
        return
    } else {
        fmt.Println(string(data))
    }
}
`

func main() {
	genCode := false

	flag.Usage = func() {
		fmt.Println("Usage:")
		fmt.Printf("  %s [option] <package>.<struct>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.BoolVar(&genCode, "gencode", false, "generate source code")
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		return
	}

	// Parse the package path and struct name
	target := args[0]
	idx := strings.LastIndex(target, ".")
	if idx < 0 {
		log.Fatalln("Invalid package format, expected <package>.<struct name>")
	}
	p := target[0:idx]
	st := path.Base(p) + "." + target[idx+1:]

	pkgs, err := packages.Load(nil, p)
	if err != nil {
		log.Fatalln(err)
	}
	if len(pkgs) != 1 {
		log.Fatalln("Package is not found", p)
	}
	pkg := pkgs[0]

	code := fmt.Sprintf(source, pkg, st)

	if genCode {
		fmt.Print(code)
		return
	}

	tmpFile := fmt.Sprintf("/tmp/%s.pkg.go", st)
	if err := ioutil.WriteFile(tmpFile, []byte(code), os.ModePerm&0660); err != nil {
		log.Fatalln(err)
	}

	out := strings.ToLower(st) + ".cfg"
	cmd := exec.Command("go", "build", "-o", out, tmpFile)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}
}
