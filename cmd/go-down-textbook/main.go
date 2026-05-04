package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/chenyb-go/go-down-textbook/internal/api"
	"github.com/chenyb-go/go-down-textbook/internal/auth"
	"github.com/chenyb-go/go-down-textbook/internal/download"
	"github.com/chenyb-go/go-down-textbook/internal/models"
	"github.com/chenyb-go/go-down-textbook/internal/pdf"
	"github.com/chenyb-go/go-down-textbook/internal/util"
)

const (
	appName    = "go-down-textbook"
	appVersion = "v0.2.0"
)

var reader = bufio.NewReader(os.Stdin)

func main() {
	fmt.Println(util.Header(appName + " " + appVersion))
	fmt.Println("国家智慧教育平台教材下载工具")
	fmt.Println()

	// 输出目录
	outputDir := "./books"
	if len(os.Args) > 1 {
		outputDir = os.Args[1]
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "创建输出目录失败: %v\n", err)
		os.Exit(1)
	}

	for {
		// Step 1: 登录
		token, err := ensureLogin()
		if err != nil {
			fmt.Fprintf(os.Stderr, "登录失败: %v\n", err)
			os.Exit(1)
		}

		client := api.NewClient(token)

		// Step 2: 展示已下载的书
		showDownloaded(outputDir)

		// Step 3: 选择教材
		books, err := selectBooks(client, outputDir)
		if err != nil {
			fmt.Println(util.Error(fmt.Sprintf("选择教材失败: %v", err)))
			fmt.Println("按 Enter 重试...")
			waitForEnter()
			continue
		}

		if len(books) == 0 {
			fmt.Println("未选择任何教材，按 Enter 返回...")
			waitForEnter()
			continue
		}

		// Step 4: 下载
		fmt.Println()
		fmt.Printf("准备下载 %d 本教材到 %s\n", len(books), util.FileLink(outputDir, outputDir))
		fmt.Println()

		downloadBooks(client, books, outputDir)

		// Step 5: 完成
		fmt.Println()
		fmt.Println(util.Success("下载完毕!"))
		fmt.Println("文件保存在: " + util.FileLink(outputDir, outputDir))
		fmt.Println()
		fmt.Println("按 Enter 返回选择教材...")
		waitForEnter()
	}
}

// ensureLogin 确保有有效的 token
func ensureLogin() (string, error) {
	// 先检查缓存的 token
	token, err := auth.GetToken()
	if err == nil && token != "" {
		fmt.Println(util.Success("已使用缓存的 Token"))
		return token, nil
	}

	// 需要登录
	fmt.Println("未检测到登录状态，即将开始登录流程...")
	return auth.LoginViaBrowser()
}

// showDownloaded 展示已下载的教材列表
func showDownloaded(dir string) {
	history, err := download.LoadHistory(dir)
	if err != nil || len(history.Books) == 0 {
		fmt.Println("暂无已下载的教材")
		return
	}

	fmt.Println(util.Header("已下载教材:"))
	for i, book := range history.Books {
		path := filepath.Join(dir, book.Filename)
		sizeMB := float64(book.Size) / 1024 / 1024
		fmt.Printf("  %d. %s  (%.1fMB)  %s\n",
			i+1,
			book.Title,
			sizeMB,
			util.FileLink(path, book.Filename),
		)
	}
	fmt.Printf("共 %d 本，保存目录: %s\n", len(history.Books), util.FileLink(dir, dir))
}

// selectBooks 交互式选择教材
func selectBooks(client *api.Client, outputDir string) ([]models.BookItem, error) {
	fmt.Println()
	fmt.Println(util.Header("正在获取教材目录..."))

	// 获取完整教材目录
	catalog, err := api.FetchCatalog(client)
	if err != nil {
		return nil, fmt.Errorf("获取教材目录失败: %w", err)
	}

	if len(catalog) == 0 {
		return nil, fmt.Errorf("教材目录为空")
	}

	fmt.Printf("  共获取到 %d 本教材\n", len(catalog))

	// 提取年级列表
	grades := api.GetGrades(catalog)
	if len(grades) == 0 {
		return nil, fmt.Errorf("未找到年级信息")
	}

	// 选择年级
	gradeIdx := selectFromList("选择年级", grades)
	if gradeIdx < 0 {
		return nil, nil
	}
	selectedGrade := grades[gradeIdx]
	fmt.Println("  已选: " + util.Bold(selectedGrade))

	// 提取学科列表
	subjects := api.GetSubjects(catalog, selectedGrade)
	if len(subjects) == 0 {
		return nil, fmt.Errorf("未找到 %s 的学科信息", selectedGrade)
	}

	// 选择学科
	subjectIdx := selectFromList("选择学科", subjects)
	if subjectIdx < 0 {
		return nil, nil
	}
	selectedSubject := subjects[subjectIdx]
	fmt.Println("  已选: " + util.Bold(selectedSubject))

	// 筛选教材
	filtered := api.FilterByGrade(catalog, selectedGrade)
	filtered = api.FilterBySubject(filtered, selectedSubject)

	if len(filtered) == 0 {
		return nil, fmt.Errorf("未找到 %s %s 的教材", selectedGrade, selectedSubject)
	}

	// 构建教材选项
	books := make([]models.BookItem, 0, len(filtered))
	for _, entry := range filtered {
		books = append(books, models.BookItem{
			ID:      entry.ID,
			Title:   entry.Title,
			Grade:   selectedGrade,
			Subject: selectedSubject,
		})
	}

	// 多选教材
	selected := multiSelectBooks("选择要下载的教材", books, outputDir)
	return selected, nil
}

