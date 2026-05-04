package models

// TagBase 标签基础信息
type TagBase struct {
	ID          string `json:"id"`
	TagName     string `json:"tag_name"`
	ParentTagID string `json:"parent_tag_id"`
	TagType     string `json:"tag_type"`
	Sort        int    `json:"sort"`
}

// TagItem 标签项，包含子标签
type TagItem struct {
	TagBase
	Children []TagItem `json:"children,omitempty"`
}

// DataVersion 教材列表版本信息
type DataVersion struct {
	URLs []string `json:"urls"`
}

// ResourceItem 资源详情
type ResourceItem struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	TiItems  []TiItem `json:"ti_items"`
	TagsList []string `json:"tags_list"`
}

// TiItem 资源文件信息
type TiItem struct {
	TiStorage     string            `json:"ti_storage"`
	TiStorages    []TiStorageRecord `json:"ti_storages"`
	TiFileFlag    string            `json:"ti_file_flag"`
	TiIsSource    bool              `json:"ti_is_source_file"`
	TiFormat      string            `json:"ti_format"`
	TiMD5         string            `json:"ti_md5"`
	TiSize        int64             `json:"ti_size"`
}

// TiStorageRecord 存储记录
type TiStorageRecord struct {
	URL    string `json:"url"`
	Domain string `json:"domain"`
	Bucket string `json:"bucket"`
}

// CatalogEntry 教材目录条目
type CatalogEntry struct {
	ContentID string   `json:"content_id"`
	Title     string   `json:"title"`
	TagList   []string `json:"tag_list"`
}

// TagByType 按类型分组标签
func TagByType(tags []TagItem, tagType string) []TagItem {
	var result []TagItem
	for _, t := range tags {
		if t.TagType == tagType {
			result = append(result, t)
		}
	}
	return result
}
