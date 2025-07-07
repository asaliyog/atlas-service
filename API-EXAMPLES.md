# Cloud Inventory API Examples

This document provides comprehensive examples of how to use the Cloud Inventory API endpoints.

## Authentication

All API endpoints require authentication via Azure Entra ID. Include your bearer token in the Authorization header:

```bash
Authorization: Bearer YOUR_JWT_TOKEN
```

## Base URL

```
http://localhost:8080/api/v1
```

## Endpoints

### GET /vms

Retrieve a paginated list of virtual machines with optional filtering and sorting.

#### Basic Usage

```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     "http://localhost:8080/api/v1/vms"
```

#### Pagination

```bash
# Get page 2 with 10 items per page
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     "http://localhost:8080/api/v1/vms?page=2&pageSize=10"
```

#### Sorting

```bash
# Sort by name in ascending order
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     "http://localhost:8080/api/v1/vms?sortBy=name&sortOrder=asc"

# Sort by creation date in descending order
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     "http://localhost:8080/api/v1/vms?sortBy=createdAt&sortOrder=desc"

# Sort by cloud-specific field
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     "http://localhost:8080/api/v1/vms?sortBy=cloudSpecificDetails.vpcId&sortOrder=asc"
```

#### Filtering

##### Single Filter Examples

```bash
# Filter by cloud type
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"cloudType\",\"operator\":\"eq\",\"value\":\"aws\"}]"

# Filter by status
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"status\",\"operator\":\"eq\",\"value\":\"running\"}]"

# Filter by name contains
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"name\",\"operator\":\"contains\",\"value\":\"web\"}]"

# Filter by instance type
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"instanceType\",\"operator\":\"eq\",\"value\":\"t2.micro\"}]"
```

##### Multiple Filters

```bash
# Filter by cloud type AND status
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"cloudType\",\"operator\":\"eq\",\"value\":\"aws\"},{\"field\":\"status\",\"operator\":\"eq\",\"value\":\"running\"}]"

# Filter by cloud account and location
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"cloudAccountId\",\"operator\":\"eq\",\"value\":\"123456789012\"},{\"field\":\"location\",\"operator\":\"eq\",\"value\":\"us-east-1\"}]"
```

##### Cloud-Specific Field Filtering

```bash
# Filter AWS VMs by VPC ID
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"cloudSpecificDetails.vpcId\",\"operator\":\"eq\",\"value\":\"vpc-12345678\"}]"

# Filter GCP VMs by machine type
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"cloudSpecificDetails.machineType\",\"operator\":\"eq\",\"value\":\"e2-standard-2\"}]"

# Filter Azure VMs by resource group
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"cloudSpecificDetails.resourceGroup\",\"operator\":\"eq\",\"value\":\"my-rg\"}]"
```

##### Advanced Filtering with Operators

```bash
# VMs created in the last 24 hours
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"createdAt\",\"operator\":\"gt\",\"value\":\"2025-01-26T00:00:00Z\"}]"

# VMs with names NOT containing "test"
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"name\",\"operator\":\"neq\",\"value\":\"test\"}]"
```

#### Combined Parameters

```bash
# Get AWS VMs, running status, sorted by name, page 1 with 5 items
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     "http://localhost:8080/api/v1/vms?page=1&pageSize=5&sortBy=name&sortOrder=asc&filter=[{\"field\":\"cloudType\",\"operator\":\"eq\",\"value\":\"aws\"},{\"field\":\"status\",\"operator\":\"eq\",\"value\":\"running\"}]"
```

## Response Format

### Success Response

```json
{
  "data": [
    {
      "id": "i-1234567890abcdef0",
      "name": "web-server-01",
      "cloudType": "aws",
      "status": "running",
      "createdAt": "2025-01-25T16:16:00Z",
      "updatedAt": "2025-01-26T16:16:00Z",
      "cloudAccountId": "123456789012",
      "location": "us-east-1",
      "instanceType": "t2.micro",
      "cloudSpecificDetails": {
        "cloudType": "aws",
        "vpcId": "vpc-12345678",
        "subnetId": "subnet-12345678",
        "securityGroupIds": ["sg-12345678", "sg-98765432"],
        "privateIpAddress": "10.0.1.100",
        "publicIpAddress": "54.123.45.67"
      }
    }
  ],
  "pagination": {
    "page": 1,
    "pageSize": 20,
    "totalItems": 100,
    "totalPages": 5
  }
}
```

### Error Response

```json
{
  "message": "Invalid filter parameter"
}
```

