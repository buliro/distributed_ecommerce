# Kubernetes Deployment Manifests

Apply manifests in order:

```bash
kubectl apply -f namespace.yaml
kubectl apply -f config.yaml
kubectl apply -f secret.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
```
