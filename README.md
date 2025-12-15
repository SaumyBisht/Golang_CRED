# Golang Microservices on Kubernetes

This project demonstrates a microservices architecture using **Go (Golang)**, **Docker**, and **Kubernetes**. It consists of two services that communicate with each other within a cluster, configured for local development using Minikube.

## 🏗 Architecture

The system consists of two main services:

1.  **Deployment Service (Frontend/Gateway):**
    * **Port:** `8082` (Container), Exposed on Port `80` (NodePort).
    * **Role:** Acts as the entry point. It receives requests and communicates with the Core Service.
2.  **Core Service (Backend):**
    * **Port:** `8081` (Container), Exposed on Port `30001` .
    * **Role:** Handles business logic and internal processing.

**Data Flow:**
```text
User Request
    │
    ▼
[ NodePort Service: deployment-service (Port 80) ]
    │
    ▼
[ Pod: deployment-service ]
    │
    │ (HTTP Request via internal DNS)
    │ "http://core-service:8081"
    │
    ▼
[ Service: core-service (ClusterIP: 8081) ]
    │
    ▼
[ Pod: core-service ]
```

---

## 🚀 Tech Stack

* **Language:** Go (Golang) 1.25+
* **Containerization:** Docker (Multistage builds)
* **Orchestration:** Kubernetes (Minikube)
* **OS:** Ubuntu/WSL2

---

## 📂 Project Structure

```text
.
├── core_service/           # Backend Logic
│   ├── main.go
│   └── Dockerfile          # Exposes Port 8081
├── deployment_service/     # Frontend/Gateway Logic
│   ├── main.go
│   └── Dockerfile          # Exposes Port 8082
├── k8s/                    # Kubernetes Manifests
│   ├── config-secrets.yaml   # ConfigMap (URLs) & Secrets (JWT)
│   ├── core-deployment.yaml  # Core Pod Definitions
│   ├── core-service.yaml     # Core Networking (NodePort 30001)
│   ├── deploy-deployment.yaml# Deploy Service Pod Definitions
│   └── deploy-service.yaml   # Deploy Service Networking (Port 80)
└── README.md
```

---

## 🛠️ Setup & Deployment Guide

### Prerequisites
* Docker installed
* Minikube installed & running (`minikube start`)
* `kubectl` configured

### Step 1: Build Docker Images
Since we are using Minikube, we must build the images and load them manually so the cluster can access them without a registry.

**1. Build Core Service**
```bash
cd core_service
docker build -t core-service:v1 .
cd ..
```

**2. Build Deployment Service**
```bash
cd deployment_service
docker build -t deployment-service:v1 .
cd ..
```

### Step 2: Load Images into Minikube
*Critically important step for WSL/Minikube users.*

```bash
minikube image load core-service:v1
minikube image load deployment-service:v1
```

### Step 3: Deploy to Kubernetes
Apply all configuration files, secrets, and deployments at once.

```bash
kubectl apply -f k8s/
```

### Step 4: Verify Deployment
Check if pods are running and services are up.

```bash
kubectl get pods
kubectl get services
```

---

## 🔌 Accessing the APIs

Since this is running in Minikube, you cannot always access `localhost` directly. Use the Minikube IP.

### 1. Access Deployment Service (Gateway)
This service is exposed on **Port 80** of the Minikube node.

```bash
# Get the direct URL
minikube service deployment-service --url
```
*Example Output:* `http://192.168.49.2:31525` (Port mapped to 80)

### 2. Access Core Service (Direct Debugging)
We exposed the core service on **NodePort 30001** for testing.

* **URL:** `http://<MINIKUBE-IP>:30001`
* **Command to get URL:**
    ```bash
    minikube service core-service --url
    ```

---

## 📝 Environment Variables

These are managed via `k8s/config-secrets.yaml`.

| Variable | Value | Source |
| :--- | :--- | :--- |
| `JWT_SECRET` | `demo-secret-sb` | **Secret:** `app-secrets` |
| `CORE_SERVICE_URL` | `http://core-service:8081` | **ConfigMap:** `app-config` |

---

## 🧹 Cleanup

To stop the application and remove resources:

```bash
# Remove all K8s resources
kubectl delete -f k8s/

# Stop Minikube
minikube stop
```