#!/bin/bash

set -e

NAMESPACE="golang-service"
IMAGE_TAG=${1:-latest}

echo "Deploying Golang Service to Kubernetes..."

# Create namespace if it doesn't exist
kubectl create namespace $NAMESPACE --dry-run=client -o yaml | kubectl apply -f -

# Apply Kubernetes manifests
echo "Applying Kubernetes manifests..."
kubectl apply -f deployments/kubernetes/

# Wait for deployment to be ready
echo "Waiting for deployment to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/golang-service -n $NAMESPACE

# Update image if specified
if [ "$IMAGE_TAG" != "latest" ]; then
    echo "Updating image to tag: $IMAGE_TAG"
    kubectl set image deployment/golang-service golang-service=golang-service:$IMAGE_TAG -n $NAMESPACE
    kubectl rollout status deployment/golang-service -n $NAMESPACE
fi

echo "Deployment completed successfully!"

# Show deployment status
kubectl get pods -n $NAMESPACE
kubectl get services -n $NAMESPACE