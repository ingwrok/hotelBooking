# Builder Stage
FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .

# Pass build args if needed, but usually ENV vars are runtime in browser
# but for VITE_API_URL it needs to be burned in or replaced at runtime.
# For simplicity, we assume generic build or env injection at runtime (which is hard for static files).
# We will use the nginx proxy /api -> backend, so VITE_API_URL can be relative usually.
RUN npm run build

# Runner Stage
FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
