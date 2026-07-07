# Terraform Plugin Framework (`thales` provider)

## Overview

This repository is a hands-on learning project for understanding and building **Custom Terraform Providers using the Terraform Plugin Framework**.

The objective is to simulate a real enterprise environment similar to the Thales project, where Terraform providers interact with proprietary APIs rather than standard cloud providers like AWS, Azure, or GCP.

This lab focuses on:

* Terraform Plugin Framework (Latest Standard)
* Custom Provider Development in Go
* Terraform Provider Protocol (gRPC)
* Resource Lifecycle Management (CRUD)
* REST API Integration
* State Management
* ImportState Implementation
* Acceptance Testing
* Module Development
* Terragrunt Concepts

---

# Architecture

```text
Terraform CLI
      │
      ▼
Terraform Core
      │
      ▼
Terraform Provider Protocol (gRPC)
      │
      ▼
terraform-provider-thales
      │
      ▼
REST API Client
      │
      ▼
Mock API Server
      │
      ▼
In-Memory Database
```

This architecture closely resembles real enterprise custom provider implementations.

---

# Technology Stack

| Component                  | Version |
| -------------------------- | ------- |
| Terraform                  | 1.15.6  |
| Go                         | 1.26.4  |
| Docker Desktop             | 4.77    |
| Docker Engine              | 29.5.3  |
| kubectl                    | 1.34.1  |
| Terraform Plugin Framework | Latest  |

---

# Project Structure

```text
terraform-plugin-framework/

├── README.md
│
├── main.go
├── go.mod
├── go.sum
│
├── provider/
│   ├── provider.go
│   ├── config.go
│   ├── client.go
│   ├── models.go
│   └── keystore_resource.go
│
├── mock-api/
│   ├── go.mod
│   ├── main.go
│   └── mock-api
│
├── examples/
│   └── basic/
│       └── main.tf
│
├── docs/
│
└── terraform-provider-thales
```

---

# Learning Objectives

The goal of this lab is to understand:

## Terraform Core

* HCL Parsing
* Execution Graph Generation
* State Management
* Plan Generation
* Drift Detection

---

## Provider Development

Understand:

```text
Metadata()

Schema()

Configure()

Resources()

DataSources()
```

---

## Resource Development

Every resource must implement:

```text
Metadata()

Schema()

Create()

Read()

Update()

Delete()

ImportState()
```

---

## Plugin Framework vs SDK v2

| SDK v2           | Plugin Framework  |
| ---------------- | ----------------- |
| schema.Provider  | provider.Provider |
| schema.Resource  | resource.Resource |
| StateFunc        | PlanModifiers     |
| DiffSuppressFunc | Validators        |
| Legacy           | Recommended       |

The Thales project primarily uses the latest Plugin Framework.

---

# Current Implementation

## Provider

Current provider:

```go
provider "thales"
```

Provider source:

```text
provider/thales
```

---

## Resource

Current resource:

```hcl
resource "thales_keystore" "payment" {

    name = "payment-keystore"

}
```

---

## Mock API

Local API server:

```text
http://localhost:8080
```

Endpoints implemented so far:

```text
POST    /keystores
GET     /keystores/{id}
```

`PUT /keystores/{id}` and `DELETE /keystores/{id}` are not implemented yet — this is the next piece of work (see Day 3/4 status below).

---

# Setup Instructions

## Clone Repository

```bash
git clone <repository>

cd terraform-plugin-framework
```

---

## Verify Tools

Terraform:

```bash
terraform version
```

---

Go:

```bash
go version
```

---

Docker:

```bash
docker version
```

---

kubectl:

```bash
kubectl version --client
```

---

# Build Provider

Build the provider:

```bash
go build -o terraform-provider-thales
```

Verify:

```bash
file terraform-provider-thales
```

Expected:

```text
Mach-O 64-bit executable
```

---

# Start Mock API

Move to:

```bash
cd mock-api
```

Build:

```bash
go build
```

Run:

```bash
./mock-api
```

Expected:

```text
Mock API running on :8080
```

---

# Test Mock API

Create resource:

