package response

import "testing"

func TestRenderTemplate(t *testing.T) {
	content := `用户请求：{{.Message}}
意图：{{.Intent}}`

	out, err := RenderTemplate(content, TemplateData{
		Message: "请分析这个系统",
		Intent:  "analysis",
	})
	if err != nil {
		t.Fatalf("RenderTemplate error = %v", err)
	}

	if out == "" {
		t.Fatal("rendered output is empty")
	}
	t.Log(out)
}
