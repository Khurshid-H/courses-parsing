package models

import (
	"bufio"
	"io"
	"strings"
)

var VariationColors = make(map[string]struct{})

// PopulateCourses ...
func PopulateCourses() {
	fs := data.FS(false)
	f, err := fs.Open("/data/variation_colors.tsv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	loadVariationColors(f)
}

func loadVariationColors(colorsFile io.Reader) {
	scanner := bufio.NewScanner(colorsFile)
	for scanner.Scan() {
		line := scanner.Text()
		cols := strings.Split(line, "\t")
		for _, color := range cols {
			VariationColors[color] = struct{}{}
		}
	}
}
