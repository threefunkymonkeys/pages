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
var layoutFileName = "application.html"

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

func SetLayoutName(layoutName string) {
	lauyoutFileName = layoutName + ".html"
}

func parseTemplates(baseDir string) (*template.Template, error) {
	var templates []string

	for _, dir := range []string{"shared", "layout", baseDir} {
		files, err := getTemplateFilenames(filepath.Join(viewsDir, dir))
		if err != nil {
			return nil, err
		}

		templates = append(templates, files...)
	}

	return template.New("").Delims(leftDelimiter, rightDelimiter).ParseFiles(templates...)
}

func getTemplateFilenames(dir string) ([]string, error) {
	var filenames []string

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		filename := entry.Name()
		if strings.HasSuffix(filename, ".html") {
			filenames = append(filenames, filepath.Join(dir, filename))
		}
	}

	return filenames, nil
}

func Render(writer io.Writer, page Page, tplDir string) error {
	buf := bytes.NewBuffer([]byte{})

	templates, err := parseTemplates(tplDir)

	if err != nil {
		return err
	}

	templates.ExecuteTemplate(buf, page.Template, nil)

	content := layoutContent{Page: page, Content: template.HTML(buf.String())}

	templates.ExecuteTemplate(writer, layoutFileName, content)

	return nil
}
