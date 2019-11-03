package collection

// Strings is slice of string
type Strings []string

// Add item to Strings
func (s *Strings) Add(item string) {
	*s = append(*s, item)
}
