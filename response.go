// Copyright 2016 HeadwindFly. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found in the LICENSE file.

package cheetah

import (
	"fmt"
	"net/http"
)

type Response struct {
	Writer http.ResponseWriter // Response writer
	IsSent bool                // Whether the response has been sent, default as false.
	Status int                 // Http response status, default as 200.
	Body   string              // Response body
}

func NewResponse(w *http.ResponseWriter) *Response {
	return &Response{
		*w, false, http.StatusOK, "",
	}
}

func (this *Response) SetHeader(key, value string) {
	this.Writer.Header().Set(key, value)
}

func (this *Response) SetHtmlHeader() {
	this.SetHeader("Content-Type", "text/html; charset=utf-8")
}

func (this *Response) SetJsonHeader() {
	this.SetHeader("Content-Type", "application/json; charset=utf-8")
}

func (this *Response) SetJsonpHeader() {
	this.SetHeader("Content-Type", "application/javascript; charset=utf-8")
}

func (this *Response) SetXmlHeader() {
	this.SetHeader("Content-Type", "application/xml; charset=utf-8")
}

type WebResponse struct {
	*Response
}

func NewWebResponse(w *http.ResponseWriter) *WebResponse {
	return &WebResponse{
		NewResponse(w),
	}
}

// Send response to client.
func (this *WebResponse) Send() {
	// The header will only be sent once.
	if !this.IsSent {
		this.IsSent = true
		this.sendHeaders()
		this.sendBody()
	}
}

func (this *WebResponse) sendHeaders() {
	this.Writer.WriteHeader(this.Status)
}

func (this *WebResponse) sendBody() {
	fmt.Fprintf(this.Writer, this.Body)
	this.Body = ""
}

func (this *WebResponse) NotFound(data string) {
	this.Status = http.StatusNotFound
	this.Body = http.StatusText(this.Status) + ": " + data
	this.Send()

}

func (this *WebResponse) InternalServerError(data string) {
	this.Status = http.StatusInternalServerError
	this.Body = http.StatusText(this.Status) + ": " + data
	this.Send()
}

func (this *WebResponse) BadRequest(data string) {
	this.Status = http.StatusBadRequest
	this.Body = http.StatusText(this.Status) + ": " + data
	this.Send()
}

func (this *WebResponse) Forbidden(data string) {
	this.Status = http.StatusForbidden
	this.Body = http.StatusText(this.Status) + ": " + data
	this.Send()
}

func (this *WebResponse) Redirect(url string) {
	this.Status = http.StatusFound
	this.SetHeader("Location", url)
	this.Send()
}
