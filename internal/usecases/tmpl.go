package usecases

import (
	"errors"
	"html/template"
	"os"
	"path/filepath"
)

// Tmpl compiles layout with given partial
func Tmpl(layout string, partial string) (*template.Template, error) {
	lp := filepath.Join("templates", layout)
	fp := filepath.Join("templates", partial)

	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
	}

	if info.IsDir() {
		return nil, errors.New("is directory")
	}

	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}
