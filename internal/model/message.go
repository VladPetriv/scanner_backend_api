package model

type FullMessage struct {
	ID              int    `json:"id" db:"messageid"`
	Title           string `json:"title" db:"messagetitle"`
	MessageURL      string `json:"messageUrl" db:"messageurl"`
	MessageImageURL string `json:"messageImageUrl" db:"messageimageurl"`

	ChannelName     string `json:"channelName" db:"channelname"`
	ChannelTitle    string `json:"channelTitle" db:"channeltitle"`
	ChannelImageURL string `json:"channelImageUrl" db:"channelimageurl"`

	UserID       int    `json:"userId" db:"userid"`
	UserFullname string `json:"userFullname" db:"userfullname"`
	UserImageURL string `json:"userImageUrl" db:"userimageurl"`

	RepliesCount int          `json:"repliesCount" db:"count"`
	Replies      []FullReplie `json:"replies"`
}
