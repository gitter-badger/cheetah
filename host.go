package cheetah

import (
	"reflect"
	"path"
	"os"
	"strings"
	"github.com/julienschmidt/httprouter"
	"fmt"
	"net/http"
)

type Host struct {
	router *Router
	routes Routes
}

func (this *Host) SetNotFoundHandler(handler http.Handler) {
	this.router.NotFound = handler
}

func (this *Host) SetMethodNotAllowedHandler(handler http.Handler) {
	this.router.MethodNotAllowed = handler
}

func (this *Host) SetPanicHandler(handler func(http.ResponseWriter, *http.Request, interface{})) {
	this.router.PanicHandler = handler
}

func (this *Host) RegisterWebController(baseRoute string, controller ControllerInterface) {
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
			this.routes[_routes[i]] = &RouteInfo{
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

func (this *Host) generateRouteHandle() {
	for key, route := range this.routes {
		handle := generateRouteHandle(route.Route, route.ControllerType, route.ControllerInfo)
		for i := 0; i < len(route.AllowMethods); i++ {
			this.router.Handle(route.AllowMethods[i], route.Route, handle)
		}
		this.routes[key].Handle = &handle
	}

	// Register default route.
	if route, ok := this.routes[App.Config.defaultRoute]; ok {
		this.router.Handle("GET", "/", *route.Handle)
	} else {
		this.router.Handle("GET", "/", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			fmt.Fprint(w, "Hello Cheetah.")
		})
	}
}

func (this *Host) RegisterResources(route, path string) {
	this.router.ServeFiles("/" + route + "/*filepath", http.Dir(path))
}

type Hosts map[string]*Host

func (this Hosts) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get domain from host.
	host := strings.Split(r.Host, ":")
	if host, ok := this[host[0]]; ok {
		host.router.ServeHTTP(w, r)
	} else {
		App.defaultHost.router.ServeHTTP(w, r)
	}
}