# Reusable API Utilities for Filtering, Sorting, and Pagination

This document explains how to use the reusable utilities in `internal/utils/` to add filtering, sorting, and pagination to any API endpoint that returns arrays.

## Overview

The utilities provide a standardized way to handle query parameters, filtering, sorting, and pagination across all API endpoints. This ensures consistency and reduces code duplication.

## Core Components

### 1. QueryParams Structure

```go
type QueryParams struct {
    Page      int           `json:"page"`
    PageSize  int           `json:"pageSize"`
    SortBy    string        `json:"sortBy"`
    SortOrder string        `json:"sortOrder"`
    Filters   []QueryFilter `json:"filters"`
}

type QueryFilter struct {
    Field    string      `json:"field"`
    Operator string      `json:"operator"`
    Value    interface{} `json:"value"`
}
```

### 2. Standard Query Parameter Format

The utilities support a standard query parameter format:

- **Pagination**: `page=1&pageSize=20`
- **Sorting**: `sortBy=name&sortOrder=asc`
- **Filtering**: `field=value` (equals) or `field_operator=value` (other operators)

## Supported Filter Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `eq` | Equals | `name=web-server` |
| `ne` | Not equals | `status_ne=running` |
| `gt` | Greater than | `cpu_cores_gt=4` |
| `gte` | Greater than or equal | `memory_gb_gte=8` |
| `lt` | Less than | `storage_gb_lt=100` |
| `lte` | Less than or equal | `created_at_lte=2024-01-01` |
| `contains` | Contains substring | `name_contains=web` |
| `starts_with` | Starts with | `name_starts_with=prod` |
| `ends_with` | Ends with | `name_ends_with=01` |
| `in` | In array | `status_in=running,stopped` |
| `not_in` | Not in array | `cloudType_not_in=aws,azure` |
| `is_null` | Is null | `public_ip_is_null=true` |
| `is_not_null` | Is not null | `private_ip_is_not_null=true` |
| `like` | SQL LIKE | `name_like=%web%` |
| `ilike` | Case-insensitive LIKE | `name_ilike=%WEB%` |
| `between` | Between values | `created_at_between=2024-01-01,2024-12-31` |

## How to Use in Any Handler

### Step 1: Import the utilities

```go
import (
    "atlas-service/internal/utils"
    "net/http"
    "reflect"
    "sort"
    "strings"
    "github.com/gin-gonic/gin"
)
```

### Step 2: Create a handler with the standard pattern

```go
type MyHandler struct {
    db *gorm.DB
}

func (h *MyHandler) GetMyResources(c *gin.Context) {
    // 1. Parse query parameters
    params, err := utils.ParseQueryParams(c)
    if err != nil {
        utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }

    // 2. Validate query parameters
    if err := utils.ValidateQueryParams(params); err != nil {
        utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }

    // 3. Get data from database
    var resources []MyResource
    if err := h.db.Find(&resources).Error; err != nil {
        utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to fetch resources")
        return
    }

    // 4. Apply filters
    filteredResources := utils.ApplyInMemoryFilters(resources, params.Filters, h.extractField)

    // 5. Apply sorting
    sortedResources := h.applySorting(filteredResources, params.SortBy, params.SortOrder)

    // 6. Apply pagination
    totalItems := len(sortedResources)
    paginatedResources := utils.ApplyPagination(sortedResources, params.Page, params.PageSize)

    // 7. Send response
    utils.SendPaginatedResponse(c, paginatedResources, params.Page, params.PageSize, totalItems)
}
```

### Step 3: Implement field extraction function

```go
// extractField extracts field value from your struct for filtering
func (h *MyHandler) extractField(resource MyResource, field string) interface{} {
    v := reflect.ValueOf(resource)
    t := v.Type()

    for i := 0; i < v.NumField(); i++ {
        fieldType := t.Field(i)
        fieldValue := v.Field(i)

        // Check JSON tag first, then field name
        jsonTag := fieldType.Tag.Get("json")
        if jsonTag != "" {
            jsonField := strings.Split(jsonTag, ",")[0]
            if jsonField == field {
                return fieldValue.Interface()
            }
        }

        // Check field name
        if strings.EqualFold(fieldType.Name, field) {
            return fieldValue.Interface()
        }
    }

    return nil
}
```

### Step 4: Implement sorting function

```go
// applySorting applies sorting to resources
func (h *MyHandler) applySorting(resources []MyResource, sortBy, sortOrder string) []MyResource {
    if sortBy == "" {
        return resources
    }

    sortedResources := make([]MyResource, len(resources))
    copy(sortedResources, resources)

    sort.Slice(sortedResources, func(i, j int) bool {
        valI := h.extractField(sortedResources[i], sortBy)
        valJ := h.extractField(sortedResources[j], sortBy)

        // Handle nil values
        if valI == nil && valJ == nil {
            return false
        }
        if valI == nil {
            return sortOrder == "asc"
        }
        if valJ == nil {
            return sortOrder == "desc"
        }

        // Convert to string for comparison
        strI := strings.ToLower(fmt.Sprintf("%v", valI))
        strJ := strings.ToLower(fmt.Sprintf("%v", valJ))

        if sortOrder == "desc" {
            return strI > strJ
        }
        return strI < strJ
    })

    return sortedResources
}
```

