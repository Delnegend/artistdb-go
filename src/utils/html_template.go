package utils

import (
	"bytes"
	"html/template"
	"log/slog"
	"net/http"
)

type HTMLTemplate struct {
	Tmpl *template.Template
}

// Read file and set it as the template
func (tmpl *HTMLTemplate) Read(filePath string) error {
	var err error
	tmpl.Tmpl, err = template.ParseFiles(filePath)
	if err != nil {
		return err
	}
	return nil
}

// Execute the template with the data and write it to the response writer
func (tmpl *HTMLTemplate) Execute(w http.ResponseWriter, data interface{}) {
	if err := tmpl.Tmpl.Execute(w, data); err != nil {
		slog.Error(err.Error())
	}
}

// Render the template with the data and return it as a string
func (tmpl *HTMLTemplate) RenderAsHTML(data interface{}) template.HTML {
	// var result strings.Builder
	var buf bytes.Buffer
	if err := tmpl.Tmpl.Execute(&buf, data); err != nil {
		slog.Error(err.Error())
	}
	return template.HTML(buf.String())
}

type IndexPageFields struct {
	Title         string
	DefaultAvatar string
	ArtistAvatar  string
	DisplayName   string
	Links         []template.HTML
}

type LinkPageFields struct {
	IsSpecial   bool
	Link        string
	Description string
}
