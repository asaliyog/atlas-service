# Cloud Inventory API Examples

## Authentication

### Local Development (Auth Bypassed)
```bash
# No authentication required for local development
curl -X GET "http://localhost:8080/api/v1/vms"
```

### Production (Azure Entra ID)
```bash
# Get a token using client credentials flow
TOKEN=$(curl -X POST "https://login.microsoftonline.com/{tenant-id}/oauth2/v2.0/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=client_credentials&client_id={client-id}&client_secret={client-secret}&scope=https://graph.microsoft.com/.default" \
  | jq -r '.access_token')

# Use the token in API requests
curl -X GET "http://localhost:8080/api/v1/vms" \
  -H "Authorization: Bearer $TOKEN"
```

## Basic VM Operations

### Get All VMs
```bash
curl -X GET "http://localhost:8080/api/v1/vms"
```

### Get VMs with Pagination
```bash
curl -X GET "http://localhost:8080/api/v1/vms?page=1&pageSize=20"
```

### Sort VMs by Creation Date (Descending)
```bash
curl -X GET "http://localhost:8080/api/v1/vms?sortBy=createdAt&sortOrder=desc"
```

## Advanced Filtering Examples

### 1. Filter by Cloud Type
```bash
# Get only AWS instances
curl -X GET "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"cloudType\",\"operator\":\"eq\",\"value\":\"aws\"}]"

# Get non-AWS instances
curl -X GET "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"cloudType\",\"operator\":\"ne\",\"value\":\"aws\"}]"
```

### 2. Filter by Status
```bash
# Get only running VMs
curl -X GET "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"status\",\"operator\":\"eq\",\"value\":\"running\"}]"

# Get stopped or terminated VMs
curl -X GET "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"status\",\"operator\":\"in\",\"value\":[\"stopped\",\"terminated\"]}]"
```

### 3. Filter by Instance Type
```bash
# Get specific instance types
curl -X GET "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"instanceType\",\"operator\":\"in\",\"value\":[\"t2.micro\",\"t3.small\",\"e2-standard-2\"]}]"

# Exclude small instances
curl -X GET "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"instanceType\",\"operator\":\"nin\",\"value\":[\"t2.nano\",\"t2.micro\"]}]"
```

### 4. Search by Name
```bash
# Case-sensitive partial match
curl -X GET "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"name\",\"operator\":\"like\",\"value\":\"web-server\"}]"

# Case-insensitive partial match
curl -X GET "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"name\",\"operator\":\"ilike\",\"value\":\"WEB\"}]"
```

### 5. Filter by Location
```bash
# Get VMs in specific regions
curl -X GET "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"location\",\"operator\":\"in\",\"value\":[\"us-east-1\",\"us-west-2\",\"East US\"]}]"
```

### 6. Date Range Filtering
```bash
# Get VMs created in the last 30 days
curl -X GET "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"createdAt\",\"operator\":\"gte\",\"value\":\"2023-11-01T00:00:00Z\"}]"

# Get VMs created between specific dates
curl -X GET "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"createdAt\",\"operator\":\"between\",\"value\":[\"2023-01-01T00:00:00Z\",\"2023-12-31T23:59:59Z\"]}]"
```

### 7. Null/Empty Value Filtering
```bash
# Get VMs without public IP addresses
curl -X GET "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"publicIpAddress\",\"operator\":\"null\",\"value\":true}]"

# Get VMs with public IP addresses
curl -X GET "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"publicIpAddress\",\"operator\":\"null\",\"value\":false}]"
```

### 8. Complex Multi-Filter Queries
```bash
# Running AWS instances in specific regions
curl -X GET "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"cloudType\",\"operator\":\"eq\",\"value\":\"aws\"},{\"field\":\"status\",\"operator\":\"eq\",\"value\":\"running\"},{\"field\":\"location\",\"operator\":\"in\",\"value\":[\"us-east-1\",\"us-west-2\"]}]"

# Non-production instances (exclude instances with "prod" in name)
curl -X GET "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"name\",\"operator\":\"ilike\",\"value\":\"prod\"}]"

# Large instances across all clouds
curl -X GET "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"instanceType\",\"operator\":\"in\",\"value\":[\"m5.large\",\"Standard_D4s_v3\",\"n1-standard-4\"]}]"
```

