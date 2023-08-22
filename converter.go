package wkhtmltoimage

/*
#cgo LDFLAGS: -lwkhtmltox
#include <stdio.h>
#include <stdlib.h>
#include <wkhtmltox/image.h>
*/
import "C"

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"unsafe"
)

type Config struct {

	//web configs
	Background       bool
	LoadImages       bool
	EnableJavascript bool
	MinimumFontSize  uint64
	UserStyleSheet   string

	// image configs
	CropLeft    uint64
	CropTop     uint64
	CropWidth   uint64
	CropHeight  uint64
	Transparent bool
	Fmt         string
	ScreenWidth uint64
	SmartWidth  uint64
	Quality     uint64

	// loading configs
	Username             string
	Password             string
	JsDelay              uint64
	ZoomFactor           float64
	BlockLocalFileAccess bool
	StopSlowScript       bool
	LoadErrorHandling    string
	Proxy                string
}

func setOption(settings *C.wkhtmltoimage_global_settings, name, value string) error {
	if name = strings.TrimSpace(name); name == "" {
		return errors.New("converter option name cannot be empty")
	}

	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))
	v := C.CString(value)
	defer C.free(unsafe.Pointer(v))

	if errCode := C.wkhtmltoimage_set_global_setting(settings, n, v); errCode != 1 {
		return fmt.Errorf("could not set converter option `%s` to `%s`: code %d", name, value, errCode)
	}

	return nil
}

func setOptions(settings *C.wkhtmltoimage_global_settings, config Config) error {
	setter := func(name, value string) error {
		return setOption(settings, name, value)
	}

	opts := []*setOp{
		// web settings
		newSetOp("web.background", config.Background, optTypeBool, setter, false),
		newSetOp("web.loadImages", config.LoadImages, optTypeBool, setter, false),
		newSetOp("web.enableJavascript", config.EnableJavascript, optTypeBool, setter, false),
		newSetOp("web.minimumFontSize", config.MinimumFontSize, optTypeUint, setter, false),
		newSetOp("web.userStyleSheet", config.UserStyleSheet, optTypeString, setter, false),
		// image settings
		newSetOp("crop.left", config.CropLeft, optTypeUint, setter, false),
		newSetOp("crop.top", config.CropTop, optTypeUint, setter, false),
		newSetOp("crop.width", config.CropWidth, optTypeUint, setter, false),
		newSetOp("crop.height", config.CropHeight, optTypeUint, setter, false),
		newSetOp("transparent", config.Transparent, optTypeBool, setter, false),
		newSetOp("screenWidth", config.ScreenWidth, optTypeUint, setter, false),
		newSetOp("smartWidth", config.SmartWidth, optTypeUint, setter, false),
		newSetOp("fmt", config.Fmt, optTypeString, setter, false),
		newSetOp("quality", config.Quality, optTypeUint, setter, false),
		// load settings
		newSetOp("load.username", config.Username, optTypeString, setter, false),
		newSetOp("load.password", config.Password, optTypeString, setter, false),
		newSetOp("load.zoomFactor", config.ZoomFactor, optTypeFloat, setter, false),
		newSetOp("load.blockLocalFileAccess", config.BlockLocalFileAccess, optTypeBool, setter, false),
		newSetOp("load.stopSlowScript", config.StopSlowScript, optTypeBool, setter, false),
		newSetOp("load.loadErrorHandling", config.LoadErrorHandling, optTypeString, setter, false),
		newSetOp("load.proxy", config.Proxy, optTypeString, setter, false),
	}

	for _, opt := range opts {
		if err := opt.execute(); err != nil {
			return err
		}
	}

	return nil
}

// applied to all converted objects.
type Converter struct {
	config *Config
}

// NewConverter returns a new converter instance.
func NewConverter(config *Config) (*Converter, error) {
	if config == nil {
		return nil, errors.New("config cannot be nil")
	}

	if err := Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize wkhtmltoimage library: %w", err)
	}

	return &Converter{
		config: config,
	}, nil
}

// Run executes the conversion and copies the output to the provided writer.
func (c *Converter) Run(input string, w io.Writer) error {
	if w == nil {
		return errors.New("the provided writer cannot be nil")
	}

	// create new settings instance
	settings := C.wkhtmltoimage_create_global_settings()
	if settings == nil {
		return errors.New("could not create converter settings")
	}

	// set config options to c struct global settings
	if err := setOptions(settings, *c.config); err != nil {
		return err
	}

	// create converter with global settings and input html
	converter := C.wkhtmltoimage_create_converter(settings, C.CString(input))
	if converter == nil {
		return errors.New("could not create converter")
	}

	// cleanup converter
	defer func() {
		C.wkhtmltoimage_destroy_converter(converter)
	}()

	// Convert objects.
	if C.wkhtmltoimage_convert(converter) != 1 {
		return errors.New("could not convert given html")
	}

	// Get conversion output buffer.
	var output *C.uchar
	size := C.wkhtmltoimage_get_output(converter, &output)
	if size == 0 {
		return errors.New("could not retrieve the converted file")
	}

	// Copy output to the provided writer.
	buf := bytes.NewBuffer(C.GoBytes(unsafe.Pointer(output), C.int(size)))
	if _, err := io.Copy(w, buf); err != nil {
		return err
	}

	return nil
}

func (c *Converter) RunOnHTMLFragment(input string, w io.Writer) error {
	if input == "" {
		return errors.New("input is empty")
	}

	hasHTMLTags, err := regexp.Match(`<\s*([^ >]+)[^>]*>.*?<\s*/\s*\1\s*>`, []byte(input))
	if err != nil {
		// we will assume that it is HTML
		hasHTMLTags = true
	} else if !hasHTMLTags {
		// if it doesn't have HTML tags, it's a plain text,
		// we need wrap it in <p> for the wkhtmltoimage to process it correctly
		input = `<p>` + input + `</p>`
	}

	if !strings.HasPrefix(input, "<html>") {
		if !strings.HasPrefix(input, "<head>") {
			if !strings.HasPrefix(input, "<body>") {
				input = `<body>` + input + `</body>`
			}
			input = `<head></head>` + input
		}
		input = `<html>` + input + `</html>`
	}

	return c.Run(input, w)
}
