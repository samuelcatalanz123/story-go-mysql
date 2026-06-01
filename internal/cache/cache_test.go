package cache

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
)

// newTestCache starts an in-memory Redis (miniredis) and returns a Redis cache
// pointed at it, so tests need no real Redis server.
func newTestCache(t *testing.T) *Redis {
	t.Helper()
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(mr.Close)
	c, err := NewRedis(context.Background(), mr.Addr(), "")
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestSetAndGet(t *testing.T) {
	c := newTestCache(t)
	ctx := context.Background()

	if _, ok, _ := c.Get(ctx, "missing"); ok {
		t.Fatal("una clave inexistente no debería encontrarse")
	}

	if err := c.Set(ctx, "k", []byte("hola"), time.Minute); err != nil {
		t.Fatal(err)
	}
	got, ok, err := c.Get(ctx, "k")
	if err != nil || !ok || string(got) != "hola" {
		t.Fatalf("esperaba 'hola', obtuve %q ok=%v err=%v", got, ok, err)
	}
}

func TestIncrCountsUp(t *testing.T) {
	c := newTestCache(t)
	ctx := context.Background()
	for want := int64(1); want <= 3; want++ {
		got, err := c.Incr(ctx, "counter", time.Minute)
		if err != nil || got != want {
			t.Fatalf("incr %d: obtuve %d (err %v)", want, got, err)
		}
	}
}

func TestDelByPrefix(t *testing.T) {
	c := newTestCache(t)
	ctx := context.Background()
	_ = c.Set(ctx, "characters:list:1", []byte("a"), time.Minute)
	_ = c.Set(ctx, "characters:list:2", []byte("b"), time.Minute)
	_ = c.Set(ctx, "other:key", []byte("c"), time.Minute)

	if err := c.DelByPrefix(ctx, "characters:list:"); err != nil {
		t.Fatal(err)
	}
	if _, ok, _ := c.Get(ctx, "characters:list:1"); ok {
		t.Fatal("la clave con prefijo debería haberse borrado")
	}
	if _, ok, _ := c.Get(ctx, "other:key"); !ok {
		t.Fatal("las claves de otro prefijo no deberían borrarse")
	}
}
