package controller

import (
	"fmt"
	"strings"

	appsv1alpha1 "github.com/mdshihabuddin01/simple-ms-with-apigw/operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// Helper functions for naming conventions
func (r *ApplicationReconciler) getConfigMapName(app *appsv1alpha1.Application) string {
	return fmt.Sprintf("%s-config", app.Name)
}

func (r *ApplicationReconciler) getSecretName(app *appsv1alpha1.Application) string {
	return fmt.Sprintf("%s-secret", app.Name)
}

func (r *ApplicationReconciler) getDeploymentName(app *appsv1alpha1.Application) string {
	return fmt.Sprintf("%s-deployment", app.Name)
}

func (r *ApplicationReconciler) getServiceName(app *appsv1alpha1.Application) string {
	return fmt.Sprintf("%s-service", app.Name)
}

func (r *ApplicationReconciler) getIngressName(app *appsv1alpha1.Application) string {
	return fmt.Sprintf("%s-ingress", app.Name)
}

// Helper functions
func (r *ApplicationReconciler) getLabels(app *appsv1alpha1.Application) map[string]string {
	return map[string]string{
		"app":     app.Name,
		"managed": "application-operator",
	}
}

func (r *ApplicationReconciler) buildResourceRequirements(resources *appsv1alpha1.ResourcesSpec) corev1.ResourceRequirements {
	req := corev1.ResourceRequirements{}

	if resources == nil {
		return req
	}

	req.Requests = corev1.ResourceList{}
	req.Limits = corev1.ResourceList{}

	if resources.Requests != nil {
		if resources.Requests.CPU != "" {
			req.Requests[corev1.ResourceCPU] = resource.MustParse(resources.Requests.CPU)
		}
		if resources.Requests.Memory != "" {
			req.Requests[corev1.ResourceMemory] = resource.MustParse(resources.Requests.Memory)
		}
	}

	if resources.Limits != nil {
		if resources.Limits.CPU != "" {
			req.Limits[corev1.ResourceCPU] = resource.MustParse(resources.Limits.CPU)
		}
		if resources.Limits.Memory != "" {
			req.Limits[corev1.ResourceMemory] = resource.MustParse(resources.Limits.Memory)
		}
	}

	return req
}

func strPtr(s string) *string {
	return &s
}

func getSecretData(envVars []appsv1alpha1.EnvVar, prefix string) map[string][]byte {
	data := make(map[string][]byte)
	for _, envVar := range envVars {
		if strings.HasPrefix(envVar.Name, prefix) {
			data[strings.TrimPrefix(envVar.Name, prefix)] = []byte(envVar.Value)
		}
	}
	return data
}

func getConfigMapData(envVars []appsv1alpha1.EnvVar, prefix string) map[string]string {
	data := make(map[string]string)
	for _, envVar := range envVars {
		if strings.HasPrefix(envVar.Name, prefix) {
			data[strings.TrimPrefix(envVar.Name, prefix)] = envVar.Value
		}
	}
	return data
}

func issuerName(rescourceName string) string {
	return strings.Join([]string{"app", rescourceName, "issuer"}, "-")
}

func ingressClassName() string {
	return "kong"
}
