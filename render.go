package pages

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Page struct {
	Title    string
	Template string
	Subject  interface{}
	Context  interface{}
}

type Renderer struct {
	leftDelimiter  string
	rightDelimiter string
	viewsDir       string
	scripts        []string
}

type layoutContent struct {
	Title   string
	Page    Page
	Content template.HTML
}

func NewRenderer() *Renderer {
	return &Renderer{viewsDir: "./views", leftDelimiter: "{{", rightDelimiter: "}}"}
}

func (r *Renderer) SetEscapeStrings(left, right string) {
	r.leftDelimiter = left
	r.rightDelimiter = right
}

func (r *Renderer) SetViewsDir(dirname string) {
	r.viewsDir = dirname
}

func tplFunctions(r *Renderer) template.FuncMap {
	return template.FuncMap{
		"add_script": func(path string) string {
			r.scripts = append(r.scripts, path)

			return ""
		},
		"render_scripts": func() template.HTML {
			content := ""

			for _, path := range r.scripts {
				content += fmt.Sprintf("<script src=\"%s\"></script>\n", path)
			}

			return template.HTML(content)
		},
	}
}

func (r *Renderer) parseTemplates(baseDir string) (*template.Template, error) {
	var templates []string

	for _, dir := range []string{"shared", "layout", baseDir} {
		files, err := getTemplateFilenames(filepath.Join(r.viewsDir, dir))
		if err != nil {
			return nil, err
		}

		templates = append(templates, files...)
	}

	return template.New("").Delims(r.leftDelimiter, r.rightDelimiter).Funcs(tplFunctions(r)).ParseFiles(templates...)
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

func (r *Renderer) Render(writer io.Writer, page Page, tplDir string) error {
	buf := bytes.NewBuffer([]byte{})

	templates, err := r.parseTemplates(tplDir)

	if err != nil {
		return err
	}

	templates.ExecuteTemplate(buf, page.Template, nil)

	content := layoutContent{Page: page, Content: template.HTML(buf.String())}

	err = templates.ExecuteTemplate(writer, "application.html", content)

	return err
}
