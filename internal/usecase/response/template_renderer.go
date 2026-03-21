package response

import (
	"bytes"
	"text/template"
)

// TemplateData 表示传给 prompt 模板的变量
type TemplateData struct {
	Message               string
	Intent                string
	StepResultsText       string
	EvidencesText         string
	CapabilityResultsText string
}

// RenderTemplate 渲染模板
func RenderTemplate(content string, data TemplateData) (string, error) {
	tpl, err := template.New("prompt").Parse(content)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
