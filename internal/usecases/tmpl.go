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

	tmpl := template.New(layout)
	tmpl = tmpl.Funcs(template.FuncMap{"assetUrl": assetUrl})

	return tmpl.ParseFiles(lp, fp)
}

func assetUrl(assetPath string) string {
	return os.Getenv("STATIC_ROOT_URL") + assetPath
}
