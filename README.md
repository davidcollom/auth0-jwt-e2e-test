# Auth0 JWT E2E Test

> **⚠️ This is a "Hack Together Demo"** - A quick and dirty end-to-end testing framework for Auth0 JWT authentication flows. Use at your own risk and adapt as needed!

## What is this?

This repository contains a Go-based end-to-end testing framework that demonstrates how to:

1. **Programmatically create Auth0 applications** using the Auth0 Management API
2. **Generate JWT tokens** using OAuth2 client credentials flow
3. **Test authentication flows** using automated browser testing with ChromeDP
4. **Clean up resources** automatically after testing

The demo uses [httpbin.org](https://httpbin.org) as a test endpoint to validate JWT token authentication in a real browser environment.

## How it works

### Test Flow
1. **Setup**: Creates a machine-to-machine Auth0 application with client credentials
2. **Token Generation**: Requests a JWT token using OAuth2 client credentials flow
3. **Browser Testing**: Uses ChromeDP to navigate to `/bearer` endpoint with the JWT token
4. **Validation**: Verifies the authentication response and token presence
5. **Cleanup**: Automatically deletes the created Auth0 application

### Key Components
- `main_test.go` - Main test orchestration and Auth0 app lifecycle management
- `helpers_test.go` - OAuth2 token generation utilities
- `browser_test.go` - ChromeDP browser automation for authentication testing
- `main.go` - Placeholder main function (tests handle everything)

## Prerequisites

### Auth0 Setup
Before running these tests, you need to configure Auth0 Management API access:

#### 1. Create a Machine-to-Machine Application
1. Go to your [Auth0 Dashboard](https://manage.auth0.com/)
2. Navigate to **Applications** → **Create Application**
3. Choose **Machine to Machine Applications**
4. Select your **Auth0 Management API**
5. Grant the following scopes:
   - `read:clients`
   - `create:clients`
   - `delete:clients`
   - `create:client_grants`
   - `read:client_grants`

#### 2. Get Your Credentials
From your newly created M2M application, copy:
- **Domain** (e.g., `dev-xxxxx.auth0.com`)
- **Client ID**
- **Client Secret**

### Local Environment
- **Go 1.24.4+** (check with `go version`)
- **Chrome/Chromium** browser (for ChromeDP)

## Installation & Setup

1. **Clone and navigate to the repository:**
   ```bash
   git clone <repository-url>
   cd auth0-jwt-e2e-test
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Create environment file:**
   ```bash
   cp .env.example .env  # if you have an example, or create manually
   ```

4. **Configure environment variables in `.env`:**
   ```env
   AUTH0_DOMAIN=your-tenant.auth0.com
   AUTH0_MANAGEMENT_CLIENT_ID=your_management_client_id
   AUTH0_MANAGEMENT_CLIENT_SECRET=your_management_client_secret
   ```

## Usage

### Run All Tests
```bash
go test -v
```

### Run Specific Tests
```bash
# Test only the browser authentication flow
go test -v -run TestHomepage

# Run with more detailed logging
go test -v -count=1
```

### What You'll See
The tests will output logs showing:
- ✅ Auth0 Management Client creation
- 🧪 Test application creation with client credentials
- 🔑 JWT token generation
- 🌐 Browser navigation and authentication testing
- 🧹 Automatic cleanup of created resources

## Expected Output

```
=== RUN   TestMain
Auth0 Management Client created successfully
Checking for existing clients...
Existing clients: X
🧪 Created app: ID=abc123, secret=xyz789
Creating client grant for the app...
JWT Token: eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIs...
OAuth2 token retrieved successfully
Auth0 App created. Starting tests...

=== RUN   TestHomepage
Starting homepage test...
Beginning capture setup
EventRequestWillBeSent: request-123: https://httpbin.org/bearer
EventLoadingFinished: request-123
✅ JSON body from /bearer: {"authenticated":true,"token":"eyJ0eXAi..."}

Tests completed. Exiting...
✅ Deleted test client
```

## Configuration

### Environment Variables
| Variable | Description | Example |
|----------|-------------|---------|
| `AUTH0_DOMAIN` | Your Auth0 tenant domain | `dev-abc123.auth0.com` |
| `AUTH0_MANAGEMENT_CLIENT_ID` | Management API client ID | `AbC123dEf456GhI789` |
| `AUTH0_MANAGEMENT_CLIENT_SECRET` | Management API client secret | `your-secret-here` |

### Customization
- **Change test endpoint**: Modify `BaseURL` in `browser_test.go`
- **Add more scopes**: Update the `Scope` array in `main_test.go`
- **Extend browser tests**: Add more test functions in `browser_test.go`

## Troubleshooting

### Common Issues

1. **"failed to create management client"**
   - Check your Auth0 domain format (should include `.auth0.com`)
   - Verify client ID and secret are correct

2. **"failed to create client"**
   - Ensure your Management API application has `create:clients` scope
   - Check if you've hit Auth0 application limits

3. **"failed to get OAuth2 token"**
   - Verify the created app has proper audience configuration
   - Check if client grants were created successfully

4. **Browser tests fail**
   - Ensure Chrome/Chromium is installed and accessible
   - Check if the target endpoint is reachable

### Debug Mode
Add more verbose logging:
```bash
export DEBUG=1
go test -v -count=1
```

## Security Considerations

⚠️ **Important Security Notes:**
- Never commit `.env` files to version control
- Use different Auth0 tenants for testing vs production
- Regularly rotate Management API credentials
- Consider using Auth0 Deploy CLI for production automation

## Contributing

This is a hack demo, but improvements are welcome! Feel free to:
- Add more comprehensive test scenarios
- Improve error handling
- Add CI/CD pipeline examples
- Extend browser automation coverage

## License

[Add your license here]

---

**Remember: This is a demonstration project.** Adapt the code patterns and security practices to fit your specific production requirements.
