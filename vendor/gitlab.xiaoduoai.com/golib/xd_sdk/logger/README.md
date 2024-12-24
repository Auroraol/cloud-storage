# logger

本库对[logrus](https://github.com/sirupsen/logrus)进行了简单封装，希望达到如下两个目的：
- 节约开发对日志库进行再次封装的成本，同时又能满足当前阶段的所有要求(譬如错误日志需要另写一份)
- 方便其它系统的接入(譬如Trace系统接入以后希望在日志里面输出对应的TraceID，这个可以统一实现)
- 按照新的[日志规范](https://doc.xiaoduoai.com/pages/viewpage.action?pageId=272834643) 格式二， 打印日志。

### 直接使用标准logger
可以不用任何显式的初始化步骤，直接使用标准logger输出：
```golang
package main

import (
	"gitlab.xiaoduoai.com/golib/xd_sdk/logger"
)

func main() {
	logger.Info(ctx, "this is a info message")
}
```
默认输出在stderr：
```shell
time="2019-07-03T14:07:17+08:00" level=info msg="this is a info message" file=main.go line=8
```

### 重新设置标准logger
现实情况是直接使用标准logger很可能不能满足需要，我们可以重新设置标准logger：
```golang
package main

import (
	"gitlab.xiaoduoai.com/golib/xd_sdk/logger"
)

func main() {
	_ = logger.ResetStandard(
		logger.WithLevel("debug"),               // 重置日志等级，默认是info
		logger.WithFile("/tmp/test.log"),        // 重置日志输出，默认是stderr
		logger.WithErrFile("/tmp/test.err.log"), // 重置日志错误复制，默认不复制
	)
	// 如下写法等价
	// _ = logger.ResetStandardWithOptions(logger.Options{
	// 	Level:   "debug",
	// 	File:    "/tmp/test.log",
	// 	ErrFile: "/tmp/test.err.log",
	// })
	logger.Debug(ctx, "this is a debug message")
	logger.Error(ctx, "this is a error message")
}
```
这时候`/tmp/test.log`里面会有：
```shell
time="2019-07-03T14:34:57+08:00" level=debug msg="this is a debug message" file=main.go line=19
time="2019-07-03T14:34:57+08:00" level=error msg="this is a error message" file=main.go line=20
```
而`/tmp/test.err.log`里面会有：
```shell
time="2019-07-03T14:34:57+08:00" level=error msg="this is a error message" file=main.go line=20
```

### 创建自定义logger
有时候你可能希望将某些日志单独打到某个地方，那么直接创建一个自定义logger即可：
```golang
package main

import (
	"gitlab.xiaoduoai.com/golib/xd_sdk/logger"
)

func main() {
	l, _ := logger.NewLogger(
		logger.WithFile("/tmp/test.api.log"),
	)
	// 如下写法等价
	// l, _ := logger.NewLoggerWithOptions(logger.Options{
	// 	File: "/tmp/test.api.log",
	// })
	l.Info(ctx, "this is a info message")
}
```
这时候`/tmp/test.api.log`里面会有：
```shell
time="2019-07-03T14:47:34+08:00" level=info msg="this is a info message" file=main.go line=15
```

### 其它的一些常见用法
你可以在打日志的时候指定一些具体的字段，这样比起自己format更好(便于搜索引擎处理)：
```golang
package main

import (
	"gitlab.xiaoduoai.com/golib/xd_sdk/logger"
)

func main() {
	logger.WithFields(logger.Fields{"foo": 1, "bar": 2}).Info(ctx, "this is a info message")
	logger.WithField("baz", 3).Info(ctx, "this is another info message")
}
```
这时候stderr会输出：
```shell
time="2019-07-03T16:42:56+08:00" level=info msg="this is a info message" bar=2 foo=1 file=main.go line=8
time="2019-07-03T16:42:56+08:00" level=info msg="this is another info message" baz=3 file=main.go line=9
```
还有比如你想接入Trace系统，在日志里面输出对应的TraceID，只需要将上下文传进去即可：
```golang
package main

import (
	"context"

	"gitlab.xiaoduoai.com/golib/xd_sdk/logger"
	"go.opencensus.io/trace"
)

func main() {
	// 这一行代码只是为了得到一个携带Trace信息的上下文，便于演示。
	ctx, _ := trace.StartSpan(context.Background(), "test")
	logger.Info(ctx, "this is a info message")
}
```
这时候stderr会输出：
```shell
time="2019-07-03T16:47:47+08:00" level=info msg="this is a info message" trace=c0a8ef1710f553ce9ed0c97a30422f70 file=main.go line=13
```
