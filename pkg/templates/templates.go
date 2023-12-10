package templates

import "text/template"

type Constructor struct {
	Provider string
	templ    template.Template
}

func (c *Constructor) InitTemplate() {}

// added line here
