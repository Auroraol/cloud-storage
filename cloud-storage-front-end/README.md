# 简介

🥳 `Electron` + `Vue3` + `Vite` + `Pinia` + `Element Plus` + `TypeScript`

## 运行项目

```bash
# 配置
1. 一键安装 .vscode 目录中推荐的插件
2. node 版本 18.x 或 20+
3. pnpm 版本 8.x 或最新版

# 克隆项目
git clone https://github.com/Auroraol/cloud-storage/tree/main/cloud-storage-front-end

# 进入项目目录
cd cloud-storage

# 安装依赖
pnpm i

# 启动服务
pnpm dev
```

## 打包

```bash

# 根据当前系统环境构建
pnpm build

# 打包成解压后的目录
pnpm build:dir

# 构建 linux 安装包, 已设置构建 AppImage 与 deb 文件
pnpm build:linux

# 构建 MacOS 安装包 (只有在 MacOS 系统上打包), 已设置构建 dmg 文件
pnpm build:macos

# 构建 x64 位 exe
pnpm build:win-x64

# 构建 x32 位 exe
pnpm build:win-x32
```

## 目录结构

```tree
├── script              主进程源码
├   ├── core            主窗口、系统菜单与托盘、本地日志等模块
├   ├── tool            一些工具类方法
├   ├── index.ts
├
├── src                 渲染进程源码
├   ├── api
├   ├── assets
├   ├── ......
├
├── static              静态资源
├   ├── icons           系统图标
```

## 站在巨人的肩膀上

- [electron-vite-vue](https://github.com/electron-vite/electron-vite-vue)
- [electron-vue-admin](https://github.com/PanJiaChen/electron-vue-admin)
- [fast-vue3](https://github.com/study-vue3/fast-vue3)