### 9. Comprehensive Query with Pagination and Sorting
```bash
# Get running web servers, sorted by creation date, with pagination
curl -X GET "http://localhost:8080/api/v1/vms?page=1&pageSize=10&sortBy=createdAt&sortOrder=desc&filter=[{\"field\":\"status\",\"operator\":\"eq\",\"value\":\"running\"},{\"field\":\"name\",\"operator\":\"ilike\",\"value\":\"web\"}]"
```

## Cloud-Specific Filtering

### AWS-Specific Filters
```bash
# Filter by VPC ID (AWS only)
curl -X GET "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"cloudType\",\"operator\":\"eq\",\"value\":\"aws\"},{\"field\":\"vpcId\",\"operator\":\"eq\",\"value\":\"vpc-12345678\"}]"

# Filter by availability zone
curl -X GET "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"cloudType\",\"operator\":\"eq\",\"value\":\"aws\"},{\"field\":\"availabilityZone\",\"operator\":\"like\",\"value\":\"us-east-1\"}]"
```

### Azure-Specific Filters
```bash
# Filter by resource group (Azure only)
curl -X GET "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"cloudType\",\"operator\":\"eq\",\"value\":\"azure\"},{\"field\":\"resourceGroup\",\"operator\":\"eq\",\"value\":\"production\"}]"

# Filter by VM size
curl -X GET "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"cloudType\",\"operator\":\"eq\",\"value\":\"azure\"},{\"field\":\"vmSize\",\"operator\":\"like\",\"value\":\"Standard_D\"}]"
```

### GCP-Specific Filters
```bash
# Filter by project ID (GCP only)
curl -X GET "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"cloudType\",\"operator\":\"eq\",\"value\":\"gcp\"},{\"field\":\"projectId\",\"operator\":\"eq\",\"value\":\"my-project-123456\"}]"

# Filter by zone
curl -X GET "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"cloudType\",\"operator\":\"eq\",\"value\":\"gcp\"},{\"field\":\"zone\",\"operator\":\"like\",\"value\":\"us-central1\"}]"
```

## Error Examples

### Invalid Filter Operator
```bash
curl -X GET "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"status\",\"operator\":\"invalid\",\"value\":\"running\"}]"
# Response: 400 Bad Request
# {"message":"Query parameter validation failed","code":"VALIDATION_ERROR","details":"filter 0: unsupported operator: invalid"}
```

### Invalid Page Size
```bash
curl -X GET "http://localhost:8080/api/v1/vms?pageSize=2000"
# Response: 400 Bad Request
# {"message":"Invalid query parameters","code":"INVALID_PARAMS","details":"pageSize cannot exceed 1000"}
```

### Invalid JSON Filter
```bash
curl -X GET "http://localhost:8080/api/v1/vms?filter=invalid-json"
# Response: 400 Bad Request
# {"message":"Invalid query parameters","code":"INVALID_PARAMS","details":"invalid filter JSON: invalid character 'i' looking for beginning of value"}
```

## Response Format

