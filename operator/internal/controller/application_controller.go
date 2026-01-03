/*
Copyright 2025.

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

package controller

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	certmanagerclient "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	appsv1alpha1 "github.com/mdshihabuddin01/simple-ms-with-apigw/operator/api/v1alpha1"
)

// ApplicationReconciler reconciles a Application object
type ApplicationReconciler struct {
	client.Client
	Scheme               *runtime.Scheme
	CertManagerClientSet *certmanagerclient.Clientset
	//KubeClients          *KubeClients
}

// +kubebuilder:rbac:groups=apps.example.com,resources=applications,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.example.com,resources=applications/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps.example.com,resources=applications/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ApplicationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := logf.FromContext(ctx)

	// Fetch the Application instance
	var app appsv1alpha1.Application
	if err := r.Get(ctx, req.NamespacedName, &app); err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Application resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get Application")
		return ctrl.Result{}, err
	}

	// Create or update ConfigMap
	if err := r.reconcileConfigMap(ctx, &app); err != nil {
		logger.Error(err, "Failed to reconcile ConfigMap")
		return ctrl.Result{}, err
	}

	// Create or update Secret
	if err := r.reconcileSecret(ctx, &app); err != nil {
		logger.Error(err, "Failed to reconcile Secret")
		return ctrl.Result{}, err
	}

	// Create or update Deployment
	if err := r.reconcileDeployment(ctx, &app); err != nil {
		logger.Error(err, "Failed to reconcile Deployment")
		return ctrl.Result{}, err
	}

	// Create or update Service
	if err := r.reconcileService(ctx, &app); err != nil {
		logger.Error(err, "Failed to reconcile Service")
		return ctrl.Result{}, err
	}

	// Create or update Ingress
	if err := r.reconcileIngress(ctx, &app); err != nil {
		logger.Error(err, "Failed to reconcile Ingress")
		return ctrl.Result{}, err
	}

	if app.Spec.TLS != nil {
		if app.Spec.TLS.Enable {
			certManagerIssuer, err := r.ReconcileIssuer(ctx, &app)
			if err != nil {
				if err.Error() != found {
					logger.Error(err, fmt.Sprintf("Failed to create certmanager issuer: %s/%s", certManagerIssuer.Name, certManagerIssuer.Namespace))
					return ctrl.Result{}, nil
				}
			} else {
				logger.Info(fmt.Sprintf("Successfully created certmanager issuer: %s/%s", certManagerIssuer.Name, certManagerIssuer.Namespace))
			}
			// TODO:
			time.Sleep(50 * time.Second)

			if len(certManagerIssuer.Status.Conditions) > 0 {
				if certManagerIssuer.Status.Conditions[0].Status == trueValue {
					//l.Info("Cert manager is ready")
					updatedIngress, err := r.updateIngressForTLS(ctx, &app)
					if err != nil {
						if err.Error() != found {
							logger.Error(err, fmt.Sprintf("Failed to update ingress for tls: %s/%s", updatedIngress.Name, updatedIngress.Namespace))
							return ctrl.Result{}, nil
						}
					} else {
						logger.Info(fmt.Sprintf("Successfully updated ingress for tls: %s/%s", updatedIngress.Name, updatedIngress.Namespace))
					}
				}
			}
		}
	}

	if err := r.ensureApplicationFinalizer(ctx, &app, logger); err != nil {
		return ctrl.Result{}, err
	}

	// Update status
	if err := r.updateStatus(ctx, &app); err != nil {
		logger.Error(err, "Failed to update status")
		return ctrl.Result{}, err
	}

	logger.Info("Successfully reconciled Application")
	return ctrl.Result{}, nil
}
