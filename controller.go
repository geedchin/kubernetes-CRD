package main

import (
	"k8s.io/client-go/kubernetes"
	"crd-demo1/pkg/client/clientset/versioned"
	"crd-demo1/pkg/client/informers/externalversions/samplecrd/v1"
	v12 "crd-demo1/pkg/client/listers/samplecrd/v1"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/client-go/tools/record"
	"k8s.io/apimachinery/pkg/util/runtime"
	networkScheme "crd-demo1/pkg/client/clientset/versioned/scheme"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog/glog"
	v13 "k8s.io/client-go/kubernetes/typed/core/v1"
	v14 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
	"fmt"
	"k8s.io/apimachinery/pkg/util/wait"
	"time"
	"k8s.io/apimachinery/pkg/api/errors"
)

const (
	controllerAgentName   = "network-controller"
	SuccessSynced         = "Synced"
	MessageResourceSynced = "Network synced successfully"
)

type Controller struct {
	kubeClient     kubernetes.Interface
	networkClient  versioned.Interface
	networkLister  v12.NetworkLister
	networksSynced cache.InformerSynced
	workqueue      workqueue.RateLimitingInterface
	recorder       record.EventRecorder
}

func NewController(
	kubeClientset kubernetes.Interface,
	networkClientset versioned.Interface,
	networkInformer v1.NetworkInformer) *Controller {

	runtime.Must(networkScheme.AddToScheme(scheme.Scheme))
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&v13.EventSinkImpl{
		Interface: kubeClientset.CoreV1().Events(""),
	})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, v14.EventSource{
		Component: controllerAgentName,
	})
	controller := &Controller{
		kubeClient:     kubeClientset,
		networkClient:  networkClientset,
		networkLister:  networkInformer.Lister(),
		networksSynced: networkInformer.Informer().HasSynced,
		workqueue:      workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Networks"),
		recorder:       recorder,
	}

	networkInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			fmt.Println("add")
			controller.enqueueNetwork(obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			fmt.Println("update")
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println("del")
			controller.enqueueNetworkForDelete(obj)
		},
	})
	return controller
}

func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {

	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()

	if ok := cache.WaitForCacheSync(stopCh, c.networksSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	glog.Info("started workers")
	<-stopCh
	glog.Info("shutting down workers")
	return nil
}

func (c *Controller) runWorker() {
	for c.processNextWorkItem() {

	}
}
func (c *Controller) processNextWorkItem() bool {

	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}
	err := func(obj interface{}) error {
		defer c.workqueue.Done(obj)
		var key string
		var ok bool

		if key, ok = obj.(string); !ok {
			c.workqueue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		if err := c.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing %s : %s", key, err.Error())
		}
		c.workqueue.Forget(obj)
		glog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)
	if err != nil {
		runtime.HandleError(err)
	}
	return true
}
func (c *Controller) syncHandler(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)

	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}
	network, err := c.networkLister.Networks(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			glog.Warningf("Network: %s/%s does not exist in local cache, will delete it from Neutron ...",
				namespace, name)

			glog.Infof("[Neutron] Deleting network: %s/%s ...", namespace, name)

			// FIX ME: call Neutron API to delete this network by name.
			//
			// neutron.Delete(namespace, name)
			return nil
		}

		runtime.HandleError(fmt.Errorf("failed to list network by: %s/%s", namespace, name))
		return err
	}
	glog.Infof("[Neutron] Try to process network: %#v ...", network)

	// FIX ME: Do diff().
	//
	// actualNetwork, exists := neutron.Get(namespace, name)
	//
	// if !exists {
	// 	neutron.Create(namespace, name)
	// } else if !reflect.DeepEqual(actualNetwork, network) {
	// 	neutron.Update(namespace, name)
	// }

	c.recorder.Event(network, v14.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

func (c *Controller) enqueueNetwork(obj interface{}) {

	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
		return
	}
	c.workqueue.AddRateLimited(key)
}

func (c *Controller) enqueueNetworkForDelete(obj interface{}) {

	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
	}
	c.workqueue.AddRateLimited(key)
}
