package typdocker

type (
	// Composer responsible to compose docker
	Composer interface {
		ComposeV3() (*Recipe, error)
	}
	// ComposeFn function
	ComposeFn    func() (*Recipe, error)
	composerImpl struct {
		fn ComposeFn
	}
)

//
// composerImpl
//

var _ Composer = (*composerImpl)(nil)

// NewCompose return new instance of composer
func NewCompose(fn ComposeFn) Composer {
	return &composerImpl{fn: fn}
}

func (i *composerImpl) ComposeV3() (*Recipe, error) {
	return i.fn()
}
