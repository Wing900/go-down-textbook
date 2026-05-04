package api

import (
	"encoding/json"
	"fmt"
)

// Bookmark 教材书签（目录节点）
type Bookmark struct {
	Title    string     `json:"title"`
	Page     int        `json:"page"`
	Level    int        `json:"level"`
	Children []Bookmark `json:"children,omitempty"`
}

// EbookMapping 电子书映射数据
type EbookMapping struct {
	EbookID string    `json:"ebook_id"`
	Mapping []Mapping `json:"mappings"`
}

// Mapping 页面映射
type Mapping struct {
	NodeID     string `json:"node_id"`
	PageNumber int    `json:"page_number"`
}

// TreeNode 目录树节点
type TreeNode struct {
	ID       string     `json:"id"`
	Title    string     `json:"title"`
	Children []TreeNode `json:"children,omitempty"`
}

// FetchBookmarks 获取教材的书签（目录）
// 1. 从资源详情中找到 ebook_mapping 文件
// 2. 下载 ebook_mapping.txt 获取 ebook_id 和页面映射
// 3. 获取目录树
// 4. 合并目录树和页面映射，生成书签列表
func FetchBookmarks(client *Client, contentID string) ([]Bookmark, error) {
	// Step 1: 获取资源详情
	detail, err := FetchResourceDetails(client, contentID)
	if err != nil {
		return nil, err
	}

	// Step 2: 找到 ebook_mapping 文件
	var ebookMappingURL string
	for _, item := range detail.TiItems {
		if item.TiFileFlag == "ebook_mapping" && item.TiStorage != "" {
			ebookMappingURL = ResolveCDNURL(item.TiStorage, client.Token != "")
			break
		}
	}

	if ebookMappingURL == "" {
		return nil, fmt.Errorf("未找到 ebook_mapping 文件")
	}

	// Step 3: 下载 ebook_mapping
	mappingData, err := client.GetJSON(ebookMappingURL)
	if err != nil {
		return nil, fmt.Errorf("下载 ebook_mapping 失败: %w", err)
	}

	var mapping EbookMapping
	if err := json.Unmarshal(mappingData, &mapping); err != nil {
		return nil, fmt.Errorf("解析 ebook_mapping 失败: %w", err)
	}

	if mapping.EbookID == "" {
		return nil, fmt.Errorf("ebook_mapping 中无 ebook_id")
	}

	// Step 4: 获取目录树
	treeURL := TreeURL(mapping.EbookID)
	treeData, err := client.GetJSON(treeURL)
	if err != nil {
		return nil, fmt.Errorf("获取目录树失败: %w", err)
	}

	var treeNodes []TreeNode
	if err := json.Unmarshal(treeData, &treeNodes); err != nil {
		return nil, fmt.Errorf("解析目录树失败: %w", err)
	}

	// Step 5: 构建 node_id -> page_number 映射
	pageMap := make(map[string]int)
	for _, m := range mapping.Mapping {
		pageMap[m.NodeID] = m.PageNumber
	}

	// Step 6: 递归构建书签树
	bookmarks := buildBookmarks(treeNodes, pageMap, 1)

	return bookmarks, nil
}

// buildBookmarks 递归构建书签列表
func buildBookmarks(nodes []TreeNode, pageMap map[string]int, level int) []Bookmark {
	var bookmarks []Bookmark

	for _, node := range nodes {
		page := pageMap[node.ID]
		bm := Bookmark{
			Title: node.Title,
			Page:  page,
			Level: level,
		}

		if len(node.Children) > 0 {
			bm.Children = buildBookmarks(node.Children, pageMap, level+1)
		}

		bookmarks = append(bookmarks, bm)
	}

	return bookmarks
}

// FlattenBookmarks 将嵌套的书签展平为列表
func FlattenBookmarks(bookmarks []Bookmark) []Bookmark {
	var result []Bookmark
	for _, bm := range bookmarks {
		result = append(result, Bookmark{
			Title: bm.Title,
			Page:  bm.Page,
			Level: bm.Level,
		})
		if len(bm.Children) > 0 {
			result = append(result, FlattenBookmarks(bm.Children)...)
		}
	}
	return result
}
