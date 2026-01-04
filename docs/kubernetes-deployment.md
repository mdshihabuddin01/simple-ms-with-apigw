# 🏗️ Kubernetes Infrastructure Setup Guide

This guide walks you through the step-by-step deployment of the application ecosystem, from core dependencies to the custom Operator and microservices.

---

## 🛠️ Phase 1: Preparation & Core Infrastructure

Before starting, ensure you are in the `artifacts` directory where the configuration files are located.

```bash
cd artifacts

```

### 1. Metrics Server

Optional: Install this only if Horizontal Pod Autoscaling (HPA) and advanced resource monitoring are required for your environment.
```bash
helm install metrics-server metrics-server/metrics-server \
  --namespace kube-system \
  -f helm-values/metrics-server-values.yaml

```

### 2. Security & Cert-Manager

We use `cert-manager` to handle TLS encryption and manage `Issuer` resources required by our Custom Resources.

```bash
# Install prerequisite CRDs for other helm charts like cert-manager or monitoring-operator 
helm install dependent-manifests ./dependent-manifests \
  --namespace app-engine --create-namespace -f helm-values/dependent-manifsts-values.yaml

# Install the Certificate Manager
helm install cert-manager ./cert-manager \
  --namespace app-engine --create-namespace 

```

### 3. monitoring-operator 

Install kube-prometheus-stack for APM
```bash
helm upgrade --install monitoring-stack-operator ./monitoring-operator \
  --namespace app-engine \
  --create-namespace \
  -f helm-values/monitoring-operator-values.yaml

```

---

## 🚦 Phase 2: API Gateway & Operator

We use **Kong** as our Ingress Controller and a custom **Go Operator** to manage our application lifecycle.

### 1. Kong Ingress Controller

```bash
helm repo add kong https://charts.konghq.com && helm repo update
kubectl create namespace kong

helm install kong kong/ingress -n kong \
  --values helm-values/kic-values.yaml

```

### 2. Application Operator

The brain of our system that reconciles our custom application specs.

```bash
helm install app-operator ./app-operator-helm-chart \
  --namespace app-engine --create-namespace -f helm-values/app-operator-values.yaml

```

---

## 📦 Phase 3: Application & Database Deployment

Now, we deploy the specific application instances into the `my-app-ns` namespace.

### 1. Namespace & Database Setup

Navigate to the prerequisites directory and verify your `storageclass` settings before applying.

```bash
cd manifests/prerequisites
kubectl apply -f .

```

> **Note:** This step initializes the `my-app-ns` namespace and sets up the **MySQL** database backend.

### 2. Deploy Microservices

Deploy the core services: the **Authentication Service (Auth)** and the **Order Management System (OMS)**. You can use the default images from the manifests.

```bash
cd .. # Back to the parent prerequisites folder
kubectl apply -f auth-service-app.yaml
kubectl apply -f order-service-app.yaml

```
**Note:** We will use `konghq.com/plugins: validate-via-auth-service` annotation in order-service-app.yaml's service section, and 
`cert-manager.io/cluster-issuer: "cluster-issuer-name"` in auth-service.yaml's ingress section. In the env section, prefix `CM` sets 
the env as configmap data, `SEC` sets as secret data but removes the prefix before placing.
Go to `manifests/kong-plugins` directory. Deploy kong plugins for proper application routing and security
```bash
kubectl apply -f validate-via-auth-service-plugin.yaml
```
If you want to enable tls, uncomment tls section with appropriate values (including your host in the ingress section) and reapply the manifest. Update the ingress section annotation with the real `clusterissuer` name what 
you can find by using  `k get ClusterIssuer -A`
### 3. Deploy Monitoring

Finally, deploy the prometheus services monitor

```bash
kubectl apply -f auth-app-monitor.yaml
kubectl apply -f order-app-monitor.yaml
```
Service monitor discovers service by service label, in our example we used `app: auth-service-application` which is the label
of `auth-service-application-setvice`. prometheus scraps metrics from `/metrics` and listens in the application port

Grafana dashboard jsons were exported in the helm chart, so you can visit the dashboard right after the deployment, you can explore them 
in the `backend/monitoring/dashboards` directory for watching or importing manually. Search dashboards with `auth` or `order` keywords in grafana. As it will not be exposed by ingress, you can port-forward the prometheus and grafana service or change it to load-balancer or node-port.
## ✅ Verification Checklist

| Resource | Namespace | Command to Verify |
| --- | --- | --- |
| **Cert-Manager** | `app-engine` | `kubectl get pods -n app-engine` |
| **Kong Gateway** | `kong` | `kubectl get svc -n kong` |
| **MySQL DB** | `my-app-ns` | `kubectl get pvc -n my-app-ns` |
| **Services** | `my-app-ns` | `kubectl get apps -n my-app-ns` |
