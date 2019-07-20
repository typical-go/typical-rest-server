package app

import (
	logrusmiddleware "github.com/bakatz/echo-logrusmiddleware"
	log "github.com/sirupsen/logrus"
)

func initLogger(s *Server) {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.JSONFormatter{
		PrettyPrint: true,
	})

	s.Logger = logrusmiddleware.Logger{Logger: log.StandardLogger()}
}
