#!/bin/bash

# Define variables
CONTAINER_NAME="expense_tracker_backend"
NETWORK_NAME="opensip-apps"

echo "🚀 Starting deployment for $CONTAINER_NAME..."

# Step 1: Pull the latest code (uncomment if you need the script to do this)
# git pull origin main

# Step 2: Stop and remove the existing containers
echo "🛑 Stopping existing containers..."
docker-compose down

# Step 3: Build and start the new containers
echo "🏗️ Building and starting new containers..."
docker-compose up -d --build

# Step 4: Connect the backend container to the external network
echo "🔌 Connecting $CONTAINER_NAME to $NETWORK_NAME network..."
# Ignore the error if it's already attached
docker network connect $NETWORK_NAME $CONTAINER_NAME || true

# Step 5: Restart the backend container so it can connect to the database on the new network
echo "🔄 Restarting $CONTAINER_NAME..."
docker restart $CONTAINER_NAME

echo "✅ Deployment complete! Run 'docker logs -f $CONTAINER_NAME' to view the logs."
