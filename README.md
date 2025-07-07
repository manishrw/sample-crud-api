# Multi-Language CRUD API Examples

## Overview
This repository demonstrates CRUD (Create, Read, Update, Delete) API implementations across multiple programming languages and frameworks. The goal is to showcase idiomatic patterns, best practices, and architectural decisions for each technology stack while maintaining consistent functionality.

## Purpose
- Provide reference implementations for developers learning new technologies
- Demonstrate language-specific idioms and patterns
- Highlight differences in error handling, validation, and data persistence approaches
- Showcase testing strategies across different ecosystems
- Illustrate deployment and containerization practices

## Core Requirements
Each implementation must provide:

### API Endpoints
- `GET /items` - List all items with pagination
- `GET /items/{id}` - Retrieve a specific item
- `POST /items` - Create a new item
- `PUT /items/{id}` - Update an existing item
- `DELETE /items/{id}` - Delete an item

### Common Features
- Request validation
- Error handling with appropriate HTTP status codes
- Database integration
- Logging
- API documentation
- Unit and integration tests
- Docker support
- Load testing scripts

## Implementations

### Go (with Gin)
- Uses clean architecture principles
- Leverages Go's strong typing and interfaces
- Implements middleware for logging and authentication
- Uses `sqlx` for type-safe database operations
- Demonstrates Go's error handling patterns

### Node.js (with Express)
- Showcases async/await patterns
- Uses TypeScript for type safety
- Implements repository pattern
- Demonstrates middleware composition
- Uses Jest for testing

### Python (with FastAPI)
- Leverages Python's type hints
- Uses Pydantic for request/response validation
- Demonstrates async database operations
- Implements dependency injection
- Uses pytest for testing

### Java (with Spring Boot)
- Demonstrates Spring's dependency injection
- Uses JPA for database operations
- Implements DTO pattern
- Shows Spring's validation framework
- Uses JUnit 5 for testing

## Architecture

### Common Patterns
Each implementation follows these architectural principles: 