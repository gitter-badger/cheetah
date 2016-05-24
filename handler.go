// Copyright 2016 HeadwindFly. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found in the LICENSE file.

package cheetah

import (
	"fmt"
	"net/http"
	"runtime"
	"runtime/debug"
	"strings"
	"github.com/hoisie/mustache"
)

type NotFoundHandler struct {
}

func (this *NotFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	App.errorHandler(w, r, 404, http.StatusText(404), 0)
}

type MethodNotAllowedHandler struct {
}

func (handler *MethodNotAllowedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	App.errorHandler(w, r, 405, http.StatusText(405), 0)
}

func panicHandler(w http.ResponseWriter, r *http.Request, v interface{}) {
	App.errorHandler(w, r, 500, v, 5)
}

type ErrorHandler func(http.ResponseWriter, *http.Request, int, interface{}, int)

func defaultErrorHandler(w http.ResponseWriter, r *http.Request, status int, v interface{}, callDepth int) {
	w.WriteHeader(status)

	title := http.StatusText(status)
	body := fmt.Sprintf("<h1>%d %s</h1>", status, http.StatusText(status))

	if App.mode == ModeDev {
		if _, file, line, ok := runtime.Caller(callDepth); ok {
			body += fmt.Sprintf(`<hr><div class="info">%s: %d</div><br><div class="info">%s</div>`, file, line, v)
		}
		body += `<br><hr><h2>STACK INFO:</h2><hr><div class="stack">`;
		stack := string(debug.Stack())
		stack = strings.Replace(stack, "\n", `<hr>`, -1)
		stack = strings.Replace(stack, "\t", `&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`, -1)
		body += stack + "</div>"
	}
	html := mustache.Render(`
	<html>
<head>
    <title>{{title}}</title>
    <style>
    	h1,h2 {
    	    text-align: center;
    	}
        hr {
            border: 1px dotted;
            color: rgba(3, 169, 244, 0.12);
            clear: both;
        }

        .info {
            color: red;
            font-weight: bold;
            text-align: center;
        }

        .stack {
            margin: 20px 30px;
        }
    </style>
</head>
<body>
{{{body}}}
</body>
</html>
	`, map[string]string{"title":title, "body":body})
	fmt.Fprint(w, html)
}