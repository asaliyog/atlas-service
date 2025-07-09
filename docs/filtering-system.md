# Configurable Filtering System Documentation

## Overview

The configurable filtering system provides a flexible, type-safe way to filter API responses based on field-specific operators and validation rules. It supports multiple data types, various comparison operators, and comprehensive validation to ensure data integrity.

## Architecture

### Core Components

1. **Filter Configuration** (`internal/config/filter_config.go`)
   - Defines field types and allowed operators
   - Validates filter parameters
   - Parses query parameters

2. **Filter Utilities** (`internal/utils/filter.go`)
   - Applies filters to data structures
   - Handles different data types and operators
   - Provides field value extraction

3. **Handler Integration** (`internal/handlers/vms.go`)
   - Integrates filtering with API endpoints
   - Handles pagination and sorting
   - Returns appropriate error responses

## How It Works

### Query Parameter Format

Filters are specified using the format: `field_operator=value`

**Examples:**
```
GET /api/v1/vms?status_eq=running
GET /api/v1/vms?name_contains=server&cloudType_in=aws,azure
GET /api/v1/vms?createdAt_gte=2024-01-01&page=1&pageSize=10
```

### Filter Processing Flow

1. **Parse Query Parameters**: Extract filter parameters from URL
2. **Validate Filters**: Check field existence, operator compatibility, and value format
3. **Apply Filters**: Filter the data based on validated parameters
4. **Return Results**: Apply pagination and sorting, then return response

## Field Types and Operators

### Supported Field Types

| Type | Description | Example Values |
|------|-------------|----------------|
| `string` | Text data | "running", "web-server" |
| `int` | Integer numbers | 1, 100, -5 |
| `float` | Decimal numbers | 1.5, 100.0, -5.2 |
| `bool` | Boolean values | true, false, 1, 0 |
| `date` | Date/time values | "2024-01-01", "2024-01-01T10:00:00Z" |
| `array` | Comma-separated lists | "aws,azure,gcp" |

### Supported Operators

#### String Operators
- `eq` - Equals (case-insensitive)
- `ne` - Not equals (case-insensitive)
- `contains` - Contains substring (case-insensitive)
- `starts_with` - Starts with substring (case-insensitive)
- `ends_with` - Ends with substring (case-insensitive)
- `like` - Pattern matching with % wildcard
- `ilike` - Case-insensitive pattern matching
- `in` - Value is in comma-separated list
- `not_in` - Value is not in comma-separated list
- `is_null` - Field is null/empty
- `is_not_null` - Field is not null/empty

#### Numeric/Date Operators
- `eq` - Equals
- `ne` - Not equals
- `gt` - Greater than
- `gte` - Greater than or equal
- `lt` - Less than
- `lte` - Less than or equal
- `between` - Value is between two values (comma-separated)
- `is_null` - Field is null/empty
- `is_not_null` - Field is not null/empty

### Field Configuration Example

```go
func VMsFilterConfig() FilterConfig {
    return FilterConfig{
        Fields: map[string]FieldConfig{
            "name": {
                Type: FieldTypeString,
                Operators: []FilterOperator{
                    OperatorEquals, OperatorNotEquals, OperatorContains,
                    OperatorStartsWith, OperatorEndsWith, OperatorLike,
                    OperatorILike, OperatorIn, OperatorNotIn,
                    OperatorIsNull, OperatorIsNotNull,
                },
            },
            "status": {
                Type: FieldTypeString,
                Operators: []FilterOperator{
                    OperatorEquals, OperatorNotEquals, OperatorIn,
                    OperatorNotIn, OperatorIsNull, OperatorIsNotNull,
                },
            },
            "createdAt": {
                Type: FieldTypeDate,
                Operators: []FilterOperator{
                    OperatorEquals, OperatorNotEquals, OperatorGreaterThan,
                    OperatorGreaterEqual, OperatorLessThan, OperatorLessEqual,
                    OperatorBetween, OperatorIsNull, OperatorIsNotNull,
                },
            },
        },
    }
}
```

## Setting Up Filtering for a New Endpoint

### Step 1: Define Filter Configuration

Create a filter configuration function in `internal/config/filter_config.go`:

