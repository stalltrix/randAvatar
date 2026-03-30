# randAvatar

随机头像API

Golang 实现的随机头像生成器。  
通过字符串 seed 唯一的头像，适用于用户默认头像、评论系统、论坛、内部工具等场景。

## 特性

- 相同的 seed 会生成相同的头像
- 支持 SVG 矢量格式，体积小、清晰不失真
- HTTP 服务方式直接使用
- 可作为 Go 库集成到项目中
- 无状态生成，无需存储头像文件

---

## 效果

<p align="center">
    <img src="img/demo1.jpg" width="160">
    <img src="img/demo2.jpg" width="160">
    <img src="img/demo3.jpg" width="160">
    <img src="img/demo4.jpg" width="160">
    <img src="img/demo5.jpg" width="160">
</p>

---


## 直接运行（HTTP 服务）

启动服务：

```
./randAvatar 0.0.0.0:8080
```

启动后访问：

```
http://localhost:8080/avatar?seed=user1234
```


即可获得一个头像。

---

## 伪静态

```
/avatar/{seed}.svg
```

等效于seed

---

## 作为 Go 库使用

在项目中引入：

```
go get github.com/stalltrix/randAvatar/strx_avatar
```

参考代码：

```go
package main

import (
    "fmt"
    "github.com/stalltrix/randAvatar/strx_avatar"
)

func main() {
    svg, err := strx_avatar.NewAvatar("user1234")
    if err != nil {
        panic(err)
    }

    fmt.Println(svg)
}
```


---

## 使用场景

- 用户默认头像
- 评论系统头像
- 自动生成用户图标
- 开发测试环境

---
