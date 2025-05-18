# 前端

```shell
# 启动frp客户端
cd frp_0.61.1_windows_amd64
# cmd 
frpc -c ./frpc.toml
```

```shell
# 进入项目目录
cd cloud-storage-front-end

# 安装依赖
pnpm i

# 启动服务
pnpm dev
```

# 后端

```shell
# 启动etcd
cd etcd-v3.5.17-windows-amd64
# 双击etcd.exe
```

```shell
# 进入项目目录
cd cloud-storage-back-end

#  RPC 服务启动
1. 上传服务 RPC
2. 日志服务 RPC
3. 用户中心 RPC

# API 服务启动
1. 用户中心 API
2. 上传服务 API
3. 日志服务 API
4. 分享服务 API
```





```
@startuml
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Container.puml

!define DEVICONS https://raw.githubusercontent.com/tupadr3/plantuml-icon-font-sprites/master/devicons
!define FONTAWESOME https://raw.githubusercontent.com/tupadr3/plantuml-icon-font-sprites/master/font-awesome-5
!include DEVICONS/go.puml
!include FONTAWESOME/users.puml
!include FONTAWESOME/server.puml
!include FONTAWESOME/database.puml
!include FONTAWESOME/cloud.puml

LAYOUT_WITH_LEGEND()

title 云存储系统 - 微服务架构详细设计


Person(authenticatedUser, "用户", "拥有账户的系统用户")

System_Boundary(cloudstorage, "云存储系统") {    
    Container(apiGateway, "API网关", "Nginx, Lua", "提供统一接入、负载均衡、安全防护")
    
    System_Boundary(userService, "用户中心服务") {
        Container(userApi, "用户中心API", "Go, Gin", "提供用户相关HTTP接口")
        Container(userRpc, "用户中心RPC", "Go, gRPC", "提供用户相关内部服务调用")
        ContainerDb(userDb, "用户数据库", "MySQL", "存储用户信息、权限配置")
        Container(userCache, "用户缓存", "Redis", "缓存用户会话、权限信息")
        
        Rel(userApi, userRpc, "调用", "gRPC")
        Rel(userRpc, userDb, "读写", "GORM")
        Rel(userRpc, userCache, "读写", "go-redis")
    }
    
    System_Boundary(uploadService, "上传服务") {
        Container(uploadApi, "上传服务API", "Go, Gin", "提供文件上传下载HTTP接口")
        Container(uploadRpc, "上传服务RPC", "Go, gRPC", "提供文件存储内部服务调用")
        ContainerDb(fileMetaDb, "文件元数据库", "MySQL", "存储文件元数据信息")
        Container(uploadCache, "上传缓存", "Redis", "缓存上传任务状态、分片信息")
        Container(fileProcessor, "文件处理器", "Go", "处理文件转码、压缩等任务")
        Container_Ext(objectStorage, "对象存储", "阿里云OSS", "存储文件实体数据")
        
        Rel(uploadApi, uploadRpc, "调用", "gRPC")
        Rel(uploadRpc, fileMetaDb, "读写", "GORM")
        Rel(uploadRpc, uploadCache, "读写", "go-redis")
        Rel(uploadRpc, objectStorage, "存取文件", "OSS SDK")
        Rel(fileProcessor, objectStorage, "处理文件", "OSS SDK")
        Rel(uploadRpc, fileProcessor, "触发处理", "消息队列")
    }
    
    System_Boundary(shareService, "分享服务") {
        Container(shareApi, "分享服务API", "Go, Gin", "提供文件分享HTTP接口")
        Container(shareRpc, "分享服务RPC", "Go, gRPC", "提供分享功能内部服务调用")
        ContainerDb(shareDb, "分享数据库", "MySQL", "存储分享链接、权限信息")
        Container(shareCache, "分享缓存", "Redis", "缓存热门分享信息")
        
        Rel(shareApi, shareRpc, "调用", "gRPC")
        Rel(shareRpc, shareDb, "读写", "GORM")
        Rel(shareRpc, shareCache, "读写", "go-redis")
    }
    
    System_Boundary(logService, "日志服务") {
        Container(logApi, "日志服务API", "Go, Gin", "提供日志查询HTTP接口")
        Container(logRpc, "日志服务RPC", "Go, gRPC", "提供日志记录内部服务调用")
        ContainerDb(logDb, "日志数据库", "MySQL", "存储系统操作日志")
        Container(logCollector, "日志收集器", "Zap, Loki", "收集系统日志")
        
        Rel(logApi, logRpc, "调用", "gRPC")
        Rel(logRpc, logDb, "读写", "GORM")
        Rel(logCollector, logRpc, "发送", "gRPC")
    }
    
    System_Boundary(infrastructureService, "基础设施服务") {
        Container(registryCenter, "服务注册中心", "ETCD", "服务发现与注册")
        Container(configCenter, "配置中心", "ETCD", "集中配置管理")
        Container(messageQueue, "消息队列", "Kafka", "异步消息处理")
        Container(monitorSystem, "监控系统", "Prometheus, Grafana", "系统监控与告警")
    }
    
    Rel(apiGateway, userApi, "转发请求", "HTTP")
    Rel(apiGateway, uploadApi, "转发请求", "HTTP")
    Rel(apiGateway, shareApi, "转发请求", "HTTP")
    Rel(apiGateway, logApi, "转发请求", "HTTP")
    
    
    Rel(uploadApi, userRpc, "验证用户", "gRPC")
    Rel(shareApi, userRpc, "验证用户", "gRPC")
    Rel(shareApi, uploadRpc, "获取文件信息", "gRPC")
    
    Rel_U(userApi, logRpc, "记录日志", "gRPC")
    Rel_U(uploadApi, logRpc, "记录日志", "gRPC")
    Rel_U(shareApi, logRpc, "记录日志", "gRPC")
    
    Rel_D(userRpc, registryCenter, "注册/发现", "gRPC")
    Rel_D(uploadRpc, registryCenter, "注册/发现", "gRPC")
    Rel_D(shareRpc, registryCenter, "注册/发现", "gRPC")
    Rel_D(logRpc, registryCenter, "注册/发现", "gRPC")
    
    Rel_L(userRpc, configCenter, "获取配置", "gRPC")
    Rel_L(uploadRpc, configCenter, "获取配置", "gRPC")
    Rel_L(shareRpc, configCenter, "获取配置", "gRPC")
    Rel_L(logRpc, configCenter, "获取配置", "gRPC")
    
    Rel(uploadRpc, messageQueue, "发送消息", "生产者")
    Rel(fileProcessor, messageQueue, "接收消息", "消费者")
    
    Rel_U(userApi, monitorSystem, "监控指标", "Prometheus")
    Rel_U(uploadApi, monitorSystem, "监控指标", "Prometheus")
    Rel_U(shareApi, monitorSystem, "监控指标", "Prometheus")
    Rel_U(logApi, monitorSystem, "监控指标", "Prometheus")
}


Rel(authenticatedUser, apiGateway, "使用", "HTTPS")


@enduml
```





