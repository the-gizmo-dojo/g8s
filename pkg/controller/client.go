package controller

import (
	clientset "github.com/the-gizmo-dojo/g8s/pkg/generated/clientset/versioned"
	informers "github.com/the-gizmo-dojo/g8s/pkg/generated/informers/externalversions/api.g8s.io/v1alpha1"
	listers "github.com/the-gizmo-dojo/g8s/pkg/generated/listers/api.g8s.io/v1alpha1"
	secretinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	secretslisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
)

type Client struct {
	// kubeclientset is a standard kubernetes clientset
	kubeClientset kubernetes.Interface

	// g8sclientset is a clientset for our own API group
	g8sClientset clientset.Interface

	// Informers for each type, just to expose for access
	passwordInformer informers.PasswordInformer
	secretInformer   secretinformers.SecretInformer

	// listers for our custom type
	passwordsLister listers.PasswordLister
	passwordsSynced cache.InformerSynced

	// listers for Secret types owned by our custom Password types
	secretsLister secretslisters.SecretLister
	secretsSynced cache.InformerSynced

	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder
}
