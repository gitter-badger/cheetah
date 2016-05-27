# Cheetah
**Cheetah** is a **simple** web framework for **Go**.
 This project aims to become a powerful web development tool and **make developer easily to build a high-performance, secure and stronger web application**.

## Cheetah's 1.0.0a is released.


# Feature
- **High Performance**
See also the following **Benchmark**
- **Simple**
Easy to configure by INI, See also the following **Configuration**
- **Secure**
Cheetah prevents the common web attacks effectively, such as **CSRF**, **XSS** etc.
- **Components**
Cheetah provides logger, session and cache components.


# Benchmark
#### Hardware
Personal laptop.
- **CPU**: Intel(R) Core(TM) i7-4720HQ CPU @ 2.60GHz, **octa-core** CPUs.
- **Hard Disk**: 256G SSD.

#### Text Benchmark
```
ab -c 100 -n 100000 -k http://127.0.0.1:8080/
```
| No.  | Requests per second |
| -----| ------------------- |
| 1    |       40854.46      |
| 2    |       41404.42      |
| 3    |       40264.88      |
| 4    |       40559.56      |
| 5    |       40669.29      |
| 6    |       39863.05      |
| 7    |       40161.26      |
| 8    |       41587.30      |
| 9    |       41873.52      |
| 10   |       40782.22      |

#### View file Benchmark
```
ab -c 100 -n 100000 -k http://127.0.0.1:8080/index/login
```
| No.  | Requests per second |
| -----| ------------------- |
| 1    |       33634.76      |
| 2    |       35143.18      |
| 3    |       34987.50      |
| 4    |       35515.66      |
| 5    |       34852.71      |
| 6    |       35382.51      |
| 7    |       35573.30      |
| 8    |       35500.40      |
| 9    |       35450.94      |
| 10   |       35333.09      |


# Quick start
### Installation
```
go get github.com/HeadwindFly/cheetah
```

### Usage
#### Application Structure
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

#### Example
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
Visit your application on [http://127.0.0.1:8080](http://127.0.0.1:8080/).


# Documentation
See also http://www.headwindfly.com (comming soon...)


# Configuration
See also **CONFIGURATION.ini**