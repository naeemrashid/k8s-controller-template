package controller

import (
	"github.com/golang/glog"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
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

// This example is about deployment, but it could be any kubernetes resource e.g statefulset, ingress, secrets, configmaps etc.
// Controller is establish a basic bare bone structure to implement more complex logic to control and respond to kubernetes
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
			// enqueue object for processing
			// controller.enqueue(obj)
		},
		UpdateFunc: func(old, new interface{}) {
			newDep := new.(*appsv1.Deployment)
			oldDep := new.(*appsv1.Deployment)
			if newDep.ResourceVersion == oldDep.ResourceVersion {
				return
			}
			// enqueue object for processing
			// controller.enqueue(new)
		},
		DeleteFunc: func(obj interface{}) {
			// object is being deleted perform custom logic on it
			// controller.enqueue(obj)
		},
	})
	return controller
}
