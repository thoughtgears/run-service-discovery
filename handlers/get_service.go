package handlers

import (
	"fmt"
	"net/http"

	"github.com/thoughtgears/run-service-discovery/db"

	"github.com/gin-gonic/gin"
)

// GetService retrieves a service from Firestore by its name
// If the service does not exist, an error will be returned
func GetService(firestoreDB *db.FirestoreDB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		name := ctx.Param("name")
		env := ctx.Query("environment")

		if name == "" || env == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "/service/:name?environment=env is required"})
			return
		}

		ID := db.SetID(fmt.Sprintf("%s-%s", name, env))

		service, err := firestoreDB.GetService(ctx, ID)
		if err != nil || service == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
			return
		}

		ctx.JSON(http.StatusOK, service)
	}
}
