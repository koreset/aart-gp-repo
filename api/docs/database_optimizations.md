# Database Query Optimization

This document outlines the database query optimizations implemented to improve application performance.

## Overview of Changes

The following optimizations have been implemented:

1. **Query Caching**: Added an in-memory cache for frequently accessed data
2. **Context with Timeouts**: Added context with timeouts to prevent long-running queries
3. **Database Connection Pooling**: Optimized connection pool settings
4. **Database Indexes**: Added indexes to frequently queried tables
5. **Query Optimization**: Refactored inefficient queries to reduce database load

## Detailed Changes

### 1. Query Caching

A simple in-memory cache has been implemented in `services/query_cache.go`. This cache:

- Stores query results with a configurable TTL (Time-To-Live)
- Automatically cleans up expired entries
- Is thread-safe for concurrent access
- Provides a simple API for caching query results

Example usage:

```go
cacheKey := fmt.Sprintf("group_pricing_quote_%s", id)
err = QueryWithContextAndCache(cacheKey, &result, 5*time.Minute, func(ctx context.Context) error {
    // Database query logic here
    return nil
})
```

### 2. Context with Timeouts

Database queries now use context with timeouts to prevent long-running queries from blocking the application. This is implemented in `services/db_utils.go`.

Example usage:

```go
err := QueryWithContext(5*time.Second, func(ctx context.Context) error {
    return DB.WithContext(ctx).Where("id = ?", id).First(&result).Error
})
```

### 3. Database Connection Pooling

The database connection pool settings have been optimized in `services/db.go`:

- Increased max open connections from 120 to 150
- Increased max idle connections from 10 to 25
- Added connection max lifetime of 1 hour
- Added connection max idle time of 30 minutes

### 4. Database Indexes

Indexes have been added to frequently queried tables to improve query performance. This is implemented in the `CreateDatabaseIndexes` function in `services/db.go`.

The following tables now have indexes:

- `aggregated_projections`: Indexes on run_id, product_code, sp_code, projection_month
- Group pricing tables: Indexes on quote_id
- User management tables: Indexes on email, subject
- Activity tracking: Indexes on user_email, date, object_type/id

### 5. Query Optimization

Several inefficient queries have been refactored to improve performance:

- `GetGroupPricingQuote`: Replaced multiple separate queries with a single query using subqueries
- `FindOrgUsers`: Added filtering to only query users for the specific organization
- `GetAggregations`: Simplified query construction and added caching

## Performance Impact

These optimizations should result in:

- Reduced database load
- Faster query response times
- Better application scalability
- Improved user experience

## Future Improvements

Additional optimizations that could be implemented in the future:

1. Distributed caching using Redis or Memcached
2. Query result pagination for large datasets
3. Asynchronous processing for long-running operations
4. Database query monitoring and profiling
5. Implement batch operations for bulk data processing