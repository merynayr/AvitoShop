package shop

// // Get - отправляет запрос в сервисный слой на получение данных о пингах
// func (api *API) Get(ctx *gin.Context) {
// 	shopObj, err := api.shopService.Get(ctx)
// 	if err != nil {
// 		_ = ctx.Error(err)
// 		ctx.JSON(int(codes.InternalServerError), gin.H{
// 			"error": fmt.Sprintf("Failed to fetch shop data: %v", err),
// 		})
// 		return
// 	}

// 	ctx.JSON(int(codes.OK), gin.H{
// 		"data": shopObj,
// 	})
// }
