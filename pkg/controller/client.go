package controller

import (
	clientset "github.com/the-gizmo-dojo/g8s/pkg/generated/clientset/versioned"
	informers "github.com/the-gizmo-dojo/g8s/pkg/generated/informers/externalversions/api.g8s.io/v1alpha1"
	listers "github.com/the-gizmo-dojo/g8s/pkg/generated/listers/api.g8s.io/v1alpha1"
	jobinformers "k8s.io/client-go/informers/batch/v1"
	secretinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	joblisters "k8s.io/client-go/listers/batch/v1"
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
	rotationInformer informers.RotationInformer
	secretInformer   secretinformers.SecretInformer
	jobInformer      jobinformers.JobInformer

	// listers for our custom types
	passwordsLister listers.PasswordLister
	passwordsSynced cache.InformerSynced
	rotationsLister listers.RotationLister
	rotationsSynced cache.InformerSynced

	// listers for k8s types owned by our custom types
	secretsLister secretslisters.SecretLister
	secretsSynced cache.InformerSynced
	jobsLister    joblisters.JobLister
	jobsSynced    cache.InformerSynced

	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder
}
