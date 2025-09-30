#!/bin/bash
# Setup script to install and run both backend and frontend

echo "Installing backend dependencies..."
cd backend
go mod download

echo "Building backend..."
go build -o bin/server ./cmd/server

echo "Starting backend server..."
./bin/server &

cd ..

echo "Installing frontend dependencies..."
cd frontend
npm install --legacy-peer-deps

echo "Starting frontend development server..."
npx expo start &

echo "Setup complete! Backend running in background (http://localhost:8080), frontend should open in browser or Expo Go app."