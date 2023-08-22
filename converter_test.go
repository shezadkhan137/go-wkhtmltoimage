package wkhtmltoimage_test

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/shezadkhan137/go-wkhtmltoimage"
)

func TestConverterRunOnHTMLFragment(t *testing.T) {
	cfg := &wkhtmltoimage.Config{
		LoadImages:       false,
		EnableJavascript: false,
		Fmt:              "jpeg",
	}
	c, err := wkhtmltoimage.NewConverter(cfg)
	if err != nil {
		t.Fatal(err)
	}

	testcases := []struct {
		Name  string
		Input string
	}{
		{
			Name:  "MainHTMLTagMissing",
			Input: `<head/><body><p>Hello, this is me.</p><p>Please be kind to me.</p></body>`,
		},
		{
			Name:  "HeadAndBodyTagsMissing",
			Input: `<p>Hello, this is me.</p><p>Please be kind to me.</p>`,
		},
		{
			Name:  "BodyTagWithHeadTagMissing",
			Input: `<body><p>Hello, this is me.</p><p>Please be kind to me.</p></body>`,
		},
		{
			Name:  "FontFormatting",
			Input: `<p><strong>bold</strong></p><p><em>italic</em></p><p><u>underlined</u></p><p><s>strike</s></p><ol><li>aufzählung</li></ol><ul><li>bullets</li></ul><p class="ql-indent-4">einrückung</p><p>link (nicht als link eingefügt): https://orf.at/stories/3327795/</p><p>link (als link eingefügt): <a href="https://orf.at/stories/3327795/" rel="noopener noreferrer" target="_blank">https://orf.at/stories/3327795/</a></p><p><br></p><p>langer text am ende:</p><p>Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.</p>`,
		},
		{
			Name:  "Listing",
			Input: `<ol><li>sfasfdasdas</li><li>sfasfdasdas22</li></ol><p><br></p>`,
		},
		{
			Name:  "ParagraphsBreakLines",
			Input: `<p> dasdasdasdas a sdas das d as das d as da sd as d asd as d asd as d asd as d asd asd as da sd asd as asd </p><p><br></p><p><br></p><p><br></p><p><br></p><p><br></p><p><br></p><p><br></p><p><br></p><p><br></p><p><br></p><p><br></p><p><br></p><p><br></p><p><br></p><p>d</p><p>as</p><p>d</p><p>asd</p><p>a</p><p>d</p><p><br></p><p>ads</p><p><br></p><p>asd</p><p>a</p><p>sd</p><p>asd</p><p><br></p><p>asd</p><p>as</p><p>das</p><p>d</p><p>asd</p><p>as</p><p>d</p><p>asd</p><p>asd</p><p>as</p><p>sda</p><p><br></p>`,
		},
	}

	for _, tc := range testcases {
		byteBuf := new(bytes.Buffer)
		out := bufio.NewWriter(byteBuf)

		if err := c.RunOnHTMLFragment(tc.Input, out); err != nil {
			t.Fatal(err, tc.Name)
		}

		if out.Size() <= 0 {
			t.Fail()
		}

		f, err := os.Create(tc.Name + ".jpeg")
		if err != nil {
			t.Fatal(err)
		}

		if _, err := f.Write(byteBuf.Bytes()); err != nil {
			t.Fatal(err)
		}
	}
}

func TestConverterRunOnHTMLFragmentPlainText(t *testing.T) {
	cfg := &wkhtmltoimage.Config{
		LoadImages:       false,
		EnableJavascript: false,
		Fmt:              "jpeg",
	}
	c, err := wkhtmltoimage.NewConverter(cfg)
	if err != nil {
		t.Fatal(err)
	}

	byteBuf := new(bytes.Buffer)
	out := bufio.NewWriter(byteBuf)

	err = c.RunOnHTMLFragment("I'm just a plain text", out)
	if !errors.Is(err, wkhtmltoimage.ErrNotHTML) {
		t.Fatal("expected ErrNotHTML, got:", err)
	}

	// input is empty
	err = c.RunOnHTMLFragment("", out)
	if !errors.Is(err, wkhtmltoimage.ErrNotHTML) {
		t.Fatal("expected ErrNotHTML, got:", err)
	}

}

func TestConverterRunInvalidHTML(t *testing.T) {
	cfg := &wkhtmltoimage.Config{
		Background:       false,
		EnableJavascript: false,
		Transparent:      true,
	}
	c, err := wkhtmltoimage.NewConverter(cfg)
	if err != nil {
		t.Fatal(err)
	}

	byteBuf := new(bytes.Buffer)
	out := bufio.NewWriter(byteBuf)

	err = c.Run("no html tags here", out)
	if err == nil {
		t.Fatal("should produce an error because it's not HTML")
	}
}

func TestConverterInitAndRun(t *testing.T) {
	cfg := &wkhtmltoimage.Config{
		Background:       false,
		EnableJavascript: false,
		Transparent:      true,
	}
	c, err := wkhtmltoimage.NewConverter(cfg)
	if err != nil {
		t.Fatal(err)
	}

	byteBuf := new(bytes.Buffer)
	out := bufio.NewWriter(byteBuf)

	c.Run("test", out)
}

func TestConverterEmptyConfig(t *testing.T) {
	_, err := wkhtmltoimage.NewConverter(nil)
	if err == nil {
		t.Fatal(err)
	}
}
