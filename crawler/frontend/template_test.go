package frontend

import (
	"html/template"
	"testing"
)

func TestTemplate(t *testing.T) {
	template.ParseFiles("template.html")
}
