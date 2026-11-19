# Login Authentication Unit Testing

## Overview

This document describes the unit testing implementation for the Login feature in the Katarina Daily backend.

The tests are implemented using Go testing tools and focus on validating the authentication process in AuthService.Login().

## Test Scenarios

Four test scenarios are covered:

1. Successful login with normalized email.
2. Login with an unknown email.
3. Login with an incorrect password.
4. Login when the repository returns an error.

## Testing Techniques

### Test Doubles

Two Test Doubles are used:

- Stub: Provides predefined repository responses to simulate different login conditions.
- Spy: Records repository interactions so that the test can verify how the service communicates with its dependency.

### Arrange, Act, Assert (AAA)

Each test follows three stages:

1. Arrange: Prepare test data and dependencies.
2. Act: Execute the Login function.
3. Assert: Verify the returned result and expected behavior.

## Running the Tests

Run all backend tests:

```bash
go test ./...
```

Run Login-specific tests:

```bash
go test ./internal/services -run Login -v
```

## Statement Coverage

Generate the coverage profile:

```bash
go test ./internal/services -coverprofile=coverage-login.out
```

Display function coverage:

```bash
go tool cover -func=coverage-login.out
```

The Login function achieved 100% statement coverage in the executed test suite.

Note: The overall services package coverage is lower because other service functions are outside the scope of this Login testing task.

## Test Results

All four Login test scenarios passed successfully.

The test results and coverage report are documented separately as testing evidence for Module 4.