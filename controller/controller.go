package controller

import (
	"fmt"
	"sync"
	"time"

	"github.com/onkarbanerjee/crd-operator/handler"
	log "github.com/sirupsen/logrus"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const (
	CREATED Event_type = "created"
	UPDATED Event_type = "updated"
	DELETED Event_type = "deleted"
)

// Event_type denotes whether the event was a create/update/delete event
type Event_type string

type DeletedItems struct {
	sync.RWMutex
	M map[string]interface{}
}

// Item denotes an item being tracked in queue
type Item struct {
	Key string
	Event_type
}

// Controller struct defines how a controller should encapsulate
// logging, client connectivity, informing (list and watching)
// queueing, and handling of resource changes
type Controller struct {
	logger       *log.Entry
	name         string
	clientset    kubernetes.Interface
	queue        workqueue.RateLimitingInterface
	informer     cache.SharedIndexInformer
	handler      handler.Handler
	deletedItems *DeletedItems
}

func (d *DeletedItems) Get(key string) interface{} {
	d.RLock()
	defer d.RUnlock()
	return d.M[key]
}

func (d *DeletedItems) Add(key string, obj interface{}) {
	d.Lock()
	defer d.Unlock()
	d.M[key] = obj
}

func (d *DeletedItems) Delete(key string) {
	d.Lock()
	defer d.Unlock()
	delete(d.M, key)
}

func New(name string, client kubernetes.Interface, informer cache.SharedIndexInformer, queue workqueue.RateLimitingInterface, handler handler.Handler, deletedItems *DeletedItems) *Controller {
	return &Controller{
		logger:       log.NewEntry(log.New()),
		name:         name,
		clientset:    client,
		informer:     informer,
		queue:        queue,
		handler:      handler,
		deletedItems: deletedItems,
	}
}

// Run is the main path of execution for the controller loop
func (c *Controller) Run(stopCh <-chan struct{}) {
	// handle a panic with logging and exiting
	defer utilruntime.HandleCrash()
	// ignore new items in the queue but when all goroutines
	// have completed existing items then shutdown
	defer c.queue.ShutDown()

	c.logger.Info("Controller.Run: initiating")

	// run the informer to start listing and watching resources
	go c.informer.Run(stopCh)

	// do the initial synchronization (one time) to populate resources
	if !cache.WaitForNamedCacheSync(c.name, stopCh, c.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("Error syncing cache"))
		return
	}
	c.logger.Info("Controller.Run: cache sync complete")

	// run the runWorker method every second with a stop channel
	wait.Until(c.runWorker, time.Second, stopCh)
}

// HasSynced allows us to satisfy the Controller interface
// by wiring up the informer's HasSynced method to it
func (c *Controller) HasSynced() bool {
	return c.informer.HasSynced()
}

// runWorker executes the loop to process new items added to the queue
func (c *Controller) runWorker() {
	log.Info("Controller.runWorker: starting")

	// invoke processNextItem to fetch and consume the next change
	// to a watched or listed resource
	for c.processNextItem() {
		log.Info("Controller.runWorker: processing next item")
	}

	log.Info("Controller.runWorker: completed")
}

// processNextItem retrieves each queued item and takes the
// necessary handler action based off of if the item was
// created or deleted
func (c *Controller) processNextItem() bool {
	log.Info("Controller.processNextItem: start")

	// fetch the next item (blocking) from the queue to process or
	// if a shutdown is requested then return out of this to stop
	// processing
	key, quit := c.queue.Get()

	// stop the worker loop from running as this indicates we
	// have sent a shutdown message that the queue has indicated
	// from the Get method
	if quit {
		return false
	}

	defer c.queue.Done(key)

	i, ok := key.(*Item)
	log.Info("i ok is", ok)
	log.Info("i is", i.Key, i.Event_type)

	// assert the string out of the key (format `namespace/name`)
	// keyRaw := key.(string)

	// take the string key and get the object out of the indexer
	//
	// item will contain the complex object for the resource and
	// exists is a bool that'll indicate whether or not the
	// resource was created (true) or deleted (false)
	//
	// if there is an error in getting the key from the index
	// then we want to retry this particular queue key a certain
	// number of times (5 here) before we forget the queue key
	// and throw an error
	item, exists, err := c.informer.GetIndexer().GetByKey(i.Key)

	// log.Info(item.(*v1.CustomConfig))

	log.Info(exists, err)
	if err != nil {
		if c.queue.NumRequeues(key) < 5 {
			c.logger.Errorf("Controller.processNextItem: Failed processing item with key %s with error %v, retrying", key, err)
			c.queue.AddRateLimited(key)
		} else {
			c.logger.Errorf("Controller.processNextItem: Failed processing item with key %s with error %v, no more retries", key, err)
			c.queue.Forget(key)
			utilruntime.HandleError(err)
		}
	}

	switch i.Event_type {
	case CREATED:
		c.logger.Infof("Controller.processNextItem: object created detected: %s", i.Key)
		c.handler.ObjectCreated(item)
		c.queue.Forget(key)
	case UPDATED:
		c.logger.Infof("Controller.processNextItem: object updated detected: %s", i.Key)
		c.handler.ObjectUpdated(item)
		c.queue.Forget(key)
	case DELETED:
		c.logger.Infof("Controller.processNextItem: object deleted detected: %s", i.Key)
		c.handler.ObjectDeleted(c.deletedItems.Get(i.Key))
		c.queue.Forget(key)
	}
	// if the item doesn't exist then it was deleted and we need to fire off the handler's
	// ObjectDeleted method. but if the object does exist that indicates that the object
	// was created (or updated) so run the ObjectCreated method
	//
	// after both instances, we want to forget the key from the queue, as this indicates
	// a code path of successful queue key processing
	// if !exists {
	// 	c.logger.Infof("Controller.processNextItem: object deleted detected: %s", i.key)
	// 	c.handler.ObjectDeleted(item)
	// 	c.queue.Forget(key)
	// } else {
	// 	c.logger.Infof("Controller.processNextItem: object created detected: %s", i.key)
	// 	c.handler.ObjectCreated(item)
	// 	c.queue.Forget(key)
	// }

	// keep the worker loop running by returning true
	return true
}