## Filter Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `eq` | Equal to | `{"field":"cloudType","operator":"eq","value":"aws"}` |
| `neq` | Not equal to | `{"field":"status","operator":"neq","value":"stopped"}` |
| `contains` | Contains substring | `{"field":"name","operator":"contains","value":"web"}` |
| `lt` | Less than | `{"field":"createdAt","operator":"lt","value":"2025-01-01T00:00:00Z"}` |
| `gt` | Greater than | `{"field":"createdAt","operator":"gt","value":"2025-01-01T00:00:00Z"}` |
| `le` | Less than or equal | `{"field":"createdAt","operator":"le","value":"2025-01-01T00:00:00Z"}` |
| `ge` | Greater than or equal | `{"field":"createdAt","operator":"ge","value":"2025-01-01T00:00:00Z"}` |

## Cloud-Specific Fields

### AWS Fields
- `cloudSpecificDetails.vpcId` - VPC ID
- `cloudSpecificDetails.subnetId` - Subnet ID
- `cloudSpecificDetails.securityGroupIds` - Security group IDs (array)
- `cloudSpecificDetails.privateIpAddress` - Private IP address
- `cloudSpecificDetails.publicIpAddress` - Public IP address

### GCP Fields
- `cloudSpecificDetails.machineType` - Machine type
- `cloudSpecificDetails.network` - Network name
- `cloudSpecificDetails.region` - Region
- `cloudSpecificDetails.privateIpAddress` - Private IP address
- `cloudSpecificDetails.publicIpAddress` - Public IP address

### Azure Fields
- `cloudSpecificDetails.resourceGroup` - Resource group name
- `cloudSpecificDetails.vmSize` - VM size
- `cloudSpecificDetails.privateIpAddress` - Private IP address
- `cloudSpecificDetails.publicIpAddress` - Public IP address

## HTTP Status Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 400 | Bad Request (invalid parameters) |
| 401 | Unauthorized (missing or invalid token) |
| 500 | Internal Server Error |

## Rate Limiting

- Maximum page size: 100 items
- Default page size: 20 items
- Pages are 1-indexed (first page is page 1)

## JavaScript Examples

### Using Fetch API

```javascript
async function getVMs(filters = [], page = 1, pageSize = 20) {
  const params = new URLSearchParams({
    page: page.toString(),
    pageSize: pageSize.toString()
  });
  
  if (filters.length > 0) {
    params.append('filter', JSON.stringify(filters));
  }
  
  const response = await fetch(`/api/v1/vms?${params}`, {
    headers: {
      'Authorization': `Bearer ${YOUR_JWT_TOKEN}`,
      'Content-Type': 'application/json'
    }
  });
  
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }
  
  return await response.json();
}

// Usage examples
const allVMs = await getVMs();
const awsVMs = await getVMs([{field: 'cloudType', operator: 'eq', value: 'aws'}]);
const runningVMs = await getVMs([{field: 'status', operator: 'eq', value: 'running'}]);
```

### Using Axios

```javascript
import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8080/api/v1',
  headers: {
    'Authorization': `Bearer ${YOUR_JWT_TOKEN}`
  }
});

// Get all VMs
const response = await api.get('/vms');

// Get filtered VMs
const awsVMs = await api.get('/vms', {
  params: {
    filter: JSON.stringify([{field: 'cloudType', operator: 'eq', value: 'aws'}])
  }
});
```

## Python Examples

### Using requests

```python
import requests
import json

def get_vms(token, filters=None, page=1, page_size=20, sort_by='id', sort_order='asc'):
    url = 'http://localhost:8080/api/v1/vms'
    headers = {'Authorization': f'Bearer {token}'}
    params = {
        'page': page,
        'pageSize': page_size,
        'sortBy': sort_by,
        'sortOrder': sort_order
    }
    
    if filters:
        params['filter'] = json.dumps(filters)
    
    response = requests.get(url, headers=headers, params=params)
    response.raise_for_status()
    return response.json()

# Usage examples
token = 'YOUR_JWT_TOKEN'
all_vms = get_vms(token)
aws_vms = get_vms(token, filters=[{'field': 'cloudType', 'operator': 'eq', 'value': 'aws'}])
```

## Testing the API

### Running Tests

```bash
# Run all tests
make test

# Run only VM handler tests
go test -v ./internal/handlers -run TestVMHandler
```

### Seeding Test Data

```bash
# Seed the database with sample VM data
make seed
```

### Manual Testing with curl

```bash
# Test basic endpoint
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     "http://localhost:8080/api/v1/vms"

# Test with filters
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"cloudType\",\"operator\":\"eq\",\"value\":\"aws\"}]"
```