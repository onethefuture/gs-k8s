package main

import (
	"flag"
	"gs-k8s/internal/service"
)

func main() {
	port := flag.Int("p", 8080, "HTTP server port")
	flag.Parse()
	//cli := kubeClient.KubeConf()
	//cli.GetImageTag("spider-man")
	service.Route(*port)
}
