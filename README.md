
## 疑问

router
1. 为什么数据存储在engine中？方便分离router吗

engine
1. ServeHTTP 执行时，handlers如何捕获错误呢？某个中间的handler执行错误，应该中断才对

server
1. 如何设置conn read/write的ddl

支持http2, http3需要做什么

## todo

router
1. 支持 group

支持http2, http3

支持更多的bind

context
1. 支持next
2. 处理中间错误

开发net包的组件
