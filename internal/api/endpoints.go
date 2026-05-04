package api

import (
	"math/rand"
)

const (
	// BaseCDN CDN 基础地址模板
	BaseCDN = "https://%s.ykt.cbern.com.cn"

	// TagsURLPath 标签（分类树）路径
	TagsURLPath = "/zxx/ndrs/tags/tch_material_tag.json"

	// DataVersionPath 教材列表版本路径
	DataVersionPath = "/zxx/ndrs/resources/tch_material/version/data_version.json"

	// DetailsPath 资源详情路径模板
	DetailsPath = "/zxx/ndrv2/resources/tch_material/details/%s.json"

	// TreePath 教材目录树路径模板
	TreePath = "/zxx/ndrv2/national_lesson/trees/%s.json"

	// TokenCDN 需要认证的下载域名
	TokenCDN = "https://r1-ndr-private.ykt.cbern.com.cn"

	// PublicCDN 公开下载域名
	PublicCDN = "https://c1.ykt.cbern.com.cn"

	// LoginURL 登录页面
	LoginURL = "https://auth.smartedu.cn/uias/login"

	// BasicWorkURL 基础工作台页面
	BasicWorkURL = "https://auth.smartedu.cn/uias/basic_work"
)

// ServerList CDN 服务器列表
var ServerList = []string{"s-file-1", "s-file-2", "s-file-3"}

// RandomServer 随机选择一个 CDN 服务器
func RandomServer() string {
	return ServerList[rand.Intn(len(ServerList))]
}

// ServerBase 返回随机 CDN 服务器的完整基地址
func ServerBase() string {
	return "https://" + RandomServer() + ".ykt.cbern.com.cn"
}

// TagsURL 返回标签 API 完整地址
func TagsURL() string {
	return "https://" + RandomServer() + ".ykt.cbern.com.cn" + TagsURLPath
}

// DataVersionURL 返回数据版本 API 完整地址
func DataVersionURL() string {
	return "https://" + RandomServer() + ".ykt.cbern.com.cn" + DataVersionPath
}

// DetailsURL 返回资源详情 API 完整地址
func DetailsURL(contentID string) string {
	return "https://" + RandomServer() + ".ykt.cbern.com.cn" +
		"/zxx/ndrv2/resources/tch_material/details/" + contentID + ".json"
}

// TreeURL 返回目录树 API 完整地址
func TreeURL(ebookID string) string {
	return "https://" + RandomServer() + ".ykt.cbern.com.cn" +
		"/zxx/ndrv2/national_lesson/trees/" + ebookID + ".json"
}