```go
func UsersFilterConfig() FilterConfig {
    return FilterConfig{
        Fields: map[string]FieldConfig{
            "id": {
                Type: FieldTypeString,
                Operators: []FilterOperator{
                    OperatorEquals, OperatorNotEquals, OperatorContains,
                    OperatorStartsWith, OperatorEndsWith, OperatorLike,
                    OperatorILike, OperatorIn, OperatorNotIn,
                    OperatorIsNull, OperatorIsNotNull,
                },
            },
            "email": {
                Type: FieldTypeString,
                Operators: []FilterOperator{
                    OperatorEquals, OperatorNotEquals, OperatorContains,
                    OperatorStartsWith, OperatorEndsWith, OperatorLike,
                    OperatorILike, OperatorIn, OperatorNotIn,
                    OperatorIsNull, OperatorIsNotNull,
                },
            },
            "age": {
                Type: FieldTypeInt,
                Operators: []FilterOperator{
                    OperatorEquals, OperatorNotEquals, OperatorGreaterThan,
                    OperatorGreaterEqual, OperatorLessThan, OperatorLessEqual,
                    OperatorBetween, OperatorIsNull, OperatorIsNotNull,
                },
            },
            "isActive": {
                Type: FieldTypeBool,
                Operators: []FilterOperator{
                    OperatorEquals, OperatorNotEquals, OperatorIsNull, OperatorIsNotNull,
                },
            },
        },
    }
}
```

### Step 2: Update Handler

Modify your handler to use the filtering system:

```go
func (h *UsersHandler) GetUsers(c *gin.Context) {
    // Get filter configuration for Users endpoint
    filterConfig := config.UsersFilterConfig()

    // Parse and validate filters from query parameters
    filters, err := filterConfig.ParseQueryParams(c.Request.URL.Query())
    if err != nil {
        utils.SendErrorResponse(c, http.StatusBadRequest, 
            fmt.Sprintf("Filter validation error: %s", err.Error()))
        return
    }

    // Parse pagination and sorting parameters
    page := 1
    pageSize := 10
    sortBy := ""
    sortOrder := "asc"

    if pageStr := c.Query("page"); pageStr != "" {
        if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
            page = p
        }
    }

    if pageSizeStr := c.Query("pageSize"); pageSizeStr != "" {
        if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 1000 {
            pageSize = ps
        }
    }

    if sortByParam := c.Query("sortBy"); sortByParam != "" {
        sortBy = sortByParam
    }

    if sortOrderParam := c.Query("sortOrder"); sortOrderParam != "" {
        if sortOrderParam == "desc" {
            sortOrder = "desc"
        }
    }

    // Fetch data from your data source
    users, err := h.fetchUsers()
    if err != nil {
        utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to fetch users")
        return
    }

    // Apply filters using the new configurable system
    filteredUsers := utils.ApplyFilters(users, filters)

    // Apply sorting
    sortedUsers := h.applySorting(filteredUsers, sortBy, sortOrder)

    // Calculate pagination
    totalItems := len(sortedUsers)
    paginatedUsers := utils.ApplyPagination(sortedUsers, page, pageSize)

    // Send response
    utils.SendPaginatedResponse(c, paginatedUsers, page, pageSize, totalItems)
}
```

### Step 3: Update Filter Utilities (if needed)

If your data structure is different from `models.VM`, you may need to update the `ApplyFilters` function in `internal/utils/filter.go` to work with your data type:

```go
// ApplyFilters applies filters to a slice of Users based on the filter configuration
func ApplyUserFilters(users []models.User, filters []config.FilterParam) []models.User {
    if len(filters) == 0 {
        return users
    }

    var filteredUsers []models.User

    for _, user := range users {
        include := true
        for _, filter := range filters {
            if !applyUserFilter(user, filter) {
                include = false
                break
            }
        }
        if include {
            filteredUsers = append(filteredUsers, user)
        }
    }

    return filteredUsers
}

// applyUserFilter applies a single filter to a User
func applyUserFilter(user models.User, filter config.FilterParam) bool {
    // Get the field value using reflection
    fieldValue := GetUserFieldValue(user, filter.Field)
    // ... rest of the implementation similar to VM filtering
}
```

## Usage Examples

### Basic Filtering

```bash
# Get all running VMs
GET /api/v1/vms?status_eq=running

# Get VMs with names containing "server"
GET /api/v1/vms?name_contains=server

# Get VMs from specific cloud providers
GET /api/v1/vms?cloudType_in=aws,azure
```

### Advanced Filtering

```bash
# Multiple filters (AND logic)
GET /api/v1/vms?status_eq=running&cloudType_eq=aws&name_contains=web

# Date range filtering
GET /api/v1/vms?createdAt_gte=2024-01-01&createdAt_lte=2024-12-31

# Pattern matching
GET /api/v1/vms?name_like=%server%

# Null value checking
GET /api/v1/vms?publicIp_is_null
```

### Pagination and Sorting

