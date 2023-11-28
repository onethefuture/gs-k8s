package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	apollo "github.com/philchia/agollo/v4"
	"github.com/philchia/agollo/v4/openapi"
	"log"
	"sync"
)

type Bootstrap struct {
	PolicyList []struct {
		BizOperation     string `json:"biz_operation"`
		EnginePolicyList []struct {
			Platform   string  `json:"platform"`
			AppVersion string  `json:"app_version"`
			Version    string  `json:"version"`
			Weight     float64 `json:"weight"`
		} `json:"enginePolicyList"`
	} `json:"policyList"`
}

//type Bootstrap struct {
//	Srvs *Srvs `json:"services"`
//}
//
//type Srvs struct {
//	Services []*Service `json:"srvs"`
//}
//
//type Service struct {
//	ID    string   `json:"id"`
//	Email string   `json:"email"`
//	Name  []string `json:"name"`
//}

var (
	boot        Bootstrap
	configMutex = &sync.Mutex{}
)

// 加载apollo配置
func LoadConfig() *Bootstrap {
	apollo.Start(&apollo.Conf{
		AppID:          "ycloud",
		Cluster:        "DEV",
		NameSpaceNames: []string{NamespaceName},
		MetaAddr:       "http://10.177.9.244:37237",
	})
	apollo.OnUpdate(func(event *apollo.ChangeEvent) {
		// 监听配置变更
		log.Printf("Event: %#v\n", event)
		updateConfigFromApollo()
	})
	log.Printf("初始化Apollo配置成功")

	// Load initial configuration from Apollo
	updateConfigFromApollo()

	return &boot
}

func updateConfigFromApollo() {
	configMutex.Lock()
	defer configMutex.Unlock()

	namespaceContent := apollo.GetContent(apollo.WithNamespace(NamespaceName))

	//var newConfig Bootstrap
	var boot Bootstrap
	err := yaml.Unmarshal([]byte(namespaceContent), &boot)
	if err != nil {
		log.Printf("Error unmarshaling JSON: %s", err)
		return
	}

	// Update the configuration
	//boot = newConfig

	buf, _ := json.Marshal(boot)
	//fmt.Println("buf: ", buf)
	log.Printf("Apollo configs: %s", string(buf))
}

// UpdateConfig updates the Apollo configuration dynamically
func UpdateConfig(newConfig, modifiedConfig *Bootstrap) {
	configMutex.Lock()
	defer configMutex.Unlock()

	// Update the configuration
	//boot = newConfig
	//
	//// Update the Apollo configuration
	//apollo.UpdateFileContent(apollo.WithNamespace("config.yaml"), string(newConfigYAML))

	newConfig = modifiedConfig

	newConfigYAML, err := yaml.Marshal(newConfig)
	if err != nil {
		log.Printf("Error marshaling YAML: %s", err)
		return
	}

	c := openapi.New("http://10.177.9.244:34262", Appid, Env, Cluster, Token)
	GetRelease, err := c.GetRelease(NamespaceName)
	if err != nil {
		log.Printf("Error getApolloConfig: %s", err)
	}
	//s := GetRelease.Configurations["config"]

	updateItemRequest := c.UpdateConfig(NamespaceName, "content", string(newConfigYAML), ReleaseComment, DataChangeLastModifiedBy)
	log.Println("updateItemRequest: ", updateItemRequest)
	log.Println("updateConfig: ", GetRelease)
	c.Release(NamespaceName, "test", ReleaseComment, DataChangeLastModifiedBy)

}

func PrintConfiguration(config Bootstrap) {
	fmt.Printf("Services:\n")
	for _, service := range config.PolicyList {
		fmt.Printf("  BizOperation: %s\n", service.BizOperation)
		fmt.Printf("  E: %s\n", service.EnginePolicyList)
	}
}

func GetConfig() *Bootstrap {
	apollo.Start(&apollo.Conf{
		AppID:          "ycloud",
		Cluster:        "DEV",
		NameSpaceNames: []string{NamespaceName},
		MetaAddr:       "http://10.177.9.244:37237",
	})
	apollo.OnUpdate(func(event *apollo.ChangeEvent) {
		// 监听配置变更
		log.Printf("Event: %#v\n", event)
		updateConfigFromApollo()
	})
	log.Printf("初始化Apollo配置成功")

	// Load initial configuration from Apollo
	//updateConfigFromApollo()
	namespaceContent := apollo.GetContent(apollo.WithNamespace(NamespaceName))

	//var newConfig Bootstrap
	//var boot Bootstrap
	err := yaml.Unmarshal([]byte(namespaceContent), &boot)
	if err != nil {
		log.Printf("Error unmarshaling JSON: %s", err)
		return &boot
	}

	// Update the configuration
	//boot = newConfig

	buf, _ := json.Marshal(boot)
	//fmt.Println("buf: ", buf)
	log.Printf("Apollo configs: %s", string(buf))
	return &boot
}
