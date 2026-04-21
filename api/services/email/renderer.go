package email

import (
	"bytes"
	"fmt"
	htmltmpl "html/template"
	texttmpl "text/template"

	"api/models"
)

// Rendered is the output of applying template variables to an EmailTemplate.
type Rendered struct {
	Subject string
	Body    string
}

// Render executes the template's subject (text/template, plain) and body
// (html/template, auto-escaped) against vars. The subject is deliberately
// kept out of html/template so it isn't HTML-escaped — subjects are plain text
// per RFC 5322.
func Render(tpl models.EmailTemplate, vars map[string]interface{}) (Rendered, error) {
	subject, err := renderText(tpl.Code+":subject", tpl.SubjectTemplate, vars)
	if err != nil {
		return Rendered{}, fmt.Errorf("render subject: %w", err)
	}
	body, err := renderHTML(tpl.Code+":body", tpl.BodyTemplate, vars)
	if err != nil {
		return Rendered{}, fmt.Errorf("render body: %w", err)
	}
	return Rendered{Subject: subject, Body: body}, nil
}

// Preview renders the template using the given sample vars. It is a thin
// pass-through today but exists so the CRUD preview endpoint has a stable
// entry point regardless of later pre-processing changes.
func Preview(tpl models.EmailTemplate, vars map[string]interface{}) (Rendered, error) {
	return Render(tpl, vars)
}

func renderText(name, src string, vars map[string]interface{}) (string, error) {
	t, err := texttmpl.New(name).Option("missingkey=zero").Parse(src)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, vars); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func renderHTML(name, src string, vars map[string]interface{}) (string, error) {
	t, err := htmltmpl.New(name).Option("missingkey=zero").Parse(src)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, vars); err != nil {
		return "", err
	}
	return buf.String(), nil
}
