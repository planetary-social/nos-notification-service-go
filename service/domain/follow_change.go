package domain

type FollowChange struct {
	ChangeType       string    `json:"changeType"`
	At               uint      `json:"at"`
	Follower         PublicKey `json:"follower"`
	Followee         PublicKey `json:"followee"`
	FriendlyFollowee string    `json:"friendlyFollowee"`
	FriendlyFollower string    `json:"friendlyFollower"`
}

func NewFollowChange(changeType string, follower PublicKey, friendlyFollower string, followee PublicKey, friendlyFollowee string, at uint) FollowChange {
	return FollowChange{
		ChangeType:       changeType,
		Follower:         follower,
		FriendlyFollower: friendlyFollower,
		Followee:         followee,
		FriendlyFollowee: friendlyFollowee,
		At:               at,
	}
}
