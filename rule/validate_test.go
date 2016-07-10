package rule

import (
	"math"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	r := Rule("(1, 10) netaddr /ab[cd]/ > 10")
	vlds, err := r.Parse()
	if err != nil {
		t.Fatal(err)
	}
	if len(vlds) != 4 {
		t.Fatal(vlds)
	}
	v := vlds[0]
	if err = v.Validate("5"); err != nil {
		t.Fatal(err)
	}
	if err = v.Validate("1"); err == nil {
		t.Fatal(v)
	} else {
		t.Log(err)
	}

	v = vlds[1]
	if err = v.Validate(":8804"); err != nil {
		t.Fatal(err)
	}

	if err = v.Validate("localhost:8804"); err != nil {
		t.Fatal(err)
	}

	v = vlds[2]
	if err = v.Validate("abc"); err != nil {
		t.Fatal(err)
	}
	if err = v.Validate("abd"); err != nil {
		t.Fatal(err)
	}
	if err = v.Validate("abe"); err == nil {
		t.Fatal()
	} else {
		t.Log(err)
	}

	v = vlds[3]
	if err = v.Validate("11"); err != nil {
		t.Fatal(err)
	}
	if err = v.Validate("10"); err == nil {
		t.Fatal()
	} else {
		t.Log(err)
	}

}

func TestParseInt(t *testing.T) {
	if v, _ := parseInt("1"); v != 1 {
		t.Fail()
	}
	if v, _ := parseInt("-1"); v != -1 {
		t.Fail()
	}

	if v, _ := parseInt("1s"); v != int64(time.Second) {
		t.Fail()
	}
	if v, _ := parseInt("-1s"); v != -int64(time.Second) {
		t.Fail()
	}
	if v, _ := parseInt("1m1s"); v != int64(time.Second)*61 {
		t.Fail()
	}

	if v, _ := parseInt("0x01"); v != 1 {
		t.Fail()
	}
	if v, _ := parseInt("-0x01"); v != -1 {
		t.Fail()
	}
	if v, _ := parseInt("0X01"); v != 1 {
		t.Fail()
	}
	if v, _ := parseInt("0x0A"); v != 10 {
		t.Fail()
	}

	if _, err := parseInt("abc"); err == nil {
		t.Fail()
	} else {
		t.Log(err)
	}
}

func TestParseRange(t *testing.T) {
	val := "(1, 10)"
	v, pos, err := parseRange(val, 0)
	if err != nil {
		t.Fatal(err)
	}
	if pos != len(val)-1 {
		t.Log(v)
		t.Fatal(pos)
	}
	if v.min != 1 {
		t.Fatal(v)
	}
	if v.max != 10 {
		t.Fatal(v)
	}
	if v.left != false {
		t.Fatal(v)
	}
	if v.right != false {
		t.Fatal(v)
	}

	val = "[1, 10)"
	v, pos, err = parseRange(val, 0)
	if err != nil {
		t.Fatal(err)
	}
	if pos != len(val)-1 {
		t.Fatal(pos)
	}
	if v.min != 1 {
		t.Fatal(v)
	}
	if v.max != 10 {
		t.Fatal(v)
	}
	if v.left != true {
		t.Fatal(v)
	}
	if v.right != false {
		t.Fatal(v)
	}

	val = "[1, 10]"
	v, pos, err = parseRange(val, 0)
	if err != nil {
		t.Fatal(err)
	}
	if pos != len(val)-1 {
		t.Fatal(pos)
	}
	if v.min != 1 {
		t.Fatal(v)
	}
	if v.max != 10 {
		t.Fatal(v)
	}
	if v.left != true {
		t.Fatal(v)
	}
	if v.right != true {
		t.Fatal(v)
	}

	val = "(  1  ,    10 )"
	v, pos, err = parseRange(val, 0)
	if err != nil {
		t.Fatal(err)
	}
	if pos != len(val)-1 {
		t.Log(v)
		t.Fatal(pos)
	}
	if v.min != 1 {
		t.Fatal(v)
	}
	if v.max != 10 {
		t.Fatal(v)
	}
	if v.left != false {
		t.Fatal(v)
	}
	if v.right != false {
		t.Fatal(v)
	}
}
func TestParseCompExp(t *testing.T) {
	val := ">1"
	v, pos, err := parseCompExp(val, 0)
	if err != nil {
		t.Fatal(err)
	}
	if pos != len(val) { //ended by the string length limitation
		t.Fatal(pos)
	}
	if v.min != 1 {
		t.Fatal(v)
	}
	if v.max != math.MaxInt64 {
		t.Fatal(v)
	}
	if v.left != false {
		t.Fatal(v)
	}
	if v.right != false {
		t.Fatal(v)
	}

	val = "<10"
	v, pos, err = parseCompExp(val, 0)
	if err != nil {
		t.Fatal(err)
	}
	if pos != len(val) { //ended by the string length limitation
		t.Fatal(pos)
	}
	if v.min != math.MinInt64 {
		t.Fatal(v)
	}
	if v.max != 10 {
		t.Fatal(v)
	}
	if v.left != false {
		t.Fatal(v)
	}
	if v.right != false {
		t.Fatal(v)
	}

	val = ">   1"
	v, pos, err = parseCompExp(val, 0)
	if err != nil {
		t.Fatal(err)
	}
	if pos != len(val) { //ended by the string length limitation
		t.Fatal(pos)
	}
	if v.min != 1 {
		t.Fatal(v)
	}
	if v.max != math.MaxInt64 {
		t.Fatal(v)
	}
	if v.left != false {
		t.Fatal(v)
	}
	if v.right != false {
		t.Fatal(v)
	}

}

func TestParseRegex(t *testing.T) {
	val := "/abc/"
	v, pos, err := parseRegex(val, 0)
	if err != nil {
		t.Fatal(err)
	}
	if pos != len(val)-1 {
		t.Fatal(pos)
	}
	if v.exp != "abc" {
		t.Fatal(v)
	}
	val = "   /abc/"
	v, pos, err = parseRegex(val, 0)
	if err != nil {
		t.Fatal(err)
	}
	if pos != len(val)-1 {
		t.Fatal(pos)
	}
	if v.exp != "abc" {
		t.Fatal(v)
	}
}
func TestParseNamed(t *testing.T) {
	val := "hello"
	v, pos, err := parseNamed(val, 0)
	if err != nil {
		t.Fatal(err)
	}
	if pos != len(val) {
		t.Fatal(pos)
	}
	if v.name != "hello" {
		t.Fatal(v)
	}

	val = "  hello"
	v, pos, err = parseNamed(val, 0)
	if err != nil {
		t.Fatal(err)
	}
	if pos != len(val) {
		t.Fatal(pos)
	}
	if v.name != "hello" {
		t.Fatal(v)
	}

	val = "  hello  "
	v, pos, err = parseNamed(val, 0)
	if err != nil {
		t.Fatal(err)
	}
	if pos != len(val)-2 {
		t.Fatal(pos)
	}
	if v.name != "hello" {
		t.Fatal(v)
	}

}
