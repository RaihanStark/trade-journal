#!/bin/bash

# Trade Journal Deployment Script
# This script helps deploy the application using Docker Compose

set -e

echo "🚀 Trade Journal Deployment Script"
echo "=================================="
echo ""

# Check if docker and docker-compose are installed
if ! command -v docker &> /dev/null; then
    echo "❌ Error: Docker is not installed"
    exit 1
fi

if ! command -v docker compose &> /dev/null; then
    echo "❌ Error: Docker Compose is not installed"
    exit 1
fi

# Check if .env file exists
if [ ! -f .env ]; then
    echo "⚠️  Warning: .env file not found"
    echo "📝 Creating .env from .env.example..."
    cp .env.example .env
    echo "✅ .env file created. Please edit it with your configuration."
    echo ""
    read -p "Press Enter to continue after editing .env file..."
fi

# Ask for deployment mode
echo "Select deployment mode:"
echo "1) Development (local)"
echo "2) Production"
read -p "Enter choice (1 or 2): " mode

if [ "$mode" == "1" ]; then
    echo ""
    echo "🔧 Starting in Development mode..."
    echo ""

    # Build and start services
    docker-compose up -d --build

    echo ""
    echo "✅ Services started successfully!"
    echo ""
    echo "📋 Access points:"
    echo "   - Frontend:      http://localhost:3000"
    echo "   - Backend API:   http://localhost:8080"
    echo "   - MinIO Console: http://localhost:9001"
    echo ""
    echo "📊 View logs with: docker-compose logs -f"

elif [ "$mode" == "2" ]; then
    echo ""
    echo "🏭 Starting in Production mode..."
    echo ""

    # Check if required env vars are set
    if ! grep -q "JWT_SECRET=" .env || grep -q "JWT_SECRET=your-secret-key-change-this-in-production" .env; then
        echo "❌ Error: Please set a secure JWT_SECRET in .env file"
        exit 1
    fi

    # Build and start services with production overrides
    docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build

    echo ""
    echo "✅ Services started successfully in production mode!"
    echo ""
    echo "⚠️  Important: Make sure to set up a reverse proxy (Nginx/Traefik) with SSL/TLS"
    echo ""
    echo "📊 View logs with: docker-compose logs -f"

else
    echo "❌ Invalid choice"
    exit 1
fi

echo ""
echo "🎉 Deployment complete!"
