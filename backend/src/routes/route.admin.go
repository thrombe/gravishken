package route

import (
	"common/models/admin"
	Batch "common/models/batch"
	Test "common/models/test"
	User "common/models/user"
	"fmt"
	"server/src/controllers"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(allControllers *controllers.ControllerClass, route *gin.Engine) {
	adminRoute := route.Group("/admin")

	adminRoute.POST("/login", func(ctx *gin.Context) {
		var adminModel admin.Admin
		fmt.Println("Admin login route")
		if err := ctx.ShouldBindJSON(&adminModel); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}

		allControllers.AdminLoginHandler(ctx, &adminModel)
	})

	adminRoute.POST("/register", func(ctx *gin.Context) {
		var adminModel admin.Admin
		if err := ctx.ShouldBindJSON(&adminModel); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}

		allControllers.AdminRegisterHandler(ctx, &adminModel)
	})

	// change password not working properly
	adminRoute.POST("/changepassword", func(ctx *gin.Context) {
		var adminModel admin.AdminChangePassword
		if err := ctx.ShouldBindJSON(&adminModel); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}

		allControllers.AdminChangePasswordHandler(ctx, &adminModel)
	})

	adminRoute.POST("/add_all_users", func(ctx *gin.Context) {
		var FilePathRequest struct {
			FilePath string `json:"filePath" binding:"required"`
		}

		if err := ctx.ShouldBindJSON(&FilePathRequest); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}
		allControllers.AddAllUsersBacthesToDb(ctx, FilePathRequest.FilePath)
	})

	adminRoute.POST("/add_batch", func(ctx *gin.Context) {
		var batchData Batch.Batch
		if err := ctx.ShouldBindJSON(&batchData); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}

		allControllers.AddBatchToDB(ctx, &batchData)
	})

	adminRoute.POST("/add_test", func(ctx *gin.Context) {

		var testModel Test.BatchTests

		if err := ctx.ShouldBindJSON(&testModel); err != nil {
			ctx.JSON(500, gin.H{"error": "Invalid request body"})
			return
		}
		fmt.Println("testModel: ", testModel)

		allControllers.AddTestToDB(ctx, &testModel)
	})

	adminRoute.POST("/update_user_data", func(ctx *gin.Context) {
		var userUpdateRequest User.UserUpdateRequest

		if err := ctx.ShouldBindJSON(&userUpdateRequest); err != nil {
			ctx.JSON(500, gin.H{"error": "Invalid request body"})
			return
		}

		allControllers.UpdateUserData(ctx, &userUpdateRequest)

	})

	adminRoute.POST("/increase_test_time", func(ctx *gin.Context) {
		var requestData struct {
			Param          string   `json:"param"`
			Username       []string `json:"username"`
			TimeToIncrease int64    `json:"time_to_increase"`
		}

		// Bind JSON body to requestData struct
		if err := ctx.ShouldBindJSON(&requestData); err != nil {
			ctx.JSON(400, gin.H{
				"message": "Invalid request body",
				"error":   err.Error(),
			})
			return
		}

		allControllers.Increase_Time(ctx, requestData.Param, requestData.Username, requestData.TimeToIncrease)

	})

	adminRoute.POST("/get_batchwise_data", func(ctx *gin.Context) {
		var batchData struct {
			Param       string `json:"param"`
			BatchNumber string `json:"batchNumber"`
			Ranges      []int  `json:"ranges"`
		}

		if err := ctx.ShouldBindJSON(&batchData); err != nil {
			ctx.JSON(500, gin.H{"error": "Invalid request body"})
			return
		}

		allControllers.GetBatchWiseData(ctx, batchData.Param, batchData.BatchNumber, batchData.Ranges)

	})

	adminRoute.POST("/set_user_data", func(ctx *gin.Context) {
		var userRequest struct {
			Username         string `json:"username"`
			Param            string `json:"param"`
			From             int    `json:"from"`
			To               int    `json:"to"`
			ResultDownloaded bool   `json:"resultDownloaded"`
		}

		if err := ctx.ShouldBindJSON(&userRequest); err != nil {
			ctx.JSON(500, gin.H{"error": "Invalid request body"})
			return
		}

		allControllers.SetUserData(ctx, userRequest.Param, &User.UserBatchRequestData{
			From:             userRequest.From,
			To:               userRequest.To,
			ResultDownloaded: userRequest.ResultDownloaded,
		}, userRequest.Username)
	})

	adminRoute.POST("/update_typing_test_text", func(ctx *gin.Context) {
		var UpdateTypingTestTextRequest struct {
			TypingTestText string `json:"typingTestText"`
			TestPassword   string `json:"testPassword"`
		}

		if err := ctx.ShouldBindJSON(&UpdateTypingTestTextRequest); err != nil {
			ctx.JSON(500, gin.H{"error": "Invalid request body"})
			return
		}

		allControllers.UpdateTypingTestText(ctx, UpdateTypingTestTextRequest.TypingTestText, UpdateTypingTestTextRequest.TestPassword)
	})
}
