package app

func initRoutes(s *Server) {
	s.CRUDController("book", s.bookController)
}
