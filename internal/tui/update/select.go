package update

type titled interface{ Title() string }

func IsLeftClick(action, button string) bool {
	if button != "left" {
		return false
	}
	return action == "left press" || action == "left release"
}

func SelectedItemTitle(item interface{}) string {
	if item == nil {
		return ""
	}
	if v, ok := item.(titled); ok {
		return v.Title()
	}
	return ""
}
