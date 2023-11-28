package kubeClient

import (
	"flag"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func KubeConf() *KubeClient {
	var kubeConf *string
	if home := homedir.HomeDir(); len(home) != 0 {
		kubeConf = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeConf = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	flag.Parse()

	cli, err := NewKubeClient(kubeConf)
	if err != nil {
		return nil
	}
	return cli
}
