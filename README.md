# k6 Extensions

This repository contains custom extensions and modifications for the [k6](https://k6.io/) load testing tool, specifically tailored for our organization's needs.

## Overview

This is an organization repository where we maintain our own version of k6 with custom extensions. We build and maintain these extensions to add functionality that's specific to our use cases and infrastructure.

## Extensions

### SignalFlow Extension
- **Purpose**: Enables integration with SignalFlow for real-time metrics processing and analysis
- **Features**:
    - Execute SignalFlow programs
    - Stream and process metrics in real-time
    - Handle large-scale data collection efficiently

## Building Custom k6

To build k6 with our custom extensions:

```bash
# Clone the repository
git clone [https://github.com/your-org/k6-splunk-extensions.git](https://github.com/your-org/k6-splunk-extensions.git)
cd k6-splunk-extensions

# Build k6 with all extensions
make