package api

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client 封装 HTTP 客户端，自动处理认证和重试
type Client struct {
	HTTPClient *http.Client
	Token      string
}

// NewClient 创建新的 API 客户端
func NewClient(token string) *Client {
	return &Client{
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		Token: token,
	}
}

// GetJSON 执行 GET 请求并返回响应体
// 自动添加认证头，5xx 错误自动重试（最多 3 次）
func (c *Client) GetJSON(url string) ([]byte, error) {
	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(attempt) * time.Second)
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("创建请求失败: %w", err)
		}

		// 添加认证头
		if c.Token != "" {
			req.Header.Set("X-ND-AUTH", FormatAuthHeader(c.Token))
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("请求失败: %w", err)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("读取响应失败: %w", err)
			continue
		}

		if resp.StatusCode >= 500 {
			lastErr = fmt.Errorf("服务器错误 %d: %s", resp.StatusCode, string(body[:min(len(body), 200)]))
			continue
		}

		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body[:min(len(body), 200)]))
		}

		return body, nil
	}
	return nil, lastErr
}

// DownloadStream 执行 GET 请求并返回响应体流（用于大文件下载）
func (c *Client) DownloadStream(url string) (io.ReadCloser, int64, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("创建请求失败: %w", err)
	}

	if c.Token != "" {
		req.Header.Set("X-ND-AUTH", FormatAuthHeader(c.Token))
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("请求失败: %w", err)
	}

	if resp.StatusCode != 200 {
		resp.Body.Close()
		return nil, 0, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return resp.Body, resp.ContentLength, nil
}

// FormatAuthHeader 格式化认证头
func FormatAuthHeader(token string) string {
	return fmt.Sprintf(`MAC id="%s",nonce="0",mac="0"`, token)
}

// ResolveCDNURL 将 CDN 域名模板解析为完整 URL
// ti_storage 格式: "cs_path:${ref-path}/edu_product/esp/assets/..."
// 需要替换 "cs_path:${ref-path}" 为实际的 CDN 基地址
func ResolveCDNURL(template string, hasToken bool) string {
	cdn := PublicCDN
	if hasToken {
		cdn = TokenCDN
	}
	// 替换 cs_path:${ref-path} 为 CDN 地址
	result := strings.ReplaceAll(template, "cs_path:${ref-path}", cdn)
	// 也处理只有 ${ref-path} 的情况
	result = strings.ReplaceAll(result, "${ref-path}", cdn)
	return result
}
