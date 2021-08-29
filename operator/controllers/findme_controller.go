/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"reflect"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"context"

	applicationv1alpha1 "github.com/cmwylie19/find-me/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
)

// FindmeReconciler reconciles a Findme object
type FindmeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=application.caseywylie.io,resources=findmes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=application.caseywylie.io,resources=findmes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=application.caseywylie.io,resources=findmes/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;
//+kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Findme object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.9.2/pkg/reconcile
func (r *FindmeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrllog.FromContext(ctx)

	// Log out
	log.Info("HITTING THE SERVICE")

	// Fetch the Findme instance
	findme := &applicationv1alpha1.Findme{}
	err := r.Get(ctx, req.NamespacedName, findme)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("Findme resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get Findme")
		return ctrl.Result{}, err
	}

	// Check if the serviceAccount already exists, if not create a new one
	serviceAccount := &corev1.ServiceAccount{}
	err = r.Get(ctx, types.NamespacedName{Name: findme.Name, Namespace: findme.Namespace}, serviceAccount)

	if err != nil && errors.IsNotFound(err) {
		// Define a new serviceAccount
		sa := r.serviceAccountForFindme(findme)
		log.Info("Creating a new Service Account", "Service.Namespace", sa.Namespace, "Service.Name", sa.Name)
		err = r.Create(ctx, sa)
		if err != nil {
			log.Error(err, "Failed to create a new Service Account", "Service.Namespace", sa.Namespace, "Service.Name", sa.Name)
			return ctrl.Result{}, err
		}
		//ServiceAccount created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "failed to get Service Account")
		return ctrl.Result{}, err
	}

	// Check if the deployment already exists, if not create a new one
	deployment := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: findme.Name, Namespace: findme.Namespace}, deployment)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep := r.deploymentForFindme(findme)
		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Create(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}

	// Check if the service already exists, if not create a new one
	service := &corev1.Service{}
	err = r.Get(ctx, types.NamespacedName{Name: findme.Name, Namespace: findme.Namespace}, service)

	if err != nil && errors.IsNotFound(err) {
		// Define a new service

		svc := r.serviceForFindme(findme)
		log.Info("Creating a new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
		err = r.Create(ctx, svc)
		if err != nil {
			log.Error(err, "Failed to create a new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
			return ctrl.Result{}, err
		}
		//Service created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "failed to get Service")
		return ctrl.Result{}, err
	}

	// Ensure the deployment size is the same as the spec
	size := findme.Spec.Size
	if *deployment.Spec.Replicas != size {
		deployment.Spec.Replicas = &size
		err = r.Update(ctx, deployment)
		if err != nil {
			log.Error(err, "Failed to update Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
			return ctrl.Result{}, err
		}
		// Ask to requeue after 1 minute in order to give enough time for the
		// pods be created on the cluster side and the operand be able
		// to do the next update step accurately.
		return ctrl.Result{RequeueAfter: time.Minute}, nil
	}

	// Update the Findme status with the pod names
	// List the pods for this memcached's deployment
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(findme.Namespace),
		client.MatchingLabels(labelsForFindme(findme.Name)),
	}
	if err = r.List(ctx, podList, listOpts...); err != nil {
		log.Error(err, "Failed to list pods", "Findme.Namespace", findme.Namespace, "Findme.Name", findme.Name)
		return ctrl.Result{}, err
	}
	podNames := getPodNames(podList.Items)

	// Update status.Nodes if needed
	if !reflect.DeepEqual(podNames, findme.Status.Nodes) {
		findme.Status.Nodes = podNames
		err := r.Status().Update(ctx, findme)
		if err != nil {
			log.Error(err, "Failed to update Findme status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// serviceAccountForFindme returns a findme ServiceAccount object
func (r *FindmeReconciler) serviceAccountForFindme(m *applicationv1alpha1.Findme) *corev1.ServiceAccount {
	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
	}

	// set Findme instance as the owner and controller
	ctrl.SetControllerReference(m, serviceAccount, r.Scheme)
	return serviceAccount
}

// serviceForFindme returns a findme Service object
func (r *FindmeReconciler) serviceForFindme(m *applicationv1alpha1.Findme) *corev1.Service {
	//ls := labelsForFindme(m.Name)

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
			Labels: map[string]string{
				"version":                "v1",
				"app.kubernetes.io/name": "findme",
				"app":                    "findme",
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				Name: "http",
				Port: 80,
			}},
			Selector: map[string]string{
				"app": "findme",
			},
		},
	}
	// Set Findme instance as the owner and controller
	ctrl.SetControllerReference(m, service, r.Scheme)
	return service
}

// deploymentForMemcached returns a findme Deployment object
func (r *FindmeReconciler) deploymentForFindme(m *applicationv1alpha1.Findme) *appsv1.Deployment {
	ls := labelsForFindme(m.Name)
	replicas := m.Spec.Size

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: m.Name,
					Containers: []corev1.Container{{
						Image: "docker.io/cmwylie19/find-me:latest",
						Name:  "findme",
						// Command: []string{"memcached", "-m=64", "-o", "modern", "-v"},
						Ports: []corev1.ContainerPort{{
							ContainerPort: 80,
							Name:          "http",
						}},
					}},
					RestartPolicy: corev1.RestartPolicyAlways,
				},
			},
		},
	}
	// Set Findme instance as the owner and controller
	ctrl.SetControllerReference(m, deployment, r.Scheme)
	return deployment
}

// labelsForMemcached returns the labels for selecting the resources
// belonging to the given memcached CR name.
func labelsForFindme(name string) map[string]string {
	return map[string]string{"app": "findme", "findme_cr": name}
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

// SetupWithManager sets up the controller with the Manager.
func (r *FindmeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&applicationv1alpha1.Findme{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.ServiceAccount{}).
		Owns(&appsv1.Deployment{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: 2}).
		Complete(r)
}
