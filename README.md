# 基于时间轮的定时器的GO实现

使用方式跟标准库类似。实现方式为: 单层时间轮 + 小堆 + 一个任务一个协程（这点可以改为协程池）

参考链接中有更高级的实现

```go
w := timewheel.NewWheel(time.Second, 5)

log.Println("hello")

for i := 0; i < 6; i++ {
    i := i
    w.AfterFunc(time.Second*time.Duration(i), func() {
        log.Println("hello world: ", i)
    })
    w.AfterFunc(time.Second*time.Duration(i), func() {
        log.Println("hello world: ", i)
    })
}
w.Start()

w.CloseAndWait()
```


参考链接
--------------
[1分钟实现“延迟消息”功能](http://www.10tiao.com/html/249/201703/2651959961/1.html)

[层级时间轮的 Golang 实现](http://russellluo.com/2018/10/golang-implementation-of-hierarchical-timing-wheels.html)

https://github.com/RussellLuo/timingwheel