package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vynious/gd-joinqueue-cms/logger"
	"log"
)

type Server struct {
	Router *gin.Engine
	qh     *QueueHandler
}

func SpawnServer() *Server {
	engine := gin.Default()
	prod := logger.SpawnKafkaProducer(logger.LoadKafkaConfigurations())
	qh := SpawnQueueHandler(prod)
	return &Server{
		Router: engine,
		qh:     qh,
	}
}

func (s *Server) MountHandlers() {
	api := s.Router.Group("/api/queue")

	api.Use(s.GenerateRequestID)

	api.POST("/join", s.qh.JoinQueue)
	api.GET("/upcoming", s.qh.GetUpcomingTicketsInQueue)
	api.GET("/next", s.qh.RetrieveNextInQueue)
}

func (s *Server) GenerateRequestID(c *gin.Context) {
	log.Println("Generating request ID")
	requestID := uuid.New().String()
	log.Println("Generated request ID:", requestID)
	c.Set("request_id", requestID)
	c.Header("X-Request-ID", requestID)
	c.Next()
}
