package aws

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
)

func Test_clientCache_get(t *testing.T) {
	cases := []struct {
		name  string
		creds []Credentials
		calls int
	}{
		{
			name: "same credential reused across sequential calls",
			creds: []Credentials{
				{Region: "us-east-1"},
				{Region: "us-east-1"},
				{Region: "us-east-1"},
			},
			calls: 1,
		},
		{
			name: "distinct credentials build distinct clients",
			creds: []Credentials{
				{Region: "us-east-1"},
				{Region: "us-west-2"},
				{Region: "us-east-1"},
			},
			calls: 2,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var count int64
			cache := newClientCache(func(Credentials) (Client, error) {
				atomic.AddInt64(&count, 1)
				return &AwsClient{}, nil
			})

			clients := map[Credentials]Client{}
			for _, creds := range tc.creds {
				client, err := cache.get(creds)
				if err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
				if prev, ok := clients[creds]; ok {
					if prev != client {
						t.Fatalf("expected the same client for %+v, got a different one", creds)
					}
				} else {
					clients[creds] = client
				}
			}

			if got := atomic.LoadInt64(&count); got != int64(tc.calls) {
				t.Fatalf("expected factory to be called %d times, got %d", tc.calls, got)
			}
		})
	}
}

func Test_clientCache_get_concurrent(t *testing.T) {
	const goroutines = 32

	var count int64
	release := make(chan struct{})
	var entered sync.WaitGroup
	entered.Add(goroutines)

	cache := newClientCache(func(Credentials) (Client, error) {
		atomic.AddInt64(&count, 1)
		return &AwsClient{}, nil
	})

	creds := Credentials{Region: "us-east-1"}
	clients := make([]Client, goroutines)
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := range goroutines {
		go func(i int) {
			defer wg.Done()
			entered.Done()
			<-release
			client, err := cache.get(creds)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}
			clients[i] = client
		}(i)
	}

	entered.Wait()
	close(release)
	wg.Wait()

	if got := atomic.LoadInt64(&count); got != 1 {
		t.Fatalf("expected factory to be called exactly once, got %d", got)
	}
	for i, client := range clients {
		if client != clients[0] {
			t.Fatalf("expected all goroutines to share one client, goroutine %d differs", i)
		}
	}
}

func Test_clientCache_get_doesNotCacheErrors(t *testing.T) {
	var count int64
	wantErr := errors.New("transient failure")
	cache := newClientCache(func(Credentials) (Client, error) {
		if atomic.AddInt64(&count, 1) == 1 {
			return nil, wantErr
		}
		return &AwsClient{}, nil
	})

	creds := Credentials{Region: "us-east-1"}

	if _, err := cache.get(creds); !errors.Is(err, wantErr) {
		t.Fatalf("expected first get to return %v, got %v", wantErr, err)
	}

	client, err := cache.get(creds)
	if err != nil {
		t.Fatalf("expected second get to succeed, got %v", err)
	}
	if client == nil {
		t.Fatal("expected second get to return a client")
	}
	if got := atomic.LoadInt64(&count); got != 2 {
		t.Fatalf("expected factory to be called twice, got %d", got)
	}
}
