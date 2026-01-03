package controller

import (
	"context"
	"fmt"
	"strings"

	appsv1alpha1 "github.com/mdshihabuddin01/simple-ms-with-apigw/operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// reconcileConfigMap creates or updates a ConfigMap for environment variables
func (r *ApplicationReconciler) reconcileConfigMap(ctx context.Context, app *appsv1alpha1.Application) error {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.getConfigMapName(app),
			Namespace: app.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(app, appsv1alpha1.GroupVersion.WithKind("Application")),
			},
		},
	}

	// Create ConfigMap data from environment variables
	configMapData := getConfigMapData(app.Spec.EnvVars, "CM_")

	if len(configMapData) == 0 {
		err := r.Get(ctx, client.ObjectKey{
			Name:      r.getConfigMapName(app),
			Namespace: app.Namespace,
		}, configMap)
		if err != nil {
			if !errors.IsNotFound(err) {
				return fmt.Errorf("failed to get ConfigMap: %w", err)
			}
		} else {
			return fmt.Errorf("configmap found but there is no configmap in env, failed to delete ConfigMap: %w", r.Delete(ctx, configMap))
		}
		return nil
	}

	op, err := controllerutil.CreateOrUpdate(ctx, r.Client, configMap, func() error {
		configMap.Labels = r.getLabels(app)
		// Set owner reference
		/*		if err := controllerutil.SetControllerReference(app, configMap, r.Scheme); err != nil {
				return err
			}*/
		configMap.Data = configMapData

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to reconcile ConfigMap: %w", err)
	}

	if op != controllerutil.OperationResultNone {
		log.FromContext(ctx).Info("ConfigMap reconciled", "operation", op)
	}
	return nil
}

// reconcileSecret creates or updates a Secret for sensitive data
func (r *ApplicationReconciler) reconcileSecret(ctx context.Context, app *appsv1alpha1.Application) error {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.getSecretName(app),
			Namespace: app.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(app, appsv1alpha1.GroupVersion.WithKind("Application")),
			},
		},
	}

	secretData := getSecretData(app.Spec.EnvVars, "SEC_")

	if len(secretData) == 0 {
		err := r.Get(ctx, client.ObjectKey{
			Name:      r.getSecretName(app),
			Namespace: app.Namespace,
		}, secret)
		if err != nil {
			if !errors.IsNotFound(err) {
				return fmt.Errorf("failed to get Secret: %w", err)
			}
		} else {
			return fmt.Errorf("secret found but there is no secret in env, failed to delete Secret: %w", r.Delete(ctx, secret))
		}
		return nil
	}

	op, err := controllerutil.CreateOrUpdate(ctx, r.Client, secret, func() error {
		secret.Labels = r.getLabels(app)

		secret.Data = secretData
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to reconcile Secret: %w", err)
	}

	if op != controllerutil.OperationResultNone {
		log.FromContext(ctx).Info("Secret reconciled", "operation", op)
	}

	return nil
}

// reconcileDeployment creates or updates a Deployment
func (r *ApplicationReconciler) reconcileDeployment(ctx context.Context, app *appsv1alpha1.Application) error {
	deploymentName := r.getDeploymentName(app)

	// 1. Build environment variables (Logic remains the same)
	envVars := []corev1.EnvVar{}
	for _, envVar := range app.Spec.EnvVars {
		if strings.HasPrefix(envVar.Name, "CM_") {
			cleanKey := strings.TrimPrefix(envVar.Name, "CM_")
			envVars = append(envVars, corev1.EnvVar{
				Name: cleanKey,
				ValueFrom: &corev1.EnvVarSource{
					ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{Name: r.getConfigMapName(app)},
						Key:                  cleanKey,
					},
				},
			})
		} else if strings.HasPrefix(envVar.Name, "SEC_") {
			cleanKey := strings.TrimPrefix(envVar.Name, "SEC_")
			envVars = append(envVars, corev1.EnvVar{
				Name: cleanKey,
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{Name: r.getSecretName(app)},
						Key:                  cleanKey,
					},
				},
			})
		} else {
			envVars = append(envVars, corev1.EnvVar{Name: envVar.Name, Value: envVar.Value})
		}
	}

	// 2. Define the Deployment using TypeMeta (REQUIRED for Patch)
	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: app.Namespace,
			Labels:    r.getLabels(app),
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(app, appsv1alpha1.GroupVersion.WithKind("Application")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: app.Spec.Replicas, // Uses pointer from Spec directly
			Selector: &metav1.LabelSelector{
				MatchLabels: r.getLabels(app),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: r.getLabels(app),
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:      app.Name,
							Image:     *app.Spec.Image,
							Ports:     []corev1.ContainerPort{{ContainerPort: app.Spec.ContainerPort}},
							Env:       envVars,
							Resources: r.buildResourceRequirements(app.Spec.Resources),
						},
					},
				},
			},
		},
	}

	// 4. Apply the Patch (Replaces r.Get, r.Create, and r.Update)
	// This is idempotent: it creates if missing, updates if changed, stays quiet if same.
	err := r.Patch(ctx, deployment, client.Apply,
		client.FieldOwner("application-controller"),
		client.ForceOwnership,
	)
	if err != nil {
		return fmt.Errorf("failed to patch Deployment: %w", err)
	}

	return nil
}

