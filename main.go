package main

import (
	"flag"
	"k8s.io/client-go/tools/clientcmd"
	"fmt"
	"k8s.io/client-go/kubernetes"

	"crd-demo1/pkg/client/clientset/versioned"
	"crd-demo1/pkg/client/informers/externalversions"
	"crd-demo1/pkg/signals"
	"k8s.io/klog/glog"
)

var (
	masterUrl   = "10.1.1.100:6443"
	kubeConfig  = "conf/conf.yaml"
	kubeCfgPath = "conf/conf.yaml"
)

func main() {

	flag.Parse()

	stopCh := signals.SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeCfgPath)
	if err != nil {
		fmt.Println(err)
		//return
	}
	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	networkClient, err := versioned.NewForConfig(cfg)

	networkInformerFactory := externalversions.NewSharedInformerFactory(networkClient, 0)

	controller := NewController(
		kubeClient,
		networkClient,
		networkInformerFactory.Samplecrd().V1().Networks())

	go networkInformerFactory.Start(stopCh)

	if err = controller.Run(2, stopCh); err != nil {
		glog.Fatal("Error running controller: %s", err.Error())
	}
}

func init() {
	flag.StringVar(&kubeConfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterUrl, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}
