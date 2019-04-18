package app

func initRoutes(s *server) {
	s.CRUDController("book", s.bookController)
}
