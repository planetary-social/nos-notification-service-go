package apns_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/planetary-social/go-notification-service/internal/fixtures"
	"github.com/planetary-social/go-notification-service/service/adapters/apns"
	"github.com/planetary-social/go-notification-service/service/domain"
	"github.com/stretchr/testify/require"
)

func TestFollowChangePayload_SingleFollow(t *testing.T) {
	pk1, pk1Npub := fixtures.PublicKeyAndNpub()
	pk2, pk2Npub := fixtures.PublicKeyAndNpub()

	batch := domain.FollowChangeBatch{
		Followee:         pk1,
		FriendlyFollower: "npub_someFollower",
		Follows:          []domain.PublicKey{pk2},
	}

	payload, err := apns.FollowChangePayload(batch)
	require.NoError(t, err)

	expectedAlert := "You have a new follower!"
	expectedPayload := map[string]interface{}{
		"aps": map[string]interface{}{
			"alert":              expectedAlert,
			"sound":              "default",
			"badge":              float64(1), // Convert badge to float64
			"thread-id":          pk1Npub,
			"interruption-level": "passive",
		},
		"data": map[string]interface{}{
			"follows":          []interface{}{pk2Npub}, // Use []interface{}
			"friendlyFollower": batch.FriendlyFollower,
		},
	}

	var actualPayload map[string]interface{}
	err = json.Unmarshal(payload, &actualPayload)
	require.NoError(t, err)

	require.Equal(t, expectedPayload, actualPayload)
}

func TestFollowChangePayload_MultipleFollowsUnfollows(t *testing.T) {
	pk1, pk1Npub := fixtures.PublicKeyAndNpub()
	pk2, pk2Npub := fixtures.PublicKeyAndNpub()
	pk3, pk3Npub := fixtures.PublicKeyAndNpub()

	batch := domain.FollowChangeBatch{
		Followee:         pk1,
		FriendlyFollower: "FriendlyUser",
		Follows:          []domain.PublicKey{pk2, pk3},
	}

	payload, err := apns.FollowChangePayload(batch)
	require.NoError(t, err)

	expectedAlert := "You have 2 new followers!"
	expectedPayload := map[string]interface{}{
		"aps": map[string]interface{}{
			"alert":              expectedAlert,
			"sound":              "default",
			"badge":              float64(1), // Convert badge to float64
			"thread-id":          pk1Npub,
			"interruption-level": "passive",
		},
		"data": map[string]interface{}{
			"follows": []interface{}{pk2Npub, pk3Npub}, // Use []interface{}
		},
	}

	var actualPayload map[string]interface{}
	err = json.Unmarshal(payload, &actualPayload)
	require.NoError(t, err)

	require.Equal(t, expectedPayload, actualPayload)
}

func TestFollowChangePayload_SingleFollow_WithFriendlyFollower(t *testing.T) {
	pk1, pk1Npub := fixtures.PublicKeyAndNpub()
	pk2, pk2Npub := fixtures.PublicKeyAndNpub()

	batch := domain.FollowChangeBatch{
		Followee:         pk1,
		FriendlyFollower: "John Doe",
		Follows:          []domain.PublicKey{pk2},
	}

	payload, err := apns.FollowChangePayload(batch)
	require.NoError(t, err)

	expectedAlert := "John Doe is a new follower!"
	expectedPayload := map[string]interface{}{
		"aps": map[string]interface{}{
			"alert":              expectedAlert,
			"sound":              "default",
			"badge":              float64(1), // Convert badge to float64
			"thread-id":          pk1Npub,
			"interruption-level": "passive",
		},
		"data": map[string]interface{}{
			"follows":          []interface{}{pk2Npub}, // Use []interface{}
			"friendlyFollower": batch.FriendlyFollower,
		},
	}

	var actualPayload map[string]interface{}
	err = json.Unmarshal(payload, &actualPayload)
	require.NoError(t, err)

	// jsonStr, err := json.Marshal(actualPayload)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }
	// fmt.Println(string(jsonStr))

	require.Equal(t, expectedPayload, actualPayload)
}

