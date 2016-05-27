// Copyright 2016 HeadwindFly. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found in the LICENSE file.

package cheetah

import (
	"github.com/HeadwindFly/cheetah/utils/string"
	"github.com/go-language/session"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

var (
	App Application
)

func init() {
	App = NewApplication()
}

func Init(filename string) {
	App.loadConfig(filename)
}

func Run() {
	App.run()
}

func NewHost(host string) *Host {
	return App.newHost(host)
}

func generateRouteHandle(route string, controllerType reflect.Type, info *ControllerInfo) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if App.Config.enableLog {
			info.Log = App.Logger.NewLog()
			defer info.Log.Flush()
		} else {
			info.Log = nil
		}

		v := reflect.New(controllerType)

		initArgs := []reflect.Value{
			reflect.ValueOf(info),
			reflect.ValueOf(&w),
			reflect.ValueOf(r),
		}
		initMethod := v.MethodByName("Init")
		initMethod.Call(initArgs)

		beforeActionMethod := v.MethodByName("BeforeAction")
		beforeResult := beforeActionMethod.Call([]reflect.Value{})
		canInvokeAction := true
		for _, value := range beforeResult {
			if _value, ok := value.Interface().(bool); ok {
				canInvokeAction = _value
			}
			break
		}

		if canInvokeAction {
			params := []reflect.Value{}
			if len(info.Params) > 0 {
				i := 0
				for ; i < len(ps); i++ {
					if info.Params[i] == "string" {
						params = append(params, reflect.ValueOf(ps[i].Value))
					} else {
						value, err := strconv.Atoi(ps[i].Value)
						if err != nil {
							panic("Invalid Params")
						}
						params = append(params, reflect.ValueOf(value))
					}
				}
				for ; i < len(info.Params); i++ {
					if info.Params[i] == "string" {
						params = append(params, reflect.ValueOf(""))
					} else {
						params = append(params, reflect.ValueOf(1))
					}
				}
			}

			// invoke the action.
			actionMethod := v.MethodByName(info.ActionFullName)
			actionMethod.Call(params)
		}

		// return response to client.
		responseClientMethod := v.MethodByName("ResponseClient")
		responseClientMethod.Call([]reflect.Value{})
	}
}

// Remove the controller's prefix and suffix that you set.
// If it is not a controller, false will be return.
func getControllerName(name string) (string, bool) {
	// remove prefix
	if len(App.Config.controllerPrefix) > 0 {
		if 0 != strings.Index(name, App.Config.controllerPrefix) {
			return "", false
		}

		prefixLen := len(App.Config.controllerPrefix)
		name = stringutil.SubString(name, prefixLen, len(name)-prefixLen)
	}
	// remove suffix
	if len(App.Config.controllerSuffix) > 0 {
		pos := len(name) - len(App.Config.controllerSuffix)

		if (pos == -1) || (pos != strings.Index(name, App.Config.controllerSuffix)) {
			return "", false
		}

		name = stringutil.SubString(name, 0, pos)
	}
	return name, true
}

// Remove the action's prefix and suffix that you set.
// If it is not a action, false will be return.
func getActionName(name string) (string, bool) {
	// the first character of action must be uppercase.
	if ('A' > name[0]) || (name[0] > 'Z') {
		return "", false
	}

	// remove prefix
	if len(App.Config.actionPrefix) > 0 {
		if 0 != strings.Index(name, App.Config.actionPrefix) {
			return "", false
		}

		prefixLen := len(App.Config.actionPrefix)
		name = stringutil.SubString(name, prefixLen, len(name)-prefixLen)
	}
	// remove suffix
	if len(App.Config.actionSuffix) > 0 {
		pos := len(name) - len(App.Config.actionSuffix)

		if (pos == -1) || (pos != strings.Index(name, App.Config.actionSuffix)) {
			return "", false
		}

		name = stringutil.SubString(name, 0, pos)
	}
	return name, true
}

func SetErrorHandler(handler ErrorHandler) {
	App.errorHandler = handler
}

func SetSessionStore(store session.Store) {
	App.sessionStore = store
}

func SetDefaultHost(host *Host) {
	App.defaultHost = host
}

func Name() string {
	return App.name
}

func Mode() int {
	return App.mode
}

func BasePath() string {
	return App.basePath
}

func State() int {
	return App.state
}
