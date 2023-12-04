package controller

import (
	//	"github.com/rancher/wrangler/pkg/crd"
	corev1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"

	//	g8sv1alpha1 "github.com/the-gizmo-dojo/g8s/pkg/apis/api.g8s.io/v1alpha1"
	clientset "github.com/the-gizmo-dojo/g8s/pkg/generated/clientset/versioned"
	g8sscheme "github.com/the-gizmo-dojo/g8s/pkg/generated/clientset/versioned/scheme"
	informers "github.com/the-gizmo-dojo/g8s/pkg/generated/informers/externalversions/api.g8s.io/v1alpha1"
	listers "github.com/the-gizmo-dojo/g8s/pkg/generated/listers/api.g8s.io/v1alpha1"
)

const controllerAgentName = "g8s-controller"

type Controller struct {
	kubeClientset kubernetes.Interface
	g8sClientset  clientset.Interface

	passwordLister  listers.PasswordLister
	passwordsSynced cache.InformerSynced

	workqueue workqueue.RateLimitingInterface

	recorder record.EventRecorder
}

func NewController(
	kubeClientset kubernetes.Interface,
	g8sClientset clientset.Interface,
	passwordInformer informers.PasswordInformer) *Controller {

	utilruntime.Must(g8sscheme.AddToScheme(scheme.Scheme))

	klog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		kubeClientset:   kubeClientset,
		g8sClientset:    g8sClientset,
		passwordLister:  passwordInformer.Lister(),
		passwordsSynced: passwordInformer.Informer().HasSynced,
		workqueue:       workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Ats"),
		recorder:        recorder,
	}

	klog.Info("Setting up event handlers")
	passwordInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueuePassword,
		UpdateFunc: func(old, new any) {
			controller.enqueuePassword(new)
		},
	})

	return controller
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
