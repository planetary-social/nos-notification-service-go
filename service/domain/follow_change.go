package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type FollowChange struct {
	ChangeType       string    `json:"changeType"`
	At               time.Time `json:"at"`
	Follower         PublicKey `json:"follower"`
	Followee         PublicKey `json:"followee"`
	FriendlyFollowee string    `json:"friendlyFollowee"`
	FriendlyFollower string    `json:"friendlyFollower"`
}

func NewFollowChange(changeType string, follower PublicKey, friendlyFollower string, followee PublicKey, friendlyFollowee string, at time.Time) FollowChange {
	return FollowChange{
		ChangeType:       changeType,
		Follower:         follower,
		FriendlyFollower: friendlyFollower,
		Followee:         followee,
		FriendlyFollowee: friendlyFollowee,
		At:               at,
	}
}

func (f *FollowChange) UnmarshalJSON(data []byte) error {
	var temp struct {
		ChangeType       string `json:"changeType"`
		At               int64  `json:"at"`
		Follower         string `json:"follower"`
		Followee         string `json:"followee"`
		FriendlyFollowee string `json:"friendlyFollowee"`
		FriendlyFollower string `json:"friendlyFollower"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	f.ChangeType = temp.ChangeType
	f.At = time.Unix(temp.At, 0)
	f.FriendlyFollowee = temp.FriendlyFollowee
	f.FriendlyFollower = temp.FriendlyFollower

	var err error
	f.Follower, err = NewPublicKeyFromHex(temp.Follower)
	if err != nil {
		return errors.New("invalid hex for follower: " + err.Error())
	}

	f.Followee, err = NewPublicKeyFromHex(temp.Followee)
	if err != nil {
		return errors.New("invalid hex for followee: " + err.Error())
	}

	return nil
}

func (f FollowChange) String() string {
	if f.ChangeType == "unfollowed" {
		return fmt.Sprintf("New unfollow: %s(%s) --x--> %s(%s)", f.FriendlyFollower, f.Follower.s, f.FriendlyFollowee, f.Followee.s)
	}

	return fmt.Sprintf("New follow: %s(%s) -----> %s(%s)", f.FriendlyFollower, f.Follower.s, f.FriendlyFollowee, f.Followee.s)
}
