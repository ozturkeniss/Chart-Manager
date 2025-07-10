#!/bin/bash

# Update Helm dependencies
set -e

echo "📦 Updating Helm dependencies..."

# Check if helm is installed
if ! command -v helm &> /dev/null; then
    echo "❌ helm is not installed. Please install helm first."
    exit 1
fi

# Add Bitnami repository if not already added
echo "🔗 Adding Bitnami repository..."
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# Update dependencies
echo "📦 Updating dependencies..."
helm dependency update ./helm

echo "✅ Dependencies updated successfully!"
echo ""
echo "📋 Updated dependencies:"
echo "  - PostgreSQL (Bitnami)"
echo "  - MongoDB (Bitnami)"
echo "  - Redis (Bitnami)"
echo "  - Kafka (Bitnami)"
echo ""
echo "💡 To install/upgrade the chart:"
echo "  ./scripts/deploy-k8s.sh" 