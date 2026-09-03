# Deployment Guide: Cloud Build CI/CD & Local CLI

This guide describes how to configure automated CI/CD deployments using **Google Cloud Build** and how to execute manual deployments from your local machine.

---

## 1. Automated CI/CD with Google Cloud Build

When configured, every commit or pull request pushed to your repository can automatically test, build, and deploy both the Go backend to **Cloud Run** and the React frontend to **Firebase Hosting**.

### 1.1. Required IAM Roles for Cloud Build Service Account

By default, Cloud Build uses a service account with the format:
`[PROJECT_NUMBER]@cloudbuild.gserviceaccount.com`

Grant the following roles to the Cloud Build service account in your Google Cloud Console (**IAM & Admin > IAM**):

1. **Cloud Run Admin** (`roles/run.admin`) - To create and update Cloud Run revisions.
2. **Service Account User** (`roles/iam.serviceAccountUser`) - To deploy containers acting as the default compute service account.
3. **Firebase Hosting Admin** (`roles/firebasehosting.admin`) - To upload and release static sites to Firebase Hosting.

You can grant these roles via the `gcloud` CLI:

```bash
PROJECT_ID="your-gcp-project-id"
PROJECT_NUMBER=$(gcloud projects describe "$PROJECT_ID" --format='value(projectNumber)')
CB_SA="${PROJECT_NUMBER}@cloudbuild.gserviceaccount.com"

# 1. Cloud Run Admin
gcloud projects add-iam-policy-binding "$PROJECT_ID" \
  --member="serviceAccount:$CB_SA" \
  --role="roles/run.admin"

# 2. Service Account User
gcloud projects add-iam-policy-binding "$PROJECT_ID" \
  --member="serviceAccount:$CB_SA" \
  --role="roles/iam.serviceAccountUser"

# 3. Firebase Hosting Admin
gcloud projects add-iam-policy-binding "$PROJECT_ID" \
  --member="serviceAccount:$CB_SA" \
  --role="roles/firebasehosting.admin"
```

---

### 1.2. Creating the Cloud Build Trigger

1. In Google Cloud Console, navigate to **Cloud Build > Triggers**.
2. Click **Create Trigger**.
3. Set the following fields:
   - **Name**: `deploy-calculator-on-commit`
   - **Event**: `Push to a branch`
   - **Repository**: Connect your GitHub/GitLab repository.
   - **Branch**: `^main$`
   - **Configuration**: `Cloud Build configuration file (yaml or json)`
   - **Location**: `Repository` -> `/cloudbuild.yaml`
4. Expand **Advanced > Substitution variables** and configure your parameters:

| Variable Name | Description | Default / Example Value |
| :--- | :--- | :--- |
| `_REGION` | Google Cloud region for Cloud Run | `us-central1` |
| `_SERVICE_NAME` | Cloud Run service name | `calculator-backend` |
| `_FIREBASE_PROJECT_ID` | Target Firebase Project ID | *(Leave empty to use `$PROJECT_ID`)* |
| `_BACKEND_API_URL` | Override backend API URL | *(Leave empty to auto-resolve)* |

> [!NOTE]
> If `_BACKEND_API_URL` is left blank, the pipeline automatically discovers the live Cloud Run HTTPS service URL and injects it into the frontend build dynamically.

---

## 2. Manual Local Deployment

You can deploy both layers directly from your local terminal using `gcloud` and `firebase-tools`.

### Prerequisites
- **Google Cloud SDK (`gcloud`)** installed and authenticated (`gcloud auth login`).
- **Node.js 20+** via `nvm`.
- Billing enabled on your Google Cloud project.

---

### Step 1: Deploy Go Backend to Cloud Run

From the repository root, deploy the Go container directly from source:

```bash
# Define target project and region
PROJECT_ID="your-gcp-project-id"
REGION="us-central1"
SERVICE_NAME="calculator-backend"

# Deploy backend
gcloud run deploy "$SERVICE_NAME" \
  --source ./backend \
  --project "$PROJECT_ID" \
  --region "$REGION" \
  --platform managed \
  --allow-unauthenticated \
  --port 8080 \
  --set-env-vars "ALLOWED_ORIGINS=*"
```

Upon completion, copy the **Service URL** output (e.g. `https://calculator-backend-xxxxx-uc.a.run.app`).

---

### Step 2: Build the Frontend

Load Node.js 20 using `nvm` and build the production bundle, specifying the Cloud Run backend URL:

```bash
cd frontend

# 1. Activate Node.js 20
source ~/.nvm/nvm.sh && nvm use 20

# 2. Set backend API URL and build
VITE_API_BASE_URL="https://calculator-backend-xxxxx-uc.a.run.app/api/v1" npm run build
```

The production assets will be generated in `frontend/dist`.

---

### Step 3: Deploy Frontend to Firebase Hosting

Deploy the compiled static bundle using `firebase-tools`:

```bash
# From repository root
cd /path/to/calculator

# Activate Node.js 20
source ~/.nvm/nvm.sh && nvm use 20

# Deploy hosting
npx -y firebase-tools deploy --only hosting --project "$PROJECT_ID"
```

Firebase Hosting will publish your assets to a global CDN and output your live web URL:
`https://[PROJECT_ID].web.app`

---

## 3. Local Development vs Cloud Deployment Summary

| Environment | Frontend URL | Backend URL | Setup Command |
| :--- | :--- | :--- | :--- |
| **Local Docker** | `http://localhost:3001` | `http://localhost:8080` | `docker compose up --build` |
| **Local Dev Server** | `http://localhost:5173` | `http://localhost:8080` | `go run cmd/api/main.go` & `npm run dev` |
| **Cloud Production** | `https://[PROJECT].web.app` | `https://[SERVICE].run.app` | See steps above or `cloudbuild.yaml` |
