package appctx

// TypiApp is Typical Application
type TypiApp struct {
	ConfigLoader
	Constructors []interface{}
	Action       interface{}
}
