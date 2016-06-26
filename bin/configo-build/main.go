//Use configo-build to build a conf generator
package main

import (
	"fmt"
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
        fmt.Printf(string(data))
    }
}
`

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Println("Usage:")
		fmt.Printf("  %s <package>.<struct>\n", args[0])
		return
	}

	pkg := args[1]
	idx := strings.LastIndex(pkg, ".")
	if idx < 0 {
		log.Fatalln("Invalid package format, expected <package>.<struct name>")
	}
	p := pkg[0:idx]
	st := path.Base(p) + "." + pkg[idx+1:]
	code := fmt.Sprintf(source, p, st)
	//fmt.Print(code)

	tmpFile := fmt.Sprintf("/tmp/%s.pkg.go", st)
	if err := ioutil.WriteFile(tmpFile, []byte(code), os.ModePerm&0660); err != nil {
		log.Fatalln(err)
	}

	out := strings.ToLower(st) + ".cfg"
	cmd := exec.Command("go", "build", "-o", out, tmpFile)
	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}
}
