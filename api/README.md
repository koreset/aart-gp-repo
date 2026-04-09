# ADS API

This repository contains the API for the ADS (AART) system.

## Testing projection worker concurrency

See docs/testing-projection-worker.md for step-by-step instructions to verify that only one instance can claim a queued projection job at a time and to run the provided helper tool.

## Deployment

This project uses GitHub Actions for automated building and deployment to various environments.

### Automated Deployment

The project is automatically deployed to the Performance environment when changes are pushed to the main/master branch or when a new release is published.

### Manual Deployment

You can manually deploy to any environment (Performance, Staging, Zambia, Afrihost) using the GitHub Actions workflow:

1. Go to the "Actions" tab in the GitHub repository
2. Select the "Build and Deploy" workflow
3. Click "Run workflow"
4. Select the target environment
5. Click "Run workflow"

For more details on the deployment process, see the [GitHub Actions documentation](.github/README.md).

## Local Development

To build and run the project locally:

```bash
# Build the project
go build -o aart_api

# Run the API
./aart_api
```

## Resilience and Backoff

- The DB layer now has resilient read helper with exponential backoff and jitter, short context timeouts, and a per-instance concurrency gate.
- Usage example:
  - err := services.DBReadWithResilience(ctx, func(d *gorm.DB) error { return d.Where("...").First(&out).Error })
- MySQL DSN includes I/O timeouts to fail fast on stuck connections.

## Reports

- Location: data/reports
- Description: This folder stores generated reports and exported analysis artifacts produced by the API/tools.
- Example: You can place or find report files such as CSV, Excel, or Markdown summaries here (e.g., data/reports/).
- Related logs used for investigations are under logs/downloaded/ (e.g., logs/downloaded/app1/logs/app.log).

## Deployment Script

For manual deployment, you can use the `deploy.sh` script:

```bash
# Deploy to Performance environment
./deploy.sh
```
