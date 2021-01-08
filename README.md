# delay queue
delay queue 是基于 Golang 实现的延时队列。基于Redis Zset, 以时间戳作为Score, 主动轮询小于当前的时间的元素。
## 安装
````
go get -u github.com/yasin-wu/delay-queue
````
推荐使用go.mod
<br>
````
require github.com/yasin-wu/delay-queue v1.1.0
````