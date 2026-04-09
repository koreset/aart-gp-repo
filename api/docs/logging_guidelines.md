# Comprehensive Logging Guidelines

## Overview

This document provides guidelines for implementing comprehensive logging throughout the application. Proper logging is essential for debugging, monitoring, and understanding application behavior in production environments.

## Logging Library

The application uses a custom logging package based on [logrus](https://github.com/sirupsen/logrus) located at `/api/log/logger.go`. This package provides structured logging with various log levels and context support.

## Key Logging Concepts

1. **Log Levels**: Use appropriate log levels based on the importance of the message:
   - `Debug`: Detailed information useful for debugging
   - `Info`: General information about application operation
   - `Warn`: Warning messages that don't prevent the application from functioning
   - `Error`: Error messages that may prevent a specific operation from completing
   - `Fatal`: Critical errors that prevent the application from starting or continuing
   - `Panic`: Errors that cause the application to panic

2. **Structured Logging**: Use structured logging with fields to provide context:
   ```go
   logger.WithFields(map[string]interface{}{
       "user_id": user.ID,
       "action": "login",
   }).Info("User logged in")
   ```

3. **Request Context**: Include request ID and user information in logs:
   ```go
   logger := log.WithContext(ctx)
   logger.Info("Processing request")
   ```

4. **Error Logging**: Always include error details when logging errors:
   ```go
   if err != nil {
       logger.WithField("error", err.Error()).Error("Failed to process request")
   }
   ```

## Logging in Controllers

Controllers should log:
1. The start of request processing
2. Key parameters from the request
3. Errors that occur during processing
4. Successful completion of the request

Example:
```go
func SomeController(c *gin.Context) {
    // Get request ID from context if available
    requestID, exists := c.Get("requestID")
    var ctx context.Context
    if exists {
        ctx = context.WithValue(context.Background(), log.RequestIDKey, requestID.(string))
    } else {
        ctx = context.Background()
    }

    // Get user info if available
    userEmail, emailExists := c.Get("userEmail")
    userName, nameExists := c.Get("userName")
    if emailExists && nameExists {
        ctx = log.ContextWithUserInfo(ctx, userEmail.(string), userName.(string))
    }

    logger := log.WithContext(ctx)
    
    // Log the start of request processing with key parameters
    id := c.Param("id")
    logger.WithField("id", id).Info("Processing SomeController request")
    
    // Log user information
    user := c.MustGet("user").(models.AppUser)
    logger.WithFields(map[string]interface{}{
        "user_email": user.UserEmail,
        "user_name": user.UserName,
    }).Debug("User retrieved from context")

    // Log service call
    logger.Info("Calling service method")
    result, err := services.SomeService(id, user)
    if err != nil {
        // Log error with details
        logger.WithField("error", err.Error()).Error("Failed to process request")
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    // Log successful completion
    logger.WithField("result_count", len(result)).Info("Request processed successfully")
    c.JSON(http.StatusOK, result)
}
```

## Logging in Services

Services should log:
1. The start of service method execution
2. Key parameters and context
3. Database operations
4. External API calls
5. Errors that occur during processing
6. Successful completion of the operation

Example:
```go
func SomeService(id string, user models.AppUser) ([]SomeResult, error) {
    logger := appLog.WithFields(map[string]interface{}{
        "user_email": user.UserEmail,
        "user_name": user.UserName,
        "id": id,
        "action": "SomeService",
    })

    logger.Info("Starting service operation")

    // Log database operations
    logger.Debug("Retrieving data from database")
    var results []SomeResult
    err := DB.Where("id = ?", id).Find(&results).Error
    if err != nil {
        logger.WithField("error", err.Error()).Error("Failed to retrieve data from database")
        return nil, err
    }
    
    logger.WithField("result_count", len(results)).Debug("Retrieved data from database")

    // Log processing steps
    logger.Debug("Processing data")
    // ... processing code ...

    // Log successful completion
    logger.WithField("result_count", len(results)).Info("Service operation completed successfully")
    return results, nil
}
```

## Logging Database Operations

Database operations should be logged with:
1. The SQL query (when appropriate)
2. The number of records affected
3. The execution time for slow queries
4. Any errors that occur

The application already has a custom GORM logger that logs database operations. Make sure it's properly configured in all database connections.

## Logging Middleware

Middleware should log:
1. Incoming requests with method, path, and client IP
2. Response status code and latency
3. Authentication and authorization events
4. Any errors that occur during middleware processing

The application already has a RequestLoggerMiddleware that logs requests. Make sure it's properly configured in all routes.

## Best Practices

1. **Be Consistent**: Use the same logging patterns throughout the application
2. **Don't Log Sensitive Data**: Avoid logging passwords, tokens, or other sensitive information
3. **Use Appropriate Log Levels**: Don't log everything at the same level
4. **Include Context**: Always include relevant context in logs
5. **Log Operation Boundaries**: Log the start and end of important operations
6. **Log Errors with Details**: Include error details when logging errors
7. **Use Structured Logging**: Use fields instead of string formatting
8. **Include Request IDs**: Include request IDs in all logs to trace requests across services

## Implementation Checklist

When adding logging to a new component:

- [ ] Import the logging package (`api/log`)
- [ ] Create a logger with appropriate context
- [ ] Log the start of operations
- [ ] Log key parameters and context
- [ ] Log errors with details
- [ ] Log successful completion
- [ ] Use appropriate log levels
- [ ] Use structured logging with fields
- [ ] Include request IDs when available
- [ ] Include user information when available