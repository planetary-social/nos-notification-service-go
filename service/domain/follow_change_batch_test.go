package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/planetary-social/go-notification-service/internal/fixtures"
	"github.com/planetary-social/go-notification-service/service/domain"
	"github.com/stretchr/testify/assert"
)

func TestFollowChangeBatch_UnmarshalJSON_Valid(t *testing.T) {
	pk1, pk1Npub := fixtures.PublicKeyAndNpub()
	pk2, pk2Npub := fixtures.PublicKeyAndNpub()
	pk3, pk3Npub := fixtures.PublicKeyAndNpub()

	jsonData := `{
		"followee": "` + pk1Npub + `",
		"friendlyFollower": "FriendlyUser",
		"follows": ["` + pk2Npub + `", "` + pk3Npub + `"]
	}`

	var batch domain.FollowChangeBatch
	err := json.Unmarshal([]byte(jsonData), &batch)

	assert.NoError(t, err)
	assert.Equal(t, pk1, batch.Followee)
	assert.Equal(t, "FriendlyUser", batch.FriendlyFollower)
	assert.Equal(t, []domain.PublicKey{pk2, pk3}, batch.Follows)
}

func TestFollowChangeBatch_UnmarshalJSON_InvalidFollowee(t *testing.T) {
	_, pk1Npub := fixtures.PublicKeyAndNpub()

	jsonData := `{
		"followee": "invalid",
		"friendlyFollower": "FriendlyUser",
		"follows": ["` + pk1Npub + `"]
	}`

	var batch domain.FollowChangeBatch
	err := json.Unmarshal([]byte(jsonData), &batch)
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid npub for followee: error decoding a nip19 entity: invalid bech32 string length 7")
}

func TestFollowChangeBatch_UnmarshalJSON_InvalidFollows(t *testing.T) {
	_, pk1Npub := fixtures.PublicKeyAndNpub()

	jsonData := `{
		"followee": "` + pk1Npub + `",
		"friendlyFollower": "FriendlyUser",
		"follows": ["invalid"]
	}`

	var batch domain.FollowChangeBatch
	err := json.Unmarshal([]byte(jsonData), &batch)
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid npub for follow: error decoding a nip19 entity: invalid bech32 string length 7")
}

func TestFollowChangeBatch_String_SingleFollow(t *testing.T) {
	pk1, pk1Npub := fixtures.PublicKeyAndNpub()
	pk2, _ := fixtures.PublicKeyAndNpub()

	batch := domain.FollowChangeBatch{
		Followee:         pk1,
		FriendlyFollower: "FriendlyUser",
		Follows:          []domain.PublicKey{pk2},
	}

	expected := "Follow: FriendlyUser -----> " + pk1Npub
	assert.Equal(t, expected, batch.String())
}

func TestFollowChangeBatch_String_MultipleFollows(t *testing.T) {
	pk1, pk1Npub := fixtures.PublicKeyAndNpub()
	pk2, _ := fixtures.PublicKeyAndNpub()
	pk3, _ := fixtures.PublicKeyAndNpub()

	batch := domain.FollowChangeBatch{
		Followee:         pk1,
		FriendlyFollower: "FriendlyUser",
		Follows:          []domain.PublicKey{pk2, pk3},
	}

	expected := "Follow aggregate: 2 followers for " + pk1Npub
	assert.Equal(t, expected, batch.String())
}
