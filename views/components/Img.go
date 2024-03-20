package components

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/a-h/templ"
)

// format the `src` property to include the necessary query params
// that Parcel consumes to transform the image with `sharp`
func getSrc(src, width, height, filetype string) string {
	if filetype == filepath.Ext(src) {
		filetype = ""
	} else {
		filetype = fmt.Sprintf("&as=%s", filetype)
	}
	return fmt.Sprintf("%s?quality=90&width=%s&height=%s%s", src, width, height, filetype)
}

// When using Parcel to process images, this manually builds an unescaped
// <picture> component so that Parcel can process the query params on the
// `src` to resize the image with `sharp`. USE CAREFULLY, but since we are
// generating only static HTML pages for this site, it shouldn't be dangerous
// to use unescaped values here since we are supplying all of the content
// in the build phase.
//
// Normal Templ components escape *all* attributes, which
// was causing issues with the processing of the `srcset` and `src` attributes
// that Parcel needs to consume to transform the images.
//
// relevant Parcel github issue:
// https://github.com/parcel-bundler/parcel/issues/9448#issuecomment-1891988960)
func Img(src, width, height string, attrs templ.Attributes) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		output := `<picture>`
		filetypes := [2]string{"avif", "webp"}

		// <source> tags
		for _, ft := range filetypes {
			output = output + fmt.Sprintf(
				`<source srcset="%s" type="image/%s" width="%s" height="%s" />`,
				getSrc(src, width, height, ft),
				ft,
				width,
				height,
			)
		}

		// <img> tag
		output = output + fmt.Sprintf(
			`<img src="%s" width="%s" height="%s"`,
			getSrc(src, width, height, filepath.Ext(src)),
			width,
			height,
		)
		// any string attributes we are given we should render onto the `<img />` tag
		for attr, val := range attrs {
			if val, ok := val.(string); ok {
				output = output + fmt.Sprintf(` %s="%s"`, attr, val)
			}
		}
		// set default values for necessary attributes
		for attr, defaultVal := range map[string]string{
			"alt":      "",
			"decoding": "async",
			"loading":  "lazy",
		} {
			if attrs[attr] == nil {
				output = output + fmt.Sprintf(` %s="%s"`, attr, defaultVal)
			}
		}
		output = output + " />"
		// close the <picture> element and this component
		output = output + "</picture>"

		_, err := io.WriteString(w, output)
		return err
	})
}
