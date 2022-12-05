package filenamify

import (
	"path/filepath"
	"testing"
)

type inputItem struct {
	str     string
	options Options
}

type exampleItem struct {
	input inputItem
	// output
	output string
}

func newExampleItem(inputStr string, options Options, outputStr string) exampleItem {
	return exampleItem{
		input: inputItem{
			inputStr, options,
		}, output: outputStr,
	}
}

func TestFilenamify(t *testing.T) {
	var output string

	example := []exampleItem{
		newExampleItem("foo/bar", Options{}, "foo!bar"),
		newExampleItem("foo//bar", Options{}, "foo!bar"),
		newExampleItem("//foo//bar//", Options{}, "foo!bar"),
		newExampleItem("foo\\\\\\bar", Options{}, "foo!bar"),
		//---
		newExampleItem("foo/bar", Options{
			Replacement: "🐴🐴",
		}, "foo🐴🐴bar"),
		newExampleItem("////foo////bar////", Options{
			Replacement: "((",
		}, "foo((bar"),
		//--
		newExampleItem("foo\u0000bar", Options{}, "foo!bar"),
		newExampleItem(".", Options{}, "!"),
		newExampleItem("..", Options{}, "!"),
		newExampleItem("./", Options{}, "!"),
		newExampleItem("../", Options{}, "!"),
		newExampleItem("con", Options{}, "con!"),
		newExampleItem("foo/bar/nul", Options{}, "foo!bar!nul"),

		newExampleItem("con", Options{
			Replacement: "🐴🐴",
		}, "con🐴🐴"),
		newExampleItem("c/n", Options{
			Replacement: "o",
		}, "cono"),
		newExampleItem("c/n", Options{
			Replacement: "con",
		}, "cconn"),
	}

	for index, item := range example {
		if output, _ = Filenamify(item.input.str, item.input.options); output != item.output {
			t.Errorf("i:%d input: %v, opt:%v, expect:%v, got:%v", index, item.input.str, item.input.options, item.output, output)
		} else {
			t.Log(index, "pass")
		}
	}

}

func TestFilenamifyV2(t *testing.T) {
	var output string

	input := "c/n"
	expect := "cconn"

	if output, _ = FilenamifyV2(input, func(options *Options) {
		options.Replacement = "con"
	}); output != expect {
		t.Error("expect:", expect, "got:", output)
	} else {
		t.Log("pass")
	}

	// test default replacement with no function
	expect = "c!n"
	if output, _ = FilenamifyV2(input); output != expect {
		t.Error("expect:", expect, "got:", output)
	} else {
		t.Log("pass")
	}

	// test blank replacement
	expect = "foobar"
	for idx, inp := range []string{"foo<bar", "foo>bar", "foo:bar", "foo\"bar", "foo/bar", "foo\\bar",
		"foo|bar", "foo?bar", "foo*bar"} {
		output, _ = FilenamifyV2(inp, func(options *Options) { options.Replacement = "" })
		if output != expect {
			t.Errorf("%v: expect:'%v' got:'%v'", idx, expect, output)
		} else {
			t.Logf("%v: pass", idx)
		}
	}
}

func TestFilenamifyPath(t *testing.T) {
	expect := "foo!bar"
	expect2 := "foohbar"
	inputStr, _ := filepath.Abs("foo:bar")

	if output, _ := Path(inputStr, Options{}); filepath.Base(output) != expect {
		t.Error("TestFilenamifyPath error", filepath.Base(output), expect)
	}

	if output, _ := PathV2(inputStr); filepath.Base(output) != expect {
		t.Error("TestFilenamifyPath error", filepath.Base(output), expect)
	}

	if output, _ := PathV2(inputStr, func(options *Options) {
		options.Replacement = "h"
	}); filepath.Base(output) != expect2 {
		t.Error("TestFilenamifyPath error", filepath.Base(output), expect2)
	}
}

func TestFilenamifyLength(t *testing.T) {
	// Basename length: 152
	const filename = "this/is/a/very_very_very_very_very_very_very_very_very_very_very_very_very_very_very_very_very_very_very_very_very_very_very_very_very_very_very_long_filename.txt"

	if output, _ := Filenamify(filepath.Base(filename), Options{}); output != "very_very_very_very_very_very_very_very_very_very_very_very_very_very_very_very_very_very_very_very_" {
		t.Error("TestFilenamifyLength error")
	}

	if output, _ := Filenamify(filepath.Base(filename), Options{MaxLength: 180}); output != "very_very_very_very_very_very_very_very_very_very_very_very_very_very_very_very_very_very_very_very_very_very_very_very_very_very_very_long_filename.txt" {
		t.Error("TestFilenamifyLength error")
	}

}
