package pages

import (
	"bytes"
	"html/template"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
)

var leftDelimiter = "{{"
var rightDelimiter = "}}"
var viewsDir = "./views"

type Page struct {
	Title    string
	Template string
	Subject  interface{}
	Context  interface{}
}

type layoutContent struct {
	Title   string
	Page    Page
	Content template.HTML
}

func SetEscapeStrings(left, right string) {
	leftDelimiter = left
	rightDelimiter = right
}

func SetViewsDir(dirname string) {
	viewsDir = dirname
}

func parseTemplates(baseDir string) (*template.Template, error) {
	var allFiles []string

	sharedDir := filepath.Join(viewsDir, "shared")
	layoutsDir := filepath.Join(viewsDir, "layout")
	templatesDir := filepath.Join(viewsDir, baseDir)

	if layouts, err := ioutil.ReadDir(layoutsDir); err == nil {
		for _, file := range layouts {
			filename := file.Name()
			if strings.HasSuffix(filename, ".html") {
				allFiles = append(allFiles, filepath.Join(layoutsDir, filename))
			}
		}
	} else {
		return nil, err
	}

	if shared, err := ioutil.ReadDir(sharedDir); err == nil {
		for _, file := range shared {
			filename := file.Name()
			if strings.HasSuffix(filename, ".html") {
				allFiles = append(allFiles, filepath.Join(sharedDir, filename))
			}
		}
	} else {
		return nil, err
	}

	if files, err := ioutil.ReadDir(templatesDir); err == nil {
		for _, file := range files {
			filename := file.Name()
			if strings.HasSuffix(filename, ".html") {
				allFiles = append(allFiles, filepath.Join(templatesDir, filename))
			}
		}
	} else {
		return nil, err
	}

	return template.New("").Delims(leftDelimiter, rightDelimiter).ParseFiles(allFiles...)
}

func Render(writer io.Writer, page Page, tplDir string) error {
	buf := bytes.NewBuffer([]byte{})

	templates, err := parseTemplates(tplDir)

	if err != nil {
		return err
	}

	templates.ExecuteTemplate(buf, page.Template, nil)

	content := layoutContent{Page: page, Content: template.HTML(buf.String())}

	templates.ExecuteTemplate(writer, "application.html", content)

	return nil
}
