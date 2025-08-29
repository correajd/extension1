# k6 Extensions

This repository contains custom extensions and modifications for the [k6](https://k6.io/) load testing tool, specifically tailored for our organization's needs.

## Overview

This is an organization repository where we maintain our own version of k6 with custom extensions. We build and maintain these extensions to add functionality that's specific to our use cases and infrastructure.

## Build Artifacts

We maintain the following build artifacts for CI/CD usage:

- **Pre-built Binaries**: Compiled k6 executables for various platforms (Linux, macOS, Windows)
- **Docker Images**: Containerized versions of k6 with our extensions pre-installed
  - Available for multiple architectures (amd64, arm64)
  - Tagged with both version numbers and `latest`

These artifacts are automatically built and published through our CI/CD pipeline and can be used directly in your testing environments.

## Extensions

### SignalFlow Extension
- **Purpose**: Enables integration with SignalFlow for real-time metrics processing and analysis
- **Features**:
    - Execute SignalFlow programs
    - Stream and process metrics in real-time
    - Handle large-scale data collection efficiently

## Building Custom k6

### Local Development

To build k6 with our custom extensions locally:

```bash
# Clone the repository
git clone [https://cd.splunkdev.com/observability/qe/k6-splunk-extensions.git](https://cd.splunkdev.com/observability/qe/k6-splunk-extensions.git)
cd k6-splunk-extensions

# Build k6 with all extensions
make
```

### CI/CD Integration

For CI/CD pipelines, we recommend using our pre-built artifacts:

```yaml
# Example GitHub Actions usage
jobs:
  load-test:
    runs-on: ubuntu-latest
    container:
      image: your-registry/k6-splunk:latest
    steps:
      - name: Run k6 test
        run: k6 run test.js
```

Pre-built binaries and Docker images are available in our artifact repository and container registry respectively.

