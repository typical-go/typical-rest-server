package app

import (
	logrusmiddleware "github.com/bakatz/echo-logrusmiddleware"
	log "github.com/sirupsen/logrus"
)

func initLogger(s *Server) {
	s.Logger = logrusmiddleware.Logger{Logger: log.StandardLogger()}
}
