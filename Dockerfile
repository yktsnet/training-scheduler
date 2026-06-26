# Stage 1: frontend-build
FROM node:22-alpine AS frontend-build
WORKDIR /app/frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Stage 2: backend-build
FROM golang:1.25-alpine AS backend-build
WORKDIR /app/backend
COPY backend/ ./
COPY --from=frontend-build /app/frontend/dist ./dist
RUN CGO_ENABLED=0 go build -o training-app .

# Stage 3: runtime
FROM gcr.io/distroless/static-debian12
COPY --from=backend-build /app/backend/training-app /training-app
EXPOSE 5000
ENTRYPOINT ["/training-app"]
