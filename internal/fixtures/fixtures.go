package fixtures

import (
	"context"
	"math/rand"
	"testing"

	"github.com/nbd-wtf/go-nostr"
	"github.com/planetary-social/go-notification-service/service/domain"
)

func Context(tb testing.TB) context.Context {
	ctx, cancelFunc := context.WithCancel(context.Background())
	tb.Cleanup(cancelFunc)
	return ctx
}

func somePrivateKeyHex() string {
	return nostr.GeneratePrivateKey()
}

func SomeKeyPair() (publicKey domain.PublicKey, secretKeyHex string) {
	hex := somePrivateKeyHex()

	p, err := nostr.GetPublicKey(hex)
	if err != nil {
		panic(err)
	}
	v, err := domain.NewPublicKey(p)
	if err != nil {
		panic(err)
	}
	return v, hex
}

func SomeRelayAddress() domain.RelayAddress {
	v, err := domain.NewRelayAddress(SomeString())
	if err != nil {
		panic(err)
	}
	return v
}

func SomeString() string {
	return randSeq(10)
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
