package server

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/ahmad-masud/KVStore/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func startTestServer(t *testing.T) (proto.KVStoreClient, func()) {
	t.Helper()

	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	s := NewServer()

	grpcServer := grpc.NewServer()
	proto.RegisterKVStoreServer(grpcServer, s)

	go func() {
		if err := grpcServer.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			t.Logf("server exited: %v", err)
		}
	}()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}

	client := proto.NewKVStoreClient(conn)

	return client, func() {
		conn.Close()
		grpcServer.Stop()
	}
}

func TestServer_SetAndGet(t *testing.T) {
	client, cleanup := startTestServer(t)
	defer cleanup()

	ctx := context.Background()

	// Set key
	_, err := client.Set(ctx, &proto.SetRequest{
		Key:   "foo",
		Value: "bar",
	})
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Get key
	resp, err := client.Get(ctx, &proto.GetRequest{
		Key: "foo",
	})
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if !resp.Found || resp.Value != "bar" {
		t.Fatalf("unexpected Get response: %+v", resp)
	}
}

func TestServer_Delete(t *testing.T) {
	client, cleanup := startTestServer(t)
	defer cleanup()

	ctx := context.Background()

	// Set key
	_, err := client.Set(ctx, &proto.SetRequest{
		Key:   "foo",
		Value: "bar",
	})
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Delete key
	delResp, err := client.Delete(ctx, &proto.DeleteRequest{
		Key: "foo",
	})
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if !delResp.Success {
		t.Fatalf("expected Delete to succeed")
	}

	// Try getting deleted key
	getResp, err := client.Get(ctx, &proto.GetRequest{
		Key: "foo",
	})
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if getResp.Found {
		t.Fatalf("expected key to be deleted, but found: %+v", getResp)
	}
}

func TestServer_TTLExpiration(t *testing.T) {
	client, cleanup := startTestServer(t)
	defer cleanup()

	ctx := context.Background()

	// Set key with short TTL
	_, err := client.Set(ctx, &proto.SetRequest{
		Key:        "baz",
		Value:      "qux",
		Ttl: 1, // expires in 1 second
	})
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Immediately check
	getResp, err := client.Get(ctx, &proto.GetRequest{
		Key: "baz",
	})
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if !getResp.Found || getResp.Value != "qux" {
		t.Fatalf("unexpected immediate Get: %+v", getResp)
	}

	// Wait for TTL to expire
	time.Sleep(2 * time.Second)

	getResp, err = client.Get(ctx, &proto.GetRequest{
		Key: "baz",
	})
	if err != nil {
		t.Fatalf("Get failed after sleep: %v", err)
	}
	if getResp.Found {
		t.Fatalf("expected key to expire, but found: %+v", getResp)
	}
}
