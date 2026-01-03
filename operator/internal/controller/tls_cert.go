package controller

import (
	"context"
	"fmt"

	acmeissuerv1 "github.com/cert-manager/cert-manager/pkg/apis/acme/v1"
	cermanagerv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	appsv1alpha1 "github.com/mdshihabuddin01/simple-ms-with-apigw/operator/api/v1alpha1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (r *ApplicationReconciler) ReconcileIssuer(ctx context.Context, app *appsv1alpha1.Application) (*cermanagerv1.ClusterIssuer, error) {
	certManagerClientSet := r.CertManagerClientSet
	if certManagerClientSet == nil {
		return nil, fmt.Errorf("failed to get certmanager client")
	}

	issuer, err := certManagerClientSet.CertmanagerV1().ClusterIssuers().Get(ctx, issuerName(app.Namespace), metav1.GetOptions{})
	if err == nil {
		updateClusterIssuer := false
		if app.Spec.TLS.Issuer.ACMEIssuer.Email != issuer.Spec.ACME.Email {
			issuer.Spec.ACME.Email = app.Spec.TLS.Issuer.ACMEIssuer.Email
			updateClusterIssuer = true
		}

		if app.Spec.TLS.Issuer.ACMEIssuer.Server != issuer.Spec.ACME.Server {
			issuer.Spec.ACME.Server = app.Spec.TLS.Issuer.ACMEIssuer.Server
			updateClusterIssuer = true
		}

		if updateClusterIssuer {
			issuer, err = certManagerClientSet.CertmanagerV1().ClusterIssuers().Update(ctx, issuer, metav1.UpdateOptions{})
			if err != nil {
				return nil, err
			}
			return issuer, nil
		}

		customError := fmt.Errorf(found)
		return issuer, customError
	}

	if !errors.IsNotFound(err) {
		return issuer, err
	}

	issuer = &cermanagerv1.ClusterIssuer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      issuerName(app.Namespace),
			Namespace: app.Namespace,
		},
		Spec: cermanagerv1.IssuerSpec{

			IssuerConfig: cermanagerv1.IssuerConfig{
				ACME: &app.Spec.TLS.Issuer.ACMEIssuer,
			},
		},
	}

	issuer.Spec.ACME.Solvers = []acmeissuerv1.ACMEChallengeSolver{
		{
			HTTP01: &acmeissuerv1.ACMEChallengeSolverHTTP01{
				Ingress: &acmeissuerv1.ACMEChallengeSolverHTTP01Ingress{
					IngressClassName: strPtr(ingressClassName()),
				},
			},
		},
	}

	issuer, err = certManagerClientSet.CertmanagerV1().ClusterIssuers().Create(ctx, issuer, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return issuer, nil

}

func (r *ApplicationReconciler) updateIngressForTLS(ctx context.Context, app *appsv1alpha1.Application) (*networkingv1.Ingress, error) {
	ingress := &networkingv1.Ingress{}
	err := r.Get(ctx, types.NamespacedName{Name: r.getIngressName(app), Namespace: app.Namespace}, ingress)
	if err == nil {
		updateIngress := false
		/*		if _, exists := ingress.ObjectMeta.Annotations[clusterIssuerLabel]; !exists {
					ingress.ObjectMeta.Annotations[clusterIssuerLabel] = issuerName(app.Namespace)
					updateIngress = true
				}
		*/
		if len(ingress.Spec.Rules) > 0 {
			if ingress.Spec.Rules[0].Host != app.Spec.Ingress.Host {
				ingress.Spec.Rules[0].Host = app.Spec.Ingress.Host
				updateIngress = true

			}
		}

		if len(ingress.Spec.TLS) > 0 {
			if ingress.Spec.TLS[0].Hosts[0] != app.Spec.Ingress.Host {
				ingress.Spec.TLS[0].Hosts = []string{app.Spec.Ingress.Host}
				updateIngress = true

			}

			if ingress.Spec.TLS[0].SecretName != app.Spec.TLS.Issuer.ACMEIssuer.PrivateKey.Name {
				ingress.Spec.TLS[0].SecretName = app.Spec.TLS.Issuer.ACMEIssuer.PrivateKey.Name
				updateIngress = true

			}
		} else {
			ingress.Spec.TLS = []networkingv1.IngressTLS{{
				Hosts:      []string{app.Spec.Ingress.Host},
				SecretName: app.Spec.TLS.Issuer.ACMEIssuer.PrivateKey.Name,
			},
			}
			updateIngress = true

		}

		if updateIngress {
			return ingress, r.Update(ctx, ingress)
		}
		customError := fmt.Errorf("Found")
		return ingress, customError

	}
	if !errors.IsNotFound(err) {
		return ingress, err
	}

	return ingress, err

}
