package auth

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

// LoginFunc 表示一个可复用的登录实现。
type LoginFunc func() (string, error)

// SessionManager 统一管理 token 的获取与失效刷新。
type SessionManager struct {
	login LoginFunc
	mu    sync.Mutex
}

// NewSessionManager 创建会话管理器。
func NewSessionManager(login LoginFunc) *SessionManager {
	return &SessionManager{login: login}
}

// EnsureToken 返回当前可用 token。
// 它优先复用环境变量或本地缓存，没有时再触发登录。
func (m *SessionManager) EnsureToken() (string, error) {
	token, err := GetToken()
	if err != nil {
		return "", err
	}
	if token != "" {
		return token, nil
	}
	return m.login()
}

// RefreshToken 在 token 失效后刷新登录态。
// 多个并发请求同时遇到 401 时，只会串行执行一次实际刷新。
func (m *SessionManager) RefreshToken(invalidated string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	current, err := GetToken()
	if err != nil {
		return "", err
	}

	// 若其他协程已经刷新过 token，则直接复用最新值。
	if current != "" && current != invalidated {
		return current, nil
	}

	// 环境变量 token 由外部注入，工具本身无法安全覆写。
	if envToken := strings.TrimSpace(os.Getenv("SMARTEDU_TOKEN")); envToken != "" && envToken == strings.TrimSpace(invalidated) {
		return "", fmt.Errorf("SMARTEDU_TOKEN 已失效，请更新环境变量或移除后重新登录")
	}

	if err := DeleteToken(); err != nil {
		return "", err
	}

	return m.login()
}
