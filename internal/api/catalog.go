package api

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/chenyb-go/go-down-textbook/internal/models"
)

// FetchCatalog 获取完整的教材目录
// 1. 获取数据版本信息（URL 列表）
// 2. 并发获取所有数据文件
// 3. 合并结果
func FetchCatalog(client *Client) ([]models.CatalogEntry, error) {
	// Step 1: 获取版本信息
	versionURL := DataVersionURL()
	versionData, err := client.GetJSON(versionURL)
	if err != nil {
		return nil, fmt.Errorf("获取数据版本失败: %w", err)
	}

	var version models.DataVersion
	if err := json.Unmarshal(versionData, &version); err != nil {
		return nil, fmt.Errorf("解析数据版本失败: %w", err)
	}

	if len(version.URLs) == 0 {
		return nil, fmt.Errorf("数据版本中无 URL 列表")
	}

	// Step 2: 并发获取所有数据文件
	var (
		mu      sync.Mutex
		wg      sync.WaitGroup
		allEntrys []models.CatalogEntry
		errs    []string
	)

	sem := make(chan struct{}, 5) // 最多 5 个并发

	for _, dataURL := range version.URLs {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			// 将 s-file-N 替换为随机服务器
			actualURL := url
			for _, server := range ServerList {
				if strings.Contains(url, server+".ykt.cbern.com.cn") {
					actualURL = strings.Replace(url, server+".ykt.cbern.com.cn",
						RandomServer()+".ykt.cbern.com.cn", 1)
					break
				}
			}

			data, err := client.GetJSON(actualURL)
			if err != nil {
				mu.Lock()
				errs = append(errs, fmt.Sprintf("获取 %s 失败: %v", url, err))
				mu.Unlock()
				return
			}

			var entries []models.CatalogEntry
			if err := json.Unmarshal(data, &entries); err != nil {
				mu.Lock()
				errs = append(errs, fmt.Sprintf("解析 %s 失败: %v", url, err))
				mu.Unlock()
				return
			}

			mu.Lock()
			allEntrys = append(allEntrys, entries...)
			mu.Unlock()
		}(dataURL)
	}

	wg.Wait()

	if len(allEntrys) == 0 && len(errs) > 0 {
		return nil, fmt.Errorf("获取教材目录失败: %s", strings.Join(errs, "; "))
	}

	return allEntrys, nil
}

// FilterByGrade 按年级筛选教材
func FilterByGrade(entries []models.CatalogEntry, grade string) []models.CatalogEntry {
	var result []models.CatalogEntry
	for _, entry := range entries {
		if containsTag(entry.TagList, grade) {
			result = append(result, entry)
		}
	}
	return result
}

// FilterBySubject 按学科筛选教材
func FilterBySubject(entries []models.CatalogEntry, subject string) []models.CatalogEntry {
	var result []models.CatalogEntry
	for _, entry := range entries {
		if containsTag(entry.TagList, subject) {
			result = append(result, entry)
		}
	}
	return result
}

// GetGrades 从教材目录中提取所有年级
func GetGrades(entries []models.CatalogEntry) []string {
	gradeOrder := []string{
		"一年级", "二年级", "三年级", "四年级", "五年级", "六年级",
		"七年级", "八年级", "九年级",
		"高一", "高二", "高三",
	}

	seen := make(map[string]bool)
	for _, entry := range entries {
		for _, tag := range entry.TagList {
			if isGradeTag(tag) && !seen[tag] {
				seen[tag] = true
			}
		}
	}

	var grades []string
	for _, g := range gradeOrder {
		if seen[g] {
			grades = append(grades, g)
		}
	}

	// 添加未在预定义列表中的年级
	for g := range seen {
		found := false
		for _, og := range gradeOrder {
			if g == og {
				found = true
				break
			}
		}
		if !found {
			grades = append(grades, g)
		}
	}

	return grades
}

// GetSubjects 从教材目录中提取指定年级的所有学科
func GetSubjects(entries []models.CatalogEntry, grade string) []string {
	seen := make(map[string]bool)
	var subjects []string

	for _, entry := range entries {
		hasGrade := containsTag(entry.TagList, grade)
		if !hasGrade {
			continue
		}
		for _, tag := range entry.TagList {
			if isSubjectTag(tag) && !seen[tag] {
				subjects = append(subjects, tag)
				seen[tag] = true
			}
		}
	}

	return subjects
}

// isGradeTag 判断是否为年级标签
func isGradeTag(tag string) bool {
	gradeKeywords := []string{
		"一年级", "二年级", "三年级", "四年级", "五年级", "六年级",
		"七年级", "八年级", "九年级",
		"高一", "高二", "高三",
	}
	for _, kw := range gradeKeywords {
		if strings.Contains(tag, kw) {
			return true
		}
	}
	return false
}

// isSubjectTag 判断是否为学科标签
func isSubjectTag(tag string) bool {
	subjectKeywords := []string{
		"语文", "数学", "英语", "物理", "化学", "生物",
		"历史", "地理", "政治", "道德与法治", "科学",
		"音乐", "美术", "体育", "信息技术", "劳动",
		"书法", "日语", "俄语", "艺术", "综合实践",
		"通用技术", "心理健康", "信息科技",
	}
	for _, kw := range subjectKeywords {
		if tag == kw || strings.Contains(tag, kw) {
			return true
		}
	}
	return false
}

// containsTag 检查标签列表中是否包含指定标签
func containsTag(tags []string, target string) bool {
	for _, tag := range tags {
		if strings.Contains(tag, target) || strings.Contains(target, tag) {
			return true
		}
	}
	return false
}
