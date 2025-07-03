# Azure Entra ID Setup Guide

This guide explains how to configure Azure Entra ID (formerly Azure AD) for authentication with your Golang service.

## Prerequisites

- Azure subscription
- Admin access to Azure Entra ID tenant
- Azure CLI installed (optional)

## Step 1: Register Application in Azure Entra ID

1. **Navigate to Azure Portal:**
   - Go to [Azure Portal](https://portal.azure.com)
   - Navigate to "Azure Active Directory" > "App registrations"

2. **Create New Registration:**
   - Click "New registration"
   - Enter application name: `golang-service-api`
   - Select "Accounts in this organizational directory only"
   - Leave Redirect URI blank for now
   - Click "Register"

3. **Note Important Values:**
   - **Application (client) ID** - This is your `AZURE_CLIENT_ID`
   - **Directory (tenant) ID** - This is your `AZURE_TENANT_ID`

## Step 2: Configure API Permissions

1. **Add API Permissions:**
   - Go to "API permissions"
   - Click "Add a permission"
   - Select "Microsoft Graph"
   - Choose "Delegated permissions"
   - Add: `User.Read`, `email`, `openid`, `profile`

2. **Grant Admin Consent:**
   - Click "Grant admin consent for [your tenant]"
   - Confirm the action

## Step 3: Configure Authentication

1. **Add Platform Configuration:**
   - Go to "Authentication"
   - Click "Add a platform"
   - Select "Web"
   - Add redirect URI: `http://localhost:8080/auth/callback` (for local development)
   - For production, add your actual domain

2. **Configure Token Settings:**
   - Enable "Access tokens" and "ID tokens"
   - Set logout URL if needed

## Step 4: Create Client Secret (Optional)

If you need a client secret for server-to-server authentication:

1. **Generate Secret:**
   - Go to "Certificates & secrets"
   - Click "New client secret"
   - Add description and set expiration
   - Copy the secret value immediately (you won't see it again)

## Step 5: Configure Application Manifest

1. **Update Manifest:**
   - Go to "Manifest"
   - Set `"accessTokenAcceptedVersion": 2`
   - Save changes

## Step 6: Set Up Application Roles (Optional)

For role-based access control:

1. **Define App Roles:**
   - Go to "App roles"
   - Click "Create app role"
   - Define roles like "Admin", "User", etc.

## Environment Configuration

Update your `.env` file with the values from Azure:

```bash
AZURE_TENANT_ID=your-directory-tenant-id
AZURE_CLIENT_ID=your-application-client-id
```

## Token Validation

In production, implement proper token validation:

1. **Verify Token Signature:**
   - Download public keys from: `https://login.microsoftonline.com/{tenant-id}/discovery/v2.0/keys`
   - Verify JWT signature using these keys

2. **Validate Claims:**
   - `iss` (issuer): `https://login.microsoftonline.com/{tenant-id}/v2.0`
   - `aud` (audience): Your application client ID
   - `exp` (expiration): Check token hasn't expired

## Example Token Payload

```json
{
  "iss": "https://login.microsoftonline.com/{tenant-id}/v2.0",
  "aud": "{client-id}",
  "sub": "{user-object-id}",
  "email": "user@company.com",
  "name": "John Doe",
  "exp": 1234567890,
  "iat": 1234567800
}
```

## Testing Authentication

### Using curl with Bearer Token

```bash
# Get token from Azure (example using client credentials flow)
curl -X POST https://login.microsoftonline.com/{tenant-id}/oauth2/v2.0/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=client_credentials&client_id={client-id}&client_secret={client-secret}&scope=https://graph.microsoft.com/.default"

# Use token with your API
curl -H "Authorization: Bearer {access-token}" \
     http://localhost:8080/api/v1/users
```

### Production Implementation Notes

The current implementation in `internal/middleware/auth.go` is simplified for demonstration. For production:

1. **Replace the simple JWT validation** with proper Azure Entra ID token validation
2. **Implement proper key rotation** handling
3. **Add audience and issuer validation**
4. **Implement proper error handling** and logging
5. **Consider using Azure SDK** for Go for easier integration

### Recommended Libraries

- `github.com/microsoft/kiota-authentication-azure-go` - Official Azure authentication
- `github.com/coreos/go-oidc/v3/oidc` - OIDC client library
- `github.com/golang-jwt/jwt/v5` - JWT handling

## Troubleshooting

### Common Issues

1. **Token validation fails:**
   - Check tenant ID and client ID
   - Verify token hasn't expired
   - Ensure proper audience claim

2. **Permission denied:**
   - Check admin consent has been granted
   - Verify user has required roles

3. **Invalid signature:**
   - Ensure using correct signing keys
   - Check token version compatibility

### Useful Endpoints

- Token endpoint: `https://login.microsoftonline.com/{tenant-id}/oauth2/v2.0/token`
- Authorization endpoint: `https://login.microsoftonline.com/{tenant-id}/oauth2/v2.0/authorize`
- JWKS endpoint: `https://login.microsoftonline.com/{tenant-id}/discovery/v2.0/keys`
- OpenID configuration: `https://login.microsoftonline.com/{tenant-id}/v2.0/.well-known/openid_configuration`