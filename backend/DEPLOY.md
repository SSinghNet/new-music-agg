# Deployment Guide — Google Cloud Run

## Prerequisites

- [gcloud CLI](https://cloud.google.com/sdk/docs/install) installed and authenticated
- Docker installed
- An Artifact Registry repository created in your GCP project
- `DATABASE_URL` stored in Secret Manager:
  ```sh
  gcloud secrets create database-url --data-file=- <<< "postgres://user:password@db.supabase.co:5432/postgres?sslmode=require"
  ```

---

## Building Images

The Dockerfile uses a `CMD` build arg to select which binary to build (defaults to `api`).

```sh
# API server
docker build -t new-music-agg-api ./backend

# Scraper (one-shot job)
docker build --build-arg CMD=scraper -t new-music-agg-scraper ./backend
```

---

## Deploying the API (Cloud Run Service)

```sh
# Tag and push
docker tag new-music-agg-api us-central1-docker.pkg.dev/PROJECT_ID/REPO/api:latest
docker push us-central1-docker.pkg.dev/PROJECT_ID/REPO/api:latest

# Deploy
gcloud run deploy new-music-agg-api \
  --image us-central1-docker.pkg.dev/PROJECT_ID/REPO/api:latest \
  --region us-central1 \
  --set-secrets DATABASE_URL=database-url:latest \
  --allow-unauthenticated \
  --port 8080
```

---

## Deploying the Scraper (Cloud Run Job + Cloud Scheduler)

### 1. Push the scraper image

```sh
docker tag new-music-agg-scraper us-central1-docker.pkg.dev/PROJECT_ID/REPO/scraper:latest
docker push us-central1-docker.pkg.dev/PROJECT_ID/REPO/scraper:latest
```

### 2. Create the Cloud Run Job

```sh
gcloud run jobs create new-music-agg-scraper \
  --image us-central1-docker.pkg.dev/PROJECT_ID/REPO/scraper:latest \
  --region us-central1 \
  --set-secrets DATABASE_URL=database-url:latest \
  --max-retries 1 \
  --task-timeout 300
```

### 3. Create a service account for the scheduler

```sh
gcloud iam service-accounts create scraper-scheduler \
  --display-name "Scraper Scheduler"

gcloud run jobs add-iam-policy-binding new-music-agg-scraper \
  --region us-central1 \
  --member "serviceAccount:scraper-scheduler@PROJECT_ID.iam.gserviceaccount.com" \
  --role "roles/run.invoker"
```

### 4. Schedule with Cloud Scheduler (daily at 6am UTC)

```sh
gcloud scheduler jobs create http scraper-daily \
  --schedule "0 6 * * *" \
  --uri "https://us-central1-run.googleapis.com/apis/run.googleapis.com/v1/namespaces/PROJECT_ID/jobs/new-music-agg-scraper:run" \
  --oauth-service-account-email scraper-scheduler@PROJECT_ID.iam.gserviceaccount.com \
  --location us-central1
```

### Run manually

```sh
gcloud run jobs execute new-music-agg-scraper --region us-central1
```

---

## Environment Variables

| Variable       | Required | Description                        |
|----------------|----------|------------------------------------|
| `DATABASE_URL` | Yes      | Postgres connection string         |
| `PORT`         | No       | HTTP port (defaults to `8080`)     |
