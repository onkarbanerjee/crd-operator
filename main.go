package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/onkarbanerjee/crd-operator/controller"
	"github.com/onkarbanerjee/crd-operator/handler"
	"github.com/onkarbanerjee/crd-operator/pkg/client/clientset/versioned"
	v1 "github.com/onkarbanerjee/crd-operator/pkg/client/informers/externalversions/customconfig/v1"
	log "github.com/sirupsen/logrus"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
)

// retrieve the Kubernetes cluster client from outside of the cluster
func getKubernetesClient() (kubernetes.Interface, versioned.Interface) {
	// construct the path to resolve to `~/.kube/config`
	kubeConfigPath := os.Getenv("HOME") + "/.kube/config"

	// create the config from the path
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		log.Error("getClusterConfig: %v", err)
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Fatalf("getClusterConfig: %v", err)
		}
	}

	// generate the client based off of the config
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("getClusterConfig: %v", err)
	}

	customconfigClient, err := versioned.NewForConfig(config)
	if err != nil {
		log.Fatalf("getClusterConfig: %v", err)
	}

	log.Info("Successfully constructed k8s client")
	return client, customconfigClient
}

// main code path
func main() {
	// get the Kubernetes client for connectivity
	client, customconfigClient := getKubernetesClient()

	// retrieve our custom resource informer which was generated from
	// the code generator and pass it the custom resource client, specifying
	// we should be looking through all namespaces for listing and watching
	informer := v1.NewCustomConfigInformer(
		customconfigClient,
		meta_v1.NamespaceAll,
		0,
		cache.Indexers{},
	)

	// create a new queue so that when the informer gets a resource that is either
	// a result of listing or watching, we can add an idenfitying key to the queue
	// so that it can be handled in the handler
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	deletedItems := &controller.DeletedItems{
		M: map[string]interface{}{},
	}
	// add event handlers to handle the three types of events for resources:
	//  - adding new resources
	//  - updating existing resources
	//  - deleting resources
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// convert the resource object into a key (in this case
			// we are just doing it in the format of 'namespace/name')
			key, err := cache.MetaNamespaceKeyFunc(obj)
			log.Infof("Add customconfig: %s", key)
			i := &controller.Item{
				Key:        key,
				Event_type: controller.CREATED,
			}
			if err == nil {
				// add the key to the queue for the handler to get
				queue.Add(i)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(newObj)
			log.Infof("Update customconfig: %s", key)
			i := &controller.Item{
				Key:        key,
				Event_type: controller.UPDATED,
			}
			if err == nil {
				queue.Add(i)
			}
		},
		DeleteFunc: func(obj interface{}) {
			// DeletionHandlingMetaNamsespaceKeyFunc is a helper function that allows
			// us to check the DeletedFinalStateUnknown existence in the event that
			// a resource was deleted but it is still contained in the index
			//
			// this then in turn calls MetaNamespaceKeyFunc
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			log.Infof("Delete customconfig: %s", key)
			deletedItems.Add(key, obj)
			i := &controller.Item{
				Key:        key,
				Event_type: controller.DELETED,
			}
			if err == nil {
				queue.Add(i)
			}
		},
	})

	// construct the Controller object which has all of the necessary components to
	// handle logging, connections, informing (listing and watching), the queue,
	// and the handler
	// controller := controller.Controller{
	// 	logger:    log.NewEntry(log.New()),
	// 	clientset: client,
	// 	informer:  informer,
	// 	queue:     queue,
	// 	handler: &handlers.CCHandler{
	// 		Client: client,
	// 	},
	// }

	ccHandler := &handler.CCHandler{
		Client: client,
	}

	ccController := controller.New("custom-config-controller", client, informer, queue, ccHandler, deletedItems)
	// use a channel to synchronize the finalization for a graceful shutdown
	stopCh := make(chan struct{})
	defer close(stopCh)

	// run the controller loop to process items
	go ccController.Run(stopCh)

	// use a channel to handle OS signals to terminate and gracefully shut
	// down processing
	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, syscall.SIGTERM)
	signal.Notify(sigTerm, syscall.SIGINT)
	<-sigTerm
}
