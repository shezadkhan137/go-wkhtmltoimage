# go-wkhtmltoimage

Implements wkhtmltoimage Go bindings. It can be used to convert HTML documents to images. The package does not use the wkhtmltoimage binary. Instead, it uses the wkhtmltox library directly.

Full documentation can be found at: https://godoc.org/github.com/shezadkhan137/go-wkhtmltoimage

Inspiration and some code taken from [https://github.com/adrg/go-wkhtmltopdf](https://github.com/adrg/go-wkhtmltopdf)


## Installation

```
go get github.com/shezadkhan137/go-wkhtmltoimage
```

## Usage

```go
package main

import (
	"log"
	"os"

	wk "github.com/shezadkhan137/go-wkhtmltoimage"
)

func main() {
	wk.Init()
	defer wk.Destroy()

	converter, err := wk.NewConverter(
		&wk.Config{
			Quality:          100,
			Fmt:              "png",
			EnableJavascript: false,
		})
	if err != nil {
		log.Fatal(err)
	}

	testString := "<html><body><p>This is some html</p></body></html>"

	outFile, err := os.Create("testme.png")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	err = converter.Run(testString, outFile)
	if err != nil {
		log.Fatal(err)
	}
}
```

See [libwkhtmltox](https://wkhtmltopdf.org/libwkhtmltox/) for settings and documentation

## License

This project is licensed under the [MIT license](https://opensource.org/licenses/MIT).
See [LICENSE](https://github.com/shezadkhan137/go-wkhtmltoimage/blob/master/LICENSE) for more details.
