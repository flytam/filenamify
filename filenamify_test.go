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
			t.Error(index, item.input.str, item.input.options, item.output)
		} else {
			t.Log(index, "pass")
		}
	}

}

func TestFilenamifyPath(t *testing.T) {
	expect := "foo!bar"
	inputStr, _ := filepath.Abs("foo:bar")

	if output, _ := Path(inputStr, Options{}); filepath.Base(output) != expect {
		t.Error("TestFilenamifyPath error", filepath.Base(output), expect)
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