func (r *ApplicationReconciler) reconcileService(ctx context.Context, app *appsv1alpha1.Application) error {
	if app.Spec.Service == nil {
		return nil
	}

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.getServiceName(app),
			Namespace: app.Namespace,
		},
	}

	op, err := controllerutil.CreateOrUpdate(ctx, r.Client, service, func() error {
		service.Labels = r.getLabels(app)
		service.OwnerReferences = []metav1.OwnerReference{
			*metav1.NewControllerRef(app, appsv1alpha1.GroupVersion.WithKind("Application")),
		}

		service.Spec.Selector = r.getLabels(app)
		service.Spec.Type = corev1.ServiceTypeClusterIP
		if app.Spec.Service.Type != "" {
			service.Spec.Type = corev1.ServiceType(app.Spec.Service.Type)
		}

		service.Annotations = app.Spec.Service.Annotations

		// 3. Update Ports
		service.Spec.Ports = []corev1.ServicePort{
			{
				Name:       "http",
				Port:       app.Spec.Service.Port,
				TargetPort: intstr.FromInt(int(app.Spec.ContainerPort)),
				Protocol:   corev1.ProtocolTCP,
			},
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to reconcile Service: %w", err)
	}

	// Log if something actually changed (Unchanged, Created, or Updated)
	if op != controllerutil.OperationResultNone {
		log.FromContext(ctx).Info("Service reconciled", "operation", op)
	}

	return nil
}

// reconcileIngress creates or updates an Ingress for the application
func (r *ApplicationReconciler) reconcileIngress(ctx context.Context, app *appsv1alpha1.Application) error {
	if app.Spec.Ingress == nil {
		// If ingress spec is nil, check if ingress exists and delete it
		ingress := &networkingv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      r.getIngressName(app),
				Namespace: app.Namespace,
			},
		}
		err := r.Get(ctx, client.ObjectKey{
			Name:      r.getIngressName(app),
			Namespace: app.Namespace,
		}, ingress)
		if err != nil {
			if !errors.IsNotFound(err) {
				return fmt.Errorf("failed to get Ingress: %w", err)
			}
			// Ingress doesn't exist, nothing to do
			return nil
		}
		// Ingress exists, delete it
		return fmt.Errorf("failed to delete Ingress: %w", r.Delete(ctx, ingress))
	}

	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.getIngressName(app),
			Namespace: app.Namespace,
		},
	}

	op, err := controllerutil.CreateOrUpdate(ctx, r.Client, ingress, func() error {
		ingress.Labels = r.getLabels(app)
		ingress.OwnerReferences = []metav1.OwnerReference{
			*metav1.NewControllerRef(app, appsv1alpha1.GroupVersion.WithKind("Application")),
		}

		// Set ingress class name
		ingress.Annotations = app.Spec.Ingress.Annotations
		ingress.Spec.IngressClassName = strPtr(ingressClassName())

		// Create ingress rule
		pathType := networkingv1.PathTypePrefix
		ingress.Spec.Rules = []networkingv1.IngressRule{
			{
				Host: app.Spec.Ingress.Host,
				IngressRuleValue: networkingv1.IngressRuleValue{
					HTTP: &networkingv1.HTTPIngressRuleValue{
						Paths: []networkingv1.HTTPIngressPath{
							{
								Path:     app.Spec.Ingress.Path,
								PathType: &pathType,
								Backend: networkingv1.IngressBackend{
									Service: &networkingv1.IngressServiceBackend{
										Name: r.getServiceName(app),
										Port: networkingv1.ServiceBackendPort{
											Number: app.Spec.Service.Port,
										},
									},
								},
							},
						},
					},
				},
			},
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to reconcile Ingress: %w", err)
	}

	// Log if something actually changed (Unchanged, Created, or Updated)
	if op != controllerutil.OperationResultNone {
		log.FromContext(ctx).Info("Ingress reconciled", "operation", op)
	}

	return nil
}
