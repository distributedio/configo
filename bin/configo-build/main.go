//Use configo-build to build a conf generator
package main

import (
	"flag"
	"fmt"
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

var source = `package main
import(
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/shafreeck/configo"

	"%s" // the config package
)
func main() {
    base := ""
    flag.StringVar(&base, "patch", "", "the base conf that to patch")
    flag.Parse()

    config := %s{} //the struct

    if base != "" {
        baseData, err := ioutil.ReadFile(base)
        if err != nil {
            fmt.Println(err)
            return
        }
        if data, err := configo.Patch(baseData, &config); err != nil {
            fmt.Println(err)
            return
        } else {
			if err := ioutil.WriteFile(base, data, os.ModePerm&0644); err != nil {
				fmt.Println(err)
			}
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

	// Get current working dir
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	// Import the path
	pkg, err := build.Import(p, cwd, 0)
	if err != nil {
		log.Fatalln(err)
	}

	code := fmt.Sprintf(source, pkg.ImportPath, st)

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
