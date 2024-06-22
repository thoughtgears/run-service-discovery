package main

import (
	"context"
	"fmt"

	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/thoughtgears/run-service-discovery/db"
	"github.com/thoughtgears/run-service-discovery/handlers"
)

type Config struct {
	Port         string `envconfig:"PORT" default:"8080"`
	GcpProjectID string `envconfig:"GCP_PROJECT_ID" required:"true"`
}

var config Config

func init() {
	envconfig.MustProcess("", &config)

	zerolog.LevelFieldName = "severity"
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = time.RFC3339Nano
}

func main() {
	ctx := context.Background()
	firestoreDB, err := db.NewFirestoreDB(ctx, config.GcpProjectID)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create Firestore client")
	}

	// Create a new gin router and add a logger and recovery middleware
	router := gin.New()
	router.Use(gin.Recovery(), logger.SetLogger(
		logger.WithLogger(func(_ *gin.Context, l zerolog.Logger) zerolog.Logger {
			return l.Output(gin.DefaultWriter).With().Logger()
		}),
	))

	// Register the routes for our service endpoints
	router.POST("/services", handlers.PostService(firestoreDB))
	router.PUT("/services/:name", handlers.UpdateService(firestoreDB))
	router.GET("/services/:name", handlers.GetService(firestoreDB))
	router.GET("/services", handlers.GetServices(firestoreDB))
	router.DELETE("/services/:name", handlers.DeleteService(firestoreDB))

	// Start the server
	log.Info().Msgf("server running on port %s", config.Port)
	if err := router.Run(fmt.Sprintf(":%s", config.Port)); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}
