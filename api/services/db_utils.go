package services

import (
	appLog "api/log"
	"context"
	"fmt"
	"time"
)

// QueryWithContext executes a database query with a timeout context
// If the query takes longer than the timeout, it will be canceled
func QueryWithContext(timeout time.Duration, queryFn func(ctx context.Context) error) error {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create a channel to receive the query result
	errCh := make(chan error, 1)

	// Execute the query in a goroutine
	go func() {
		errCh <- queryFn(ctx)
	}()

	// Wait for the query to complete or timeout
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("query timed out after %v", timeout)
		}
		return ctx.Err()
	}
}

// QueryWithContextAndCache executes a database query with a timeout context and caching
// If the result is found in the cache, it is returned without executing the query
// If the query takes longer than the timeout, it will be canceled
func QueryWithContextAndCache(cacheKey string, result interface{}, timeout time.Duration, queryFn func(ctx context.Context) error) error {
	logger := appLog.WithFields(map[string]interface{}{
		"cache_key": cacheKey,
		"timeout":   timeout.String(),
		"action":    "QueryWithContextAndCache",
	})

	// Try to get the result from the cache
	if cachedResult, found := QueryCacheInstance.Get(cacheKey); found {
		logger.Debug("Cache hit")
		
		// Copy the cached result to the result parameter
		switch typedResult := result.(type) {
		case *[]interface{}:
			*typedResult = cachedResult.([]interface{})
		case *interface{}:
			*typedResult = cachedResult
		default:
			// For other types, we need to use reflection or type assertions
			// This is a simplified implementation
			logger.Debug("Cache hit but result type not handled, executing query")
			goto executeQuery
		}
		
		return nil
	}

	logger.Debug("Cache miss, executing query")

executeQuery:
	// Execute the query with context
	err := QueryWithContext(timeout, func(ctx context.Context) error {
		return queryFn(ctx)
	})

	if err != nil {
		logger.WithField("error", err.Error()).Error("Query execution failed")
		return err
	}

	// Store the result in the cache
	QueryCacheInstance.Set(cacheKey, result)
	logger.Debug("Query result cached")

	return nil
}

// ExecuteBatchOperation executes a batch operation on a slice of items
// The batchSize parameter determines how many items are processed in each batch
// The processFn parameter is a function that processes a batch of items
func ExecuteBatchOperation(items interface{}, batchSize int, processFn func(batch interface{}) error) error {
	// TODO: Implement batch processing logic
	// This is a placeholder for future implementation
	return nil
}