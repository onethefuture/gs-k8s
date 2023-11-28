package kubeClient

import (
	"context"
	"errors"
	"fmt"
	"gs-k8s/internal/pkg"
	"k8s.io/api/core/v1"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/golang/glog"
	v12 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func NewKubeClient(kubeConf *string) (*KubeClient, error) {
	config, err := clientcmd.BuildConfigFromFlags("", *kubeConf)
	if err != nil {
		glog.Error(err)
		return nil, errors.New(fmt.Sprintf("Error building Flags: %s", err.Error()))
	}

	if config == nil {
		return nil, errors.New("Error building Flags: config is nil")
	}

	//创建链接
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	deploymentsClient := clientset.AppsV1().Deployments(pkg.DefaultServiceName)

	list, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	// 列出所有的Deployments
	for _, d := range list.Items {
		Image := d.Spec.Template.Spec.Containers[0].Image
		//re := regexp.MustCompile("[^:]+?$")
		tagcomit := strings.Split(Image, ":")

		//fmt.Println(Image, tagcomit)
		//ImageTag, _ := regexp.MatchString("[^:]+?$", Image)
		log.Printf(" * %s (imagetag: %s)\n", d.Name, tagcomit[1])
	}

	return &KubeClient{
		cli: clientset,
	}, nil
}

type KubeClient struct {
	cli *kubernetes.Clientset
}

func (k *KubeClient) ExpGetService() *v1.ServiceList {
	serviceClient := k.cli.CoreV1().Services(pkg.DefaultServiceName)
	w, err := serviceClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		glog.Error(err.Error())
	}
	return w
}

func (k *KubeClient) GetAllServiceName() []string {
	w := k.ExpGetService()
	var serviceName []string
	for _, d := range w.Items {
		//glog.Info(d.Name)
		serviceName = append(serviceName, d.Name)
	}
	//fmt.Println("service name []string is: ", serviceName)
	return serviceName
}

func (k *KubeClient) GetServiceName(service string) string {
	w := k.ExpGetService()
	serviceName := make(map[string]string)
	for _, d := range w.Items {
		//glog.Info(d.Name)
		serviceName[d.Name] = d.Name
	}
	//fmt.Println("service name map is: ", serviceName)
	return serviceName[service]
}

func (k *KubeClient) GetImageTag(service string) []string {
	deploymentsClient := k.cli.AppsV1().Deployments(pkg.DefaultServiceName)

	list, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		glog.Error(err.Error())
		return nil
	}

	serviceName := k.GetServiceName(service)
	fmt.Println("The serviceName is: ", serviceName)

	serviceNamePattern := "^" + serviceName
	deployNameRegex := regexp.MustCompile(serviceNamePattern)
	serviceVersion := make(map[string]string)
	for _, d := range list.Items {
		if deployNameRegex.MatchString(d.Name) {
			Image := d.Spec.Template.Spec.Containers[0].Image
			imageTag := strings.Split(Image, ":")
			lastIndex := strings.LastIndex(imageTag[1], "-")
			if lastIndex != -1 {
				serviceVersion[d.Name] = imageTag[1][:lastIndex]
			}
		}
	}
	//ServiceTag := serviceVersion[service]
	var result []string
	for _, deployName := range serviceVersion {
		result = append(result, deployName)
	}

	fmt.Println("service current version is: ", result)
	return result
}

func (k *KubeClient) WatchDeployment() {
	deploymentsClient := k.cli.AppsV1().Deployments(pkg.DefaultServiceName)

	w, err := deploymentsClient.Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		glog.Error(err.Error())
		return
	}
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	//定义一个空字典，来记录最初 Deployment 对应的 Image
	serviceImageTag := make(map[string]string)
	var before string
	list, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		glog.Error(err.Error())
		return
	}
	for _, d := range list.Items {
		serviceImageTag[d.Name] = d.Spec.Template.Spec.Containers[0].Image
	}

	for {
		var sendcall bool
		select {
		case event := <-w.ResultChan():
			switch event.Type {
			case watch.Added, watch.Modified:
				deployment := event.Object.(*v12.Deployment)
				afterImageTag := event.Object.(*v12.Deployment).Spec.Template.Spec.Containers[0].Image
				log.Printf("Deployment : %s , beforeImageTag : %s \n", deployment.Name, afterImageTag)
				//for _, container := range deployment.Spec.Template.Spec.Containers {
				//	fmt.Printf("Image: %s\n", container.Image)
				//}
				before = serviceImageTag[deployment.Name]
				if before != afterImageTag {
					sendcall = true
				}
				if sendcall {
					// 镜像tag部分
					tagcomit := strings.Split(afterImageTag, ":")
					var tag string
					genName := deployment.Name
					if len(tagcomit) > 1 && len(tagcomit[1]) > 1 {
						tag = strings.Split(tagcomit[1], "-")[0]
						genName += tag
					} else {
						tag = tagcomit[1]
						genName += tag
					}
					if err != nil {
						log.Println("获取 commitmes 失败: ", err)
						return
					}
					//message := fmt.Sprintf("【** 通知 **】聚合通知(生产)\n发布服务: %s\n发布时间: %s\n发布版本: %s\n更新内容: %s\n更多详情，请登入 ArgoCD 首页查看\n安格斯地址: %s", deployment.Name, time.Now().Format("2006-01-02 15:04:05.000"), tag, commitmes, "https://argocdtrain.shmao.net/")
					//pkg.SendWarning(deployment.Name, message)
				}
				serviceImageTag[deployment.Name] = afterImageTag
			case watch.Error:
				log.Printf("Error watching K8s: %v", event.Object)
			}
		case <-sigCh:
			w.Stop()
			log.Printf("Stopping...")
			return
		}
	}
}
