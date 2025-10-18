package service

import (
	"google.golang.org/grpc"
	"context"
	"sync"

	"google.golang.org/grpc/status"

	"google.golang.org/grpc/codes"

)

type Limiters struct{
	mu sync.Mutex
	active map[string]int
	maxUp int
	maxDo int
	maxList int
}

func NewLimiters() *Limiters{
	return &Limiters{
		active: map[string]int{},
		maxUp:  10,
		maxDo: 10,
		maxList:    100,
	}
}

func (l *Limiters) StreamInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error{
	method := info.FullMethod
	var limit int
	if methodContains(method, "Upload"){
		limit = l.maxUp
	} else if methodContains(method, "Download"){
		limit = l.maxList
	} else {
		return handler(srv, ss)
	}

	if !l.tryAcquire(method, limit){
		return status.Error(codes.ResourceExhausted, "Error Concurrent Stream Req")
	}
	defer l.release(method)

	return handler(srv, ss)
}

func (l *Limiters) tryAcquire(method string, limit int) bool{
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.active[method] >= limit{
		return false
	}
	l.active[method]++
	return true
}

func (l *Limiters) release(method string){
	l.mu.Lock()
	defer l.mu.Unlock()
	l.active[method]--
}

func containsIgnoreCase(s, substr string) bool{
	for i := 0; i + len(substr) <= len(s); i++ {
		match := true
		for j := range substr{
			c1 := s[i + j]
			c2 := substr[j]
			if c1 >= 'A' && c1 <= 'Z'{
				c1 += 32
			}
			if c2 >= 'A' && c2 <= 'Z'{
				c2 += 32
			}
			if c1 != c2{
				match = false
				break
			}
		}
		if match{
			return true
		}
	}
	return false
}

func methodContains(method, part string) bool{
	return len(method) > 0 && (method == part || containsIgnoreCase(method, part))
}

func (l *Limiters) UnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	method := info.FullMethod
	if methodContains(method, "List"){
		if !l.tryAcquire(method, l.maxList){
			return nil, status.Error(codes.ResourceExhausted, "Error Concurrent List")
		}
		defer l.release(method)
	}
	return handler(ctx, req)
}