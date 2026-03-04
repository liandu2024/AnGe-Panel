[[ 简体中文 ]](https://github.com/liandu2024/AnGe-Panel/blob/main/README.md) |
[[ English ]](https://github.com/liandu2024/AnGe-Panel/blob/main/README_EN.md)

<div align=center>

<img src="./doc/images/main-preview.jpg" width="100%" />

# AnGe-Panel

[![Github](https://img.shields.io/badge/Github-123456?logo=github&labelColor=242424)](https://github.com/liandu2024/AnGe-Panel)
[![docker](https://img.shields.io/badge/docker-123456?logo=docker&logoColor=fff&labelColor=1c7aed)](https://github.com/liandu2024/AnGe-Panel/pkgs/container/ange-panel) 
[![Telegram](https://img.shields.io/badge/Telegram-123456?logo=telegram&labelColor=229ED9)](https://t.me/angeworld2024)
<br>

</div>

A perfect website navigation + webpage bookmarks panel.

完美的网站导航 + 网页收藏夹。

---

## 主要功能

### 1) 左侧分组目录条（快速定位）
![](./doc/images/sidebar-nav.jpg)

- 页面左侧新增"分组目录条"，点击圆点可 **一键跳转** 到对应分组。
- 分组多的时候，不用一直往下滑找。

### 2) 网站 + 网页 两种收藏模式
- **网站**：适合收藏「一个站点」和它的入口（比如 NAS、博客、后台管理）。
- **网页**：适合收藏「一篇文章 / 一个页面」的链接（比如知乎文章、头条链接、教程页面）。
- 你可以在"网站 / 网页"之间切换，搜索条件也会保持，不会丢。

### 3) 网页收藏更好用
![](./doc/images/webbookmarks.jpg)

- 网页列表支持：一键置顶/取消置顶、快速编辑、快速删除。
- 置顶/新建/修改/删除 **不会整页刷新**，界面不会"闪一下/跳一下"。

### 4) 图标 / 壁纸 分开管理，支持复用
- 上传图片时区分"图标"和"壁纸"，避免混在一起不好找。
- 已上传的图标/壁纸可以在历史里 **重复选择复用**，不需要每次重新上传。

### 5) 一些贴心的体验优化
- 链接去重校验（避免重复收藏同一 URL）。
- 列表标题自动省略、需要时再悬停看全名。
- 移动端体验优化（更紧凑、更清爽）。

![](./doc/images/mobile-1.jpg)
![](./doc/images/mobile-2.jpg)

## 🐳 Docker 部署

在支持 Docker 设备上执行以下指令：

```bash
docker run -d --name ange-panel --restart=unless-stopped \
  -p 3002:3002 \
  -v /root/ange-data:/data \
  ghcr.io/liandu2024/ange-panel:latest
```

> 说明：数据持久化存储在 `/root/ange-data` 目录，首次启动后可在面板内修改配置。

## 🔐 首次登录

登录网址：http://[部署设备的IP]:3002

- **默认管理员账号**：`admin`
- **默认管理员密码**：`admin`

首次启动时会自动创建示例分组和示例网站/网页链接，方便快速体验功能。

> ⚠️ 首次登录后请尽快修改密码！

## 📋 文档

- [中文文档](https://github.com/liandu2024/AnGe-Panel/blob/main/README.md)
- [English Documentation](https://github.com/liandu2024/AnGe-Panel/blob/main/README_EN.md)

## 🏖️ 其他

- [Telegram 群](https://t.me/angeworld2024)
- [安格超市](https://blog.angeworld.cc/market)

---

## ❤️ 感谢

- 本项目基于 [Sun-Panel](https://github.com/hslr-s/sun-panel) 项目做的二次开发，感谢原作者辛勤付出！

---

[![Star History Chart](https://api.star-history.com/svg?repos=liandu2024/AnGe-Panel&type=Date)](https://star-history.com/#liandu2024/AnGe-Panel&Date)
