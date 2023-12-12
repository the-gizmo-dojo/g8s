package controller

import (
	"context"
	"time"

	"golang.org/x/time/rate"

	//	"github.com/rancher/wrangler/pkg/crd"
	corev1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"

	secretinformers "k8s.io/client-go/informers/core/v1"

	//	g8sv1alpha1 "github.com/the-gizmo-dojo/g8s/pkg/apis/api.g8s.io/v1alpha1"
	clientset "github.com/the-gizmo-dojo/g8s/pkg/generated/clientset/versioned"
	g8sscheme "github.com/the-gizmo-dojo/g8s/pkg/generated/clientset/versioned/scheme"
	informers "github.com/the-gizmo-dojo/g8s/pkg/generated/informers/externalversions/api.g8s.io/v1alpha1"
	// listers "github.com/the-gizmo-dojo/g8s/pkg/generated/listers/api.g8s.io/v1alpha1"
)

const controllerAgentName = "g8s-controller"

// Controller is the controller implementation for Password resources
type Controller struct {
	Client
	Executor
}

// NewController returns a new g8s controller
func NewController(
	ctx context.Context,
	kubeClientset kubernetes.Interface,
	g8sClientset clientset.Interface,
	passwordInformer informers.PasswordInformer,
	secretInformer secretinformers.SecretInformer) *Controller {

	logger := klog.FromContext(ctx)

	// Create event broadcaster
	// Add g8s types to the default Kubernetes Scheme so Events can be
	// logged for g8s types.
	utilruntime.Must(g8sscheme.AddToScheme(scheme.Scheme))
	logger.V(4).Info("Creating event broadcaster")

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartStructuredLogging(0)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	ratelimiter := workqueue.NewMaxOfRateLimiter(
		workqueue.NewItemExponentialFailureRateLimiter(5*time.Millisecond, 1000*time.Second),
		&workqueue.BucketRateLimiter{Limiter: rate.NewLimiter(rate.Limit(50), 300)},
	)

	controller := &Controller{
		Client: Client{
			kubeClientset:    kubeClientset,
			g8sClientset:     g8sClientset,
			passwordInformer: passwordInformer,
			secretInformer:   secretInformer,
			passwordsLister:  passwordInformer.Lister(),
			passwordsSynced:  passwordInformer.Informer().HasSynced,
			secretsLister:    secretInformer.Lister(),
			secretsSynced:    secretInformer.Informer().HasSynced,
			recorder:         recorder,
		},
		Executor: Executor{
			workqueue: workqueue.NewNamedRateLimitingQueue(ratelimiter, "Password"),
		},
	}

	logger.Info("Setting up event handlers")

	// Set up an event handler for when Password resources change
	passwordInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueuePassword,
		UpdateFunc: func(old, new any) {
			controller.enqueuePassword(new)
		},
	})

	// Set up an event handler for when Secret resources change. This
	// handler will lookup the owner of the given Secret, and if it is
	// owned by a Password resource then the handler will enqueue that Password resource for
	// processing. This way, we don't need to implement custom logic for
	// handling Secret resources. More info on this pattern:
	// https://github.com/kubernetes/community/blob/8cafef897a22026d42f5e5bb3f104febe7e29830/contributors/devel/controllers.md
	secretInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.handleObject,
		UpdateFunc: func(old, new interface{}) {
			newDepl := new.(*corev1.Secret)
			oldDepl := old.(*corev1.Secret)
			if newDepl.ResourceVersion == oldDepl.ResourceVersion {
				// Periodic resync will send update events for all known Secrets.
				// Two different versions of the same Secret will always have different ResourceVersions.
				// This section will skip calling handleObject() if they are the same.
				return
			}
			controller.handleObject(new)
		},
		DeleteFunc: controller.handleObject,
	})

	return controller
}
