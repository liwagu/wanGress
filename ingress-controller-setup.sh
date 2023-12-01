#!/bin/bash

# Install Nginx Ingress Controller
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.48.1/deploy/static/provider/cloud/deploy.yaml

# Verify that the Ingress Controller has been created
kubectl get pods -n ingress-nginx \
  -l app.kubernetes.io/name=ingress-nginx --watch

echo "Ingress Controller setup completed."
