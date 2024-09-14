package domain

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/nbd-wtf/go-nostr/nip19"
)

// This is the struct coming from the follow-change pubsub topic produced by
// the followers service. The type represents a struct with the same name
// there.
type FollowChangeBatch struct {
	Followee         PublicKey   `json:"followee"`
	FriendlyFollower string      `json:"friendlyFollower"`
	Follows          []PublicKey `json:"follows"`
}

func (f *FollowChangeBatch) UnmarshalJSON(data []byte) error {
	var temp struct {
		Followee         string   `json:"followee"`
		FriendlyFollower string   `json:"friendlyFollower"`
		Follows          []string `json:"follows"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	f.FriendlyFollower = temp.FriendlyFollower

	var err error
	f.Followee, err = NewPublicKeyFromNpub(temp.Followee)
	if err != nil {
		return errors.New("invalid npub for followee: " + err.Error())
	}

	f.Follows = make([]PublicKey, len(temp.Follows))
	for i, npub := range temp.Follows {
		f.Follows[i], err = NewPublicKeyFromNpub(npub)
		if err != nil {
			return errors.New("invalid npub for follow: " + err.Error())
		}
	}

	return nil
}

func (f FollowChangeBatch) String() string {
	friendlyFollowee, err := nip19.EncodePublicKey(f.Followee.Hex())
	if err != nil {
		friendlyFollowee = f.Followee.Hex()
	}

	if len(f.Follows) == 1 {
		return fmt.Sprintf("Follow: %s -----> %s", f.FriendlyFollower, friendlyFollowee)
	}

	return fmt.Sprintf("Follow aggregate: %d followers for %s", len(f.Follows), friendlyFollowee)
}