```
@startuml
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Container.puml

!define DEVICONS https://raw.githubusercontent.com/tupadr3/plantuml-icon-font-sprites/master/devicons
!define FONTAWESOME https://raw.githubusercontent.com/tupadr3/plantuml-icon-font-sprites/master/font-awesome-5
!include DEVICONS/go.puml
!include FONTAWESOME/users.puml
!include FONTAWESOME/server.puml
!include FONTAWESOME/database.puml
!include FONTAWESOME/cloud.puml

LAYOUT_WITH_LEGEND()


Person(authenticatedUser, "用户", "拥有账户的系统用户")

System_Boundary(cloudstorage, "云存储系统") {    
    Container(apiGateway, "API网关", "Nginx, Lua", "提供统一接入、负载均衡、安全防护")
    
    System_Boundary(userService, "用户中心服务") {
        Container(userApi, "用户中心API", "Go, Gin", "提供用户相关HTTP接口")
        Container(userRpc, "用户中心RPC", "Go, gRPC", "提供用户相关内部服务调用")
        ContainerDb(userDb, "用户数据库", "MySQL", "存储用户信息、权限配置")
        Container(userCache, "用户缓存", "Redis", "缓存用户会话、权限信息")
        
        Rel(userApi, userRpc, "调用", "gRPC")
        Rel(userRpc, userDb, "读写", "GORM")
        Rel(userRpc, userCache, "读写", "go-redis")
    }
    
    System_Boundary(uploadService, "上传服务") {
        Container(uploadApi, "上传服务API", "Go, Gin", "提供文件上传下载HTTP接口")
        Container(uploadRpc, "上传服务RPC", "Go, gRPC", "提供文件存储内部服务调用")
        ContainerDb(fileMetaDb, "文件元数据库", "MySQL", "存储文件元数据信息")
        Container(uploadCache, "上传缓存", "Redis", "缓存上传任务状态、分片信息")
        Container(fileProcessor, "文件处理器", "Go", "处理文件转码、压缩等任务")
        Container_Ext(objectStorage, "对象存储", "阿里云OSS", "存储文件实体数据")
        
        Rel(uploadApi, uploadRpc, "调用", "gRPC")
        Rel(uploadRpc, fileMetaDb, "读写", "GORM")
        Rel(uploadRpc, uploadCache, "读写", "go-redis")
        Rel(uploadRpc, objectStorage, "存取文件", "OSS SDK")
        Rel(fileProcessor, objectStorage, "处理文件", "OSS SDK")
        Rel(uploadRpc, fileProcessor, "触发处理", "消息队列")
    }
    
    System_Boundary(shareService, "分享服务") {
        Container(shareApi, "分享服务API", "Go, Gin", "提供文件分享HTTP接口")
        Container(shareRpc, "分享服务RPC", "Go, gRPC", "提供分享功能内部服务调用")
        ContainerDb(shareDb, "分享数据库", "MySQL", "存储分享链接、权限信息")
        Container(shareCache, "分享缓存", "Redis", "缓存热门分享信息")
        
        Rel(shareApi, shareRpc, "调用", "gRPC")
        Rel(shareRpc, shareDb, "读写", "GORM")
        Rel(shareRpc, shareCache, "读写", "go-redis")
    }
    
    System_Boundary(logService, "日志服务") {
        Container(logApi, "日志服务API", "Go, Gin", "提供日志查询HTTP接口")
        Container(logRpc, "日志服务RPC", "Go, gRPC", "提供日志记录内部服务调用")
        ContainerDb(logDb, "日志数据库", "MySQL", "存储系统操作日志")
        Container(logCollector, "日志收集器", "Zap, Loki", "收集系统日志")
        
        Rel(logApi, logRpc, "调用", "gRPC")
        Rel(logRpc, logDb, "读写", "GORM")
        Rel(logCollector, logRpc, "发送", "gRPC")
    }
    
    System_Boundary(infrastructureService, "基础设施服务") {
        Container(registryCenter, "服务注册中心", "ETCD", "服务发现与注册")
        Container(configCenter, "配置中心", "ETCD", "集中配置管理")
        Container(messageQueue, "消息队列", "Kafka", "异步消息处理")
        Container(monitorSystem, "监控系统", "Prometheus, Grafana", "系统监控与告警")
    }
    
    Rel(apiGateway, userApi, "转发请求", "HTTP")
    Rel(apiGateway, uploadApi, "转发请求", "HTTP")
    Rel(apiGateway, shareApi, "转发请求", "HTTP")
    Rel(apiGateway, logApi, "转发请求", "HTTP")
    
    
    Rel(uploadApi, userRpc, "验证用户", "gRPC")
    Rel(shareApi, userRpc, "验证用户", "gRPC")
    Rel(shareApi, uploadRpc, "获取文件信息", "gRPC")
    
    Rel_U(userApi, logRpc, "记录日志", "gRPC")
    Rel_U(uploadApi, logRpc, "记录日志", "gRPC")
    Rel_U(shareApi, logRpc, "记录日志", "gRPC")
    
    Rel_D(userRpc, registryCenter, "注册/发现", "gRPC")
    Rel_D(uploadRpc, registryCenter, "注册/发现", "gRPC")
    Rel_D(shareRpc, registryCenter, "注册/发现", "gRPC")
    Rel_D(logRpc, registryCenter, "注册/发现", "gRPC")
    
    Rel_L(userRpc, configCenter, "获取配置", "gRPC")
    Rel_L(uploadRpc, configCenter, "获取配置", "gRPC")
    Rel_L(shareRpc, configCenter, "获取配置", "gRPC")
    Rel_L(logRpc, configCenter, "获取配置", "gRPC")
    
    Rel(uploadRpc, messageQueue, "发送消息", "生产者")
    Rel(fileProcessor, messageQueue, "接收消息", "消费者")
    
    Rel_U(userApi, monitorSystem, "监控指标", "Prometheus")
    Rel_U(uploadApi, monitorSystem, "监控指标", "Prometheus")
    Rel_U(shareApi, monitorSystem, "监控指标", "Prometheus")
    Rel_U(logApi, monitorSystem, "监控指标", "Prometheus")
}


Rel(authenticatedUser, apiGateway, "使用", "HTTPS")
@enduml
```





```
@startuml
skinparam backgroundColor white
skinparam monochrome true
skinparam shadowing false
skinparam defaultFontName Arial
skinparam defaultFontSize 12
skinparam linetype ortho
skinparam packageStyle rectangle
skinparam componentStyle uml2
skinparam monochrome true
skinparam handDrawn true

title 基于Golang的云存储与文件审计系统后端架构设计

rectangle "API网关层" {
  [Nginx网关]
}

together {
  package "应用服务层" {
    component "用户中心服务" {
      [用户认证模块]
      [用户管理模块]
    }
    
    component "上传服务" {
      [文件上传模块]
      [文件处理模块]
    }
    
    component "分享服务" {
      [分享管理模块]
      [访问控制模块]
    }
    
    component "日志服务" {
      [操作日志模块]
      [审计模块]
    }
  }
}

rectangle "数据存储层" {
  database "关系数据库" {
    [用户数据]
    [文件元数据]
    [分享数据]
    [日志数据]
  }
  
  database "缓存服务" {
    [Redis]
  }
  
  cloud "对象存储" {
    [阿里云OSS]
  }
}

rectangle "基础架构层" {
  [服务注册中心]
  [配置中心]
  [消息队列]
  [监控系统]
}

请求 --> [Nginx网关]

[Nginx网关] --> [用户认证模块]
[Nginx网关] --> [文件上传模块]
[Nginx网关] --> [分享管理模块]
[Nginx网关] --> [操作日志模块]

[用户认证模块] --> [用户数据]
[文件上传模块] --> [文件元数据]
[分享管理模块] --> [分享数据]
[操作日志模块] --> [日志数据]

[文件上传模块] --> [阿里云OSS]

[用户认证模块] --> [Redis]
[文件上传模块] --> [Redis]
[分享管理模块] --> [Redis]

[用户认证模块] --> [服务注册中心]
[文件上传模块] --> [服务注册中心]
[分享管理模块] --> [服务注册中心]
[操作日志模块] --> [服务注册中心]

[文件上传模块] --> [消息队列]
@enduml
```





```
@startuml
skinparam backgroundColor white
skinparam monochrome true
skinparam shadowing false
skinparam defaultFontName Arial
skinparam defaultFontSize 12
skinparam linetype ortho
skinparam packageStyle rectangle
skinparam componentStyle uml2

title 基于Golang的云存储与文件审计系统前端架构设计

' 用户
actor 用户

rectangle "前端应用" {
  [Vue3应用]
}

package "展示层" {
  component "页面视图" {
    [登录/认证视图]
    [文件管理视图]
    [分享管理视图]
    [个人中心视图]
    [日志与审计视图]
    [回收站视图]
  }
  

}

package "业务逻辑层" {
  component "状态管理" {
    [Pinia存储]
    [用户状态]
    [文件状态]
    [权限状态]
    [应用设置]
  }
 
}

package "数据访问层" {
  component "API服务" {
    [用户API模块]
    [文件API模块]
    [分享API模块]
    [日志API模块]
    [审计API模块]
  }
}

' 连接关系
用户 --> [Vue3应用]

[Vue3应用] --> [登录/认证视图]
[Vue3应用] --> [文件管理视图]
[Vue3应用] --> [分享管理视图]
[Vue3应用] --> [个人中心视图]
[Vue3应用] --> [日志与审计视图]

[登录/认证视图] --> [用户API模块]
[文件管理视图] --> [文件API模块]
[分享管理视图] --> [分享API模块]
[日志与审计视图] --> [日志API模块]
[日志与审计视图] --> [审计API模块]

[用户API模块] ..> [请求拦截器]
[文件API模块] ..> [请求拦截器]
[分享API模块] ..> [请求拦截器]
[日志API模块] ..> [请求拦截器]


[页面视图] --> [Pinia存储]
[Pinia存储] --> [用户状态]
[Pinia存储] --> [文件状态]
[Pinia存储] --> [权限状态]
@enduml
```



