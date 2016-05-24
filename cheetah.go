// Copyright 2016 HeadwindFly. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found in the LICENSE file.

package cheetah

import (
	"github.com/go-language/utils/string"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"github.com/go-language/session"
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

// Register web controller
func RegisterWebController(baseRoute string, controller ControllerInterface) {
	if (len(baseRoute) == 0) || (baseRoute[0] != '/') {
		panic("The first character of route named \"" + baseRoute + "\" must be \"/\".")
	}

	if len(baseRoute) < 2 {
		panic("The length of route named \"" + baseRoute + "\" must greater than one.")
	}

	t := reflect.TypeOf(controller)
	v := reflect.ValueOf(controller)

	controllerName := t.Elem().Name() // the full name of controller

	controllerName, ok := getControllerName(controllerName)
	if !ok {
		panic("The " + t.Elem().Name() + "'s name is invalid, The controller's prefix and suffix must be '" + App.Config.controllerPrefix + "' and '" + App.Config.controllerSuffix + "'.")
	}

	// get method filter
	methodFilter := MethodFilter{}
	methodFilterMethod := v.MethodByName("MethodFilter")
	if methodFilterMethod.IsValid() {
		values := methodFilterMethod.Call([]reflect.Value{})
		for _, value := range values {
			if _value, ok := value.Interface().(MethodFilter); ok {
				methodFilter = _value
			}
			break
		}
	}

	// get package path
	pkgPath := path.Join(os.Getenv("GOPATH"), "src", v.Elem().Type().PkgPath())

	// set view path
	viewPath := path.Join(path.Dir(pkgPath), App.Config.viewDir, BuildPrettyRoute(controllerName))

	// get layout
	viewLayout := ""
	layoutMethod := v.MethodByName("GetLayout")
	if layoutMethod.IsValid() {
		values := layoutMethod.Call([]reflect.Value{})
		for _, value := range values {
			if _value, ok := value.Interface().(string); ok && !strings.EqualFold(_value, "0") {
				viewLayout = _value
			}
			break
		}
	}

	if 0 == strings.Compare("0", viewLayout) {
		viewLayout = App.Config.viewLayout
	}

	for j := 0; j < t.NumMethod(); j++ {
		_routes := []string{}

		method := t.Method(j)

		actionName, ok := getActionName(method.Name)
		if !ok {
			continue
		}

		// get allow methods.
		allowMethods := []string{"GET", "POST"}
		if len(methodFilter) > 0 {
			if value, ok := methodFilter[actionName]; ok {
				allowMethods = value
			}
		}

		// if the current action is the default action, add route
		if strings.EqualFold(actionName, App.Config.defaultAction) {
			_routes = append(_routes, baseRoute)
		}

		methodType := reflect.TypeOf(v.Method(j).Interface())

		actionRoute := BuildPrettyRoute(actionName)

		route := baseRoute + "/" + actionRoute

		_routes = append(_routes, route) // add route

		params := []string{}

		if methodType.NumIn() > 0 {
			routeWithParams := route
			for k := 0; k < methodType.NumIn(); k++ {
				paramKind := methodType.In(k).Kind().String()
				// the param'kind must be string or int
				if (paramKind != "string") && (paramKind != "int") {
					panic("The type of " + v.Elem().Type().String() + "." + method.Name + "()" + "'s params must be string or int.")
				}
				params = append(params, paramKind)

				routeWithParams += "/:" + string(rune(97 + k))

				_routes = append(_routes, routeWithParams) // add route
			}
		}
		// add route to the routes map.
		for i := 0; i < len(_routes); i++ {
			App.routes[_routes[i]] = &RouteInfo{
				Route:          _routes[i],
				AllowMethods:   allowMethods,
				ControllerType: v.Elem().Type(),
				ControllerInfo: &ControllerInfo{
					PkgPath:        pkgPath,
					FullName:       t.Elem().Name(),
					Name:           controllerName,
					ViewPath:       viewPath,
					ActionFullName: method.Name,
					ActionName:     actionName,
					Layout:         viewLayout,
					Params:         params,
				},
			}
		}
	}
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
		name = stringutil.SubString(name, prefixLen, len(name) - prefixLen)
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
		name = stringutil.SubString(name, prefixLen, len(name) - prefixLen)
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

func SetPanicHandler(handler func(http.ResponseWriter, *http.Request, interface{})) {
	App.router.PanicHandler = handler
}

func SetSessionStore(store session.Store) {
	App.SetSessionStore(store)
}
