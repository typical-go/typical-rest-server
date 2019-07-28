package app

func initRoutes(s *Server) {

	s.GET("book", s.bookController.List)
	s.POST("book", s.bookController.Create)
	s.GET("book/:id", s.bookController.Get)
	s.PUT("book", s.bookController.Update)
	s.DELETE("book/:id", s.bookController.Delete)

	s.Any("application/health-check", s.applicationController.Health)
}
