package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thoughtgears/run-service-discovery/db"
)

// DeleteService deletes a service from Firestore by its name
func DeleteService(firestoreDB *db.FirestoreDB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		name := ctx.Param("name")
		env := ctx.Query("environment")

		ID := db.SetID(fmt.Sprintf("%s-%s", name, env))

		if err := firestoreDB.DeleteService(ctx, ID); err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Error deleting service not found"})
			return
		}

		ctx.Status(http.StatusAccepted)
	}
}
