# Integration Tests

These tests run against a real Rails API server and verify end-to-end functionality.

## Prerequisites

1. Rails API server must be running at `http://localhost:3000`
2. Test user must exist: `dev@ai-memoria.com` with password `dev123`

## Running the Tests

\`\`\`bash
# Start Rails API server (if not running)
cd ../api
rails server
\`\`\`

In another terminal:
\`\`\`bash
# Run all integration tests
make test-integration

# Run specific test
go test -v ./tests/integration/... -run TestIntegrationLogin
\`\`\`

## Environment Variables

You can override the default test credentials:

\`\`\`bash
export AI_MEMORIA_API_URL=http://localhost:3000
export AI_MEMORIA_TEST_EMAIL=dev@ai-memoria.com
export AI_MEMORIA_TEST_PASSWORD=dev123
make test-integration
\`\`\`

## What's Tested

- ✅ Login authentication
- ✅ Whoami (current user info)
- ✅ User creation
- ✅ Duplicate user handling
- ✅ API status check
- ✅ Logout (token revocation)
- ✅ JSON output format
- ✅ Invalid login handling
- ✅ Unauthenticated access

## Test Isolation

Each test runs with an isolated configuration directory (`~/.ai-memoria` is overridden via `HOME`), so tests don't interfere with your development config.

## Notes

- Tests create real users in the database (with unique emails)
- No cleanup is performed - users are left for manual verification
- Tests use the development database, so be aware of data created
