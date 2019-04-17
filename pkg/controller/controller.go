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

// This example is about k8s Deployment, but it could be any kubernetes resource e.g statefulset, ingress, secrets, configmaps etc.
// Controller establishes a basic bare bone structure to implement more complex logic to control and respond to kubernetes
// object lifecycle events.

const controllerAgentName = "ControllerTemplate"

const (
	SuccessSynced         = "Synced"
	ErrResourceExists     = "ErrResourceExists"
	MessageResourceExists = "Resource %q already exists and is not managed by <controllerAgentName> controller"
	MessageResourceSynced = "Resource<name> synced successfully"
)

type Controller struct {
	name             string
	namespace        string
	kubeclientset    kubernetes.Interface
	deploymentLister appslisters.DeploymentLister
	deploymentSynced cache.InformerSynced
	workqueue        workqueue.RateLimitingInterface
	recorder         record.EventRecorder
}
type Event struct {
	key          string
	eventType    string
	namespace    string
	resourceType string
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
		AddFunc:    controller.Create,
		UpdateFunc: controller.Update,
		DeleteFunc: controller.Delete,
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

func (c *Controller) Create(obj interface{}) {
	var event Event
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		return
	}
	event.key = key
	event.eventType = "create"
	c.enqueue(event)
}

func (c *Controller) Delete(obj interface{}) {
	var event Event
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		return
	}
	event.key = key
	event.eventType = "delete"
	c.enqueue(event)
}
func (c *Controller) Update(old, new interface{}) {
	var event Event
	newDep := new.(*appsv1.Deployment)
	oldDep := new.(*appsv1.Deployment)
	if newDep.ResourceVersion == oldDep.ResourceVersion {
		return
	}
	key, err := cache.MetaNamespaceKeyFunc(new)
	if err != nil {
		return
	}
	event.key = key
	event.eventType = "update"
	c.enqueue(event)

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
		var event Event
		var ok bool
		if event, ok = obj.(Event); !ok {
			c.workqueue.Forget(event)
			runtime.HandleError(fmt.Errorf("expected event (of type Event) in workqueue but got %#v", event))
			return nil
		}
		if err := c.takeAction(event); err != nil {
			return fmt.Errorf("error syncing '%s': %s", event, err.Error())
		}
		c.workqueue.Forget(event)
		glog.Infof("Successfully synced '%s'", event)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

func (c *Controller) takeAction(event Event) error {

	switch event.eventType {
	case "create":
		c.ObjectCreated(event.key)
	case "update":
		c.ObjectUpdated(event.key)
	case "delete":
		c.ObjectDeleted(event.key)
	}
	return nil
}

// ObjectCreated performs logic to create event.
func (c *Controller) ObjectCreated(key string) error {
	return nil
}

// ObjectUpdated performs logic to create event.
func (c *Controller) ObjectUpdated(key string) error {
	return nil
}

// ObjectDeleted performs logic to create event.
func (c *Controller) ObjectDeleted(key string) error {
	return nil
}
