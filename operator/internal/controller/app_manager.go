package controller

import (
	"context"
	"fmt"

	cermanagerv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	"github.com/go-logr/logr"
	appsv1alpha1 "github.com/mdshihabuddin01/simple-ms-with-apigw/operator/api/v1alpha1"
	"github.com/pkg/errors"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	authorizationv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ApplicationReconciler) updateStatus(ctx context.Context, app *appsv1alpha1.Application) error {
	// Update status with conditions
	condition := metav1.Condition{
		Type:               "Available",
		Status:             metav1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Reason:             "Reconciled",
		Message:            "Application resources successfully reconciled",
	}

	if len(app.Status.Conditions) == 0 || app.Status.Conditions[0].Type != condition.Type {
		app.Status.Conditions = []metav1.Condition{condition}
		return r.Status().Update(ctx, app)
	}

	return nil
}

func (r *ApplicationReconciler) watchClusterIssuer(b *builder.Builder) *builder.Builder {
	return b.Watches(
		&cermanagerv1.ClusterIssuer{},
		handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, object client.Object) []reconcile.Request {
			attachedApplications := &appsv1alpha1.ApplicationList{}
			listOps := &client.ListOptions{
				Namespace: "",
			}
			err := r.List(context.TODO(), attachedApplications, listOps)
			if err != nil {
				return []reconcile.Request{}
			}

			requests := make([]reconcile.Request, len(attachedApplications.Items))
			for i, item := range attachedApplications.Items {
				requests[i] = reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      item.GetName(),
						Namespace: item.GetNamespace(),
					},
				}
			}
			return requests
		}),
		builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		builder.WithPredicates(predicate.Funcs{
			DeleteFunc: func(e event.DeleteEvent) bool {
				// trigger reconciliation for deleted mysqlCR resources
				return true
			},
		}),
		builder.WithPredicates(predicate.Funcs{
			CreateFunc: func(e event.CreateEvent) bool {
				// trigger reconciliation for deleted mysqlCR resources
				return true
			},
		}),
		builder.WithPredicates(predicate.Funcs{
			UpdateFunc: func(e event.UpdateEvent) bool {
				// trigger reconciliation for deleted mysqlCR resources
				return true
			},
		}),
	)
}

// SetupWithManager sets up the controller with the Manager.
func (r *ApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	appBuilder := ctrl.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.Application{}).
		Named("application").
		WithOptions(controller.Options{MaxConcurrentReconciles: 2})

	appBuilder = r.watchClusterIssuer(appBuilder)
	return appBuilder.Owns(&corev1.Namespace{}).
		Owns(&corev1.ServiceAccount{}).
		Owns(&authorizationv1.ClusterRole{}).
		Owns(&authorizationv1.ClusterRoleBinding{}).
		Owns(&corev1.ServiceAccount{}).
		Owns(&authorizationv1.Role{}).
		Owns(&authorizationv1.RoleBinding{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.Service{}).
		Owns(&appsv1.Deployment{}).
		Owns(&batchv1.Job{}).
		Owns(&corev1.Secret{}).
		Owns(&networkingv1.IngressClass{}).
		Owns(&admissionregistrationv1.ValidatingWebhookConfiguration{}).
		Owns(&networkingv1.NetworkPolicy{}).
		Owns(&networkingv1.Ingress{}).
		Complete(r)
}

func (r *ApplicationReconciler) ensureApplicationFinalizer(ctx context.Context, app *appsv1alpha1.Application, l logr.Logger) error {
	// examine DeletionTimestamp to determine if object is under deletion
	if app.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// to registering our finalizer.
		if !controllerutil.ContainsFinalizer(app, applicationFinalizer) {
			controllerutil.AddFinalizer(app, applicationFinalizer)
			if err := r.Update(ctx, app); err != nil {
				return errors.Wrap(err, "Failed to add finalizer with app")
			}
			l.Info("-------------------Added finalizer in app -------------------")
		}
	} else {
		// The object is being deleted
		// remove our finalizer from the list and update it.
		var err error
		if err = r.deleteApplicationDependencies(ctx, app); err == nil {

			if controllerutil.RemoveFinalizer(app, applicationFinalizer) {
				l.Info("Removed finalizer")
			}
			if err = r.Update(ctx, app); err != nil {
				return errors.Wrap(err, "Failed to remove finalizer from app")
			}
			l.Info("-------------------Removed finalizer in app -------------------")

			return nil
		}
		return fmt.Errorf("couldn't delete resources: %w", err)
	}

	return nil
}

func (r *ApplicationReconciler) deleteApplicationDependencies(ctx context.Context, app *appsv1alpha1.Application) error {
	certManagerClientSet := r.CertManagerClientSet
	if certManagerClientSet == nil {
		return fmt.Errorf("failed to get certmanager client")
	}

	issuer, err := certManagerClientSet.CertmanagerV1().ClusterIssuers().Get(ctx, issuerName(app.Namespace), metav1.GetOptions{})

	if kerr.IsNotFound(err) {
		//l.Info("Issuer not found")
		return nil
	}

	err = certManagerClientSet.CertmanagerV1().ClusterIssuers().Delete(ctx, issuer.Name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
