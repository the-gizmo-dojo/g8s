package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/the-gizmo-dojo/g8s/pkg/apis/api.g8s.io/v1alpha1"
	"github.com/the-gizmo-dojo/g8s/pkg/g8s"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

const (
	// SuccessSynced is used as part of the Event 'reason' when a Password is synced
	SuccessSynced = "Synced"
	// ErrResourceExists is used as part of the Event 'reason' when a Password fails
	// to sync due to a Secret of the same name already existing.
	ErrResourceExists = "ErrResourceExists"

	// MessageResourceExists is the message used for Events when a resource
	// fails to sync due to a Secret already existing
	MessageResourceExists = "Resource %q already exists and is not managed by Password"
	// MessageResourceSynced is the message used for an Event fired when a Password
	// is synced successfully
	MessageResourceSynced = "Password synced successfully"
)

type Executor struct {
	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workqueue workqueue.RateLimitingInterface
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(ctx context.Context, workers int) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()
	logger := klog.FromContext(ctx)

	// Start the informer factories to begin populating the informer caches
	logger.Info("Starting Password controller")

	// Wait for the caches to be synced before starting workers
	logger.Info("Waiting for informer caches to sync")

	if ok := cache.WaitForCacheSync(ctx.Done(), c.passwordsSynced, c.secretsSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	logger.Info("Starting workers", "count", workers)
	// Launch two workers to process At resources
	for i := 0; i < workers; i++ {
		go wait.UntilWithContext(ctx, c.runWorker, time.Second)
	}

	logger.Info("Started workers")
	<-ctx.Done()
	logger.Info("Shutting down workers")

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *Controller) runWorker(ctx context.Context) {
	for c.processNextWorkItem(ctx) {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextWorkItem(ctx context.Context) bool {
	obj, shutdown := c.workqueue.Get()
	logger := klog.FromContext(ctx)

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// Password resource to be synced.
		if err := c.syncHandler(ctx, key); err != nil {
			// Put the item back on the workqueue to handle any transient errors.
			c.workqueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(obj)
		logger.Info("Successfully synced", "resourceName", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the Password resource
// with the current status of the resource.
func (c *Controller) syncHandler(ctx context.Context, key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	logger := klog.LoggerWithValues(klog.FromContext(ctx), "resourceName", key)

	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the Password resource with this namespace/name
	password, err := c.passwordsLister.Passwords(namespace).Get(name)
	if err != nil {
		// The Password resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("password '%s' in work queue no longer exists", key))
			return nil
		}

		return err
	}

	backendName := password.ObjectMeta.Name
	historyName := password.ObjectMeta.Name + "-history"
	// Get the backend Secret and history Secret with the name specified in Password.ObjectMeta.Name
	backend, berr := c.secretsLister.Secrets(password.Namespace).Get(backendName)
	history, herr := c.secretsLister.Secrets(password.Namespace).Get(historyName)

	// If the backend and history resources don't exist, create them
	if errors.IsNotFound(berr) && errors.IsNotFound(herr) {
		logger.V(4).Info("Create backend and history Secret resources")
		pwstr := g8s.GeneratePassword(password)
		backend, err = c.Client.kubeClientset.CoreV1().Secrets(password.Namespace).Create(ctx, newBackend(password, pwstr), metav1.CreateOptions{})
		history, err = c.Client.kubeClientset.CoreV1().Secrets(password.Namespace).Create(ctx, newHistory(password, pwstr), metav1.CreateOptions{})
	} else if errors.IsNotFound(berr) { // backend dne but history does, rebuild backend from history
		logger.V(4).Info("Create backend Secret resources from history")
		pwbyte := history.Data["password"]
		backend, err = c.Client.kubeClientset.CoreV1().Secrets(password.Namespace).Create(ctx, newBackend(password, string(pwbyte)), metav1.CreateOptions{})
	} else if errors.IsNotFound(herr) { // backend exists but history dne, rebuild history from backend
		logger.V(4).Info("Create history Secret resources from backend")
		pwbyte := backend.Data["password"]
		history, err = c.Client.kubeClientset.CoreV1().Secrets(password.Namespace).Create(ctx, newHistory(password, string(pwbyte)), metav1.CreateOptions{})
	} else {
		logger.V(4).Info("Secret resources for history and backend exist")
	}

	// If an error occurs during Get/Create, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return err
	}

	// If the Secret is not controlled by this Password resource, we should log
	// a warning to the event recorder and return error msg.
	if !metav1.IsControlledBy(backend, password) {
		msg := fmt.Sprintf(MessageResourceExists, backend.Name)
		c.recorder.Event(password, corev1.EventTypeWarning, ErrResourceExists, msg)
		return fmt.Errorf("%s", msg)
	}

	// Finally, we update the status block of the Password resource to reflect the
	// current state of the world
	err = c.updatePasswordStatus(password, backend)
	fmt.Println("updateStatusErr: ", err)
	if err != nil {
		return err
	}

	c.recorder.Event(password, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

func (c *Controller) updatePasswordStatus(password *v1alpha1.Password, secret *corev1.Secret) error {
	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	passwordCopy := password.DeepCopy()
	passwordCopy.Status.Ready = true
	// If the CustomResourceSubresources feature gate is not enabled,
	// we must use Update instead of UpdateStatus to update the Status block of the Password resource.
	// UpdateStatus will not allow changes to the Spec of the resource,
	// which is ideal for ensuring nothing other than resource status has been updated.
	_, err := c.Client.g8sClientset.ApiV1alpha1().Passwords(password.Namespace).Update(context.TODO(), passwordCopy, metav1.UpdateOptions{})
	return err
}

// enqueuePassword takes a Password resource and converts it into a namespace/name
// string which is then put onto the workqueue. This method should *not* be
// passed resources of any type other than Password.
func (c *Controller) enqueuePassword(obj any) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.workqueue.Add(key)
}

// handleObject will take any resource implementing metav1.Object and attempt
// to find the Password resource that 'owns' it. It does this by looking at the
// objects metadata.ownerReferences field for an appropriate OwnerReference.
// It then enqueues that Password resource to be processed. If the object does not
// have an appropriate OwnerReference, it will simply be skipped.
func (c *Controller) handleObject(obj interface{}) {
	var object metav1.Object
	var ok bool
	logger := klog.FromContext(context.Background())
	if object, ok = obj.(metav1.Object); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding object, invalid type"))
			return
		}
		object, ok = tombstone.Obj.(metav1.Object)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding object tombstone, invalid type"))
			return
		}
		logger.V(4).Info("Recovered deleted object", "resourceName", object.GetName())
	}
	logger.V(4).Info("Processing object", "object", klog.KObj(object))
	if ownerRef := metav1.GetControllerOf(object); ownerRef != nil {
		// If this object is not owned by a Password, we should not do anything more
		// with it.
		if ownerRef.Kind != "Password" {
			return
		}

		password, err := c.passwordsLister.Passwords(object.GetNamespace()).Get(ownerRef.Name)
		if err != nil {
			logger.V(4).Info("Ignore orphaned object", "object", klog.KObj(object), "password", ownerRef.Name)
			return
		}

		c.enqueuePassword(password)
		return
	}
}

// Secret.Immutable requires a *bool, helper func to return that
func boolPtr(b bool) *bool {
	return &b
}

// newSecret creates a new Secret for a Password resource which contains the actual password.
// It also sets the appropriate OwnerReferences on the resource so handleObject can discover
// the Password resource that 'owns' it.
func newBackend(pw *v1alpha1.Password, pwstr string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pw.ObjectMeta.Name,
			Namespace: pw.ObjectMeta.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(pw, v1alpha1.SchemeGroupVersion.WithKind("Password")),
			},
			Annotations: map[string]string{
				"controller": "g8s",
			},
		},
		Immutable: boolPtr(true),
		StringData: map[string]string{
			"password": pwstr,
		},
		Type: "Opaque",
	}
}

// newHistory creates a new Secret for a Password resource which contains the password's history.
// It also sets the appropriate OwnerReferences on the resource so handleObject can discover
// the Password resource that 'owns' it.
func newHistory(pw *v1alpha1.Password, pwstr string) *corev1.Secret {

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pw.ObjectMeta.Name + "-history",
			Namespace: pw.ObjectMeta.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(pw, v1alpha1.SchemeGroupVersion.WithKind("Password")),
			},
			Annotations: map[string]string{
				"controller": "g8s",
			},
		},
		Immutable: boolPtr(true),
		StringData: map[string]string{
			"password": pwstr,
		},
		Type: "Opaque",
	}
}
