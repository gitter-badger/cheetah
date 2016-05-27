// Copyright 2016 HeadwindFly. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found in the LICENSE file.

package cheetah

import (
	"net/http"
	"strings"
)

type Context struct {
	Request       *http.Request
	csrfToken     string
	trueCsrfToken string
}

func NewContext(w *http.ResponseWriter, r *http.Request) *Context {
	r.ParseForm()
	return &Context{
		Request: r,
	}
}

// Return current request's CSRF token.
func (this *Context) CsrfToken() string {
	return this.csrfToken
}

// Validate CSRF Token.
// True is returned if the token is valid, false is returned otherwise.
func (this *Context) ValidateCsrfToken() bool {
	if strings.EqualFold("GET", this.Request.Method) || strings.EqualFold("HEAD", this.Request.Method) {
		return true
	}
	if len(this.trueCsrfToken) == 0 {
		return false
	}

	if ValidateCsrfToken(App.Config.csrfMaskLength, this.getCsrfTokenFromForm(), this.trueCsrfToken) ||
		ValidateCsrfToken(App.Config.csrfMaskLength, this.getCsrfTokenFromHeader(), this.trueCsrfToken) {
		return true
	}
	return false
}

// Get CSRF token from request's header.
func (this *Context) getCsrfTokenFromHeader() string {
	return this.Request.Header.Get("X-CSRF-Token")
}

// Get CSRF token from post form.
func (this *Context) getCsrfTokenFromForm() string {
	return this.Request.PostFormValue("_csrf")
}

// Returns a boolean indicating whether this is a GET request.
func (this *Context) IsGet() bool {
	return strings.EqualFold("GET", this.Request.Method)
}

// Returns a boolean indicating whether this is a POST request.
func (this *Context) IsPost() bool {
	return strings.EqualFold("POST", this.Request.Method)
}

// Returns a boolean indicating whether this is a PUT request.
func (this *Context) IsPut() bool {
	return strings.EqualFold("PUT", this.Request.Method)
}

// Returns a boolean indicating whether this is a HEAD request.
func (this *Context) IsHead() bool {
	return strings.EqualFold("HEAD", this.Request.Method)
}

// Returns a boolean indicating whether this is a DELETE request.
func (this *Context) IsDelete() bool {
	return strings.EqualFold("DELETE", this.Request.Method)
}

// Returns a boolean indicating whether this is a PATCH request.
func (this *Context) IsPatch() bool {
	return strings.EqualFold("PATCH", this.Request.Method)
}

// Returns a boolean indicating whether this is a OPTIONS request.
func (this *Context) IsOptions() bool {
	return strings.EqualFold("OPTIONS", this.Request.Method)
}

// Returns a boolean indicating whether this is a AJAX request.
func (this *Context) IsAjax() bool {
	header := this.Request.Header.Get("X-Requested-With")
	return strings.Compare("XMLHttpRequest", header) == 0
}
