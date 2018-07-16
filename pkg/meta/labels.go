package meta

type Labels map[string]string

type Selector map[string]string

func SelectorMatches(selector Selector, labels Labels) bool {
	for key, v1 := range selector {
		if v2, ok := labels[key]; !ok || v1 != v2 {
			return false
		}
	}
	return true
}
