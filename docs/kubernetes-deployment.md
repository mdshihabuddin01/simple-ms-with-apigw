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
  -f helm-values/metrics-values.yaml

```

### 2. Security & Cert-Manager

We use `cert-manager` to handle TLS encryption and manage `Issuer` resources required by our Custom Resources.

```bash
# Install base dependencies
helm install dependent-manifests ./dependent-manifests \
  --namespace app-engine --create-namespace

# Install the Certificate Manager
helm install cert-manager ./cert-manager \
  --namespace app-engine --create-namespace 

```

### 3. monitoring-operator 

Install kube-prometheus-stack for APM
```bash
helm upgrade -i monitoring-stack-operator ./monitoring-operator \                     
  --namespace app-engine --create-namespace -f helm-values/monitoring-operator-values.yaml

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
  --namespace app-engine --create-namespace

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

Finally, deploy the core services: the **Authentication Service (Auth)** and the **Order Management System (OMS)**.

```bash
cd .. # Back to the parent prerequisites folder
kubectl apply -f auth-service-app.yaml
kubectl apply -f order-service-app.yaml

```

### 3. Deploy Monitoring

Finally, deploy the prometheus services monitor

```bash
kubectl apply -f auth-app-monitor.yaml
kubectl apply -f order-app-monitor.yaml
```
Service monitor discovers service by service label, in our example we used `app: auth-service-application` which is the label
of `auth-service-application-setvice`. prometheus scraps metrics from `/metrics` and listens in the same port as application for 
`auth-service` port is 8081 and `order-service` port is 8082

Grafana dashboard jsons were exported in the helm chart, so you can visit the dashboard right after the deployment, you can explore them 
in the `backend/monitoring/dashboards` directory for watching or importing manually.

As it will not be exposed by ingress, you can port-forward the prometheus and grafana service or change it to load-balancer or node-port.
## ✅ Verification Checklist

| Resource | Namespace | Command to Verify |
| --- | --- | --- |
| **Cert-Manager** | `app-engine` | `kubectl get pods -n app-engine` |
| **Kong Gateway** | `kong` | `kubectl get svc -n kong` |
| **MySQL DB** | `my-app-ns` | `kubectl get pvc -n my-app-ns` |
| **Services** | `my-app-ns` | `kubectl get apps -n my-app-ns` |
