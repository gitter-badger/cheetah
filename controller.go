// Copyright 2016 HeadwindFly. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found in the LICENSE file.

package cheetah

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	log "github.com/go-language/logger"
	"github.com/go-language/utils/string"
	"github.com/hoisie/mustache"
	"net/http"
	"path"
	"github.com/go-language/session"
)

// Controller Interface.
type ControllerInterface interface {
	Init(config *ControllerInfo, w *http.ResponseWriter, r *http.Request)
	BeforeAction() bool
	BeforeResponse()
	ResponseClient()
}

// Controller Config.
type ControllerInfo struct {
	PkgPath        string   // package path of the controller.
	Name           string   // controller's name.
	FullName       string   // controller's full name.
	ViewPath       string   // view's path.
	ActionFullName string   // method' full name.
	ActionName     string   // action name.
	Params         []string // params of action,such as {"string","int"} means that the first param type of string,the second param type of int.
	Layout         string   // layout's name.
	Log            *log.Log // log
}

type WebController struct {
	Name     string           // controller's name.
	FullName string           // controller's full name.
	PkgPath  string           // controller's package path.
	Action   string           // controller's action name
	ViewPath string           // view's path
	Layout   string           // layout's name, if empty means that do not use layout.
	Context  *Context         // Context
	Response *WebResponse     // web response
	Session  *session.Session // session
	Log      *log.Log         // log
}

func (this *WebController) Init(info *ControllerInfo, w *http.ResponseWriter, r *http.Request) {
	this.Name = info.Name
	this.PkgPath = info.PkgPath
	this.ViewPath = info.ViewPath
	this.Action = info.ActionName
	this.Layout = info.Layout
	this.Log = info.Log

	this.Context = NewContext(w, r)

	this.Response = NewWebResponse(w)

	this.getSession(r)

	this.validateCsrfToken()
}

// Do something before invling the action.
// If true was returned, the action will be invoked.
// But on the contrary, it will not invoke the action, just response client directly.
func (this *WebController) BeforeAction() bool {
	return true
}

func (this *WebController) BeforeResponse() {

}

// Response client.
func (this *WebController) ResponseClient() {
	this.saveSession()
	this.Response.Send()
}

func (this *WebController) validateCsrfToken() {
	if App.Config.enableCsrfValidation && !this.Context.ValidateCsrfToken() {
		this.Response.BadRequest("Unable to verify your data submission.")
	}
}

func (this *WebController) getSession(r *http.Request) {
	if App.Config.enableSession {
		var err error
		this.Session, err = App.sessionStore.Get(r, App.Config.sessionName)
		if err != nil {
			this.Response.InternalServerError(err.Error())
			return
		}

		this.Context.trueCsrfToken = this.getTrueCsrfToken()
		this.Context.csrfToken = GenerateCsrfToken(App.Config.csrfMaskLength, []byte(this.Context.trueCsrfToken))
	}
}

func (this *WebController) getTrueCsrfToken() string {
	token, ok := this.Session.Values["_csrf"]

	if !ok {
		token = stringutil.GenerateRandomString(32)
		this.Session.Values["_csrf"] = token
	}
	return token.(string)
}

func (this *WebController) saveSession() {
	if App.Config.enableSession {
		if err := this.Session.Save(this.Response.Writer); err != nil {
			this.Response.InternalServerError(fmt.Sprintf("Error saving session: %v", err))
		}
	}
}

func (this *WebController) Render(context ...interface{}) {
	this.RenderFile("", context...)
}

func (this *WebController) RenderData(data string, context ...interface{}) {
	this.Response.SetHtmlHeader()
	this.Response.Body = mustache.Render(data, context...)
}

// @param name the view file name
func (this *WebController) RenderFile(name string, context ...interface{}) {
	this.Response.SetHtmlHeader()

	if len(name) == 0 {
		name = BuildPrettyRoute(this.Action) + App.Config.viewSuffix
	} else {
		name = name + App.Config.viewSuffix
	}
	file := this.getViewFile(name)

	if len(this.Layout) > 0 {
		this.Response.Body = mustache.RenderFileInLayout(file, this.getLayoutFile(), context...)
	} else {
		this.Response.Body = mustache.RenderFile(file, context...)
	}
}

func (this *WebController) RenderPartial(context ...interface{}) {
	this.RenderPartialFile("", context...)
}

func (this *WebController) RenderPartialFile(name string, context ...interface{}) {
	if len(name) == 0 {
		name = BuildPrettyRoute(this.Action) + App.Config.viewSuffix
	} else {
		name = name + App.Config.viewSuffix
	}
	file := this.getViewFile(name)

	this.Response.Body = mustache.RenderFile(file, context...)
}

func (this *WebController) getLayoutFile() string {
	return path.Join(path.Dir(this.ViewPath), App.Config.viewLayoutDir, this.Layout)
}

// the v will be responsed directly if type of v is string
func (this *WebController) RenderJson(v interface{}) {
	this.Response.SetJsonHeader()

	if value, ok := v.(string); ok {
		this.Response.Body = value
	} else {
		json, err := json.Marshal(v)
		if err != nil {
			this.Response.InternalServerError(err.Error())
			return
		}
		this.Response.Body += string(json)
	}
}

// the v will be responsed directly if type of v is string
func (this *WebController) RenderJsonp(v interface{}, callback string) {
	this.Response.SetJsonpHeader()

	if value, ok := v.(string); ok {
		this.Response.Body = value
	} else {

		json, err := json.Marshal(v)
		if err != nil {
			this.Response.InternalServerError(err.Error())
			return
		}
		this.Response.Body += callback + "(" + string(json) + ")"
	}
}

func (this *WebController) RenderText(text string) {
	this.Response.SetHtmlHeader()
	this.Response.Body = text
}

// the v will be responsed directly if type of v is string
func (this *WebController) RenderXml(v interface{}, header string) {
	this.Response.SetXmlHeader()

	if value, ok := v.(string); ok {
		this.Response.Body = value
	} else {
		byteXML, err := xml.MarshalIndent(v, "", `   `)
		if err != nil {
			this.Response.InternalServerError(err.Error())
			return
		}

		if len(header) == 0 {
			header = xml.Header
		}

		this.Response.Body = header + string(byteXML)
	}
}

func (this *WebController) getViewFile(name string) string {
	return path.Join(this.ViewPath, name)
}

// Get layout name
// If it is set as empty string means that disabled the layout.
// Returns zero string '0' default, means that use the global layout.
// Layout enable default.
func (this *WebController) GetLayout() string {
	return "0"
}
