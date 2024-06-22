package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thoughtgears/run-service-discovery/db"
)

// GetServices retrieves a list of services from Firestore
func GetServices(firestoreDB *db.FirestoreDB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		services, err := firestoreDB.GetServices(ctx)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err})
			return
		}

		ctx.JSON(http.StatusOK, services)
	}
}
