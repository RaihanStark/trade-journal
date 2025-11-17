    FROM node:20-alpine AS builder

    # Enable pnpm
    RUN corepack enable && corepack prepare pnpm@latest --activate
    
    WORKDIR /app
    
    # Copy only package files first (great for Docker caching)
    COPY package.json pnpm-lock.yaml ./
    
    # Install deps using pnpm
    RUN pnpm install --frozen-lockfile
    
    # Copy the entire project
    COPY . .
    
    # Build argument for API URL
    ARG VITE_API_URL=http://localhost:8080
    ENV VITE_API_URL=$VITE_API_URL
    
    # Increase memory for build (helps SvelteKit/Vite)
    ENV NODE_OPTIONS="--max-old-space-size=2048"
    
    # Build the app (pnpm version)
    RUN pnpm build
    
    
    # ---- Production Stage ----
    FROM nginx:alpine
    
    # Copy custom nginx config
    COPY nginx.conf /etc/nginx/conf.d/default.conf
    
    # Copy built assets
    COPY --from=builder /app/build /usr/share/nginx/html
    
    EXPOSE 80
    CMD ["nginx", "-g", "daemon off;"]
    