// selectFromList 从列表中选择一项
func selectFromList(prompt string, items []string) int {
	fmt.Println()
	fmt.Println(util.Header(prompt + ":"))
	for i, item := range items {
		fmt.Printf("  %d. %s\n", i+1, item)
	}
	fmt.Println()

	for {
		fmt.Printf("请输入序号 (1-%d, 0=返回): ", len(items))
		input := readLine()
		input = strings.TrimSpace(input)

		if input == "0" || input == "" {
			return -1
		}

		idx, err := strconv.Atoi(input)
		if err == nil && idx >= 1 && idx <= len(items) {
			return idx - 1
		}

		fmt.Println(util.Error("输入无效，请重新输入"))
	}
}

// multiSelectBooks 多选教材
func multiSelectBooks(prompt string, books []models.BookItem, outputDir string) []models.BookItem {
	selected := make([]bool, len(books))
	allSelected := false

	for {
		fmt.Println()
		fmt.Println(util.Header(prompt + ":"))
		fmt.Println(util.Dim("  输入 a=全选/取消全选, s=确认, 0=返回"))
		fmt.Println(util.Dim("  输入序号切换选中状态，多个用逗号分隔"))
		fmt.Println()

		for i, book := range books {
			mark := "  "
			if selected[i] {
				mark = util.Success("✓ ")
			}
			downloaded := ""
			if download.HasBook(outputDir, book.ID) {
				downloaded = util.Dim(" [已下载]")
			}
			fmt.Printf("  %s%d. %s%s\n", mark, i+1, book.Title, downloaded)
		}

		count := 0
		for _, s := range selected {
			if s {
				count++
			}
		}
		fmt.Println()
		fmt.Printf("已选 %d 本\n", count)
		fmt.Printf(">>> ")

		input := readLine()
		input = strings.TrimSpace(input)

		switch input {
		case "0":
			return nil
		case "s", "S", "确认":
			var result []models.BookItem
			for i, s := range selected {
				if s {
					result = append(result, books[i])
				}
			}
			return result
		case "a", "A":
			allSelected = !allSelected
			for i := range selected {
				selected[i] = allSelected
			}
		default:
			parts := strings.Split(input, ",")
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if part == "" {
					continue
				}
				idx, err := strconv.Atoi(part)
				if err == nil && idx >= 1 && idx <= len(books) {
					selected[idx-1] = !selected[idx-1]
				}
			}
		}
	}
}

// downloadBooks 下载选中的教材
func downloadBooks(client *api.Client, books []models.BookItem, outputDir string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	manager := download.NewManager(client, 3)
	events := manager.DownloadBooks(ctx, books, outputDir)

	var doneCount int32

	for event := range events {
		switch event.Type {
		case download.EventStart:
			fmt.Printf("  ⬇ %s ...", event.BookTitle)

		case download.EventProgress:
			if event.TotalBytes > 0 {
				pct := float64(event.BytesRead) / float64(event.TotalBytes) * 100
				fmt.Printf("\r  ⬇ %s %.1f%%", event.BookTitle, pct)
			}

		case download.EventDone:
			atomic.AddInt32(&doneCount, 1)
			sizeMB := float64(event.BytesRead) / 1024 / 1024
			fmt.Printf("\r  ✓ %s (%.1fMB)\n", event.BookTitle, sizeMB)

			// 尝试添加书签
			go addBookmarksIfNeeded(client, filepath.Join(outputDir, event.Filename), event.BookID)

		case download.EventError:
			atomic.AddInt32(&doneCount, 1)
			fmt.Printf("\r  ✗ %s: %v\n", event.BookTitle, event.Error)
		}
	}
}

// addBookmarksIfNeeded 如果有书签数据则添加到 PDF
func addBookmarksIfNeeded(client *api.Client, pdfPath, contentID string) {
	if !pdf.HasPDFCPU() {
		return
	}

	bookmarks, err := api.FetchBookmarks(client, contentID)
	if err != nil || len(bookmarks) == 0 {
		return
	}

	if err := pdf.AddBookmarks(pdfPath, bookmarks); err != nil {
		fmt.Fprintf(os.Stderr, "  添加书签失败: %v\n", err)
		return
	}

	fmt.Printf("  📑 已添加目录书签: %s\n", filepath.Base(pdfPath))
}

// readLine 从标准输入读取一行
func readLine() string {
	line, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}
	return strings.TrimRight(line, "\r\n")
}

// waitForEnter 等待用户按 Enter
func waitForEnter() {
	reader.ReadString('\n')
}
