# Kubernetes Deployment

This folder contains production-friendly Kubernetes manifests for FeedPulse_Cloud_Native_Platform microservices.

## Structure

- `base/` - shared manifests for all environments
- `overlays/atlas/` - MongoDB Atlas-first deployment (recommended)
- `overlays/local-mongo/` - local in-cluster MongoDB deployment

## Prerequisites

- Kubernetes cluster (minikube, kind, EKS, GKE, AKS, etc.)
- Ingress controller installed (NGINX Ingress expected by `ingressClassName: nginx`)
- Docker images published and referenced in `k8s/base/*.yaml`

## Important setup

1. Update image names in:
   - `k8s/base/auth-service.yaml`
   - `k8s/base/feedback-service.yaml`
   - `k8s/base/ai-service.yaml`
   - `k8s/base/api-gateway.yaml`
2. Update secrets in overlay patch files:
   - `k8s/overlays/atlas/secret-patch.yaml` or
   - `k8s/overlays/local-mongo/secret-patch.yaml`
3. Update ingress host in `k8s/base/ingress.yaml` from `feedpulse.local` to your DNS.

## Deploy (Atlas)

```bash
kubectl apply -k k8s/overlays/atlas
kubectl -n feedpulse get pods
kubectl -n feedpulse get svc
kubectl -n feedpulse get ingress
```

## Deploy (Local Mongo in cluster)

```bash
kubectl apply -k k8s/overlays/local-mongo
kubectl -n feedpulse get pods
kubectl -n feedpulse get svc
kubectl -n feedpulse get statefulset
```

## Rollout checks

```bash
kubectl -n feedpulse rollout status deploy/auth-service
kubectl -n feedpulse rollout status deploy/ai-service
kubectl -n feedpulse rollout status deploy/feedback-service
kubectl -n feedpulse rollout status deploy/api-gateway
```

## Access

- Through ingress host configured in `k8s/base/ingress.yaml`
- Or port-forward gateway:

```bash
kubectl -n feedpulse port-forward svc/api-gateway 8080:80
```

Then open `http://localhost:8080`.

## Cleanup

```bash
kubectl delete -k k8s/overlays/atlas
# or
kubectl delete -k k8s/overlays/local-mongo
```

