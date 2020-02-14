package typpostgres

func (m *Module) reset(cfg Config) (err error) {
	if err = m.drop(cfg); err != nil {
		return
	}
	if err = m.create(cfg); err != nil {
		return
	}
	if err = m.migrate(cfg); err != nil {
		return
	}
	if err = m.seed(cfg); err != nil {
		return
	}
	return
}
