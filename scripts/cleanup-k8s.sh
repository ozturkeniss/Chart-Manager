#!/bin/bash

# Cleanup Kubernetes deployment
set -e

echo "🧹 Cleaning up Rancher Manager from Kubernetes..."

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl is not installed."
    exit 1
fi

# Check if helm is installed
if ! command -v helm &> /dev/null; then
    echo "❌ helm is not installed."
    exit 1
fi

NAMESPACE="rancher-manager"
RELEASE_NAME="rancher-manager"

echo "🗑️  Uninstalling Helm release: ${RELEASE_NAME}"
helm uninstall ${RELEASE_NAME} -n ${NAMESPACE} --wait --timeout=5m

echo "🗑️  Deleting namespace: ${NAMESPACE}"
kubectl delete namespace ${NAMESPACE} --wait --timeout=5m

echo "✅ Cleanup completed successfully!"
echo ""
echo "💡 If you want to keep the data, you can:"
echo "  1. Backup PVCs before cleanup"
echo "  2. Use 'helm uninstall' without deleting namespace"
echo "  3. Manually delete specific resources" 