package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/thoughtgears/run-service-discovery/pkg/db"
)

// GetService retrieves a service from Firestore by its name
// If the service does not exist, an error will be returned
func GetService(firestoreDB *db.FirestoreDB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		name := ctx.Param("name")

		ID := db.SetID(name)

		service, err := firestoreDB.GetService(ctx, ID)
		if err != nil || service == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
			return
		}

		ctx.JSON(http.StatusOK, service)
	}
}
