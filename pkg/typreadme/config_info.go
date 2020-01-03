package typreadme

// ConfigInfos is collection of configuration information
type ConfigInfos []*ConfigInfo

// ConfigInfo is configuration information
type ConfigInfo struct {
	Name     string
	Type     string
	Default  string
	Required string
}

// Append config info
func (c *ConfigInfos) Append(detail *ConfigInfo) {
	*c = append(*c, detail)
}
