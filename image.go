package wkhtmltoimage

/*
#cgo LDFLAGS: -lwkhtmltox
#include <stdlib.h>
#include <wkhtmltox/image.h>
*/
import "C"
import "errors"

// Init initializes the library, allocating all necessary resources.
func Init() error {
	if C.wkhtmltoimage_init(0) != 1 {
		return errors.New("could not initialize library")
	}

	return nil
}

// Version returns the version of the library.
func Version() string {
	return C.GoString(C.wkhtmltoimage_version())
}

// Destroy releases all the resources used by the library.
func Destroy() {
	C.wkhtmltoimage_deinit()
}
