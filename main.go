package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/thoughtgears/run-service-discovery/handlers"
	"github.com/thoughtgears/run-service-discovery/pkg/config"
	"github.com/thoughtgears/run-service-discovery/pkg/db"
)

func main() {
	// Set the logger configuration for GCP stackdriver
	zerolog.LevelFieldName = "severity"
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = time.RFC3339Nano

	// Load the configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	// Create a new Firestore client
	ctx := context.Background()
	firestoreDB, err := db.NewFirestoreDB(ctx, cfg.GcpProjectID)
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

	// Start the server
	log.Info().Msgf("server running on port %s", cfg.Port)
	if err := router.Run(fmt.Sprintf(":%s", cfg.Port)); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}
