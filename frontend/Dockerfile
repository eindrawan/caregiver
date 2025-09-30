# Use multi-stage build to keep the final image small
FROM node:20-alpine

# Set working directory
WORKDIR /app

# Install build dependencies for Expo and native modules
RUN apk add --no-cache python3 make g++ libc6-compat

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy source code
COPY .

# Install Expo CLI globally
RUN npm install -g @expo/cli

# Build the web app using Expo export
RUN npx expo export

# Install nginx
RUN apk add --no-cache nginx

# Remove default nginx files
RUN rm -rf /usr/share/nginx/html/*

# Copy the built web app to nginx directory
# The export command creates a dist directory
COPY dist/ /usr/share/nginx/html/

# Copy custom nginx configuration
COPY nginx.conf /etc/nginx.conf

EXPOSE 80

# Start nginx server
CMD ["nginx", "-g", "daemon off;"]