package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/thoughtgears/run-service-discovery/pkg/db"
)

func UpdateService(firestoreDB *db.FirestoreDB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		name := ctx.Param("name")

		var service db.Service
		if err := ctx.ShouldBindJSON(&service); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		service.ID = db.SetID(name)

		existingService, err := firestoreDB.GetService(ctx, service.ID)
		if err != nil || existingService == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
			return
		}

		if err := firestoreDB.UpdateServiceURL(ctx, service.ID, service.URL); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Service updated successfully"})
	}
}
