package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	kubeClient "gs-k8s/internal/k8s"
	"gs-k8s/internal/pkg"
	"net/http"
)

func Route(port int) {
	// 创建Gin路由
	router := gin.Default()
	cli := kubeClient.KubeConf()
	// 定义处理POST请求的路由
	router.GET("/gstrain/service", func(c *gin.Context) {
		// 调用cli.GetAllServiceName方法
		serviceNames := cli.GetAllServiceName()

		// 返回调用结果
		c.JSON(http.StatusOK, gin.H{"serviceNames": serviceNames})
	})
	router.GET("/gstrain/getconf", func(c *gin.Context) {
		getConfig := pkg.GetConfig()
		// 返回调用结果
		c.JSON(http.StatusOK, getConfig)
	})

	router.POST("/gstrain/updateconf", func(c *gin.Context) {
		var requestData *pkg.Bootstrap
		config := pkg.LoadConfig()
		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 获取servicename
		//modifiedConfig, exists := requestData["policyList"]
		//if !exists {
		//	c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'modifiedPrefix' in JSON"})
		//	return
		//}
		pkg.UpdateConfig(config, requestData)

		c.JSON(http.StatusOK, gin.H{"StatusOK": http.StatusOK})
	})
	router.POST("/gstrain/version", func(c *gin.Context) {
		var requestData map[string]string

		// 解析JSON数据
		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 获取servicename
		serviceName, exists := requestData["servicename"]
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'servicename' in JSON"})
			return
		}

		// 调用cli.GetImageTag方法
		imageTag := cli.GetImageTag(serviceName)

		// 返回调用结果
		c.JSON(http.StatusOK, gin.H{"imageTag": imageTag})
	})

	// 启动HTTP服务，监听8080端口
	addr := fmt.Sprintf(":%d", port)
	err := router.Run(addr)
	if err != nil {
		fmt.Println("Failed to start the server:", err)
	}
}
