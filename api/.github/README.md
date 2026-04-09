git # GitHub Actions Workflows

This directory contains GitHub Actions workflows for automating various processes in the project.

## Deploy API Workflow

The `deploy.yml` workflow automates the deployment of the API to the production environment. It builds the Go binary, handles migrations, and deploys to three services (app1, app2, app3) on the remote server.

### Workflow Triggers

The workflow is triggered by:
- Pushing to the `main` or `master` branch
- Manual trigger using the "Run workflow" button in the GitHub Actions UI

### Required Secrets

The workflow requires the following secrets to be set in the GitHub repository:

- `SSH_PRIVATE_KEY`: The private SSH key for connecting to the remote server. This should be the content of the `aart_key_access.pem` file.

To add these secrets:
1. Go to your GitHub repository
2. Click on "Settings"
3. Click on "Secrets and variables" > "Actions"
4. Click on "New repository secret"
5. Add the secret with the appropriate name and value

### Deployment Process

The workflow performs the following steps:

1. **Build**: Builds the Go binary for Linux
2. **Migrations**: Ensures the migrations directory exists and copies it to the server
3. **App1 Deployment**:
   - Creates migrations directory on the server
   - Copies migrations to the server
   - Stops app1 service and backs up the existing binary
   - Copies the new binary to app1
   - Sets executable permissions for app1
4. **App2 Deployment**:
   - Stops app2 service and backs up the existing binary
   - Copies the new binary to app2
   - Sets executable permissions for app2
5. **App3 Deployment**:
   - Stops app3 service and backs up the existing binary
   - Copies the new binary to app3
   - Sets executable permissions for app3
6. **Service Start**: Starts all three services
7. **Cleanup**: Removes the local binary

### Error Handling

The workflow includes comprehensive error handling:
- Each step is conditional on the success of previous steps
- Explicit error handling for critical steps
- Detailed status notifications at the end of the workflow

### Deployment Status

At the end of the workflow, a detailed deployment status summary is provided, including:
- Start and end times
- Status of each major step in the workflow
- Overall success or failure of the deployment

This information is valuable for troubleshooting if the deployment fails.