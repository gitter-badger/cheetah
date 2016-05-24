// Copyright 2016 HeadwindFly. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found in the LICENSE file.

package cheetah

// Configuration of application.
type Config struct {
	// Server Configuration
	serverAddr                   string
	serverProtocol               string
	serverCertFile               string
	serverKeyFile                string

	// Controller Configuration
	controllerPrefix             string
	controllerSuffix             string

	// Action Configuration
	actionPrefix                 string
	actionSuffix                 string
	defaultAction                string

	// View Configuration
	viewLayout                   string
	viewLayoutDir                string
	viewDir                      string
	viewSuffix                   string

	// Session Configuration
	enableSession                bool
	sessionName                  string
	sessionStore                 string
	sessionMaxAge                int

	// Log Configuration
	enableLog                    bool
	logLevel                     int
	logFlag                      int
	// FileTartget
	logFileLevel                 int
	logFileDir                   string
	logFileName                  string
	logFilePath                  string
	logFileMaxSize               int64
	logFileInterval              int
	// EmailTarget
	logMailLevel                 int
	logMailHost                  string
	logMailPort                  string
	logMailUser                  string
	logMailPassword              string
	logMailFrom                  string
	logMailTo                    string
	logMailSubject               string

	// CSRF Configuration
	enableCsrfValidation         bool
	csrfMaskLength               int
	csrfSessionParam             string
	csrfHeaderParam              string
	csrfFormParam                string

	// Router Configuration
	defaultRoute                 string
	routerRedirectTrailingSlash  bool
	routerRedirectFixedPath      bool
	routerHandleMethodNotAllowed bool
	routerHandleOPTIONS          bool

	// Redis Configuration
	enableCache                  bool
	redisNetwork                 string
	redisAddress                 string
	redisPassword                string
	redisDb                      string
	redisMaxIdle                 int
	redisIdleTimeout             int
}

func (this *Config) ServerAddr() string {
	return this.serverAddr
}

func (this *Config) ServerProtocol() string {
	return this.serverProtocol
}

func (this *Config) ServerCertFile() string {
	return this.serverCertFile
}

func (this *Config) ServerKeyFile() string {
	return this.serverKeyFile
}

func (this *Config) ControllerPrefix() string {
	return this.controllerPrefix
}

func (this *Config) ControllerSuffix() string {
	return this.controllerSuffix
}

func (this *Config) ActionPrefix() string {
	return this.actionPrefix
}

func (this *Config) ActionSuffix() string {
	return this.actionSuffix
}

func (this *Config) DefaultAction() string {
	return this.defaultAction
}

func (this *Config) ViewLayout() string {
	return this.viewLayout
}

func (this *Config) ViewDir() string {
	return this.viewDir
}

func (this *Config) ViewSuffix() string {
	return this.viewSuffix
}

func (this *Config) EnableSession() bool {
	return this.enableSession
}

func (this *Config) SessionName() string {
	return this.sessionName
}

func (this *Config) EnableLog() bool {
	return this.enableLog
}

func (this *Config) LogFileDir() string {
	return this.logFileDir
}

func (this *Config) LogFilePath() string {
	return this.logFilePath
}

func (this *Config) LogFileName() string {
	return this.logFileName
}

func (this *Config) LogLevel() int {
	return this.logLevel
}

func (this *Config) LogFlag() int {
	return this.logFlag
}

func (this *Config) LogFileMaxSize() int64 {
	return this.logFileMaxSize
}

func (this *Config) LogFileInterval() int {
	return this.logFileInterval
}

func (this *Config) EnableCsrfValidation() bool {
	return this.enableCsrfValidation
}

func (this *Config) CsrfHeaderParam() string {
	return this.csrfHeaderParam
}

func (this *Config) CsrfFormParam() string {
	return this.csrfFormParam
}

func (this *Config) DefaultRoute() string {
	return this.defaultRoute
}
