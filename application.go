// Copyright 2016 HeadwindFly. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found in the LICENSE file.

package cheetah

import (
	"fmt"
	"github.com/go-language/rediscache"
	log "github.com/go-language/logger"
	"github.com/go-language/utils/file"
	"github.com/go-language/utils/ini"
	"net/http"
	"path"
	"strings"
	"github.com/go-language/session"
	"net/smtp"
)

const (
	ModeDev = iota
	ModePro
)

const (
	StateUninitialized = iota
	StateInitialized
	StateRuning
)

const (
	ServerPort = "8080"
	ServerProtocol = "HTTP"

	ControllerPrefix = ""
	ControllerSuffix = "Controller"

	ActionPrefix = "Action"
	ActionSuffix = ""
	DefaultAction = "Index"

	ViewDir = "views"
	ViewSuffix = ".html"
	ViewLayout = "layout.html"
	ViewLayoutDir = "layouts"

	EnableSession = true
	SessionName = "GOSESSION"

	LogDir = "logs"
	LogName = "app.log"

	EnableCsrfValidation = true
	CsrfMaskLength = 8
	CsrfSessionParam = "_csrf"
	CsrfHeaderParam = "X-CSRF-Token"
	CsrfFormParam = "_csrf"

	DefaultRoute = "/index"
)

type Application struct {
	name         string
	mode         int
	basePath     string
	state        int
	language     string
	port         string
	hosts        Hosts
	defaultHost  *Host
	Config       *Config
	errorHandler ErrorHandler
	sessionStore session.Store
	Logger       *log.Logger
	Cache        *rediscache.RedisCache
	resources    map[string]string
}

func NewApplication() Application {
	return Application{
		state:    StateUninitialized,
		name:     "Cheetah Application",
		basePath: "",
		mode:     ModePro,
		language:"en",
		hosts:make(Hosts),
		defaultHost:nil,
		Config: &Config{
			// Server configuration
			serverPort:                         ServerPort,
			serverProtocol:                     ServerProtocol,
			serverCertFile:                     "",
			serverKeyFile:                      "",

			// Controller configuration
			controllerPrefix:             ControllerPrefix,
			controllerSuffix:             ControllerSuffix,

			// Action configuration
			actionPrefix:                 ActionPrefix,
			actionSuffix:                 ActionSuffix,
			defaultAction:                DefaultAction,

			// View configuration
			viewLayout:                   ViewLayout,
			viewLayoutDir:ViewLayoutDir,
			viewDir:                      ViewDir,
			viewSuffix:                   ViewSuffix,

			// Session configuration
			enableSession:                EnableSession,
			sessionName:                   SessionName,
			sessionStore:                 "REDIS",
			sessionMaxAge:                10 * 24 * 3600,

			// CSRF configuration
			enableCsrfValidation:         EnableCsrfValidation,
			csrfMaskLength:               CsrfMaskLength,
			csrfSessionParam:             CsrfSessionParam,
			csrfHeaderParam:              CsrfHeaderParam,
			csrfFormParam:                CsrfFormParam,

			// Log configuration
			enableLog:                    true,
			logFlag:                      log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile,
			logLevel:                     log.LevelDebug | log.LevelInfo | log.LevelWarn | log.LevelError | log.LevelFatal,
			logFileLevel:log.LevelInfo | log.LevelWarn | log.LevelError | log.LevelFatal,
			logFileDir:                       LogDir,
			logFilePath:                      "",
			logFileName:                      LogName,
			logFileMaxSize:                   int64(20 * 1024 * 1024),
			logFileInterval:                  3600,
			logMailLevel:log.LevelError | log.LevelFatal,
			logMailHost:"",
			logMailPort:"",
			logMailUser:"",
			logMailPassword:"",
			logMailFrom:"",
			logMailTo:"",
			logMailSubject:"Application Log",

			// Route configuration
			defaultRoute:                 DefaultRoute,
			routerRedirectTrailingSlash:  true,
			routerRedirectFixedPath:      true,
			routerHandleMethodNotAllowed: true,
			routerHandleOPTIONS:          true,

			// Cache configuration
			enableCache:true,
			redisNetwork:                 "tcp",
			redisAddress:                 ":6379",
			redisPassword:                "",
			redisDb:                      "0",
			redisMaxIdle:                 1000,
			redisIdleTimeout:             300,
		},
		errorHandler: defaultErrorHandler,
		sessionStore: nil,
		Logger:       nil,
		resources:    make(map[string]string, 0),
	}
}

