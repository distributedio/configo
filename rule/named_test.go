package rule

import (
	"testing"
)

func TestValidatePath(t *testing.T) {
	if err := ValidatePath("/abc"); err != nil {
		t.Fatal(err)
	}
	if err := ValidatePath("./abc"); err != nil {
		t.Fatal(err)
	}
	if err := ValidatePath("../abc"); err != nil {
		t.Fatal(err)
	}
	if err := ValidatePath(".../abc"); err == nil {
		t.Fatal()
	} else {
		t.Log(err)
	}
	if err := ValidatePath("./../abc"); err != nil {
		t.Fatal(err)
	}
	if err := ValidatePath(".././../abc"); err != nil {
		t.Fatal(err)
	}
	if err := ValidatePath("./"); err != nil {
		t.Fatal(err)
	}
	if err := ValidatePath("/"); err != nil {
		t.Fatal(err)
	}
	if err := ValidatePath("."); err != nil {
		t.Fatal(err)
	}
	if err := ValidatePath("abc/"); err != nil {
		t.Fatal(err)
	}
	if err := ValidatePath("abc"); err != nil {
		t.Fatal(err)
	}
	if err := ValidatePath("abc.log"); err != nil {
		t.Fatal(err)
	}
	if err := ValidatePath("abc..log"); err == nil {
		t.Fatal()
	} else {
		t.Log(err)
	}
}
