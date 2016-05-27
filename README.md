# Cheetah
**Cheetah** is a **simple** web framework for **Go**.
 This project aims to become a powerful web development tool and **make developer easily to build a high-performance, secure and stronger web application**.

## Cheetah is still under development, but it will come soon.


# Feature
## High Performance
See also the following **Benchmark**
## Simple
Easy to configure by INI, See also the following **Configuration**
## Secure
It prevents the common web attacks effectively, such as **CSRF**, **XSS** etc.


# Benchmark
## Hardware
Personal laptop.
- **CPU**: Intel(R) Core(TM) i7-4720HQ CPU @ 2.60GHz, **octa-core** CPUs.
- **Hard Disk**: 256G SSD.

## Text Benchmark
```
ab -c 100 -n 100000 -k http://127.0.0.1:8080/
```
RPS
| No. | Requests per second |
| ------------- |:-------------:|
| 1   | 48000 |
| 1   | 48000 |

### View file Benchmark 
```
ab -c 100 -n 100000 -k http://127.0.0.1:8080/index/login
```


# Quick start
## Installation
```
go get github.com/HeadwindFly/cheetah
```

## Usage
### Application Structure
***src/***
　　***app/***
　　　　***config/***
　　　　　　**main.ini** // application's configuration.

　　　　***controllers/***
　　　　　　**index.go** // IndexController

　　　　***views/*** // view files.
　　　　　　***index*** // the IndexController's view files.
  　　　　　　　**index.html**
  
　　　　***resources/***
　　　　　　***css/***
　　　　　　***image/***
　　　　　　***js/***

　　　　***logs/***
　　　　　　**app.log** // the log file.

　　　　**app.go**

### Example
```
package main

import (
	"github.com/HeadwindFly/cheetah"
	"runtime"
	"path"
	"os"
	"headwindfly.com/user/controllers"
)

func main() {
	config := path.Join(os.Getenv("GOPATH"), "src", "app", "config", "main.ini")
	cheetah.Init(config)

	if runtime.NumCPU() > 2 {
		runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	}

	host := cheetah.NewHost("www.headwindfly.com")
	host.RegisterWebController("/index", &controllers.IndexController{})
	
	host.RegisterResources("resources","/path/to/resources")

	cheetah.Run()
}
```
http://127.0.0.1:8080/


# Documentation
See also http://www.headwindfly.com


# Configuration
See also **CONFIGURATION.ini**