### Successful Response
```json
{
  "data": [
    {
      "id": "i-1234567890abcdef0",
      "name": "web-server-01",
      "cloudType": "aws",
      "status": "running",
      "createdAt": "2023-11-15T10:30:00Z",
      "updatedAt": "2023-11-15T12:00:00Z",
      "cloudAccountId": "123456789012",
      "location": "us-east-1",
      "instanceType": "t2.micro",
      "cloudSpecificDetails": {
        "id": "i-1234567890abcdef0",
        "name": "web-server-01",
        "status": "running",
        "createdAt": "2023-11-15T10:30:00Z",
        "updatedAt": "2023-11-15T12:00:00Z",
        "location": "us-east-1",
        "instanceType": "t2.micro",
        "accountId": "123456789012",
        "vpcId": "vpc-12345678",
        "subnetId": "subnet-12345678",
        "securityGroupIds": ["sg-12345678", "sg-98765432"],
        "privateIpAddress": "10.0.1.100",
        "publicIpAddress": "54.123.45.67",
        "keyName": "my-key-pair",
        "imageId": "ami-0abcdef1234567890",
        "launchTime": "2023-11-15T10:30:00Z",
        "availabilityZone": "us-east-1a",
        "architecture": "x86_64",
        "virtualizationType": "hvm",
        "platform": "linux",
        "rootDeviceType": "ebs",
        "monitoringState": "enabled",
        "ebsOptimized": true,
        "tags": {
          "Environment": "production",
          "Owner": "team-web"
        }
      }
    }
  ],
  "pagination": {
    "page": 1,
    "pageSize": 20,
    "totalItems": 42,
    "totalPages": 3
  }
}
```

### Error Response
```json
{
  "message": "Query parameter validation failed",
  "code": "VALIDATION_ERROR",
  "details": "filter 0: unsupported operator: invalid"
}
```

## Health Check

### Check Service Health
```bash
curl -X GET "http://localhost:8080/health"
```

### Expected Response
```json
{
  "status": "healthy",
  "timestamp": "2023-11-15T12:00:00Z",
  "version": "1.0.0",
  "uptime": "2h30m15s",
  "database": {
    "status": "healthy",
    "ping": "2.5ms",
    "connections": 5
  },
  "services": {
    "api": "healthy"
  }
}
```

## Filter Operator Reference

| Operator | Description | Example Value | Use Case |
|----------|-------------|---------------|-----------|
| `eq` | Equal | `"running"` | Exact match |
| `ne` | Not Equal | `"stopped"` | Exclude specific values |
| `lt` | Less Than | `"2023-01-01"` | Date/number comparisons |
| `lte` | Less Than or Equal | `"2023-01-01"` | Date/number comparisons |
| `gt` | Greater Than | `"2023-01-01"` | Date/number comparisons |
| `gte` | Greater Than or Equal | `"2023-01-01"` | Date/number comparisons |
| `in` | In List | `["t2.micro", "t3.small"]` | Multiple options |
| `nin` | Not in List | `["stopped", "terminated"]` | Exclude multiple values |
| `like` | Partial Match | `"web-server"` | Contains pattern |
| `ilike` | Case-Insensitive LIKE | `"WEB"` | Case-insensitive contains |
| `between` | Between Two Values | `["2023-01-01", "2023-12-31"]` | Range queries |
| `null` | Is Null / Not Null | `true` / `false` | Empty value checks |

## Common Use Cases

### 1. DevOps Monitoring
```bash
# Get all non-running instances for alerting
curl -X GET "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"status\",\"operator\":\"ne\",\"value\":\"running\"}]"
```

### 2. Cost Optimization
```bash
# Find large instances that might be over-provisioned
curl -X GET "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"instanceType\",\"operator\":\"in\",\"value\":[\"m5.xlarge\",\"Standard_D8s_v3\",\"n1-standard-8\"]}]"
```

### 3. Security Auditing
```bash
# Find instances with public IP addresses
curl -X GET "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"publicIpAddress\",\"operator\":\"null\",\"value\":false}]"
```

### 4. Environment Management
```bash
# Get all production instances
curl -X GET "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"name\",\"operator\":\"ilike\",\"value\":\"prod\"}]"
```

### 5. Multi-Cloud Inventory
```bash
# Compare instance counts across clouds
curl -X GET "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"cloudType\",\"operator\":\"eq\",\"value\":\"aws\"}]" | jq '.pagination.totalItems'
curl -X GET "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"cloudType\",\"operator\":\"eq\",\"value\":\"azure\"}]" | jq '.pagination.totalItems'
curl -X GET "http://localhost:8080/api/v1/vms?filter=[{\"field\":\"cloudType\",\"operator\":\"eq\",\"value\":\"gcp\"}]" | jq '.pagination.totalItems'
```