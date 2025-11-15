# Contributing to Trade Journal

Thank you for wanting to help! Please follow these simple rules.

## Important Rule: Always Test Your Code

**Every new feature must have tests.** This is very important!

### Types of Tests

1. **Unit Tests** - Test small parts of code
2. **Integration Tests** - Test how parts work together
3. **E2E Tests** - Test the complete feature

### Where to Put Tests

- Unit tests: Same folder as your code, file name ends with `_test.go`
- Integration tests: Same folder, file name ends with `_integration_test.go`
- E2E tests: Put in `test/e2e/` folder

### Running Tests

Before you submit your code, run these commands:

```bash
# Run all tests
make test

# Run tests and see details
make test-verbose

# Check test coverage
make test-coverage
```

**Your tests must pass before submitting!**

## How to Contribute

1. **Fork** this project
2. **Create a branch** for your feature: `git checkout -b my-new-feature`
3. **Write your code**
4. **Write tests** for your code (required!)
5. **Run tests** to make sure everything works
6. **Commit** your changes: `git commit -m "Add new feature"`
7. **Push** to your branch: `git push origin my-new-feature`
8. **Create a Pull Request**

## Code Style

- Write clean, simple code
- Add comments to explain difficult parts
- Follow the existing code style in the project

## Questions?

If you need help, create an issue and ask!
