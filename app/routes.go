package app

func initRoutes(s *Server) {
	s.CRUDController("book", s.bookController)

	s.Any("application/health-check", s.applicationController.Health)
}
