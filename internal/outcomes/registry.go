package outcomes

import (
	"fmt"
	"sync"
)

// Registry manages all available outcomes
type Registry struct {
	mu       sync.RWMutex
	outcomes map[string]*Outcome
}

// NewRegistry creates a new outcome registry
func NewRegistry() *Registry {
	return &Registry{
		outcomes: make(map[string]*Outcome),
	}
}

// Register adds an outcome to the registry
func (r *Registry) Register(outcome *Outcome) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if outcome.ID == "" {
		return fmt.Errorf("outcome ID cannot be empty")
	}

	if _, exists := r.outcomes[outcome.ID]; exists {
		return fmt.Errorf("outcome with ID %s already registered", outcome.ID)
	}

	r.outcomes[outcome.ID] = outcome
	return nil
}

// Get retrieves an outcome by ID
func (r *Registry) Get(id string) (*Outcome, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	outcome, exists := r.outcomes[id]
	if !exists {
		return nil, fmt.Errorf("outcome not found: %s", id)
	}

	return outcome, nil
}

// List returns all registered outcomes
func (r *Registry) List() []*Outcome {
	r.mu.RLock()
	defer r.mu.RUnlock()

	outcomes := make([]*Outcome, 0, len(r.outcomes))
	for _, outcome := range r.outcomes {
		outcomes = append(outcomes, outcome)
	}

	return outcomes
}

// ListByCategory returns outcomes filtered by category
func (r *Registry) ListByCategory(category string) []*Outcome {
	r.mu.RLock()
	defer r.mu.RUnlock()

	outcomes := make([]*Outcome, 0)
	for _, outcome := range r.outcomes {
		if outcome.Category == category {
			outcomes = append(outcomes, outcome)
		}
	}

	return outcomes
}

// Categories returns all unique categories
func (r *Registry) Categories() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	categories := make(map[string]bool)
	for _, outcome := range r.outcomes {
		if outcome.Category != "" {
			categories[outcome.Category] = true
		}
	}

	result := make([]string, 0, len(categories))
	for cat := range categories {
		result = append(result, cat)
	}

	return result
}