```
@startuml
skinparam backgroundColor white
skinparam monochrome true
skinparam shadowing false
skinparam defaultFontName SimSun
skinparam defaultFontSize 12
skinparam linetype ortho
skinparam rectangleBorderThickness 1
skinparam arrowColor black
skinparam nodesep 40
skinparam ranksep 30

' 设置文字竖向显示
skinparam rectangleTextDirection vertical
skinparam classAttributeTextDirection vertical

' 顶层节点
rectangle "云\n存\n储\n系\n统" as main

' 中间层节点
rectangle "用\n户\n模\n块" as user
rectangle "文\n件\n模\n块" as file
rectangle "分\n享\n模\n块" as share
rectangle "日\n志\n模\n块" as log
rectangle "审\n计\n模\n块" as audit

' 底层节点 - 用户模块
rectangle "登\n录\n认\n证" as login
rectangle "注\n册" as register
rectangle "权\n限\n管\n理" as permission
rectangle "修\n改\n信\n息" as userInfo
rectangle "账\n号\n管\n理" as account

' 底层节点 - 文件模块
rectangle "上\n传\n下\n载" as transfer
rectangle "分\n类\n存\n储" as category
rectangle "文\n件\n预\n览" as preview
rectangle "文\n件\n搜\n索" as search
rectangle "在\n线\n编\n辑" as edit

' 底层节点 - 分享模块
rectangle "权\n限\n设\n置" as sharePermission
rectangle "分\n享\n管\n理" as sharing
rectangle "协\n作\n记\n录" as collabRecord

' 底层节点 - 日志模块
rectangle "操\n作\n日\n志" as opLog
rectangle "系\n统\n日\n志" as sysLog
rectangle "登\n录\n日\n志" as loginLog

' 底层节点 - 审计模块
rectangle "操\n作\n审\n计" as opAudit
rectangle "统\n计\n分\n析" as statistics
rectangle "数\n据\n导\n出" as export

' 主连接
main -down-> user
main -down-> file
main -down-> share
main -down-> log
main -down-> audit

' 用户模块连接
user -down-> login
user -down-> register
user -down-> permission
user -down-> userInfo
user -down-> account

' 文件模块连接
file -down-> transfer
file -down-> category
file -down-> preview
file -down-> search
file -down-> edit

' 分享模块连接
share -down-> sharePermission
share -down-> sharing
share -down-> collabRecord

' 日志模块连接
log -down-> opLog
log -down-> sysLog
log -down-> loginLog

' 审计模块连接
audit -down-> opAudit
audit -down-> statistics
audit -down-> export

' 布局调整 - 强制水平排列
login -[hidden]right-> register
register -[hidden]right-> permission
permission -[hidden]right-> userInfo
userInfo -[hidden]right-> account

transfer -[hidden]right-> category
category -[hidden]right-> preview
preview -[hidden]right-> search
search -[hidden]right-> edit

sharePermission -[hidden]right-> sharing
sharing -[hidden]right-> collabRecord

opLog -[hidden]right-> sysLog
sysLog -[hidden]right-> loginLog

opAudit -[hidden]right-> statistics
statistics -[hidden]right-> export

' 强制二级模块水平排列
user -[hidden]right-> file
file -[hidden]right-> share
share -[hidden]right-> log
log -[hidden]right-> audit

hide empty members
hide circle
hide methods
@enduml
```





```
@startuml
skinparam monochrome true
skinparam ClassBackgroundColor white

class 用户表 {
    id
    version
    username
    password
    mobile
    nickname
    gender
    avatar
    birthday
    email
    brief
    info
    del_state
    delete_time
    status
    admin
    now_volume
    total_volume
    create_time
    update_time
}

note left of 用户表::id
  id: 用户ID
  version: 版本号
  username: 用户名
  password: 密码
  mobile: 手机号
  nickname: 昵称
  gender: 性别
  avatar: 用户头像
  birthday: 生日
  email: 电子邮箱
endnote

note right of 用户表::brief
  brief: 简介/个性签名
  info: 附加信息
  del_state: 删除状态
  delete_time: 删除时间
  status: 状态
  admin: 是否管理员
  now_volume: 当前存储容量
  total_volume: 最大存储容量
  create_time: 创建时间
  update_time: 更新时间
endnote
@enduml

@startuml
skinparam monochrome true
skinparam ClassBackgroundColor white

class 用户存储库表 {
    id
    user_id
    parent_id
    repository_id
    name
    status
    create_time
    update_time
}

note left of 用户存储库表::id
  id: 主键ID
  user_id: 用户ID
  parent_id: 父文件夹ID
  repository_id: 存储项ID
endnote

note right of 用户存储库表::name
  name: 名称
  status: 文件状态
  create_time: 创建时间
  update_time: 更新时间
endnote

@enduml

@startuml
skinparam monochrome true
skinparam ClassBackgroundColor white

class 上传历史表 {
    id
    user_id
    file_name
    size
    repository_id
    status
    create_time
    update_time
}

note left of 上传历史表::id
  id: 主键ID
  user_id: 用户ID
  file_name: 文件名
  size: 文件大小
endnote

note right of 上传历史表::repository_id
  repository_id: 文件ID
  status: 上传状态
  create_time: 创建时间
  update_time: 更新时间
endnote

@enduml

@startuml
skinparam monochrome true
skinparam ClassBackgroundColor white

class 存储池表 {
    id
    identity
    oss_key
    hash
    ext
    size
    path
    name
    create_time
    update_time
}

note left of 存储池表::id
  id: 主键ID
  identity: 文件ID
  oss_key: OSS键
  hash: 文件哈希值
  ext: 文件扩展名
endnote

note right of 存储池表::size
  size: 文件大小
  path: 文件路径
  name: 文件名
  create_time: 创建时间
  update_time: 更新时间
endnote

@enduml

@startuml
skinparam monochrome true
skinparam ClassBackgroundColor white

class 分享表 {
    id
    user_id
    repository_id
    user_repository_id
    expired_time
    click_num
    create_time
    code
    update_time
    deleted_at
}

note left of 分享表::id
  id: 主键ID
  user_id: 用户ID
  repository_id: 公共池中的唯一标识
  user_repository_id: 用户池中的唯一标识
  expired_time: 失效时间
endnote

note right of 分享表::click_num
  click_num: 点击次数
  create_time: 创建时间
  code: 提取码
  update_time: 更新时间
  deleted_at: 删除时间
endnote

@enduml

@startuml
skinparam monochrome true
skinparam ClassBackgroundColor white

class 操作日志表 {
    id
    user_id
    content
    flag
    repository_id
    file_name
    create_time
}

note left of 操作日志表::id
  id: 操作记录ID
  user_id: 用户ID
  content: 操作内容
  flag: 操作类型
endnote

note right of 操作日志表::repository_id
  repository_id: 文件ID
  file_name: 文件名
  create_time: 创建时间戳
endnote

@enduml

@startuml
skinparam monochrome true
skinparam ClassBackgroundColor white

class 日志文件表 {
    id
    user_id
    name
    host
    path
    create_time
    comment
    monitor_choice
}

note left of 日志文件表::id
  id: 日志文件ID
  user_id: 用户ID
  name: 日志文件名
  host: 主机信息
endnote

note right of 日志文件表::path
  path: 日志文件路径
  create_time: 创建时间
  comment: 备注
  monitor_choice: 监控选择
endnote

@enduml

@startuml
skinparam monochrome true
skinparam ClassBackgroundColor white

class SSH信息表 {
    id
    user_id
    host
    port
    username
    password
    created_at
    updated_at
}

note left of SSH信息表::id
  id: 主键ID
  user_id: 用户ID
  host: 主机地址
  port: 端口号
endnote

note right of SSH信息表::username
  username: 用户名
  password: 密码
  created_at: 创建时间
  updated_at: 更新时间
endnote

@enduml
```





