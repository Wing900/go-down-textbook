package api

import (
	"encoding/json"
	"fmt"

	"github.com/chenyb-go/go-down-textbook/internal/models"
)

// FetchResourceDetails 获取资源详情
func FetchResourceDetails(client *Client, contentID string) (*models.ResourceItem, error) {
	url := DetailsURL(contentID)
	data, err := client.GetJSON(url)
	if err != nil {
		return nil, fmt.Errorf("获取资源详情失败 [%s]: %w", contentID, err)
	}

	var item models.ResourceItem
	if err := json.Unmarshal(data, &item); err != nil {
		return nil, fmt.Errorf("解析资源详情失败: %w", err)
	}

	return &item, nil
}

// ParseResourceURL 从 ti_items 中解析下载链接
// 优先选择 PDF 格式的文件，然后是其他格式
func ParseResourceURL(tiItems []models.TiItem, hasToken bool) (downloadURL, backupURL string, size int64) {
	// 优先查找 PDF 文件
	for _, item := range tiItems {
		if item.TiFormat == "pdf" || item.TiFormat == "PDF" {
			return resolveItemURLs(item, hasToken)
		}
	}

	// 查找 ebook_mapping 标记的文件（排除）
	// 然后选择第一个非 ebook_mapping 的文件
	for _, item := range tiItems {
		if item.TiFileFlag != "ebook_mapping" && item.TiStorage != "" {
			return resolveItemURLs(item, hasToken)
		}
	}

	// 最后选择第一个有存储信息的文件
	for _, item := range tiItems {
		if item.TiStorage != "" {
			return resolveItemURLs(item, hasToken)
		}
	}

	return "", "", 0
}

// resolveItemURLs 解析单个 TiItem 的下载 URL
func resolveItemURLs(item models.TiItem, hasToken bool) (url, backup string, size int64) {
	url = ResolveCDNURL(item.TiStorage, hasToken)
	size = item.TiSize

	// 从 ti_storages 中找备用链接
	if len(item.TiStorages) > 0 {
		for _, s := range item.TiStorages {
			if s.URL != "" {
				backup = ResolveCDNURL(s.URL, hasToken)
				break
			}
		}
	}

	// 如果没有备用链接，用主链接替换域名生成
	if backup == "" && url != "" {
		backup = replaceCDNServer(url)
	}

	return url, backup, size
}

// replaceCDNServer 替换 URL 中的 CDN 服务器域名
func replaceCDNServer(url string) string {
	for _, server := range ServerList {
		domain := server + ".ykt.cbern.com.cn"
		if containsStr(url, domain) {
			alt := RandomServer() + ".ykt.cbern.com.cn"
			if alt == domain {
				// 选一个不同的
				for _, s := range ServerList {
					if s+".ykt.cbern.com.cn" != domain {
						alt = s + ".ykt.cbern.com.cn"
						break
					}
				}
			}
			return replaceDomain(url, domain, alt)
		}
	}
	return url
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsSubstring(s, sub))
}

func containsSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func replaceDomain(url, old, new string) string {
	result := ""
	i := 0
	for i <= len(url)-len(old) {
		if url[i:i+len(old)] == old {
			result += new
			i += len(old)
		} else {
			result += string(url[i])
			i++
		}
	}
	result += url[i:]
	return result
}
