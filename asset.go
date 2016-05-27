package cheetah

import (
	"fmt"
	"strings"
)

type Asset interface {
	Output() string
	Option(key, value string) Asset
	Condition(condition string) Asset
}

type CssAsset struct {
	Href      string
	Rel       string
	Type      string
	condition string
	Options   map[string]string
}

func NewCssAsset(href string) *CssAsset {
	return &CssAsset{
		Href:      href,
		Rel:       "stylesheet",
		Type:      "text/css",
		condition: "",
		Options:   make(map[string]string, 0),
	}
}

func (this *CssAsset) Output() string {
	args := make([]interface{}, 0)
	format := "<link rel=\"%s\" type=\"%s\" href=\"%s\"/>"
	args = append(args, this.Rel, this.Type, this.Href)

	options := make([]string, 0)
	if len(this.Options) > 0 {
		format = "<link rel=\"%s\" type=\"%s\" href=\"%s\" %s/>"
		for k, v := range this.Options {
			options = append(options, k+"=\""+v+"\"")
		}
		args = append(args, strings.Join(options, " "))
	}

	asset := fmt.Sprintf(format, args...)

	if len(this.condition) > 0 {
		asset = "<!--[if " + this.condition + "]> -->" + asset + "<!-- <![endif]-->"
	}
	return asset
}

func (this *CssAsset) Option(key, value string) *CssAsset {
	this.Options[key] = value
	return this
}

func (this *CssAsset) Condition(condition string) *CssAsset {
	this.condition = condition
	return this
}

type JsAsset struct {
	Src       string // it links a  javascript source file  If src is not empty.
	Script    string // it is a javascript if script is not empty.
	Type      string
	condition string
	Options   map[string]string
}

func NewJsAsset(src, script string) *JsAsset {
	return &JsAsset{
		Src:       src,
		Script:    script,
		Type:      "text/javascript",
		condition: "",
		Options:   make(map[string]string, 0),
	}
}

func (this *JsAsset) Output() string {
	args := make([]interface{}, 0)
	format := ""

	if len(this.Src) > 0 {
		args = append(args, this.Type, this.Src)
		if len(this.Options) > 0 {
			format = "<script type=\"%s\" src=\"%s\" %s></script>"
		} else {
			format = "<script type=\"%s\" src=\"%s\"></script>"
		}
	} else {
		args = append(args, this.Type, this.Script)
		if len(this.Options) > 0 {
			format = "<script type=\"%s\" %s></script>"
		} else {
			format = "<script type=\"%s\">%s</script>"
		}
	}

	options := make([]string, 0)
	if len(this.Options) > 0 {
		format = "<link rel=\"%s\" type=\"%s\" href=\"%s\" %s/>"
		for k, v := range this.Options {
			options = append(options, k+"=\""+v+"\"")
		}
		args = append(args, strings.Join(options, " "))
	}

	asset := fmt.Sprintf(format, args...)

	if len(this.condition) > 0 {
		asset = "<!--[if " + this.condition + "]> -->" + asset + "<!-- <![endif]-->"
	}
	return asset
}

func (this *JsAsset) Option(key, value string) *JsAsset {
	this.Options[key] = value
	return this
}

func (this *JsAsset) Condition(condition string) *JsAsset {
	this.condition = condition
	return this
}
