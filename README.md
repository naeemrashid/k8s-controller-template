# k8s-controller-template
A template to quick-start with kubernetes controller/operator implementation or more precisely kubernetes watcher to watch for changes in resources and apply custom logic. Template does not introduces new CRDs(Custom Resource Definitions) but that is intentional, have a look at kubernetes [sample-controller](https://github.com/kubernetes/sample-controller) for refrences.
### Specify type of resources to watch for
Below is a watch loop for changes in Deployment resources
```
deploymentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.Create,
		UpdateFunc: controller.Update,
		DeleteFunc: controller.Delete,
	})
```
### Implement your custom logic Here

```
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

```
### Build and Run

```
// vendor the project using
go get github.com/thenaeem/k8s-controller-template
// Build container
make container
// Deploy to Kuberentes as a single replica of Deployment.
```
