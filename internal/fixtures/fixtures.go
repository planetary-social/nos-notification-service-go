package fixtures

import (
	"context"
	cryptorand "crypto/rand"
	"encoding/hex"
	"fmt"
	"math/rand"
	"testing"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
	"github.com/planetary-social/go-notification-service/internal"
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
	v, err := domain.NewPublicKeyFromHex(p)
	if err != nil {
		panic(err)
	}
	return v, hex
}

func SomeRelayAddress() domain.RelayAddress {
	protocol := internal.RandomElement([]string{"ws", "wss"})
	address := fmt.Sprintf("%s://%s", protocol, SomeString())

	v, err := domain.NewRelayAddress(address)
	if err != nil {
		panic(err)
	}
	return v
}

func SomeString() string {
	return randSeq(10)
}

func SomeAPNSToken() domain.APNSToken {
	v, err := domain.NewAPNSTokenFromHex(SomeHexBytesOfLen(10 + rand.Intn(10)))
	if err != nil {
		panic(err)
	}
	return v
}

func SomeHexBytesOfLen(l int) string {
	b := make([]byte, l)
	n, err := cryptorand.Read(b)
	if n != len(b) {
		panic("short read")
	}
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func PublicKeyAndNpub() (domain.PublicKey, string) {
	pk, _ := SomeKeyPair()
	npub, _ := nip19.EncodePublicKey(pk.Hex())
	return pk, npub
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
