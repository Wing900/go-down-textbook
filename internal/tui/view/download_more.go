package view

type DownloadFooterData struct {
	Done   bool
	Failed []string
}

func downloadTitle(done bool) string {
	if done {
		return "下载完成"
	}
	return "下载中"
}

func footerText(d DownloadData) string {
	if d.Done {
		return "O 打开目录   Enter 返回选书   q 退出"
	}
	if len(d.Failed) > 0 {
		return "O 打开目录   R 回选书页   q 退出"
	}
	return "O 打开目录   q 退出"
}
