package walker

import (
	"regexp"
	"strings"
)

// DocTags is tag in godoc
type DocTags []string

// ParseDocTag to parse tag from document tag
func ParseDocTag(doc string) (tags DocTags) {
	r, _ := regexp.Compile("\\[(.*?)\\]")
	for _, s := range r.FindAllString(doc, -1) {
		tags = append(tags, strings.ToLower(s[1:len(s)-1]))
	}
	return
}

// Contain to check is tag avaible
func (t DocTags) Contain(tag string) bool {
	tag = strings.ToLower(tag)
	for _, docTag := range t {
		if docTag == tag {
			return true
		}
	}
	return false
}
