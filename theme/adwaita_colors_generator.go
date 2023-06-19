//go:build ignore
// +build ignore

/*
This tool will visit the Adwaita color page and generate a Go file with all the colors.
To add a new color, just add it to the colorToGet map. The key is the name of the color for Fyne, and the color name
is the name of the color in the Adwaita page without the "@".
*/

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"image/color"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"text/template"
)

const (
	adwaitaColorPage = "https://gnome.pages.gitlab.gnome.org/libadwaita/doc/1.0/named-colors.html"
	output           = "adwaita_colors.go"
	sourceTpl        = `package theme

// This file is generated by adwaita_colors_generator.go
// Please do not edit manually, use:
// go generate ./theme/
//
// The colors are taken from: https://gnome.pages.gitlab.gnome.org/libadwaita/doc/1.0/named-colors.html

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var adwaitaDarkScheme = map[fyne.ThemeColorName]color.Color{
{{- range $key, $value := .DarkScheme }}
    {{$key}}: {{$value}},
{{- end }}
}

var adwaitaLightScheme = map[fyne.ThemeColorName]color.Color{
{{- range $key, $value := .LightScheme }}
    {{$key}}: {{$value}},
{{- end }}
}`
)

var (
	tableRowMatcher       = regexp.MustCompile(`(?s)<tr>(.*?)</tr>`)
	tableColorCellMatcher = regexp.MustCompile(`(?s)<tt>((?:rgba|#).*?)</tt>`)
	rows                  = [][]string{}
	colorToGet            = map[string]string{
		"theme.ColorNameBackground":        "window_bg_color", // or "view_bg_color"
		"theme.ColorNameForeground":        "window_fg_color", // or "view_fg_color"
		"theme.ColorNameOverlayBackground": "popover_bg_color",
		"theme.ColorNamePrimary":           "accent_bg_color",
		"theme.ColorNameInputBackground":   "view_bg_color", // or "window_bg_color"
		"theme.ColorNameError":             "destructive_color",
		"theme.ColorNameButton":            "headerbar_bg_color", // it's the closer color to the button color
		"theme.ColorNameShadow":            "shade_color",
		// some colors to ensure that the theme is almost complete
		"theme.ColorGreen":  "success_color",
		"theme.ColorYellow": "warning_color",
		"theme.ColorRed":    "destructive_color",
		"theme.ColorBlue":   "accent_color",
	}
)

func main() {

	darkScheme := map[string]string{}
	lightScheme := map[string]string{}

	reps, err := http.Get(adwaitaColorPage)
	if err != nil {
		log.Fatal(err)
	}
	defer reps.Body.Close()
	htpage, err := ioutil.ReadAll(reps.Body)
	if err != nil {
		log.Fatal(err)
	}
	rows = tableRowMatcher.FindAllStringSubmatch(string(htpage), -1)

	for colname, color := range colorToGet {
		lcol, err := getColorFor(color, "light")
		if err != nil {
			log.Fatal(err)
		}
		dcol, err := getColorFor(color, "dark")
		if err != nil {
			log.Fatal(err)
		}
		lightScheme[colname] = fmt.Sprintf("color.RGBA{0x%x, 0x%x, 0x%x, 0x%x}", lcol.R, lcol.G, lcol.B, lcol.A)
		darkScheme[colname] = fmt.Sprintf("color.RGBA{0x%x, 0x%x, 0x%x, 0x%x}", dcol.R, dcol.G, dcol.B, dcol.A)
	}

	out, err := os.Create(output)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	tpl := template.New("source")
	tpl, err = tpl.Parse(sourceTpl)
	if err != nil {
		log.Fatal(err)
	}
	// generate the source
	buffer := bytes.NewBufferString("")
	err = tpl.Execute(buffer, struct {
		LightScheme map[string]string
		DarkScheme  map[string]string
	}{
		LightScheme: lightScheme,
		DarkScheme:  darkScheme,
	})
	if err != nil {
		log.Fatal(err)
	}

	// format the file
	content := buffer.String()
	if formatted, err := format.Source([]byte(content)); err != nil {
		log.Fatal(err)
	} else {
		out.Write(formatted)
	}

}

func getColorFor(name string, variant string) (col color.RGBA, err error) {
	// get the color from the adwaita_colors.go
	// return the color
	for _, row := range rows {
		// check if the row is for "@success_color" (@ is html encoded)
		if strings.Contains(row[0], "&#64;"+name) || strings.Contains(row[0], "@"+name) {
			// the color is in the second column
			c := tableColorCellMatcher.FindAllStringSubmatch(row[0], -1)
			switch variant {
			case "light":
				col, err = stringToColor(c[0][1])
			case "dark":
				col, err = stringToColor(c[1][1])
			}
			return
		}
	}
	return
}

func stringToColor(s string) (c color.RGBA, err error) {
	c.A = 0xff
	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 9:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x%02x", &c.R, &c.G, &c.B, &c.A)
	default:
		// rgba(...) format
		var r, g, b uint8
		var a float32
		_, err = fmt.Sscanf(s, "rgba(%d, %d, %d, %f)", &r, &g, &b, &a)
		c.R = r
		c.G = g
		c.B = b
		c.A = uint8(a * 255)
	}
	return
}
