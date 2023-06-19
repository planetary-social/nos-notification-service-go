package domain

type Registration struct {
	publicKeys []PublicKeyWithRelays
	locale     Locale
	apnsToken  APNSToken
}

type PublicKeyWithRelays struct {
	publicKey PublicKey
	relays    Relays
}
