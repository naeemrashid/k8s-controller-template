package controller

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	appslisters "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
)

const controllerAgentName = "ControllerTemplate"

const (
	SuccessSynced         = "Synced"
	ErrResourceExists     = "ErrResourceExists"
	MessageResourceExists = "Resource %q already exists and is not managed by <controllerAgentName> controller"
	MessageResourceSynced = "Resource<name> synced successfully"
)

// This example is about k8s Deployment, but it could be any kubernetes resource e.g statefulset, ingress, secrets, configmaps etc.
// Controller establishes a basic bare bone structure to implement more complex logic to control and respond to kubernetes
// object lifecycle events.
type Controller struct {
	name             string
	namespace        string
	kubeclientset    kubernetes.Interface
	deploymentLister appslisters.DeploymentLister
	deploymentSynced cache.InformerSynced
	workqueue        workqueue.RateLimitingInterface
	recorder         record.EventRecorder
}

func NewController(
	kubeclientset kubernetes.Interface,
	deploymentInformer appsinformers.DeploymentInformer,

) *Controller {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})
	controller := &Controller{
		name:             controllerAgentName,
		namespace:        "default",
		kubeclientset:    kubeclientset,
		deploymentLister: deploymentInformer.Lister(),
		deploymentSynced: deploymentInformer.Informer().HasSynced,
		workqueue:        workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
		recorder:         recorder,
	}
	glog.Info("Setting up event handlers")
	deploymentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// Object is newly created.
			// Perform custom logic
			controller.createEvent(obj)
		},
		UpdateFunc: func(old, new interface{}) {
			newDep := new.(*appsv1.Deployment)
			oldDep := new.(*appsv1.Deployment)
			if newDep.ResourceVersion == oldDep.ResourceVersion {
				return
			}
			// Object is been updated.
			// Perform custom logic
			controller.updateEvent(new)
		},
		DeleteFunc: func(obj interface{}) {
			// object is been deleted.
			// Perform custom logic
			controller.deleteEvent(obj)
		},
	})
	return controller
}
func (c *Controller) enqueue(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	c.workqueue.AddRateLimited(key)
}

// Perform custom logic on object creation event.
func (c *Controller) createEvent(obj interface{}) {

	// enqueue object after performing operations.
	c.enqueue(obj)
}

// Perform custom logic on object delete event.
func (c *Controller) deleteEvent(obj interface{}) {

	// enqueue object after performing operations.
	c.enqueue(obj)
}

// perform custom logic on object update event.
func (c *Controller) updateEvent(obj interface{}) {

	// enqueue object after performing operations.
	c.enqueue(obj)
}

func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()
	if ok := cache.WaitForCacheSync(stopCh, c.deploymentSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}
	<-stopCh
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
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}
		c.workqueue.Forget(obj)
		glog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

// Sync resource, this function should be called after update, delete, create event.
func (c *Controller) syncHandler(key string) error {
	return nil
}
