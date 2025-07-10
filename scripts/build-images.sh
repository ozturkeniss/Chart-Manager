#!/bin/bash

# Build Docker images for Kubernetes deployment
set -e

echo "🚀 Building Docker images for Kubernetes deployment..."

# Set image registry (change this to your registry)
REGISTRY="rancher-manager"
TAG="latest"

# Build AuthService
echo "📦 Building AuthService image..."
docker build -f docker/Dockerfile.authservice -t ${REGISTRY}-authservice:${TAG} .

# Build ItemService
echo "📦 Building ItemService image..."
docker build -f docker/Dockerfile.itemservice -t ${REGISTRY}-itemservice:${TAG} .

# Build InventoryService
echo "📦 Building InventoryService image..."
docker build -f docker/Dockerfile.inventoryservice -t ${REGISTRY}-inventoryservice:${TAG} .

# Build FiberGateway
echo "📦 Building FiberGateway image..."
docker build -f docker/Dockerfile.fibergateway -t ${REGISTRY}-fibergateway:${TAG} .

echo "✅ All images built successfully!"
echo ""
echo "📋 Built images:"
echo "  - ${REGISTRY}-authservice:${TAG}"
echo "  - ${REGISTRY}-itemservice:${TAG}"
echo "  - ${REGISTRY}-inventoryservice:${TAG}"
echo "  - ${REGISTRY}-fibergateway:${TAG}"
echo ""
echo "💡 To push to registry, run:"
echo "  docker tag ${REGISTRY}-authservice:${TAG} your-registry/${REGISTRY}-authservice:${TAG}"
echo "  docker push your-registry/${REGISTRY}-authservice:${TAG}"
echo "  (repeat for other services)" 