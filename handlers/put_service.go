package handlers

import (
	"fmt"
	"net/http"

	db2 "github.com/thoughtgears/run-service-discovery/db"

	"github.com/gin-gonic/gin"
)

func UpdateService(firestoreDB *db2.FirestoreDB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		name := ctx.Param("name")

		var service db2.Service
		if err := ctx.ShouldBindJSON(&service); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := service.ValidateURL(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		service.ID = db2.SetID(fmt.Sprintf("%s-%s", name, service.Environment))

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
