#!/bin/bash

# Deploy the Kubernetes cluster
echo "Deploying Kubernetes cluster..."
bash cluster-setup.sh

# Deploy the Ingress Controller
echo "Deploying Ingress Controller..."
bash ingress-controller-setup.sh

# Apply the Ingress Controller customization
echo "Applying Ingress Controller customization..."
kubectl apply -f ingress-controller-customization.yaml

# Deploy the service
echo "Deploying the service..."
kubectl apply -f service-definition.yaml

# Deploy the Ingress
echo "Deploying the Ingress..."
kubectl apply -f ingress-definition.yaml

echo "Deployment completed."
