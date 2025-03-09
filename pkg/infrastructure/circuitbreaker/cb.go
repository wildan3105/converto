package circuitbreaker

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type State string

const (
	Closed   State = "Closed"
	Open     State = "Open"
	HalfOpen State = "Half-Open"
)

var ErrCircuitBreakerOpen = errors.New("circuit breaker is open")

type CircuitBreaker struct {
	mu               sync.Mutex
	state            State
	failures         int
	failureThreshold int
	cooldownTime     time.Duration
	lastFailureTime  time.Time
}

// NewCircuitBreaker initializes a new circuit breaker
func NewCircuitBreaker(failureThreshold int, cooldownTime time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:            Closed,
		failureThreshold: failureThreshold,
		cooldownTime:     cooldownTime,
	}
}

// Execute runs a function with circuit breaker protection
func (cb *CircuitBreaker) Execute(fn func() error) error {
	cb.mu.Lock()

	switch cb.state {
	case Open:
		if time.Since(cb.lastFailureTime) > cb.cooldownTime {
			cb.state = HalfOpen
			fmt.Println("Circuit breaker transitioning to half-open state.")
		} else {
			cb.mu.Unlock()
			return ErrCircuitBreakerOpen
		}
	case HalfOpen:
		fmt.Println("Circuit breaker is in half-open state. Trying limited requests.")
	}

	cb.mu.Unlock()

	// Execute the function
	err := fn()

	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failures++
		cb.lastFailureTime = time.Now()

		if cb.failures >= cb.failureThreshold {
			cb.state = Open
			fmt.Println("Circuit breaker opened!")
		}
		return err
	}

	// On success, reset the circuit breaker
	cb.reset()
	return nil
}

// reset brings the circuit breaker to its initial state
func (cb *CircuitBreaker) reset() {
	cb.failures = 0
	cb.state = Closed
	fmt.Println("Circuit breaker reset to closed state.")
}
