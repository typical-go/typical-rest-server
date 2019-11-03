package collection

// Interfaces is slice of interface{}
type Interfaces []interface{}

// Add item to Interfaces
func (i *Interfaces) Add(item interface{}) {
	*i = append(*i, item)
}
