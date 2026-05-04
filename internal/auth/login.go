package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/chenyb-go/go-down-textbook/internal/util"
)

// LoginViaBrowser 通过浏览器登录获取 token
// 使用本地 HTTP 服务器 + JS 代码段的方式：
// 1. 启动本地 HTTP 服务器接收 token
// 2. 打印 JS 代码段供用户在浏览器 DevTools 中执行
// 3. JS 代码将 token POST 到本地服务器
func LoginViaBrowser() (string, error) {
	fmt.Println()
	fmt.Println(util.Header("=== 登录国家智慧教育平台 ==="))
	fmt.Println()

	// 创建一个 channel 接收 token
	tokenCh := make(chan string, 1)
	errCh := make(chan error, 1)

	// 启动本地 HTTP 服务器
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", fmt.Errorf("启动本地服务器失败: %w", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port

	mux := http.NewServeMux()
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var body struct {
			Token string `json:"token"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if body.Token == "" {
			http.Error(w, "Token is empty", http.StatusBadRequest)
			return
		}

		// 允许 CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))

		tokenCh <- body.Token
	})

	// 处理 CORS preflight
	mux.HandleFunc("/token/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusOK)
			return
		}
		mux.ServeHTTP(w, r)
	})

	server := &http.Server{Handler: mux}
	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	// 尝试打开浏览器
	browserOpened := tryOpenBrowser("https://auth.smartedu.cn/uias/basic_work")

	// 打印说明
	fmt.Println("请按以下步骤操作：")
	fmt.Println()
	if browserOpened {
		fmt.Println("  1. 浏览器已打开，请登录您的账号")
	} else {
		fmt.Println("  1. 请在浏览器中打开: " + util.HTTPLink("https://auth.smartedu.cn/uias/basic_work", "https://auth.smartedu.cn/uias/basic_work"))
	}
	fmt.Println("  2. 登录成功后，按 F12 打开开发者工具")
	fmt.Println("  3. 切换到 Console（控制台）标签")
	fmt.Println("  4. 粘贴以下代码并按回车：")
	fmt.Println()

	jsCode := fmt.Sprintf(`fetch("http://127.0.0.1:%d/token",{method:"POST",headers:{"Content-Type":"application/json"},body:JSON.stringify({token:(function(){var k=Object.keys(localStorage).find(k=>k.startsWith("ND_UC_AUTH"));if(!k)return"";var d=JSON.parse(localStorage.getItem(k));return JSON.parse(d.value).access_token})()})}).then(r=>r.json()).then(d=>console.log("Token 已发送！",d))`, port)

	fmt.Println(util.Dim("─────────────────────────────────────────────────────────"))
	fmt.Println(util.Bold(jsCode))
	fmt.Println(util.Dim("─────────────────────────────────────────────────────────"))
	fmt.Println()

	// 倒计时等待
	fmt.Println("等待 Token... (5 分钟超时)")

	select {
	case token := <-tokenCh:
		fmt.Println(util.Success("✓ Token 获取成功!"))
		// 保存 token
		if err := SaveToken(token); err != nil {
			fmt.Fprintf(os.Stderr, "警告: 保存 token 失败: %v\n", err)
		}
		// 关闭服务器
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		server.Shutdown(ctx)
		return token, nil

	case err := <-errCh:
		return "", fmt.Errorf("服务器错误: %w", err)

	case <-time.After(5 * time.Minute):
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		server.Shutdown(ctx)
		return "", fmt.Errorf("等待超时，请重试")
	}
}

// tryOpenBrowser 尝试打开浏览器
func tryOpenBrowser(url string) bool {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}

	if cmd != nil {
		cmd.Stdout = nil
		cmd.Stderr = nil
		if err := cmd.Start(); err != nil {
			return false
		}
		return true
	}
	return false
}

// isChromeRunning 检查 Chrome 是否正在运行（Windows）
func isChromeRunning() bool {
	if runtime.GOOS != "windows" {
		return false
	}

	cmd := exec.Command("tasklist", "/FI", "IMAGENAME eq chrome.exe")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.Contains(string(output), "chrome.exe")
}

// findChrome 查找 Chrome 可执行文件路径（Windows）
func findChrome() string {
	if runtime.GOOS != "windows" {
		return ""
	}

	// 常见的 Chrome 安装路径
	paths := []string{
		os.Getenv("PROGRAMFILES") + `\Google\Chrome\Application\chrome.exe`,
		os.Getenv("PROGRAMFILES(X86)") + `\Google\Chrome\Application\chrome.exe`,
		os.Getenv("LOCALAPPDATA") + `\Google\Chrome\Application\chrome.exe`,
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	return ""
}
