/* This file is part of bbcrawl, ©2020 Jörg Walter
 *  This software is licensed under the "GNU General Public License version 3" */

package cmdline

import (
	"fmt"
	"strings"
	"testing"
)

func TestBoolean(t *testing.T) {
	var tests = map[string]bool{
		"true":  true,
		"TRUE":  true,
		"false": false,
		"FALSE": false,
	}
	for k, v := range tests {
		test := new(Boolean)
		if err := test.Set(k); err != nil {
			t.Error(err)
		}
		if bool(*test) != v {
			t.Errorf("expected %t, got %t", v, bool(*test))
		}
	}

	var errors = []string{"", "abc", "trve", "1", "0"}
	for _, v := range errors {
		test := new(Boolean)
		if err := test.Set(v); err == nil {
			t.Error("Error expected")
		}
	}

	trve := new(Boolean)
	*trve = true
	if trve.String() != "true" {
		t.Errorf("expected %s, got %s", "true", trve.String())
	}
	f := new(Boolean)
	if f.String() != "false" {
		t.Errorf("expected %s, got %s", "false", f.String())
	}
}

func TestURLCollection(t *testing.T) {
	input := "https://www.google.com,ftp://example.com,relative/url/example.html"
	input_split := strings.Split(input, ",")
	urls := &URLCollection{}
	if urls.String() != "" {
		t.Logf("%s: Expected empty string.", t.Name())
		t.Fail()
	}
	if err := urls.Set(input); err != nil {
		t.Logf("%s: error in method \"Set\": %v", t.Name(), err)
		t.FailNow()
	}
	for i, url := range urls.URLs {
		if input_split[i] != url.String() {
			t.Logf("%s: expected %s==%s", t.Name(), input_split[i], url.String())
			t.FailNow()
		}
	}
	if input != urls.String() {
		t.Logf("%s: expected %s==%s", t.Name(), input, urls.String())
		t.FailNow()
	}
}

func TestSingleURL(t *testing.T) {
	input := "https://www.google.com"
	url := &SingleURL{}
	if url.String() != "" {
		t.Logf("%s: Expected empty string.", t.Name())
		t.Fail()
	}
	if err := url.Set(input); err != nil {
		t.Logf("%s: error in method \"Set\": %v", t.Name(), err)
		t.FailNow()
	}
	if input != url.String() {
		t.Logf("%s: expected %s==%s", t.Name(), input, url.String())
		t.FailNow()
	}
}

func TestIntRange(t *testing.T) {
	ir := &IntRange{}
	if ir.String() != "0,0" {
		t.Errorf("%s: Expected \"0,0\", got %q", t.Name(), ir.String())
	}
	shouldWork := func(first, second int, positive bool) {
		input := fmt.Sprintf("%d,%d", first, second)
		if positive {
			if err := ir.Set(input); err != nil {
				t.Errorf("input \"%s\": %v", input, err)
				return
			}
		} else {
			if err := ir.Set(input); err == nil {
				t.Errorf("input \"%s\" should have caused an error", input)
				return
			}
		}
		if ir.String() != input {
			t.Errorf("%s: Expected \"%s\", got %q", t.Name(), input, ir.String())
		}
	}
	//legal combinations
	shouldWork(3, 3, true)
	shouldWork(23, 42, true)
	shouldWork(-1, 2, true)
	//illegal combinations
	shouldWork(4, 3, false)
}

func TestFSDirectory(t *testing.T) {
	dir := "/var"
	nodir := "allyourbasearebelongtous"
	fsdir := &FSDirectory{}
	nofsdir := &FSDirectory{}
	if err := fsdir.Set(dir); err != nil {
		t.Logf("%s: %v.", t.Name(), err)
		t.Fail()
	}
	if err := nofsdir.Set(nodir); err == nil {
		t.Logf("%s: Expected an error.", t.Name())
		t.Fail()
	}
}

func TestStartPage(t *testing.T) {
	page := StartPage(0)
	if err := page.Set("0"); err == nil {
		t.Logf("%s: Set called with 0, no error thrown.", t.Name())
		t.Fail()
	}
	if err := page.Set("23"); err != nil {
		t.Logf("%s: %v", t.Name(), err)
		t.Fail()
	}
	if int(page) != 23 {
		t.Logf("%s: StartPage expected to be 23, got %d.", t.Name(), int(page))
		t.Fail()
	}
}

func TestEndPage(t *testing.T) {
	s := StartPage(23)
	e := NewEndPage(&s)
	if err := e.Set("21"); err == nil {
		t.Logf("%s: StartPage (%d) is greater than EndPage (%d), expected error.", t.Name(), int(s), e.End)
		t.Fail()
	}
	err := e.Set("24")
	if err != nil {
		t.Logf("%s: %v.", t.Name(), err)
		t.Fail()
	}
	if e.End != 24 {
		t.Logf("%s: EndPage is expected to be 21, but is %d.", t.Name(), e.End)
		t.Fail()
	}
	if err := e.Set("0"); err == nil {
		t.Logf("%s: EndPage is less than 1, expected error.", t.Name())
		t.Fail()
	}
}

func TestAttrs(t *testing.T) {
	const (
		key_test    = "test"
		key_style   = "style"
		key_numbers = "numbers"
		val_style   = "height:20px;width:30px"
	)
	input := fmt.Sprintf("%s=high,low/%s=%s/%s=1,2,3,4,5", key_test, key_style, val_style, key_numbers)
	a := make(Attrs)
	if err := a.Set(input); err != nil {
		t.Logf("%s: %s.", t.Name(), err)
		t.FailNow()
	}
	if len(a) != 3 {
		t.Errorf("%s: Expected %d pairs, got %d.", t.Name(), 3, len(a))
	}
	if len(a[key_test]) != 2 {
		t.Errorf("%s: Expected key %q to have %d values.", t.Name(), key_test, 2)
	}
	if a[key_style][0] != val_style {
		t.Errorf("%s: Key %q: Expected %q, got %q.", t.Name(), key_style, val_style, a[key_style][0])
	}
	if len(a[key_numbers]) != 5 {
		t.Errorf("%s: Expected key %q to have %d values.", t.Name(), key_numbers, 5)
	}
	t.Logf("%s:\n\tInput: %q\n\tString(): %q\n", t.Name(), input, a.String())
}