## Complete Example: Users API

Here's a complete example of how the Users API uses these utilities:

```go
package handlers

import (
    "atlas-service/internal/models"
    "atlas-service/internal/utils"
    "fmt"
    "net/http"
    "reflect"
    "sort"
    "strings"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type UsersHandler struct {
    db *gorm.DB
}

func NewUsersHandler(db *gorm.DB) *UsersHandler {
    return &UsersHandler{db: db}
}

func (h *UsersHandler) GetUsers(c *gin.Context) {
    // Parse query parameters using the reusable utility
    params, err := utils.ParseQueryParams(c)
    if err != nil {
        utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }

    // Validate query parameters
    if err := utils.ValidateQueryParams(params); err != nil {
        utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }

    // Get users from database
    var users []models.User
    if err := h.db.Find(&users).Error; err != nil {
        utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to fetch users")
        return
    }

    // Apply filters using the reusable utility
    filteredUsers := utils.ApplyInMemoryFilters(users, params.Filters, h.extractUserField)

    // Apply sorting
    sortedUsers := h.applySorting(filteredUsers, params.SortBy, params.SortOrder)

    // Calculate pagination
    totalItems := len(sortedUsers)
    paginatedUsers := utils.ApplyPagination(sortedUsers, params.Page, params.PageSize)

    // Send response using the reusable utility
    utils.SendPaginatedResponse(c, paginatedUsers, params.Page, params.PageSize, totalItems)
}

func (h *UsersHandler) extractUserField(user models.User, field string) interface{} {
    v := reflect.ValueOf(user)
    t := v.Type()

    for i := 0; i < v.NumField(); i++ {
        fieldType := t.Field(i)
        fieldValue := v.Field(i)

        jsonTag := fieldType.Tag.Get("json")
        if jsonTag != "" {
            jsonField := strings.Split(jsonTag, ",")[0]
            if jsonField == field {
                return fieldValue.Interface()
            }
        }

        if strings.EqualFold(fieldType.Name, field) {
            return fieldValue.Interface()
        }
    }

    return nil
}

func (h *UsersHandler) applySorting(users []models.User, sortBy, sortOrder string) []models.User {
    if sortBy == "" {
        return users
    }

    sortedUsers := make([]models.User, len(users))
    copy(sortedUsers, users)

    sort.Slice(sortedUsers, func(i, j int) bool {
        valI := h.extractUserField(sortedUsers[i], sortBy)
        valJ := h.extractUserField(sortedUsers[j], sortBy)

        if valI == nil && valJ == nil {
            return false
        }
        if valI == nil {
            return sortOrder == "asc"
        }
        if valJ == nil {
            return sortOrder == "desc"
        }

        strI := strings.ToLower(fmt.Sprintf("%v", valI))
        strJ := strings.ToLower(fmt.Sprintf("%v", valJ))

        if sortOrder == "desc" {
            return strI > strJ
        }
        return strI < strJ
    })

    return sortedUsers
}
```

## API Usage Examples

### Basic pagination
```
GET /api/v1/users?page=1&pageSize=10
```

### Sorting
```
GET /api/v1/users?sortBy=name&sortOrder=asc
GET /api/v1/users?sortBy=createdAt&sortOrder=desc
```

### Filtering
```
GET /api/v1/users?name=john
GET /api/v1/users?email_contains=gmail
GET /api/v1/users?isActive=true
GET /api/v1/users?createdAt_gte=2024-01-01
```

### Combined queries
```
GET /api/v1/users?page=1&pageSize=20&sortBy=name&sortOrder=asc&isActive=true&email_contains=@company.com
```

## Response Format

All endpoints return a standardized paginated response:

```json
{
  "data": [
    {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "isActive": true,
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
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

## Error Responses

Error responses follow a standard format:

```json
{
  "error": "Invalid query parameters: page must be greater than 0"
}
```

## Benefits

1. **Consistency**: All APIs follow the same pattern for filtering, sorting, and pagination
2. **Reusability**: Write once, use everywhere
3. **Maintainability**: Changes to query logic only need to be made in one place
4. **Type Safety**: Generic functions ensure type safety
5. **Performance**: Efficient in-memory filtering and sorting
6. **Flexibility**: Support for multiple operators and complex queries

## Best Practices

1. **Field Extraction**: Always implement a proper field extraction function that handles JSON tags
2. **Error Handling**: Always validate query parameters before processing
3. **Performance**: For large datasets, consider implementing database-level filtering
4. **Documentation**: Document supported fields and operators for each endpoint
5. **Testing**: Test various filter combinations and edge cases

## Migration Guide

To migrate an existing endpoint to use these utilities:

1. Replace custom query parameter parsing with `utils.ParseQueryParams()`
2. Replace custom validation with `utils.ValidateQueryParams()`
3. Replace custom filtering with `utils.ApplyInMemoryFilters()`
4. Replace custom pagination with `utils.ApplyPagination()`
5. Replace custom response formatting with `utils.SendPaginatedResponse()`
6. Implement the required field extraction and sorting functions

This approach ensures that all your APIs have consistent, powerful, and maintainable filtering, sorting, and pagination capabilities. 