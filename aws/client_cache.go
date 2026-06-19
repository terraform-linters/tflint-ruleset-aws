package aws

import "sync"

// clientFactory builds a Client from credentials. Defaults to NewClient;
// overridable in tests to count and control client construction.
type clientFactory func(Credentials) (Client, error)

type clientEntry struct {
	once   sync.Once
	client Client
	err    error
}

// clientCache builds at most one Client per distinct Credentials value and
// shares it across all runners. Safe for concurrent use.
type clientCache struct {
	factory clientFactory
	mu      sync.Mutex
	entries map[Credentials]*clientEntry
}

func newClientCache(factory clientFactory) *clientCache {
	return &clientCache{factory: factory, entries: map[Credentials]*clientEntry{}}
}

func (c *clientCache) get(creds Credentials) (Client, error) {
	c.mu.Lock()
	entry, ok := c.entries[creds]
	if !ok {
		entry = &clientEntry{}
		c.entries[creds] = entry
	}
	c.mu.Unlock()

	entry.once.Do(func() {
		entry.client, entry.err = c.factory(creds)
	})

	if entry.err != nil {
		// Do not cache failures: drop the entry so a later call retries.
		c.mu.Lock()
		if c.entries[creds] == entry {
			delete(c.entries, creds)
		}
		c.mu.Unlock()
		return nil, entry.err
	}
	return entry.client, nil
}