```bash
# With pagination
GET /api/v1/vms?status_eq=running&page=1&pageSize=5

# With sorting
GET /api/v1/vms?status_eq=running&sortBy=name&sortOrder=desc

# Combined filtering, pagination, and sorting
GET /api/v1/vms?cloudType_eq=aws&page=2&pageSize=10&sortBy=createdAt&sortOrder=desc
```

## Error Handling

### Validation Errors

The system returns HTTP 400 Bad Request for validation errors:

```json
{
  "error": "Filter validation error: operator 'gte' is not allowed for field 'status' of type 'string'"
}
```

### Common Error Scenarios

1. **Invalid Field**: Field not configured for filtering
   ```
   GET /api/v1/vms?invalidField_eq=value
   Response: "field 'invalidField' is not allowed for filtering"
   ```

2. **Invalid Operator**: Operator not allowed for field type
   ```
   GET /api/v1/vms?status_gte=running
   Response: "operator 'gte' is not allowed for field 'status' of type 'string'"
   ```

3. **Invalid Value Format**: Value doesn't match field type
   ```
   GET /api/v1/vms?createdAt_gte=invalid-date
   Response: "value 'invalid-date' is not a valid date for field 'createdAt'"
   ```

4. **Invalid Parameter Format**: Wrong query parameter format
   ```
   GET /api/v1/vms?status=running
   Response: "invalid filter parameter format: 'status'. Expected format: field_operator"
   ```

## Best Practices

### 1. Field Configuration

- **Be Specific**: Only allow operators that make sense for each field type
- **Consider Performance**: String operators like `contains` and `like` can be slower on large datasets
- **Document Fields**: Clearly document what each field represents and its expected values

### 2. Error Messages

- **Be Descriptive**: Provide clear, actionable error messages
- **Include Context**: Mention the field name and type in error messages
- **Suggest Alternatives**: When possible, suggest valid alternatives

### 3. Performance Considerations

- **Index Fields**: Ensure database fields used for filtering are properly indexed
- **Limit Page Size**: Set reasonable limits on page size (default: 1000)
- **Use Appropriate Operators**: Prefer exact matches (`eq`) over pattern matching (`like`) when possible

### 4. Security

- **Validate Input**: All filter parameters are validated for type safety
- **Sanitize Output**: Ensure filtered data doesn't expose sensitive information
- **Rate Limiting**: Consider implementing rate limiting for expensive filter operations

## Extending the System

### Adding New Field Types

1. Add the new type to `FieldType` enum in `filter_config.go`
2. Implement validation logic in `validateValue` function
3. Add comparison logic in `compareValues` function in `filter.go`

### Adding New Operators

1. Add the new operator to `FilterOperator` enum
2. Implement the operator logic in `applyFilter` function
3. Update field configurations to include the new operator where appropriate

### Adding New Data Types

1. Create a new filter configuration function
2. Implement `ApplyFilters` function for the new data type
3. Update handlers to use the new filtering system

## Testing

### Unit Tests

Create unit tests for filter validation and application:

```go
func TestFilterValidation(t *testing.T) {
    config := config.VMsFilterConfig()
    
    // Test valid filter
    err := config.ValidateFilter("status", "eq", "running")
    assert.NoError(t, err)
    
    // Test invalid operator
    err = config.ValidateFilter("status", "gte", "running")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "operator 'gte' is not allowed")
}
```

### Integration Tests

Test the complete filtering flow:

```go
func TestVMsFiltering(t *testing.T) {
    // Setup test data
    vms := []models.VM{...}
    
    // Test filtering
    filters := []config.FilterParam{
        {Field: "status", Operator: config.OperatorEquals, Value: "running"},
    }
    
    result := utils.ApplyFilters(vms, filters)
    assert.Len(t, result, expectedCount)
}
```

## Troubleshooting

### Common Issues

1. **Filters Not Working**: Check that field names match JSON tags in your model
2. **Type Mismatches**: Ensure field types in configuration match actual data types
3. **Performance Issues**: Consider adding database indexes for frequently filtered fields
4. **Memory Usage**: Large datasets with complex filters may consume significant memory

### Debug Tips

1. **Enable Logging**: Add debug logs to see which filters are being applied
2. **Test Incrementally**: Test filters one at a time to isolate issues
3. **Check Data Types**: Verify that your data matches the configured field types
4. **Validate Configuration**: Ensure all fields in your model are properly configured

## Conclusion

The configurable filtering system provides a robust, type-safe, and extensible solution for API filtering. It enforces data integrity through comprehensive validation while maintaining flexibility for different use cases. By following the patterns established in this documentation, you can easily add filtering capabilities to any endpoint in your API. 