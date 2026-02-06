# AI Coding Agents

This document describes the AI coding agents available for the showheaders project and how to use them effectively.

## Overview

AI coding agents are specialized assistants that can help with various development tasks in this project. They have deep understanding of the codebase structure, Go best practices, and can assist with implementation, testing, debugging, and documentation.

## Available Agents

### Code Implementation Agent
**Purpose**: Implements new features and functionality

**Best used for**:
- Adding new HTTP endpoints
- Implementing middleware
- Creating new handlers
- Extending configuration options

**Example prompts**:
- "Add a new endpoint that shows only specific headers"
- "Implement request body logging in the middleware"
- "Add support for custom response headers"

### Testing Agent
**Purpose**: Creates and maintains comprehensive test coverage

**Best used for**:
- Writing unit tests
- Creating table-driven tests
- Implementing integration tests
- Improving test coverage

**Example prompts**:
- "Add tests for the new middleware"
- "Create table-driven tests for the headers handler"
- "Add benchmarks for the logging middleware"

### Documentation Agent
**Purpose**: Maintains clear and comprehensive documentation

**Best used for**:
- Updating README.md
- Adding code comments
- Creating API documentation
- Writing usage examples

**Example prompts**:
- "Update the README with the new feature"
- "Add godoc comments to the handler package"
- "Create examples for the configuration options"

### Debugging Agent
**Purpose**: Identifies and fixes issues in the codebase

**Best used for**:
- Analyzing error messages
- Investigating failing tests
- Fixing bugs
- Performance optimization

**Example prompts**:
- "Why is this test failing?"
- "Fix the race condition in the middleware"
- "Optimize the headers rendering"

### Refactoring Agent
**Purpose**: Improves code structure and maintainability

**Best used for**:
- Restructuring code
- Improving error handling
- Simplifying complex functions
- Applying Go best practices

**Example prompts**:
- "Refactor the server setup code"
- "Improve error handling in the handlers"
- "Simplify the configuration logic"

## Working with Agents

### General Guidelines

1. **Be specific**: Provide clear, detailed descriptions of what you want to accomplish
2. **Provide context**: Reference specific files, functions, or error messages
3. **Review changes**: Always review and test agent-generated code
4. **Iterate**: Agents work best with feedback - refine requests as needed

### Project-Specific Context

When working with agents on this project, keep in mind:

- **Go version**: This project uses Go 1.21+
- **Logging**: Uses Zap for structured logging
- **Testing**: Follows standard Go testing patterns with table-driven tests
- **HTTP framework**: Uses standard `net/http` package
- **Project structure**: Internal packages for organization

### Best Practices

#### Do's ✅
- Ask agents to write tests alongside new features
- Request explanations for complex changes
- Have agents check for existing patterns in the codebase
- Ask for multiple implementation options when unsure

#### Don'ts ❌
- Don't skip testing agent-generated code
- Don't accept changes you don't understand
- Don't ignore warnings or deprecation notices
- Don't assume agents know the latest project changes

## Example Workflows

### Adding a New Feature
1. Ask agent to outline the implementation plan
2. Review the plan and provide feedback
3. Have agent implement the feature
4. Request tests for the new feature
5. Ask for documentation updates

### Fixing a Bug
1. Share the error message or failing test with the agent
2. Ask agent to identify the root cause
3. Request a fix with explanation
4. Have agent add tests to prevent regression

### Improving Code Quality
1. Ask agent to analyze a specific file or package
2. Request suggestions for improvements
3. Have agent implement refactoring
4. Ensure all tests still pass

## Tips for Effective Agent Use

### Context Matters
Provide relevant information:
```
"In internal/handlers/headers.go, the ServeHTTP method needs to..."
```

### Be Specific About Requirements
Good: "Add a configuration option for log level with environment variable support"
Bad: "Make logging better"

### Request Explanations
"Explain why you chose this approach over alternatives"

### Verify Understanding
"Show me the changes you'll make before implementing them"

## Continuous Improvement

This document should evolve as we discover better ways to work with AI coding agents. Contributions and suggestions are welcome!

## Related Resources

- [Go Documentation](https://go.dev/doc/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Zap Logger Documentation](https://pkg.go.dev/go.uber.org/zap)

## Contributing

If you discover effective patterns for working with agents on this project, please update this document with your findings.
