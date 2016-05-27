// Copyright 2016 HeadwindFly. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found in the LICENSE file.

package cheetah

import (
	"github.com/julienschmidt/httprouter"
	"reflect"
	"strings"
)

type Router struct {
	*httprouter.Router
}

func NewRouter(redirectTrailingSlash, redirectFixedPath, handleMethodNotAllowed, handleOPTIONS bool) *Router {
	return &Router{
		&httprouter.Router{
			RedirectTrailingSlash:  redirectTrailingSlash,
			RedirectFixedPath:      redirectFixedPath,
			HandleMethodNotAllowed: handleMethodNotAllowed,
			HandleOPTIONS:          handleOPTIONS,
			NotFound:               new(NotFoundHandler),
			MethodNotAllowed:       new(MethodNotAllowedHandler),
			PanicHandler:           panicHandler,
		},
	}
}

type Routes map[string]*RouteInfo

type RouteInfo struct {
	Route          string             // route, such as "/", "/user" "/user/post" etc.
	AllowMethods   []string           // allowed methods.
	ControllerType reflect.Type       // controller's reflect.Type.
	ControllerInfo *ControllerInfo    // controller's info
	Handle         *httprouter.Handle // route handle.
}

// PostController.CommentAdd()'s route will be formated as "/post/comment-add"
func BuildPrettyRoute(route string) string {
	if len(route) == 0 {
		return ""
	}
	prettyRoute := strings.ToLower(string(route[0]))
	for i := 1; i < len(route); i++ {
		c := route[i]
		if ('A' <= c) && (c <= 'Z') {
			prettyRoute += "-" + string(rune(int(c)+32))
		} else {
			prettyRoute += string(c)
		}
	}
	return prettyRoute
}