```bash
curl -X POST localhost:8080/keystores
```

Response:

```json
{
    "id":"ks-001",
    "name":"payment-keystore"
}
```

---

Read resource:

```bash
curl localhost:8080/keystores/ks-001
```

Response:

```json
{
    "id":"ks-001",
    "name":"payment-keystore"
}
```

---

# Terraform Provider Internal Flow

Terraform does not directly communicate with infrastructure.

Internal execution:

```text
terraform apply

↓

Terraform Core

↓

Provider Discovery

↓

Launch Provider Binary

↓

gRPC Handshake

↓

Configure()

↓

Resource Create()

↓

REST API Call

↓

State Update
```

---

# Resource Lifecycle

The Terraform resource lifecycle consists of:

```text
Create()

↓

Read()

↓

Update()

↓

Delete()

↓

ImportState()
```

---

## Create

Responsible for:

```text
POST /resource
```

Creates infrastructure.

---

## Read

Responsible for:

```text
GET /resource/{id}
```

Refreshes Terraform state.

Used for:

* Drift Detection
* terraform plan
* terraform refresh

---

## Update

Responsible for:

```text
PUT /resource/{id}
```

Modifies infrastructure.

---

## Delete

Responsible for:

```text
DELETE /resource/{id}
```

Destroys infrastructure.

---

## ImportState

Responsible for:

```bash
terraform import
```

Allows onboarding of manually created resources.

---

# Development Roadmap

## Day 1

Environment Setup

* Terraform
* Go
* Docker
* Mock API

Status:

```text
COMPLETED
```

---

## Day 2

Plugin Framework Fundamentals

* Provider Skeleton
* Resource Skeleton
* CRUD Interfaces

Status:

```text
COMPLETED
```

---

## Day 3

Real REST Integration

Implement:

```text
Create()

Read()

Update()

Delete()
```

using:

```text
localhost:8080
```

Status:

```text
IN PROGRESS
```

* Create() / Read() — done, wired through a shared *APIClient injected via provider Configure() → resource Configure() (standard Plugin Framework dependency-injection pattern)
* Update() / Delete() — still stubs; mock API has no PUT/DELETE routes yet

---

## Day 4

State Management

Topics:

* ImportState
* Drift Detection
* 404 Handling
* Resource Removal

Status:

```text
IN PROGRESS
```

* 404 handling — done: `GetKeystore` returns a sentinel `ErrNotFound` on a 404 response, and `Read()` calls `resp.State.RemoveResource(ctx)` instead of erroring
* ImportState — not started

---

## Day 5

Plan Modifiers

Topics:

* RequiresReplace
* Validators
* Immutable Attributes
* Sensitive Fields

---

## Day 6

Modules and Terragrunt

Topics:

* Module Design
* Reusability
* Multi-environment Structure
* Terragrunt Integration

---

## Day 7

Testing and CI/CD

Topics:

* Acceptance Tests
* Unit Tests
* GitHub Actions
* Provider Releases

---

# Interview Preparation Notes

## Common Questions

### How does Terraform communicate with Providers?

Answer:

Terraform Core communicates with Providers using the Terraform Provider Protocol over gRPC.

---

### Why is Go used for Terraform Providers?

Answer:

HashiCorp officially supports Go for Provider development because the Plugin Framework and SDK are implemented in Go.

---

### Difference between SDK v2 and Plugin Framework?

Answer:

Plugin Framework is the latest recommended approach and provides:

* Strong typing
* Better validation
* Plan Modifiers
* Improved state handling
* Long-term support

---

### How do custom providers work?

Answer:

Custom providers act as an abstraction layer between Terraform and proprietary platforms by implementing CRUD operations over REST APIs or SDKs.

---

# Future Enhancements

Planned improvements (none of these are done yet):

```text
☐ ImportState

☐ Update API

☐ Delete API

☐ Acceptance Tests

☐ Docker Deployment

☐ Terragrunt Integration

☐ Provider Versioning

☐ Release Automation

☐ Multi-resource Support
```

---

# Author

Vikash Jaiswal

Terraform Plugin Framework

2026
