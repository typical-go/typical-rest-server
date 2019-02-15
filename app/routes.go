package app

func initRoutes(s *server) {
	s.CRUD("book", s.bookController)
}
