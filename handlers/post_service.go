package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thoughtgears/run-service-discovery/db"
)

// PostService registers a new service in Firestore
// If the service already exists, an error will be returned
// The service requires a name and a URL
// The ID is a sha256 hash of the service name
func PostService(firestoreDB *db.FirestoreDB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var service db.Service
		if err := ctx.ShouldBindJSON(&service); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := service.Validate(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		service.ID = db.SetID(fmt.Sprintf("%s-%s", service.Name, service.Environment))

		existingService, err := firestoreDB.GetService(ctx, service.ID)
		if err == nil && existingService != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Service already exists"})
			return
		}

		if err := firestoreDB.SetService(ctx, service.ID, service); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "Service registered successfully"})
	}
}
