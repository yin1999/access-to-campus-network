# 河海大学校园网登录认证工具

由于学校网络认证登录站点会自动跳转到 `http://www.baidu.com/`，在用户未获得访问互联网权限时，会劫持网页内容，使浏览器跳转到认证登录页面。但用户如果在近期使用浏览器访问过百度，浏览器会因为 [`HSTS策略`](https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Headers/Strict-Transport-Security) 自动由 `http` 协议跳转到 `https` 协议，从而使认证系统劫持的网页无法正常被用户访问（通常浏览器会提示有关 **不安全** 的警告）。

该项目用于方便用户为了在不更换浏览器访问的情况（或针对服务器无GUI的情况）下能够方便地完成校园网认证。

## 使用说明

前往 [Release](https://gitee.com/allo123/access-to-campus-network/releases) 下载对应平台的二进制文件，将其放置到方便自己管理的路径下。可采用两种方式启动程序：配置文件、传递参数，为方便后续使用，推荐使用创建配置文件的方法。

## 配置文件

在程序二进制文件所在的目录下，新建配置文件 `config.ini`。打开并配置以下选项：

```ini
user=username     # 用户名
password=password # 密码
count=5           # 最大重试认证次数（可选，默认为 5）
net=out-campus    # 网络接入类型（可选[ out-campus | cmcc ]，默认为 out-campus）
```

使用时，仅需直接打开程序二进制文件，即可根据配置文件自动尝试网络接入认证（Windows系统可在桌面建立快捷方式）。

关于网络接入类型两个选项的说明如下（与校园网系统中的一致）：

| 接入类型 | 说明 |
| ------- | ---- |
| out-campus | 校园外网服务(out-campus NET) |
| cmcc | 中国移动(CMCC NET) |

## 传递参数

程序也可使用传递参数的方式启动：

```bash
./access-network -u username -p password -c 5 -net out-campus
```

传递参数的各项参数与 配置文件 方法中的各项配置对应关系如下（与配置文件相同，`-c`、`-net` 参数均为可选项，且有相同的默认值）：

| 参数名 | 对应配置文件 |
| ----- | ---------- |
| -u | user |
| -p | password |
| -c | count |
| -net | net |

上述参数及其使用方法可通过以下命令查询：

```bash
./access-network -h
```
