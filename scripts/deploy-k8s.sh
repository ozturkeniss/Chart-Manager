#!/bin/bash

# Deploy to Kubernetes using Helm
set -e

echo "🚀 Deploying Rancher Manager to Kubernetes..."

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl is not installed. Please install kubectl first."
    exit 1
fi

# Check if helm is installed
if ! command -v helm &> /dev/null; then
    echo "❌ helm is not installed. Please install helm first."
    exit 1
fi

# Check if we're connected to a cluster
if ! kubectl cluster-info &> /dev/null; then
    echo "❌ Not connected to a Kubernetes cluster. Please connect to a cluster first."
    exit 1
fi

# Create namespace if it doesn't exist
NAMESPACE="rancher-manager"
echo "📦 Creating namespace: ${NAMESPACE}"
kubectl create namespace ${NAMESPACE} --dry-run=client -o yaml | kubectl apply -f -

# Install/upgrade the Helm chart
RELEASE_NAME="rancher-manager"
echo "📦 Installing/upgrading Helm chart: ${RELEASE_NAME}"

# Check if release exists
if helm list -n ${NAMESPACE} | grep -q ${RELEASE_NAME}; then
    echo "🔄 Upgrading existing release..."
    helm upgrade ${RELEASE_NAME} ./helm -n ${NAMESPACE} --wait --timeout=10m
else
    echo "🆕 Installing new release..."
    helm install ${RELEASE_NAME} ./helm -n ${NAMESPACE} --wait --timeout=10m
fi

echo "✅ Deployment completed successfully!"
echo ""
echo "📋 Check deployment status:"
echo "  kubectl get pods -n ${NAMESPACE}"
echo "  kubectl get services -n ${NAMESPACE}"
echo "  kubectl get ingress -n ${NAMESPACE}"
echo ""
echo "🌐 Access the application:"
echo "  kubectl port-forward -n ${NAMESPACE} svc/${RELEASE_NAME}-fibergateway 8080:8080"
echo "  Then visit: http://localhost:8080"
echo ""
echo "📊 View logs:"
echo "  kubectl logs -n ${NAMESPACE} -l component=gateway"
echo "  kubectl logs -n ${NAMESPACE} -l service=authservice"
echo "  kubectl logs -n ${NAMESPACE} -l service=itemservice"
echo "  kubectl logs -n ${NAMESPACE} -l service=inventoryservice" 