[[ 简体中文 ]](https://sun-panel-doc.enianteam.com/zh_cn/introduce/project.html) |
[[ English ]](https://sun-panel-doc.enianteam.com/introduce/project.html)

<div align=center>

<img src="./doc/images/logo.png" width="100" height="100" />

# AnGe-Panel

[![Github](https://img.shields.io/badge/Github-123456?logo=github&labelColor=242424)](https://github.com/hslr-s/sun-panel)
[![Gitee](https://img.shields.io/badge/Gitee-123456?logo=gitee&labelColor=c71d23)](https://gitee.com/hslr/sun-panel)
[![docker](https://img.shields.io/badge/docker-123456?logo=docker&logoColor=fff&labelColor=1c7aed)](https://hub.docker.com/r/hslr/sun-panel) 
[![Bilibili](https://img.shields.io/badge/Bilibili-123456?logo=bilibili&logoColor=fff&labelColor=fb7299)](https://space.bilibili.com/27407696/channel/collectiondetail?sid=2023810)
[![YouTube](https://img.shields.io/badge/YouTube-123456?logo=youtube&labelColor=ff0000)](https://www.youtube.com/channel/UCKwbFmKU25R602z6P2fgPYg)
<br>
[![GitHub User's stars](https://img.shields.io/github/stars/hslr-s%2Fsun-panel?style=flat&logo=github)](https://github.com/hslr-s/sun-panel)
[![github downloads](https://img.shields.io/github/downloads/hslr-s/sun-panel/total.svg?logo=github)](https://github.com/hslr-s/sun-panel/releases)
[![docker pulls](https://img.shields.io/docker/pulls/hslr/sun-panel.svg?logo=docker)](https://hub.docker.com/r/hslr/sun-panel)

[[ 中文文档 ]](https://sun-panel-doc.enianteam.com/zh_cn) |
[[ Document ]](https://sun-panel-doc.enianteam.com) |
[[ Demo ]](http://sunpaneldemo.enianteam.com) 

A perfect website navigation + webpage bookmarks panel.
<br>
完美的网站导航 + 网页收藏夹。

</div>


![](./doc/images/main-dark.png)

> [!IMPORTANT]
> In order to maintain the livelihood, the author added some [`PRO`] (https://pro.sun-panel.top) function, so the project temporarily entered a closed source state.; At present, the latest version of the open source is `v1.3.0`, [Please see the latest version of closed source](https://github.com/hslr-s/sun-panel/releases).; When the modular technology is developed, the separation of the PRO and the programs will be opened again, and the closed source will have no effect on ordinary users.; Let's look forward to open source again, and at the same time, we are welcome to supervise and review the security of the program.
> 
> 作者为了维持生计，增加了一些 [`PRO`](https://pro.sun-panel.top) 功能，所以项目暂时进入闭源状态。目前开源最新版本为`v1.3.0`，[闭源最新版本请查看](https://github.com/hslr-s/sun-panel/releases)。待开发出模块化技术，然后对PRO和主程序进行分离会再次开源，闭源对普通用户没有任何影响。我们一起期待再次开源吧，同时也欢迎各位大佬对程序的安全性进行监督和审查。

## 📌 Fork / 二次开发说明（本仓库改动声明）

> **AnGe-Panel** 是基于 Sun‑Panel（MIT License）的二次开发版本，非官方版本。

## 🔐 首次登录

- **默认管理员账号**：`admin`
- **默认管理员密码**：`admin`

首次启动时会自动创建示例分组和示例网站/网页链接，方便快速体验功能。

> ⚠️ 首次登录后请尽快修改密码！

本仓库基于 **Sun‑Panel 开源版本（MIT License）** 做了面向“日常导航 + 信息收集”的二次开发。

- 本仓库 **不是官方仓库**，与原作者/官方无隶属关系；如需官方版本/PRO 功能请前往原项目。
- 原项目地址：<https://github.com/hslr-s/sun-panel>
- License：保留原项目 MIT License（见 `LICENSE`）。

### ✅ 这个 Fork 做了什么（给普通人看的）

你可以把它理解成：在“导航主页”之外，多加了一个更顺手的 **网页收藏 / 信息流** 模式，并把使用体验做得更像 App。

#### 1) 左侧分组目录条（快速定位）
- 页面左侧新增“分组目录条”，点击圆点可 **一键跳转** 到对应分组。
- 分组多的时候，不用一直往下滑找。

#### 2) 网站 + 网页 两种收藏模式
- **网站**：适合收藏「一个站点」和它的入口（比如 NAS、博客、后台管理）。
- **网页**：适合收藏「一篇文章 / 一个页面」的链接（比如知乎文章、头条链接、教程页面）。
- 你可以在“网站 / 网页”之间切换，搜索条件也会保持，不会丢。

#### 3) 网页收藏更好用
- 网页列表支持：一键置顶/取消置顶、快速编辑、快速删除。
- 置顶/新建/修改/删除 **不会整页刷新**，界面不会“闪一下/跳一下”。

#### 4) 图标 / 壁纸 分开管理，支持复用
- 上传图片时区分“图标”和“壁纸”，避免混在一起不好找。
- 已上传的图标/壁纸可以在历史里 **重复选择复用**，不需要每次重新上传。

#### 5) 一些贴心的体验优化
- 链接去重校验（避免重复收藏同一 URL）。
- 列表标题自动省略、需要时再悬停看全名。
- 移动端体验优化（更紧凑、更清爽）。

> 备注：本 Fork 也包含一些视觉与布局微调（间距、悬停高亮等），但 README 不展开这些细节，避免影响阅读。

## 😎 Features

- 🍉 Clean interface, powerful functionality, low resource consumption
- 🍊 Easy to use, visual operation, zero-code usage
- 🍠 One-click switch between internal and external network modes
- 🍵 Supports Docker deployment (compatible with Arm systems)
- 🎪 Supports multi-account isolation
- 🎏 Supports viewing system status
- 🫙 Supports custom JS, CSS
- 🍻 Simple usage without the need to connect to an external database
- 🍾 Rich icon styles for free combination, supports [Iconify icon library](https://icon-sets.iconify.design/)
- 🚁 Supports opening small windows in the webpage (some third-party websites may block this feature)

## 🖼️ Preview Screenshots

**Various styles, freely combined**

![](./doc/images/icon-small-new.png)
![](./doc/images/transparent-info.png)
![](./doc/images/transparent-small.png)
![](./doc/images/solid-color-info.png)
![](./doc/images/full-color-small.jpg)

**Built-in small windows**

![](./doc/images/window-ssh.png)
![](./doc/images/window-xunlei.png)



## 🐳 Docker（推荐）

本仓库已内置 `Dockerfile`，并提供 GitHub Actions 自动构建镜像（发布到 GHCR）：

- 镜像地址（默认）：`ghcr.io/<你的GitHub账号>/<仓库名>:<tag>`
- 端口：`3002`

示例：

```bash
docker run -d --name ange-panel -p 3002:3002 \
  ghcr.io/<owner>/<repo>:latest
```

> 说明：首次发布后，你可以在 GitHub → Packages 里看到镜像；如需公开下载，请把该 Package 设为 Public。

## 🐳 Deployment tutorial
[Deployment Tutorial](https://sun-panel-doc.enianteam.com/usage/quick_deploy.html)

## 🍵 Donate

> Open-source development is not easy. If you feel that my project has helped you, you are welcome to [donate](./doc/donate.md) or buy me a cup of tea☕ (please leave your nickname or name in the note if possible). Your support is my motivation, thank you.


<a href="https://www.paypal.me/hslrs">
<img height="60" src="./doc/images/donate/paypal.png" target="_blank"></img> 
</a>


|   |   |
| ------------ | ------------ |
| <img height="300" src="./doc/images/donate/weixin.png"/> |  <img height="300" src="./doc/images/donate/alipay.png" /> |

## 🏖️ Communication group & community

Author：**[红烧猎人](https://blog.enianteam.com/u/sun/content/11)**

[Github Discussions](https://github.com/hslr-s/sun-panel/discussions)

QQ交流群，进不去可以点上方连接联系作者

<img src="./doc/images/qq_group_qr2.png"  height="350" />

## ❤️ Thanks

- [Roc](https://github.com/RocCheng)
- [jackloves111](https://github.com/jackloves111)
- [Rock.L](https://github.com/gitlyp)


---

[![Star History Chart](https://api.star-history.com/svg?repos=hslr-s/sun-panel&type=Date)](https://star-history.com/#hslr-s/sun-panel&Date)
