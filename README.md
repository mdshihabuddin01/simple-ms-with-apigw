# Simple Mini-service Application with API Gateway

This repository contains a simple mini-service application with an API Gateway, and all necessary Kubernetes deployment configurations, including Helm charts and an operator.

## Getting Started

To get the application running, follow these general steps:

### 1. Backend Setup

The `backend` directory contains the services (`auth-service`, `order-service`) and the API Gateway configuration for running the application locally .

-   **Local Development:** Refer to `backend/docker-compose.yml` to set up and run the services using Docker Compose. More detailed instructions might be in `docs/running-backend-app.md`.

### 2. Kubernetes Deployment

The `manifests` and `artifacts` directories contain everything needed for Kubernetes deployment.

-   **Manifests:** Basic Kubernetes YAML files are located in `manifests/`.
-   **Helm Charts:** The `artifacts/` directory contains Helm charts for the application operator, Cert-Manager, Metrics Server, and Monitoring Operator.
-   **Operator:** The `operator/` directory holds the source code for a Kubernetes operator, likely for managing the application's lifecycle within the cluster.
-   **Documentation:** Consult `docs` for detailed deployment instructions.

## Monitoring

The `backend/monitoring` and `artifacts/monitoring-operator` directories provide configurations for Prometheus and Grafana to monitor the services and the Kubernetes cluster.

## API Documentation

Refer to `docs/api-documentation.md` for details on the available API endpoints and how to interact with the services.