func TestFollowChangePayload_BatchedFollow_WithNoFriendlyFollower(t *testing.T) {
	pk1, pk1Npub := fixtures.PublicKeyAndNpub()
	pk2, pk2Npub := fixtures.PublicKeyAndNpub()
	pk3, pk3Npub := fixtures.PublicKeyAndNpub()

	batch := domain.FollowChangeBatch{
		Followee: pk1,
		Follows:  []domain.PublicKey{pk2, pk3},
	}

	payload, err := apns.FollowChangePayload(batch)
	require.NoError(t, err)

	expectedAlert := "You have 2 new followers!"
	expectedPayload := map[string]interface{}{
		"aps": map[string]interface{}{
			"alert":              expectedAlert,
			"sound":              "default",
			"badge":              float64(1), // Convert badge to float64
			"thread-id":          pk1Npub,
			"interruption-level": "passive",
		},
		"data": map[string]interface{}{
			"follows": []interface{}{pk2Npub, pk3Npub}, // Use []interface{}
		},
	}

	var actualPayload map[string]interface{}
	err = json.Unmarshal(payload, &actualPayload)
	require.NoError(t, err)

	// jsonStr, err := json.Marshal(actualPayload)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }
	// fmt.Println(string(jsonStr))

	require.Equal(t, expectedPayload, actualPayload)
}
func TestFollowChangePayload_Exceeds4096Bytes_With59TotalNpubs(t *testing.T) {
	pk1, _ := fixtures.PublicKeyAndNpub()

	batch := domain.FollowChangeBatch{
		Followee:         pk1,
		FriendlyFollower: "npub_someFollower_wont_be_added",
		Follows:          []domain.PublicKey{},
	}

	for i := 0; i < 59; i++ {
		follow, _ := fixtures.PublicKeyAndNpub()
		batch.Follows = append(batch.Follows, follow)
	}

	payload, err := apns.FollowChangePayloadWithValidation(batch, false)
	require.NoError(t, err)

	// Ensure 59 is the maximum size we can get
	payloadSize := len(payload)
	t.Logf("Payload size with 59 total follows and unfollows: %d bytes", payloadSize)
	require.True(t, payloadSize > 4096, fmt.Sprintf("Payload size should exceed 4096 bytes, but was %d bytes", payloadSize))
}

func TestFollowChangePayload_ValidPayload_With58TotalNpubs_IsValid(t *testing.T) {
	pk1, _ := fixtures.PublicKeyAndNpub()

	batch := domain.FollowChangeBatch{
		Followee:         pk1,
		FriendlyFollower: "npub_someFollower_wont_be_added",
		Follows:          []domain.PublicKey{},
	}

	for i := 0; i < 58; i++ {
		follow, _ := fixtures.PublicKeyAndNpub()
		batch.Follows = append(batch.Follows, follow)
	}

	payload, err := apns.FollowChangePayloadWithValidation(batch, true) // With validation
	require.NoError(t, err)

	// Ensure 58 is the maximum size we can get
	payloadSize := len(payload)
	t.Logf("Payload size with 58 total follows and unfollows: %d bytes", payloadSize)
	require.True(t, payloadSize <= 4096, fmt.Sprintf("Payload size should be within 4096 bytes, but was %d bytes", payloadSize))
}

func TestFollowChangePayload_InvalidPayload_With59TotalNpubs_Fails_With_Validation(t *testing.T) {
	pk1, _ := fixtures.PublicKeyAndNpub()

	batch := domain.FollowChangeBatch{
		Followee:         pk1,
		FriendlyFollower: "npub_someFollower_wont_be_added",
		Follows:          []domain.PublicKey{},
	}

	for i := 0; i < 59; i++ { // 29 follows
		follow, _ := fixtures.PublicKeyAndNpub()
		batch.Follows = append(batch.Follows, follow)
	}

	payload, err := apns.FollowChangePayload(batch) // This always validates
	require.Error(t, err)
	require.Nil(t, payload)

	expectedError := fmt.Sprintf("FollowChangeBatch for followee %s has too many npubs (59). MAX_TOTAL_NPUBS is 58", pk1.Hex())
	require.EqualError(t, err, expectedError)
}
