# go-down-textbook

国家智慧教育平台教材下载工具 — 批量下载中小学教材 PDF。

## 技术原理

```
用户 → CLI → API 客户端 → 智慧教育平台 API → CDN 下载 PDF
         ↓
    chromedp 自动登录获取 Token
```

1. **登录** — 调用 chromedp 启动本地 Chrome，用户登录后通过 CDP 协议从 `localStorage` 自动提取 `access_token`，无需手动复制粘贴
2. **选书** — 调用平台开放 API 获取教材目录，按年级→学科→教材名三级筛选，支持多选
3. **下载** — 通过 X-ND-AUTH 鉴权头访问私有 CDN，并发下载（默认 3 路），失败自动切备用 CDN 域名重试
4. **书签** — 可选依赖 `pdfcpu`，自动从平台拉取目录树生成 PDF 书签

## 特点

- **小白开箱即用** — 下载即用，自动打开浏览器登录，无需配环境变量、无需手动复制 Token
- **内存极简** — 流式下载 + 32KB 固定缓冲区，不将 PDF 整体加载到内存，理论上可处理任意大小文件
- **跨平台** — Windows / macOS / Linux 全支持
- **断点续感知** — 已下载的教材自动跳过，不重复下载

## 教程

### 1. 下载

从 [Releases](https://github.com/Wing900/go-down-textbook/releases) 下载对应系统的可执行文件或压缩包。

| 文件 | 平台 |
|------|------|
| `go-down-textbook-windows-amd64.exe` | Windows 10/11 64位 |
| `go-down-textbook-darwin-amd64` | macOS Intel |
| `go-down-textbook-darwin-arm64` | macOS Apple Silicon |
| `go-down-textbook-linux-amd64` | Linux 64位 |

### 2. 运行

直接双击或终端运行：

```bash
# 下载到程序所在目录下的 books/
go-down-textbook.exe

# 或指定输出目录
go-down-textbook.exe D:\我的教材
```

### 3. 登录

程序自动打开 Chrome 浏览器 → 在国家智慧教育平台登录您的账号 → 登录成功后自动获取 Token，无需手动操作。

### 4. 选书与下载

按提示选择年级 → 学科 → 教材（支持多选），确认后自动并发下载。下载完成后按 `o` 可在文件管理器中打开目录。

### 5. 书签（可选）

如需自动生成 PDF 书签（目录），安装 [pdfcpu](https://github.com/pdfcpu/pdfcpu) 后放在 `PATH` 即可：

```bash
go install github.com/pdfcpu/pdfcpu/cmd/pdfcpu@latest
```

## 开发打包

Windows 分发包会使用仓库根目录的 `logo.jpg` 自动生成程序图标，并将图标嵌入 `go-down-textbook.exe`：

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\build-windows.ps1
```

也可以直接执行：

```bash
make package-windows
```

产物输出到 `dist/`：

- `dist/go-down-textbook-windows-amd64/go-down-textbook.exe`
- `dist/go-down-textbook-windows-amd64.zip`
