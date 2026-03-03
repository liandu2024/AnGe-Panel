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

---

## Main Features

### 1) Left Sidebar Navigation (Quick Jump)
![](./doc/images/sidebar-nav.jpg)

- New "group sidebar" on the left side of the page. Click dots to **jump instantly** to the corresponding group.
- No more scrolling when you have many groups.

### 2) Two Collection Modes: Website + Webpage
- **Website**: Perfect for收藏「a site」and its entry points (e.g., NAS, blog, admin panel).
- **Webpage**: Perfect for收藏「an article / a page」link (e.g., Zhihu articles, news links, tutorial pages).
- Switch between "Website / Webpage" freely - search conditions are preserved.

### 3) Webpage Collection Enhanced
![](./doc/images/webbookmarks.jpg)

- Webpage list supports: one-click pin/unpin, quick edit, quick delete.
- Pin/Create/Edit/Delete **won't refresh the entire page** - no flashing or jumping.

### 4) Icons & Wallpapers - Separate Management with Reuse
- Distinguish between "icons" and "wallpapers" when uploading images.
- Previously uploaded icons/wallpapers can be **reused** from history - no need to re-upload every time.

### 5) Thoughtful UX Improvements
- Duplicate link detection (avoid saving the same URL twice).
- Auto-ellipsis for list titles, hover to see full names.
- Mobile-optimized experience (more compact, cleaner).

![](./doc/images/mobile-1.jpg)
![](./doc/images/mobile-2.jpg)

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

## 🔐 First Login

- **Default Admin Username**: `admin`
- **Default Admin Password**: `admin`

On first startup, sample groups and example website/webpage links will be created automatically for quick testing.

> ⚠️ Please change your password after first login!

## 🐳 Docker Deployment (Recommended)

```bash
docker run -d --name ange-panel --restart=unless-stopped \
  -p 3002:3002 \
  -v /root/ange-data:/data \
  ghcr.io/liandu2024/ange-panel:latest
```

> Note: Data is persisted in `/root/ange-data` directory. You can modify settings in the panel after first startup.

## 📋 Documentation

- [中文文档](https://github.com/liandu2024/AnGe-Panel/blob/main/README.md)
- [English Documentation](https://github.com/liandu2024/AnGe-Panel/blob/main/README_EN.md)

## 🏖️ Community

- [Telegram Group](https://t.me/angeworld2024)
- [AnGe Market](https://blog.angeworld.cc/market)

---

## ❤️ Thanks

- [Sun-Panel](https://github.com/hslr-s/sun-panel) - Original project