func (this *Application) loadConfig(filename string) {
	var err error

	config := ini.NewConfig(filename)

	section, err := config.GetSection()
	if err != nil {
		panic(err)
	}

	// Set the base configuration
	name, err := section.GetString("name")
	if err == nil {
		this.name = name
	}
	mode, err := section.GetString("mode")
	if strings.EqualFold(mode, "DEV") {
		this.mode = ModeDev
	} else {
		this.mode = ModePro
	}
	basePath, err := section.GetString("base_path")
	if err == nil {
		this.basePath = basePath
	}
	// Set basePath as the configuration's parent directory if it is not specific.
	if len(this.basePath) == 0 {
		this.basePath = path.Dir(path.Dir(filename))
	}

	// Set server configuration
	port, err := section.GetString("server.port")
	if err == nil {
		this.Config.serverPort = port
	}
	protocol, err := section.GetString("server.protocol")
	if err == nil {
		if !strings.EqualFold("HTTP", protocol) && !strings.EqualFold("HTTPS", protocol) {
			panic("The protocol is not supported: " + protocol + ", only support HTTP and HTTPS")
		}
		this.Config.serverProtocol = protocol
	}
	certFile, err := section.GetString("server.cert_file")
	if err == nil {
		this.Config.serverCertFile = certFile
	}
	keyFile, err := section.GetString("server.key_file")
	if err == nil {
		this.Config.serverKeyFile = keyFile
	}

	// Set controller configuration
	controllerPrefix, err := section.GetString("controller.prefix")
	if err == nil {
		this.Config.controllerPrefix = controllerPrefix
	}
	controllerSuffix, err := section.GetString("controller.suffix")
	if err == nil {
		this.Config.controllerSuffix = controllerSuffix
	}

	// Set action configuration
	actionPrefix, err := section.GetString("action.prefix")
	if err == nil {
		this.Config.actionPrefix = actionPrefix
	}
	actionSuffix, err := section.GetString("action.suffix")
	if err == nil {
		this.Config.actionSuffix = actionSuffix
	}
	defaultAction, err := section.GetString("action.default")
	if err == nil {
		this.Config.defaultAction = defaultAction
	}

	// Set view configuration
	viewDir, err := section.GetString("view.dir")
	if err == nil {
		this.Config.viewDir = viewDir
	}
	viewSuffix, err := section.GetString("view.suffix")
	if err == nil {
		this.Config.viewSuffix = viewSuffix
	}
	viewLayout, err := section.GetString("view.layout")
	if err == nil {
		this.Config.viewLayout = viewLayout
	}
	viewLayoutDir, err := section.GetString("view.layout_dir")
	if err == nil {
		this.Config.viewLayoutDir = viewLayoutDir
	}

	// Set session configuration
	enableSession, err := section.GetBool("session.enable")
	if err == nil {
		this.Config.enableSession = enableSession
	}
	sessionName, err := section.GetString("session.name")
	if err == nil {
		this.Config.sessionName = sessionName
	}
	sessionMaxAge, err := section.GetInt("session.max_age")
	if (err == nil) && (sessionMaxAge > 0) {
		this.Config.sessionMaxAge = sessionMaxAge
	}
	sessionStore, err := section.GetString("session.store")
	if err == nil {
		this.Config.sessionStore = sessionStore
	}

	// Set Redis Cache configuration
	redisMaxIdle, err := section.GetInt("redis.max_idle")
	if err == nil {
		this.Config.redisMaxIdle = redisMaxIdle
	}
	redisIdleTimeout, err := section.GetInt("redis.idle_timeout")
	if (err == nil) && (redisIdleTimeout > 0) {
		this.Config.redisIdleTimeout = redisIdleTimeout
	}
	redisNetwork, err := section.GetString("redis.network")
	if err == nil {
		this.Config.redisNetwork = redisNetwork
	}
	redisAddress, err := section.GetString("redis.address")
	if err == nil {
		this.Config.redisAddress = redisAddress
	}
	redisPassword, err := section.GetString("redis.password")
	if err == nil {
		this.Config.redisPassword = redisPassword
	}
	redisDb, err := section.GetString("redis.db")
	if err == nil {
		this.Config.redisDb = redisDb
	}

	// Set CSRF configuration
	enableCsrfValidation, err := section.GetBool("csrf.enable_validation")
	if err == nil {
		this.Config.enableCsrfValidation = enableCsrfValidation
	}
	csrfMaskLength, err := section.GetInt("csrf.mask_length")
	if err == nil {
		this.Config.csrfMaskLength = csrfMaskLength
	}
	csrfSessionParam, err := section.GetString("csrf.session_param")
	if err == nil {
		this.Config.csrfSessionParam = csrfSessionParam
	}
	csrfHeaderParam, err := section.GetString("csrf.header_param")
	if err == nil {
		this.Config.csrfHeaderParam = csrfHeaderParam
	}
	csrfFormParam, err := section.GetString("csrf.form_param")
	if err == nil {
		this.Config.csrfFormParam = csrfFormParam
	}

	// Set log configuration
	enableLog, err := section.GetBool("log.enable")
	if err == nil {
		this.Config.enableLog = enableLog
	}
	logLevel, err := section.GetInt("log.level")
	if err == nil {
		this.Config.logLevel = logLevel
	}
	logFlag, err := section.GetInt("log.flag")
	if err == nil {
		this.Config.logFlag = logFlag
	}
	logFileDir, err := section.GetString("log.file_dir")
	if err == nil {
		this.Config.logFileDir = logFileDir
	}
	logFileName, err := section.GetString("log.file_name")
	if err == nil {
		this.Config.logFileName = logFileName
	}
	logFilePath, err := section.GetString("log.file_path")
	if err == nil {
		this.Config.logFilePath = logFilePath
	}
	logFileMaxSize, err := section.GetInt("log.file_max_size")
	if err == nil {
		this.Config.logFileMaxSize = int64(logFileMaxSize)
	}
	logFileInterval, err := section.GetInt("log.file_interval")
	if err == nil {
		this.Config.logFileInterval = logFileInterval
	}
	logFileLevel, err := section.GetInt("log.file_level")
	if err == nil {
		this.Config.logLevel = logFileLevel
	}
	logMailLevel, err := section.GetInt("log.mail_level")
	if err == nil {
		this.Config.logMailLevel = logMailLevel
	}
	logMailHost, err := section.GetString("log.mail_host")
	if err == nil {
		this.Config.logMailHost = logMailHost
	}
	logMailPort, err := section.GetString("log.mail_port")
	if err == nil {
		this.Config.logMailPort = logMailPort
	}
	logMailUser, err := section.GetString("log.mail_user")
	if err == nil {
		this.Config.logMailUser = logMailUser
	}
	logMailPassword, err := section.GetString("log.mail_password")
	if err == nil {
		this.Config.logMailPassword = logMailPassword
	}
	logMailFrom, err := section.GetString("log.mail_from")
	if err == nil {
		this.Config.logMailFrom = logMailFrom
	}
	logMailTo, err := section.GetString("log.mail_to")
	if err == nil {
		this.Config.logMailTo = logMailTo
	}
	logMailSubject, err := section.GetString("log.mail_subject")
	if err == nil {
		this.Config.logMailSubject = logMailSubject
	}

	// Set router configuration
	defaultRoute, err := section.GetString("router.default")
	if err == nil {
		this.Config.defaultRoute = defaultRoute
	}
	routerRedirectTrailingSlash, err := section.GetBool("router.redirect_trailing_slash")
	if err == nil {
		this.Config.routerRedirectTrailingSlash = routerRedirectTrailingSlash
	}
	routerRedirectFixedPath, err := section.GetBool("router.redirect_fixed_path")
	if err == nil {
		this.Config.routerRedirectFixedPath = routerRedirectFixedPath
	}
	routerHandleMethodNotAllowed, err := section.GetBool("router.handle_method_not_allowed")
	if err == nil {
		this.Config.routerHandleMethodNotAllowed = routerHandleMethodNotAllowed
	}
	routerHandleOPTIONS, err := section.GetBool("router.handle_options")
	if err == nil {
		this.Config.routerHandleOPTIONS = routerHandleOPTIONS
	}

	// Set resources configuration
	resourcesSection, err := config.GetSection("resources")
	if err == nil {
		this.resources = resourcesSection.Params
	}

	// validate configuration
	this.validateConfig()

	// Change Application' state to StateInitialized
	this.state = StateInitialized
}

