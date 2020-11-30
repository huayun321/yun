# yun
### yun.web框架
----

#### 用途

* 该项目为学习GO HTTP框架原理所用
* 提供的功能为路由 context 前缀树路由 路由分组 中间件 模版 错误恢复


#### day 1

* 实现handler

```shell
➜  huayun321 curl http://localhost:9999/hello
Header["User-Agent"] = ["curl/7.64.1"]
Header["Accept"] = ["*/*"]
➜  huayun321 curl http://localhost:9999      
URL.Path = "/"
➜  huayun321 curl http://localhost:9999/hello
Header["User-Agent"] = ["curl/7.64.1"]
Header["Accept"] = ["*/*"]
➜  huayun321 curl http://localhost:9999/x    
404 NOT FOUND: /x

```

* 封装进yun包

```shell
➜  huayun321 curl http://localhost:9999      
URL.Path = "/"
➜  huayun321 curl http://localhost:9999/hello
Header["Accept"] = ["*/*"]
Header["User-Agent"] = ["curl/7.64.1"]
➜  huayun321 curl http://localhost:9999/x    
404 NOT FOUND: /x

```