```
请按照下述sql修改上述 CREATE TABLE IF NOT EXISTS  `user` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '用户id',
    `version` BIGINT NOT NULL,
    `username` VARCHAR(25) DEFAULT NULL COMMENT '用户名',
    `password` VARCHAR(255) DEFAULT NULL COMMENT '密码',
    `mobile` BIGINT(11) DEFAULT NULL COMMENT '手机号',
    `nickname` VARCHAR(50) NOT NULL COMMENT '昵称',
    `gender` TINYINT(1) NOT NULL DEFAULT '1' COMMENT '性别，1：男，0：女，默认为1',
    `avatar` VARCHAR(255) DEFAULT NULL COMMENT '用户头像',
    `birthday` DATE DEFAULT NULL COMMENT '生日',
    `email` VARCHAR(254) DEFAULT NULL COMMENT '电子邮箱',
    `brief` VARCHAR(255) DEFAULT NULL COMMENT '简介|个性签名',
    `info` TEXT,
    `del_state` INt COMMENT '删除状态，0: 未删除，1：已删除',
    `delete_time` TIMESTAMP COMMENT '删除时间',
    `status` TINYINT(1) NOT NULL DEFAULT '0' COMMENT '状态，0：正常，1：锁定，2：禁用，3：过期',
    `admin` TINYINT(1) NOT NULL DEFAULT '0' COMMENT '是否管理员，1：是，0：否',
    `now_volume` INT(11) NOT NULL DEFAULT '0' COMMENT '当前存储容量',
    `total_volume` INT(11) NOT NULL DEFAULT '1000000000' COMMENT '最大存储容量',
    `create_time` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_name_unique` (`username`),
    UNIQUE KEY `idx_mobile_unique` (`mobile`)
);  CREATE TABLE `user_repository`
(
    `id`                  bigint unsigned NOT NULL AUTO_INCREMENT,
    `user_id`       bigint unsigned NOT NULL DEFAULT '0',
    `parent_id`           bigint unsigned NOT NULL DEFAULT '0' ,
    `repository_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '0则为文件夹, 其他为文件id',
    `name`                varchar(255) NOT NULL DEFAULT '' COMMENT '文件夹名称',
    `status` int(11) NOT NULL DEFAULT '0' COMMENT '文件状态(0正常1已删除2禁用)',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_repository_id` (`repository_id`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 8
  DEFAULT CHARSET = utf8; CREATE TABLE IF NOT EXISTS `upload_history`
(
    `id`         bigint unsigned NOT NULL AUTO_INCREMENT,
    `user_id`       bigint unsigned NOT NULL DEFAULT '0' comment '用户id',
    `file_name` varchar(255) NOT NULL DEFAULT '' COMMENT '文件名',
    `size` int(11) NOT NULL DEFAULT '0' COMMENT '文件大小',
    `repository_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '文件id',
    `status`       tinyint(1) NOT NULL DEFAULT '0' COMMENT '上传状态，0：上传中，1：上传成功，2：上传失败',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_repository_id_unique` (`repository_id`)
    ) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8; CREATE TABLE IF NOT EXISTS `repository_pool` (
    `id`          bigint unsigned NOT NULL AUTO_INCREMENT,
    `identity`    bigint unsigned NOT NULL DEFAULT '0' COMMENT '文件id',
    `oss_key`     varchar(255)   NOT NULL DEFAULT '' COMMENT '文件在OSS中的键',
    `hash`        varchar(32)    NOT NULL DEFAULT '' COMMENT '文件的唯一标识',
    `ext`         varchar(30)    NOT NULL DEFAULT '' COMMENT '文件扩展名',
    `size`        int(11)        NOT NULL DEFAULT '0' COMMENT '文件大小',
    `path`        varchar(255)   NOT NULL DEFAULT '' COMMENT '文件url路径',
    `name`        varchar(255)   NOT NULL DEFAULT '',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `idx_identity` (`identity`),        -- 新增的普通索引
    UNIQUE KEY `idx_hash_unique` (`hash`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8; CREATE TABLE IF NOT EXISTS `share_basic` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL DEFAULT '0',
  `repository_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '公共池中的唯一标识',
  `user_repository_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '用户池子中的唯一标识',
  `expired_time` int(11) NOT NULL DEFAULT '0' COMMENT '失效时间，单位秒, 【0-永不失效】',
  `click_num` int(11) NOT NULL DEFAULT '0' COMMENT '点击次数',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `code` VARCHAR(30) NOT NULL DEFAULT '' COMMENT '提取码',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_repository_id` (`repository_id`),
  KEY `idx_user_repository_id` (`user_repository_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8; CREATE TABLE IF NOT EXISTS `audit`(
    `id` INT UNSIGNED AUTO_INCREMENT COMMENT '操作记录ID',
    `user_id` bigint unsigned NOT NULL DEFAULT '0' comment '用户id',
    `content` TEXT NOT NULL COMMENT '操作内容',
    `flag` tinyint(1) NOT NULL default 0 comment '操作类型，0：上传，1：下载，2：删除，3.恢复 4：重命名，5：移动，6：复制，7：创建文件夹，8：修改文件',
    `repository_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '0则为文件夹, 其他为文件id',
    `file_name`    varchar(255) NOT NULL DEFAULT '' COMMENT '文件夹名称',
    `create_time` int(11) NOT NULL default 0 comment '创建时间戳',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 comment '操作记录表'; CREATE TABLE IF NOT EXISTS `logfile`(
    `id` INT UNSIGNED AUTO_INCREMENT COMMENT '日志文件ID',
    `user_id` INT NOT NULL,
    `name` VARCHAR(200) NOT NULL COMMENT '日志文件名',
    `host` LONGTEXT NOT NULL COMMENT '主机信息',
    `path` VARCHAR(1024) NOT NULL COMMENT '日志文件路径',
    `create_time` DATETIME COMMENT '创建时间',
    `comment` VARCHAR(200) COMMENT '备注',
    `monitor_choice` INT COMMENT '监控选择',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8; CREATE TABLE `ssh_info` (
                            `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
                            `user_id` bigint(20) NOT NULL COMMENT '用户ID',
                            `host` varchar(255) NOT NULL COMMENT '主机地址',
                            `port` int(11) NOT NULL COMMENT '端口号',
                            `username` varchar(255) NOT NULL COMMENT '用户名',
                            `password` varchar(255) NOT NULL COMMENT '密码',
                            `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                            `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                            PRIMARY KEY (`id`),
                            UNIQUE KEY `idx_user_host` (`user_id`, `host`) -- 新增的唯一索引
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='SSH信息表';
```








(1) 用户表
用户表是系统识别和管理用户的核心，存储着丰富的用户信息。"id" 作为主键，采用 bigint 类型自增长，独一无二地标识每个用户。"version" 字段用于记录用户信息的版本号，便于数据同步和冲突处理。"username" 字段是用户登录系统的关键凭证，"password" 字段经加密存储，保障用户登录安全。"mobile" 字段作为重要联系方式，在安全验证等场景发挥作用。"nickname"、"gender"、"avatar"、"birthday" 和 "email" 等字段满足用户个性化信息展示需求。"brief" 存储用户简介或个性签名，"info" 存储附加信息。"del_state" 和 "delete_time" 用于软删除功能实现，"status" 标识用户状态，"admin" 标识是否为管理员。"now_volume" 和 "total_volume" 字段分别记录已用和最大存储空间，用于监控和管理用户存储资源。"create_time" 和 "update_time" 字段记录创建和更新时间，用于追溯和跟踪用户信息变化。
系统对用户表的设计特别注重安全性和可扩展性。通过对用户名和手机号建立唯一索引，确保用户身份的唯一性；通过版本控制机制，有效处理并发更新冲突；通过软删除机制，保障用户数据的安全性；通过存储空间管理，实现精细化的资源控制。这些设计使得用户管理模块能够灵活应对各种业务需求，同时保持数据的一致性和完整性。
表3-1 用户表
| 字段名称 | 类型 | 长度 | 允许空 | 默认值 | 备注 |
|----------|------|------|--------|--------|------|
| id | bigint | 20 | 否 | 自增 | 用户ID，主键 |
| version | int | 11 | 否 | 无 | 版本号，用于乐观锁 |
| username | varchar | 25 | 是 | null | 用户名，唯一索引 |
| password | varchar | 255 | 是 | null | 密码，加密存储 |
| mobile | varchar | 11 | 是 | null | 手机号，唯一索引 |
| nickname | varchar | 50 | 否 | 无 | 昵称 |
| gender | tinyint | 1 | 否 | 1 | 性别，1：男，0：女 |
| avatar | varchar | 255 | 是 | null | 用户头像URL |
| birthday | date | - | 是 | null | 生日 |
| email | varchar | 254 | 是 | null | 电子邮箱 |
| brief | varchar | 255 | 是 | null | 简介/个性签名 |
| info | text | - | 是 | null | 附加信息 |
| del_state | tinyint | 1 | 是 | 0 | 删除状态，0：未删除，1：已删除 |
| delete_time | timestamp | - | 是 | null | 删除时间 |
| status | tinyint | 1 | 否 | 0 | 状态，0：正常，1：锁定，2：禁用，3：过期 |
| admin | tinyint | 1 | 否 | 0 | 是否管理员，1：是，0：否 |
| now_volume | int | 11 | 否 | 0 | 当前存储容量，单位字节 |
| total_volume | int | 11 | 否 | 1000000000 | 最大存储容量，默认1GB |
| create_time | timestamp | - | 是 | CURRENT_TIMESTAMP | 创建时间 |
| update_time | timestamp | - | 是 | CURRENT_TIMESTAMP | 更新时间，自动更新 |
索引设计：
主键索引：id
唯一索引：idx_name_unique(username)
唯一索引：idx_mobile_unique(mobile)
(2) 用户存储库表
用户存储库表用于管理用户的文件和文件夹结构，是实现文件组织和管理的核心表。该表采用树形结构设计，通过 "parent_id" 字段构建文件夹的层级关系，使得系统能够支持无限层级的文件夹嵌套，满足用户对文件组织的多样化需求。
"id" 为主键，唯一确定每个存储项。"user_id" 字段关联用户表，明确存储项所有者，是实现多用户隔离的关键。"parent_id" 字段用于确定文件所属文件夹，构建文件目录结构，值为0表示位于根目录。"repository_id" 字段关联存储池表，指向实际存储的文件，值为0表示该记录是文件夹而非文件。"name" 字段是用户识别文件或文件夹的重要标识，支持中文和特殊字符。"status" 字段标识文件状态，支持正常、已删除和禁用三种状态，实现文件的软删除和访问控制。"create_time" 和 "update_time" 字段分别记录创建和更新时间，便于文件版本管理和历史追踪。
该表的设计特别注重性能和可扩展性。通过对 "user_id" 和 "repository_id" 建立索引，大幅提高文件列表查询和文件定位的效率；通过状态字段的设计，实现文件回收站功能，防止误删除导致的数据丢失；通过树形结构的设计，支持灵活的文件组织方式，满足不同场景下的使用需求。
表3-2 用户存储库表
| 字段名称 | 类型 | 长度 | 允许空 | 默认值 | 备注 |
|----------|------|------|--------|--------|------|
| id | bigint | 20 | 否 | 自增 | 主键ID |
| user_id | bigint | 20 | 否 | 0 | 用户ID，关联用户表 |
| parent_id | bigint | 20 | 否 | 0 | 父文件夹ID，0表示根目录 |
| repository_id | bigint | 20 | 否 | 0 | 存储项ID，0表示文件夹 |
| name | varchar | 255 | 否 | '' | 文件或文件夹名称 |
| status | int | 11 | 否 | 0 | 文件状态(0正常1已删除2禁用) |
| create_time | timestamp | - | 是 | CURRENT_TIMESTAMP | 创建时间 |
| update_time | timestamp | - | 是 | CURRENT_TIMESTAMP | 更新时间，自动更新 |
索引设计：
主键索引：id
普通索引：idx_user_id(user_id)
普通索引：idx_repository_id(repository_id)
业务规则：
同一用户同一目录下不允许存在同名文件或文件夹
删除文件夹时需级联处理其下的所有文件和子文件夹
移动文件时需检查目标位置是否存在同名文件
文件夹层级不应超过系统设定的最大值(通常为10层)
(3) 存储池表
存储池表用于管理系统中的实际文件存储，是实现文件去重和高效存储的关键。该表采用内容寻址存储的设计理念，通过文件哈希值唯一标识每个文件，实现基于内容的文件去重，大幅节省存储空间。
"id" 作为主键，唯一标识每个存储项。"identity" 字段是文件的唯一标识，用于系统内部引用。"oss_key" 字段存储对象存储服务中的键值，是访问云存储中实际文件的关键。"hash" 字段存储文件哈希值，通常采用MD5或SHA256算法，用于实现秒传和文件完整性校验。"ext" 字段记录文件扩展名，便于文件类型识别。"size" 字段记录文件大小，用于存储空间统计和限制检查。"path" 字段存储文件实际存储路径，支持多种存储策略。"name" 字段记录原始文件名，保留上传时的文件名信息。"create_time" 和 "update_time" 字段记录创建和更新时间，便于存储管理和清理策略实施。
该表的设计特别注重存储效率和访问性能。通过对文件哈希值建立唯一索引，实现基于内容的文件去重，同一内容的文件只会存储一份；通过对文件ID建立索引，提高文件检索效率；通过存储文件扩展名和大小信息，支持文件类型过滤和大小统计；通过记录原始文件名，在文件下载时保持文件名的一致性。这些设计使得系统能够高效管理海量文件，同时保证用户体验。
表3-3 存储池表
| 字段名称 | 类型 | 长度 | 允许空 | 默认值 | 备注 |
|----------|------|------|--------|--------|------|
| id | bigint | 20 | 否 | 自增 | 主键ID |
| identity | bigint | 20 | 否 | 0 | 文件ID，系统内部标识 |
| oss_key | varchar | 255 | 否 | '' | 文件在OSS中的键，访问路径 |
| hash | varchar | 32 | 否 | '' | 文件哈希值，用于去重和秒传 |
| ext | varchar | 30 | 否 | '' | 文件扩展名，不含点号 |
| size | int | 11 | 否 | 0 | 文件大小，单位字节 |
| path | varchar | 255 | 否 | '' | 文件URL路径，完整访问地址 |
| name | varchar | 255 | 否 | '' | 原始文件名，保留上传时的名称 |
| create_time | timestamp | - | 是 | CURRENT_TIMESTAMP | 创建时间 |
| update_time | timestamp | - | 是 | CURRENT_TIMESTAMP | 更新时间，自动更新 |
索引设计：
主键索引：id
普通索引：idx_identity(identity)
唯一索引：idx_hash_unique(hash)
业务规则：
文件上传前先计算哈希值，检查是否已存在相同内容的文件
文件实际删除需检查引用计数，确保没有用户仍在使用该文件
定期清理无引用的文件，回收存储空间
大文件采用分片上传，上传完成后合并处理
(4) 分享表
分享表记录文件和文件夹的分享信息，是实现文件共享功能的核心。该表设计支持灵活的分享策略，包括密码保护、有效期限制和访问统计等功能，满足用户对文件分享的多样化需求。
"id" 作为主键，唯一标识每次分享记录。"user_id" 字段关联用户表，确定分享发起者，便于分享管理和权限控制。"repository_id" 字段关联存储池中的唯一标识，"user_repository_id" 字段关联用户存储库中的唯一标识，这两个字段共同确定分享的具体内容，支持文件和文件夹的分享。"expired_time" 字段设置分享过期时间，值为0表示永不过期，增强分享的时效性控制。"click_num" 字段记录分享链接的点击次数，便于分享效果统计和分析。"code" 字段存储分享提取码，增强分享的安全性，防止未授权访问。"create_time" 和 "update_time" 字段记录分享创建和更新时间，"deleted_at" 字段用于软删除功能实现，支持分享记录的管理和恢复。
该表的设计特别注重安全性和灵活性。通过提取码机制，确保只有知晓密码的用户才能访问分享内容；通过有效期设置，控制分享的时间范围，过期自动失效；通过点击统计，帮助用户了解分享的传播情况；通过软删除机制，防止误操作导致分享记录丢失。这些设计使得文件分享功能既安全可控，又灵活便捷。
表3-4 分享表
| 字段名称 | 类型 | 长度 | 允许空 | 默认值 | 备注 |
|----------|------|------|--------|--------|------|
| id | bigint | 20 | 否 | 自增 | 主键ID |
| user_id | bigint | 20 | 否 | 0 | 用户ID，分享创建者 |
| repository_id | bigint | 20 | 否 | 0 | 公共池中的唯一标识 |
| user_repository_id | bigint | 20 | 否 | 0 | 用户池中的唯一标识 |
| expired_time | int | 11 | 否 | 0 | 失效时间，0表示永不过期 |
| click_num | int | 11 | 否 | 0 | 点击次数，访问统计 |
| create_time | timestamp | - | 是 | CURRENT_TIMESTAMP | 创建时间 |
| code | varchar | 30 | 否 | '' | 提取码，为空表示无需密码 |
| update_time | timestamp | - | 是 | CURRENT_TIMESTAMP | 更新时间，自动更新 |
| deleted_at | datetime | - | 是 | null | 删除时间，软删除标记 |
索引设计：
主键索引：id
普通索引：idx_user_id(user_id)
普通索引：idx_repository_id(repository_id)
普通索引：idx_user_repository_id(user_repository_id)
业务规则：
分享链接的有效性检查包括：是否过期、是否被删除、源文件是否存在
提取码验证成功后才能访问分享内容
访问分享链接时自动增加点击计数
支持批量取消分享和批量延长有效期
(5) 上传历史表
上传历史表用于跟踪用户文件上传行为，是实现上传管理和断点续传的重要依据。该表详细记录了每次上传操作的过程和结果，便于系统监控上传状态和用户查询上传历史。
"id" 作为主键，唯一标识每次上传记录。"user_id" 字段关联用户表，确定上传执行者，便于用户区分自己的上传记录。"file_name" 字段记录上传文件名，保留原始文件名信息。"size" 字段记录文件大小，用于上传进度计算和存储空间检查。"repository_id" 字段关联文件ID，指向上传成功后的存储记录。"status" 字段标识上传状态，支持上传中、成功和失败三种状态，便于系统和用户了解上传进度。"create_time" 和 "update_time" 字段分别记录上传开始和状态更新时间，用于计算上传耗时和判断上传是否超时。
该表的设计特别注重实用性和可追溯性。通过对上传状态的记录，系统能够识别中断的上传任务，支持断点续传功能；通过记录上传时间和文件大小，系统能够分析上传速度和性能瓶颈；通过对repository_id建立唯一索引，确保一个文件只有一条有效的上传记录，防止重复上传。这些设计使得上传功能更加可靠和高效。
表3-5 上传历史表
| 字段名称 | 类型 | 长度 | 允许空 | 默认值 | 备注 |
|----------|------|------|--------|--------|------|
| id | bigint | 20 | 否 | 自增 | 主键ID |
| user_id | bigint | 20 | 否 | 0 | 用户ID，上传者 |
| file_name | varchar | 255 | 否 | '' | 文件名，原始名称 |
| size | int | 11 | 否 | 0 | 文件大小，单位字节 |
| repository_id | bigint | 20 | 否 | 0 | 文件ID，关联存储池 |
| status | tinyint | 1 | 否 | 0 | 上传状态，0：上传中，1：成功，2：失败 |
| create_time | timestamp | - | 是 | CURRENT_TIMESTAMP | 创建时间，上传开始时间 |
| update_time | timestamp | - | 是 | CURRENT_TIMESTAMP | 更新时间，状态变更时间 |
索引设计：
主键索引：id
唯一索引：idx_repository_id_unique(repository_id)
业务规则：
上传开始时创建状态为"上传中"的记录
上传完成后更新状态为"成功"或"失败"
定期清理长时间处于"上传中"状态的记录，释放资源
支持查询特定时间段内的上传记录，便于统计分析
(6) 操作日志表
操作日志表是系统进行安全审计和行为分析的重要依据，记录了用户在系统中的各种操作行为，为系统安全和用户行为分析提供数据支持。该表设计注重全面性和可查询性，确保系统中的重要操作都有迹可循。
"id" 作为主键，唯一标识每条操作记录。"user_id" 字段关联用户表，确定操作执行者，便于追踪用户行为。"content" 字段记录操作内容详情，采用文本类型存储详细的操作描述。"flag" 字段标识操作类型，如上传、下载、删除等，采用数值编码便于统计和筛选。"repository_id" 字段关联文件ID，指明操作涉及的具体文件。"file_name" 字段记录操作涉及的文件名，保留操作时的文件名信息。"create_time" 字段记录操作发生时间，便于时间序列分析和事件追溯。
该表的设计特别注重审计需求和性能平衡。通过对操作类型的编码，支持高效的操作类型筛选和统计；通过记录详细的操作内容，满足安全审计的信息需求；通过关联文件信息，便于追踪文件生命周期中的各种操作。同时，该表不设置更新操作，一旦记录创建就不再修改，确保审计日志的不可篡改性，增强系统的安全性和可靠性。
表3-6 操作日志表
| 字段名称 | 类型 | 长度 | 允许空 | 默认值 | 备注 |
|----------|------|------|--------|--------|------|
| id | int | 10 | 否 | 自增 | 操作记录ID |
| user_id | bigint | 20 | 否 | 0 | 用户ID，操作执行者 |
| content | text | - | 否 | 无 | 操作内容，详细描述 |
| flag | tinyint | 1 | 否 | 0 | 操作类型(0:上传,1:下载,2:删除,3:恢复,4:重命名,5:移动,6:复制,7:创建文件夹,8:修改文件) |
| repository_id | bigint | 20 | 否 | 0 | 文件ID，操作对象 |
| file_name | varchar | 255 | 否 | '' | 文件名，操作对象名称 |
| create_time | int | 11 | 否 | 0 | 创建时间戳，操作发生时间 |
业务规则：
操作日志一旦创建不可修改，确保审计的真实性
系统重要操作必须记录日志，包括但不限于文件上传、下载、删除、分享等
支持按用户、操作类型、时间范围等多维度查询日志
定期归档历史日志，确保系统性能不受影响
(7) 日志文件表
日志文件表用于管理系统日志文件，支持系统监控和问题排查。该表设计面向系统管理员，提供对服务器日志文件的集中管理和监控功能，是系统运维的重要工具。
"id" 作为主键，唯一标识每个日志文件记录。"user_id" 字段关联用户表，确定日志所有者或管理员。"name" 字段记录日志文件名，便于识别不同的日志文件。"host" 字段记录主机信息，支持多服务器环境下的日志管理。"path" 字段存储日志文件路径，指明日志文件在服务器上的位置。"comment" 字段用于添加备注信息，增强日志管理的灵活性。"monitor_choice" 字段用于设置监控选项，支持自定义监控策略。"create_time" 字段记录日志文件创建时间，便于追踪日志文件的生命周期。
该表的设计特别注重系统监控和问题排查需求。通过记录日志文件的位置和主机信息，系统管理员能够快速定位需要查看的日志文件；通过设置监控选项，系统能够自动检测日志中的异常情况并发出警报；通过添加备注信息，管理员能够记录与日志相关的特殊情况或处理方法。这些设计使得系统日志管理更加高效和智能化。
表3-7 日志文件表
| 字段名称 | 类型 | 长度 | 允许空 | 默认值 | 备注 |
|----------|------|------|--------|--------|------|
| id | int | 10 | 否 | 自增 | 日志文件ID |
| user_id | int | 11 | 否 | 无 | 用户ID，日志所有者 |
| name | varchar | 200 | 否 | 无 | 日志文件名，便于识别 |
| host | longtext | - | 否 | 无 | 主机信息，服务器标识 |
| path | varchar | 1024 | 否 | 无 | 日志文件路径，完整位置 |
| create_time | datetime | - | 是 | null | 创建时间，日志生成时间 |
| comment | varchar | 200 | 是 | null | 备注，附加说明 |
| monitor_choice | int | 11 | 是 | null | 监控选择，监控策略配置 |
业务规则：
系统自动收集各服务器的关键日志文件信息
支持按主机、日志类型、时间范围等条件筛选日志
根据监控选项定期检查日志内容，发现异常时触发告警
支持日志文件的在线查看和下载，便于问题排查
(8) SSH信息表
SSH信息表用于存储远程连接信息，支持系统远程管理和日志采集功能。该表设计面向系统管理员，提供对远程服务器的安全访问能力，是系统运维和日志采集的基础设施。
"id" 作为主键，唯一标识每条SSH记录。"user_id" 字段关联用户表，确定SSH配置所有者，通常是系统管理员。"host" 字段记录主机地址，可以是IP地址或域名。"port" 字段记录端口号，默认为22但支持自定义。"username" 和 "password" 字段存储登录凭证，用于建立SSH连接。"created_at" 和 "updated_at" 字段记录创建和更新时间，便于追踪配置变更历史。
该表的设计特别注重安全性和可用性。通过对用户ID和主机地址建立联合唯一索引，确保一个用户对同一主机只有一条SSH配置，避免配置冗余和混淆；通过存储完整的连接参数，支持系统自动建立SSH连接，实现日志采集和远程命令执行等功能；通过记录创建和更新时间，便于审计和配置管理。这些设计使得系统能够安全高效地管理远程服务器连接，支持分布式部署和集中化管理。
表3-8 SSH信息表
| 字段名称 | 类型 | 长度 | 允许空 | 默认值 | 备注 |
|----------|------|------|--------|--------|------|
| id | bigint | 20 | 否 | 自增 | 主键ID |
| user_id | bigint | 20 | 否 | 无 | 用户ID，配置所有者 |
| host | varchar | 255 | 否 | 无 | 主机地址，IP或域名 |
| port | int | 11 | 否 | 无 | 端口号，SSH服务端口 |
| username | varchar | 255 | 否 | 无 | 用户名，SSH登录用户 |
| password | varchar | 255 | 否 | 无 | 密码，SSH登录密码 |
| created_at | datetime | - | 否 | CURRENT_TIMESTAMP | 创建时间 |
| updated_at | datetime | - | 否 | CURRENT_TIMESTAMP | 更新时



```
@startuml
skinparam backgroundColor white
skinparam monochrome true
skinparam shadowing false
skinparam defaultFontName Arial
skinparam defaultFontSize 12
skinparam linetype ortho
skinparam packageStyle rectangle
skinparam componentStyle uml2
skinparam monochrome true
skinparam handDrawn true
title 用户认证流程

actor 用户
participant "前端应用" as Frontend
participant "API网关" as Gateway
participant "用户中心服务" as UserService
database "用户数据库" as UserDB

用户 -> Frontend: 输入用户名和密码
Frontend -> Gateway: 发送登录请求
Gateway -> UserService: 转发认证请求
UserService -> UserDB: 查询用户信息
UserDB --> UserService: 返回用户数据
UserService -> UserService: 验证密码
alt 验证成功
    UserService -> UserService: 生成访问令牌和刷新令牌
    UserService --> Gateway: 返回令牌和用户信息
    Gateway --> Frontend: 返回认证结果和令牌
    Frontend -> Frontend: 存储令牌
    Frontend --> 用户: 显示登录成功，跳转到主页
else 验证失败
    UserService --> Gateway: 返回认证失败信息
    Gateway --> Frontend: 返回认证失败
    Frontend --> 用户: 显示错误信息
end
@enduml
```



```
@startuml
skinparam backgroundColor white
skinparam monochrome true
skinparam shadowing false
skinparam defaultFontName Arial
skinparam defaultFontSize 12
skinparam linetype ortho
skinparam packageStyle rectangle
skinparam componentStyle uml2
skinparam monochrome true
skinparam handDrawn true
title 文件上传流程

actor 用户
participant "前端应用" as Frontend
participant "上传服务" as UploadService
participant "存储服务" as StorageService
database "文件元数据库" as FileDB
database "对象存储" as OSS

用户 -> Frontend: 选择文件上传
Frontend -> Frontend: 计算文件哈希
Frontend -> UploadService: 发送文件哈希检查请求
UploadService -> FileDB: 查询哈希值是否存在
alt 文件已存在(秒传)
    FileDB --> UploadService: 返回文件已存在
    UploadService --> Frontend: 通知文件已存在，无需上传
    Frontend -> UploadService: 创建用户文件记录
    UploadService -> FileDB: 保存用户文件关联
    UploadService --> Frontend: 返回上传成功
    Frontend --> 用户: 显示上传成功
else 文件不存在
    FileDB --> UploadService: 返回文件不存在
    UploadService --> Frontend: 返回需要上传文件
    alt 小文件(普通上传)
        Frontend -> UploadService: 上传文件
        UploadService -> StorageService: 存储文件
        StorageService -> OSS: 保存到对象存储
        OSS --> StorageService: 返回存储结果
        StorageService --> UploadService: 返回存储路径
    else 大文件(分片上传)
        Frontend -> UploadService: 请求初始化分片上传
        UploadService -> StorageService: 初始化分片上传
        StorageService -> OSS: 创建分片上传任务
        OSS --> StorageService: 返回上传ID
        StorageService --> UploadService: 返回上传ID和分片策略
        UploadService --> Frontend: 返回分片上传参数
        loop 每个分片
            Frontend -> UploadService: 上传分片
            UploadService -> StorageService: 存储分片
            StorageService -> OSS: 上传分片
            OSS --> StorageService: 返回分片结果
            StorageService --> UploadService: 返回分片状态
            UploadService --> Frontend: 返回分片上传结果
        end
        Frontend -> UploadService: 请求完成分片上传
        UploadService -> StorageService: 完成分片上传
        StorageService -> OSS: 合并分片
        OSS --> StorageService: 返回合并结果
        StorageService --> UploadService: 返回存储路径
    end
    UploadService -> FileDB: 保存文件元数据和用户关联
    UploadService --> Frontend: 返回上传成功
    Frontend --> 用户: 显示上传成功
end

@enduml
```





```
@startuml
skinparam backgroundColor white
skinparam monochrome true
skinparam shadowing false
skinparam defaultFontName Arial
skinparam defaultFontSize 12
skinparam linetype ortho
skinparam packageStyle rectangle
skinparam componentStyle uml2
skinparam monochrome true
skinparam handDrawn true

title 文件下载流程

actor 用户
participant "前端应用" as Frontend
participant "上传服务" as UploadService
participant "存储服务" as StorageService
database "文件元数据库" as FileDB
database "对象存储" as OSS

用户 -> Frontend: 请求下载文件
Frontend -> UploadService: 发送文件下载请求
UploadService -> FileDB: 查询文件元数据
FileDB --> UploadService: 返回文件信息
UploadService -> UploadService: 验证用户权限
alt 权限验证通过
    UploadService -> StorageService: 请求文件访问URL
    StorageService -> OSS: 生成临时访问URL
    OSS --> StorageService: 返回临时URL
    StorageService --> UploadService: 返回文件访问URL
    UploadService -> UploadService: 记录下载操作
    UploadService --> Frontend: 返回文件下载链接
    Frontend --> 用户: 开始下载文件
else 权限验证失败
    UploadService --> Frontend: 返回权限错误
    Frontend --> 用户: 显示权限错误信息
end
@enduml
```



```
@startuml
skinparam backgroundColor white
skinparam monochrome true
skinparam shadowing false
skinparam defaultFontName Arial
skinparam defaultFontSize 12
skinparam linetype ortho
skinparam packageStyle rectangle
skinparam componentStyle uml2
skinparam monochrome true
skinparam handDrawn true

title 文件分享流程

actor "分享者" as Sharer
actor "接收者" as Receiver
participant "前端应用" as Frontend
participant "分享服务" as ShareService
database "分享数据库" as ShareDB
participant "上传服务" as UploadService

' 创建分享
Sharer -> Frontend: 选择文件分享
Frontend -> ShareService: 发送创建分享请求
ShareService -> ShareService: 生成分享码和链接
ShareService -> ShareDB: 保存分享记录
ShareDB --> ShareService: 返回保存结果
ShareService --> Frontend: 返回分享链接和提取码
Frontend --> Sharer: 显示分享信息

' 访问分享
Receiver -> Frontend: 访问分享链接
Frontend -> ShareService: 验证分享链接
ShareService -> ShareDB: 查询分享记录
ShareDB --> ShareService: 返回分享信息
ShareService -> ShareService: 验证分享有效性
alt 需要提取码
    ShareService --> Frontend: 请求输入提取码
    Frontend --> Receiver: 提示输入提取码
    Receiver -> Frontend: 输入提取码
    Frontend -> ShareService: 提交提取码
    ShareService -> ShareService: 验证提取码
end
alt 验证通过
    ShareService -> ShareService: 更新访问计数
    ShareService -> UploadService: 获取文件信息
    UploadService --> ShareService: 返回文件信息
    ShareService --> Frontend: 返回文件访问信息
    Frontend --> Receiver: 显示文件预览或下载选项
else 验证失败
    ShareService --> Frontend: 返回验证失败
    Frontend --> Receiver: 显示错误信息
end
@enduml
```



```
@startuml
skinparam backgroundColor white
skinparam monochrome true
skinparam shadowing false
skinparam defaultFontName Arial
skinparam defaultFontSize 12
skinparam linetype ortho
skinparam packageStyle rectangle
skinparam componentStyle uml2
skinparam monochrome true
skinparam handDrawn true

title 文件管理流程

actor 用户
participant "前端应用" as Frontend
participant "上传服务" as UploadService
database "用户存储库" as UserRepo
database "文件元数据库" as FileDB

' 文件浏览
用户 -> Frontend: 请求浏览文件
Frontend -> UploadService: 获取文件列表
UploadService -> UserRepo: 查询用户文件
UserRepo --> UploadService: 返回文件列表
UploadService --> Frontend: 返回文件信息
Frontend --> 用户: 显示文件列表

' 文件搜索
用户 -> Frontend: 输入搜索条件
Frontend -> UploadService: 发送搜索请求
UploadService -> UserRepo: 执行文件搜索
UserRepo --> UploadService: 返回搜索结果
UploadService --> Frontend: 返回匹配文件
Frontend --> 用户: 显示搜索结果

' 文件操作
用户 -> Frontend: 选择文件操作(移动/复制/重命名/删除)
alt 移动文件
    Frontend -> UploadService: 发送移动请求
    UploadService -> UserRepo: 更新文件位置
    UserRepo --> UploadService: 返回更新结果
    UploadService -> UploadService: 记录操作日志
    UploadService --> Frontend: 返回操作结果
else 复制文件
    Frontend -> UploadService: 发送复制请求
    UploadService -> UserRepo: 创建文件副本
    UserRepo --> UploadService: 返回创建结果
    UploadService -> UploadService: 记录操作日志
    UploadService --> Frontend: 返回操作结果
else 重命名文件
    Frontend -> UploadService: 发送重命名请求
    UploadService -> UserRepo: 更新文件名称
    UserRepo --> UploadService: 返回更新结果
    UploadService -> UploadService: 记录操作日志
    UploadService --> Frontend: 返回操作结果
else 删除文件
    Frontend -> UploadService: 发送删除请求
    UploadService -> UserRepo: 标记文件删除
    UserRepo --> UploadService: 返回更新结果
    UploadService -> UploadService: 记录操作日志
    UploadService --> Frontend: 返回操作结果
end
Frontend --> 用户: 显示操作结果
@enduml
```



```
@startuml
skinparam backgroundColor white
skinparam monochrome true
skinparam shadowing false
skinparam defaultFontName Arial
skinparam defaultFontSize 12
skinparam linetype ortho
skinparam packageStyle rectangle
skinparam componentStyle uml2
skinparam monochrome true
skinparam handDrawn true

title 审计日志流程

actor "用户" as User
participant "前端应用" as Frontend
participant "各微服务" as Services
participant "日志服务" as LogService
database "日志数据库" as LogDB

' 日志记录
User -> Frontend: 执行操作(上传/下载/分享等)
Frontend -> Services: 处理业务请求
Services -> LogService: 发送操作日志
LogService -> LogDB: 保存操作记录
LogDB --> LogService: 返回保存结果
LogService --> Services: 确认日志已记录
Services --> Frontend: 返回业务处理结果
Frontend --> User: 显示操作结果

' 日志查询
Frontend -> LogService: 发送日志查询请求
LogService -> LogDB: 查询日志记录
LogDB --> LogService: 返回日志数据
LogService -> LogService: 处理日志数据
LogService --> Frontend: 返回日志信息

' 日志分析
Frontend -> LogService: 发送分析请求
LogService -> LogDB: 执行聚合查询
LogDB --> LogService: 返回统计数据
LogService -> LogService: 生成分析报告
LogService --> Frontend: 返回分析结果
@enduml
```



```
sql修改上述 CREATE TABLE IF NOT EXISTS  `user` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '用户id',
    `version` BIGINT NOT NULL,
    `username` VARCHAR(25) DEFAULT NULL COMMENT '用户名',
    `password` VARCHAR(255) DEFAULT NULL COMMENT '密码',
    `mobile` BIGINT(11) DEFAULT NULL COMMENT '手机号',
    `nickname` VARCHAR(50) NOT NULL COMMENT '昵称',
    `gender` TINYINT(1) NOT NULL DEFAULT '1' COMMENT '性别，1：男，0：女，默认为1',
    `avatar` VARCHAR(255) DEFAULT NULL COMMENT '用户头像',
    `birthday` DATE DEFAULT NULL COMMENT '生日',
    `email` VARCHAR(254) DEFAULT NULL COMMENT '电子邮箱',
    `brief` VARCHAR(255) DEFAULT NULL COMMENT '简介|个性签名',
    `info` TEXT,
    `del_state` INt COMMENT '删除状态，0: 未删除，1：已删除',
    `delete_time` TIMESTAMP COMMENT '删除时间',
    `status` TINYINT(1) NOT NULL DEFAULT '0' COMMENT '状态，0：正常，1：锁定，2：禁用，3：过期',
    `admin` TINYINT(1) NOT NULL DEFAULT '0' COMMENT '是否管理员，1：是，0：否',
    `now_volume` INT(11) NOT NULL DEFAULT '0' COMMENT '当前存储容量',
    `total_volume` INT(11) NOT NULL DEFAULT '1000000000' COMMENT '最大存储容量',
    `create_time` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_name_unique` (`username`),
    UNIQUE KEY `idx_mobile_unique` (`mobile`)
);  CREATE TABLE `user_repository`
(
    `id`                  bigint unsigned NOT NULL AUTO_INCREMENT,
    `user_id`       bigint unsigned NOT NULL DEFAULT '0',
    `parent_id`           bigint unsigned NOT NULL DEFAULT '0' ,
    `repository_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '0则为文件夹, 其他为文件id',
    `name`                varchar(255) NOT NULL DEFAULT '' COMMENT '文件夹名称',
    `status` int(11) NOT NULL DEFAULT '0' COMMENT '文件状态(0正常1已删除2禁用)',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_repository_id` (`repository_id`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 8
  DEFAULT CHARSET = utf8; CREATE TABLE IF NOT EXISTS `upload_history`
(
    `id`         bigint unsigned NOT NULL AUTO_INCREMENT,
    `user_id`       bigint unsigned NOT NULL DEFAULT '0' comment '用户id',
    `file_name` varchar(255) NOT NULL DEFAULT '' COMMENT '文件名',
    `size` int(11) NOT NULL DEFAULT '0' COMMENT '文件大小',
    `repository_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '文件id',
    `status`       tinyint(1) NOT NULL DEFAULT '0' COMMENT '上传状态，0：上传中，1：上传成功，2：上传失败',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_repository_id_unique` (`repository_id`)
    ) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8; CREATE TABLE IF NOT EXISTS `repository_pool` (
    `id`          bigint unsigned NOT NULL AUTO_INCREMENT,
    `identity`    bigint unsigned NOT NULL DEFAULT '0' COMMENT '文件id',
    `oss_key`     varchar(255)   NOT NULL DEFAULT '' COMMENT '文件在OSS中的键',
    `hash`        varchar(32)    NOT NULL DEFAULT '' COMMENT '文件的唯一标识',
    `ext`         varchar(30)    NOT NULL DEFAULT '' COMMENT '文件扩展名',
    `size`        int(11)        NOT NULL DEFAULT '0' COMMENT '文件大小',
    `path`        varchar(255)   NOT NULL DEFAULT '' COMMENT '文件url路径',
    `name`        varchar(255)   NOT NULL DEFAULT '',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `idx_identity` (`identity`),        -- 新增的普通索引
    UNIQUE KEY `idx_hash_unique` (`hash`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8; CREATE TABLE IF NOT EXISTS `share_basic` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL DEFAULT '0',
  `repository_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '公共池中的唯一标识',
  `user_repository_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '用户池子中的唯一标识',
  `expired_time` int(11) NOT NULL DEFAULT '0' COMMENT '失效时间，单位秒, 【0-永不失效】',
  `click_num` int(11) NOT NULL DEFAULT '0' COMMENT '点击次数',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `code` VARCHAR(30) NOT NULL DEFAULT '' COMMENT '提取码',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_repository_id` (`repository_id`),
  KEY `idx_user_repository_id` (`user_repository_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8; CREATE TABLE IF NOT EXISTS `audit`(
    `id` INT UNSIGNED AUTO_INCREMENT COMMENT '操作记录ID',
    `user_id` bigint unsigned NOT NULL DEFAULT '0' comment '用户id',
    `content` TEXT NOT NULL COMMENT '操作内容',
    `flag` tinyint(1) NOT NULL default 0 comment '操作类型，0：上传，1：下载，2：删除，3.恢复 4：重命名，5：移动，6：复制，7：创建文件夹，8：修改文件',
    `repository_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '0则为文件夹, 其他为文件id',
    `file_name`    varchar(255) NOT NULL DEFAULT '' COMMENT '文件夹名称',
    `create_time` int(11) NOT NULL default 0 comment '创建时间戳',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 comment '操作记录表'; CREATE TABLE IF NOT EXISTS `logfile`(
    `id` INT UNSIGNED AUTO_INCREMENT COMMENT '日志文件ID',
    `user_id` INT NOT NULL,
    `name` VARCHAR(200) NOT NULL COMMENT '日志文件名',
    `host` LONGTEXT NOT NULL COMMENT '主机信息',
    `path` VARCHAR(1024) NOT NULL COMMENT '日志文件路径',
    `create_time` DATETIME COMMENT '创建时间',
    `comment` VARCHAR(200) COMMENT '备注',
    `monitor_choice` INT COMMENT '监控选择',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8; CREATE TABLE `ssh_info` (
                            `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
                            `user_id` bigint(20) NOT NULL COMMENT '用户ID',
                            `host` varchar(255) NOT NULL COMMENT '主机地址',
                            `port` int(11) NOT NULL COMMENT '端口号',
                            `username` varchar(255) NOT NULL COMMENT '用户名',
                            `password` varchar(255) NOT NULL COMMENT '密码',
                            `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                            `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                            PRIMARY KEY (`id`),
                            UNIQUE KEY `idx_user_host` (`user_id`, `host`) -- 新增的唯一索引
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='SSH信息表';
```

