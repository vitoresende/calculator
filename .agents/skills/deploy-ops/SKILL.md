---
name: deploy-ops
description: >-
  Use this skill when managing Docker containerization for Go backend, Google Cloud Run deployments,
  Firebase Hosting configurations, and Google Cloud Build CI/CD pipelines.
---

# Deploy Operations & Cloud Architecture

This skill defines the deployment workflows, containerization strategies, and CI/CD automation rules.

## 1. Backend Containerization (Docker)
- Use multi-stage Docker builds to keep the final image minimal, secure, and fast.
- Stage 1 (`builder`): Official `golang:alpine` image with caching of `go.mod` and `go.sum`. Compile with `CGO_ENABLED=0 GOOS=linux`.
- Stage 2 (`runtime`): Minimal base image like `gcr.io/distroless/static:nonroot` or `alpine:latest` with a non-root user.
- Listen on the port provided by the `$PORT` environment variable (default to `8080` for local runs and Cloud Run).

## 2. Backend Deployment (Google Cloud Run)
- Cloud Run executes the containerized Go REST API as a serverless service.
- Flags:
  - `--platform managed`
  - `--region <REGION>`
  - `--allow-unauthenticated` (for public calculator API)
  - Memory and CPU: Lightweight (e.g. `--memory 256Mi --cpu 1`).

## 3. Frontend Deployment (Firebase Hosting)
- Frontend is built into static assets (e.g., in `frontend/dist`).
- Configuration via `firebase.json`:
  - Set `public` to the build output directory (`frontend/dist`).
  - Configure single-page application rewrites to `/index.html`.
- Deployment via Firebase CLI / Cloud Build: `firebase deploy --only hosting`.

## 4. CI/CD Pipeline (Google Cloud Build: `cloudbuild.yaml`)
- Trigger: Automatic execution on commits to the primary branch (`main`).
- Pipeline Steps:
  1. **Backend Tests**: Run `go test -v -race ./...`.
  2. **Frontend Tests**: Run `npm test`.
  3. **Build Backend Container**: `docker build -t gcr.io/$PROJECT_ID/calculator-backend:$COMMIT_SHA backend/`.
  4. **Push Backend Image**: Push image to Container Registry / Artifact Registry.
  5. **Deploy Backend to Cloud Run**: `gcloud run deploy calculator-backend --image ... --region ...`.
  6. **Build Frontend**: Install dependencies and build static assets in `frontend/`.
  7. **Deploy Frontend to Firebase Hosting**: Deploy via Firebase CLI container.