func (this *Application) validateConfig() {
	// The basePath must be set
	if len(this.basePath) == 0 {
		panic("The basePath must be set.")
	}

	// Check server configuration
	if strings.EqualFold("HTTPS", this.Config.serverProtocol) {
		isCertFileExist, _ := fileutil.IsFile(this.Config.serverCertFile)
		if !isCertFileExist {
			panic("The server's cert file dose not exist: " + this.Config.serverCertFile)
		}
		isKeyFileExist, _ := fileutil.IsFile(this.Config.serverKeyFile)
		if !isKeyFileExist {
			panic("The server's key file dose not exist: " + this.Config.serverKeyFile)
		}
	}

	// Check action configuration
	if (len(this.Config.actionPrefix) == 0) && (len(this.Config.actionSuffix) == 0) {
		panic(fmt.Sprintf("actionPrefix and actionSuffix cannot be empty string at the same time."))
	}

	// Check log configuration
	if this.Config.enableLog {
		if len(this.Config.logFilePath) == 0 {
			this.Config.logFilePath = path.Join(this.basePath, this.Config.logFileDir)
		}
	}

	// Check session configuration
	if this.Config.enableSession {
		if len(this.Config.sessionName) == 0 {
			panic("The session can not be empty string.")
		}
	}

	// Check CSRF configuration
	if this.Config.enableCsrfValidation {
		if !this.Config.enableSession {
			panic("The CSRF validation depend on the session, please enable the session or disable CSRF validation.")
		}
		if this.Config.csrfMaskLength < CsrfMaskLength {
			fmt.Println("The csrfMaskLength is too short, it be set to 8.")
		}
		if len(this.Config.csrfSessionParam) == 0 {
			panic("The csrfSessionParam can not be empty string.")
		}
		if len(this.Config.csrfHeaderParam) == 0 {
			panic("The csrfHeaderParam can not be empty string.")
		}
		if len(this.Config.csrfFormParam) == 0 {
			panic("The csrfFormParam can not be empty string.")
		}
	}
}

