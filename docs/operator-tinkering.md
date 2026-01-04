# 🚀 Operator Tinkering Guide



### 🛠 Development & Code Generation
Api can be generaed using `kubebuilder create api` command. If you change any field or spec in api, before running your controller, you must ensure your API definitions and manifests are up to date.

* **`make generate`** Invokes `controller-gen` to update `zz_generated.deepcopy.go`. Run this whenever you modify your `_types.go` files to ensure your Go types implement the required `runtime.Object` interface.
* **`make manifests`** Generates Custom Resource Definitions (CRDs), RBAC roles, and Webhook configurations. This transforms your Go markers (e.g., `// +kubebuilder:rbac`) into YAML files in the `config/` directory.

---

### 💻 Local Testing

Test your operator logic against a live Kubernetes cluster (like Minikube or Kind) without building a Docker image.

* **`make install`** Applies the generated CRDs to your current Kubernetes cluster context.
* **`make run`** Starts the operator locally as a standard Go process. It uses your local `~/.kube/config` to communicate with the cluster—perfect for rapid debugging.

---

### 📦 Deployment & Packaging

When you are ready to move beyond local testing, use these commands to package your operator.

* **`make build-installer`** Generates a consolidated `install.yaml` file using Kustomize. This file contains everything needed (CRDs, RBAC, Deployment) to install your operator on any cluster with a single `kubectl apply` command.
* **Dockerization** Build a production-ready container image for your operator:
```bash
docker build -t codedapp/app-operator:1 . --no-cache

```



---

### 📋 Quick Reference Summary

| Command | Purpose | When to run it? |
| --- | --- | --- |
| **`make generate`** | Generates DeepCopy Go code | After changing API struct fields |
| **`make manifests`** | Generates CRD & RBAC YAMLs | After changing markers or types |
| **`make install`** | Registers CRDs in cluster | Before running/deploying |
| **`make run`** | Runs controller locally | During active development |
| **`make build-installer`** | Creates production YAML bundle | When releasing a new version |