func (this *Application) newHost(host string) *Host {
	this.hosts[host] = &Host{
		router:NewRouter(
			App.Config.routerRedirectTrailingSlash,
			App.Config.routerRedirectFixedPath,
			App.Config.routerHandleMethodNotAllowed,
			App.Config.routerHandleOPTIONS,
		),
		routes:make(Routes, 0),
	};
	return this.hosts[host]
}

func (this *Application) run() {
	this.registerRouteHandler()

	var err error

	if this.state == StateUninitialized {
		panic("Please initialize the Application by invoking the method: " + "cheetah.Init(\"/path/to/ini_config_file\")")
	}

	// Register logger
	if this.Config.enableLog {
		this.Logger = log.NewLogger(
			this.Config.logLevel,
			this.Config.logFlag,
		)
		defer this.Logger.Close()

		// Add FileTarget
		logFile, err := log.OpenFile(path.Join(this.Config.logFilePath, this.Config.logFileName))
		if err != nil {
			panic(err.Error())
		}

		if len(this.Config.logFilePath) > 0 {
			fileTarget := log.NewFileTarget(this.Logger, this.Config.logFileLevel, logFile)

			go fileTarget.Crontab()

			this.Logger.AddTarget(fileTarget)
		}

		if len(this.Config.logMailHost) > 0 {
			auth := smtp.PlainAuth("", this.Config.logMailUser, this.Config.logMailPassword, this.Config.logMailHost)

			mailTarget := log.NewMailTarget(
				this.Config.logMailLevel,
				this.Config.logMailHost + ":" + this.Config.logMailPort,
				this.Config.logMailFrom,
				this.Config.logMailTo,
				auth,
			)

			mailTarget.SetSubject(this.Config.logMailSubject)
			this.Logger.AddTarget(mailTarget)
		}
	}

	// Register Cache
	if this.Config.enableCache {
		redisPool := rediscache.NewRedisPool(
			this.Config.redisMaxIdle,
			this.Config.redisIdleTimeout,
			this.Config.redisNetwork,
			this.Config.redisAddress,
			this.Config.redisPassword,
			this.Config.redisDb,
		)

		defer redisPool.Close()

		this.Cache = rediscache.NewRedisCache(redisPool)
	}

	// Register session store
	if this.Config.enableSession {
		if !this.Config.enableCache {
			panic("The session depends on redis cache, please enable the cache component.")
		}

		store := session.NewRedisStore(this.Cache.GetPool(), session.Options{})

		if err != nil {
			panic(err)
		}

		store.SetMaxAge(this.Config.sessionMaxAge)

		SetSessionStore(store)
	}

	this.state = StateRuning

	fmt.Println("Application started.")

	if len(this.hosts) == 0 {
		panic("No host.")
	}

	var handler http.Handler
	if len(this.hosts) == 1 {
		for _, host := range this.hosts {
			handler = host.router
		}

	} else {
		handler = this.hosts
		if this.defaultHost == nil {
			panic("The default host must be set.")
		}
	}

	addr := ":" + this.Config.serverPort
	// If the protocol equal HTTPS
	if strings.EqualFold("HTTPS", this.Config.serverProtocol) {
		err = http.ListenAndServeTLS(
			addr,
			this.Config.serverCertFile,
			this.Config.serverKeyFile,
			handler,
		)
	} else {
		err = http.ListenAndServe(addr, handler)
	}

	if err != nil {
		panic(err.Error())
	}
}

// Register route handler.
func (this *Application) registerRouteHandler() {
	for _, host := range this.hosts {
		host.generateRouteHandle()
	}
}