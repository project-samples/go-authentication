package myprofile

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type User struct {
	UserId       string        `json:"userId,omitempty" gorm:"column:userid;primary_key" bson:"_id,omitempty" dynamodbav:"userId,omitempty" firestore:"userId,omitempty" validate:"required,max=40"`
	Username     string        `json:"username,omitempty" gorm:"column:username" bson:"username,omitempty" dynamodbav:"username,omitempty" firestore:"username,omitempty" validate:"required,username,max=100"`
	Email        string        `json:"email,omitempty" gorm:"column:email" bson:"email,omitempty" dynamodbav:"email,omitempty" firestore:"email,omitempty" validate:"email,max=100"`
	Phone        string        `json:"phone,omitempty" gorm:"column:phone" bson:"phone,omitempty" dynamodbav:"phone,omitempty" firestore:"phone,omitempty" validate:"required,phone,max=18"`
	FirstName    string        `json:"firstName,omitempty" gorm:"column:firstName" bson:"firstName,omitempty" dynamodbav:"firstName,omitempty" firestore:"firstName,omitempty" `
	LastName     string        `json:"lastName,omitempty" gorm:"column:lastName" bson:"lastName,omitempty" dynamodbav:"lastName,omitempty" firestore:"lastName,omitempty" `
	ImageURL     string        `json:"imageURL,omitempty" gorm:"column:imageurl" bson:"imageURL,omitempty" dynamodbav:"imageURL,omitempty" firestore:"imageURL,omitempty" `
	Occupation   string        `json:"occupation,omitempty" gorm:"column:occupation" bson:"occupation,omitempty" dynamodbav:"occupation,omitempty" firestore:"occupation,omitempty" `
	Company      string        `json:"company,omitempty" gorm:"column:company" bson:"company,omitempty" dynamodbav:"company,omitempty" firestore:"company,omitempty" `
	DateOfBirth  *time.Time    `json:"dateOfBirth,omitempty" gorm:"column:date_of_birth" bson:"dateOfBirth,omitempty" dynamodbav:"dateOfBirth,omitempty" firestore:"dateOfBirth,omitempty"`
	Interests    []string      `json:"interests,omitempty" gorm:"column:interests" bson:"interests,omitempty" dynamodbav:"interests,omitempty" firestore:"interests,omitempty"`
	LookingFor   []string      `json:"lookingFor,omitempty" gorm:"column:lookingfor" bson:"lookingFor,omitempty" dynamodbav:"lookingFor,omitempty" firestore:"lookingFor,omitempty"`
	Skills       []Skills      `json:"skills,omitempty" gorm:"column:skills" bson:"skills,omitempty" dynamodbav:"skills,omitempty" firestore:"skills,omitempty"`
	Achievements []Achievement `json:"achievements,omitempty" gorm:"column:achievements" bson:"achievements,omitempty" dynamodbav:"achievements,omitempty" firestore:"achievements,omitempty"`
	Settings     *Settings     `json:"settings,omitempty" gorm:"column:settings" bson:"settings,omitempty" dynamodbav:"settings,omitempty" firestore:"settings,omitempty"`
	Gallery      []UploadInfo  `json:"gallery,omitempty" gorm:"column:gallery"`
	CoverURL     string        `json:"coverURL,omitempty" gorm:"column:coverurl" bson:"coverURL,omitempty" dynamodbav:"coverURL,omitempty" firestore:"coverURL,omitempty" `
}
type Skills struct {
	Skill   string `json:"skill,omitempty" gorm:"column:skill" bson:"skill,omitempty" dynamodbav:"skill,omitempty" firestore:"skill,omitempty" validate:"required"`
	Hirable bool   `json:"hirable,omitempty" gorm:"column:hirable" bson:"hirable,omitempty" dynamodbav:"hirable,omitempty" firestore:"hirable,omitempty"`
}

type UploadInfo struct {
	Source string `json:"source,omitempty" gorm:"column:source" bson:"source,omitempty" dynamodbav:"source,omitempty" firestore:"source,omitempty"`
	Url    string `json:"url,omitempty" gorm:"column:url" bson:"url,omitempty" dynamodbav:"url,omitempty" firestore:"url,omitempty"`
	Type   string `json:"category,omitempty" gorm:"column:category" bson:"category,omitempty" dynamodbav:"category,omitempty" firestore:"category,omitempty"`
}

func (c Skills) Value() (driver.Value, error) {
	return json.Marshal(c)
}
func (c *Skills) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &c)
}

type Achievement struct {
	Subject     string `json:"subject,omitempty" gorm:"column:subject" bson:"subject,omitempty" dynamodbav:"subject,omitempty" firestore:"subject,omitempty" validate:"required"`
	Description string `json:"description,omitempty" gorm:"column:description" bson:"description,omitempty" dynamodbav:"description,omitempty" firestore:"description,omitempty" validate:"required"`
}

func (c Achievement) Value() (driver.Value, error) {
	return json.Marshal(c)
}
func (c *Achievement) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &c)
}

type Settings struct {
	Language                        string `json:"language,omitempty" gorm:"column:language" bson:"language,omitempty" dynamodbav:"language,omitempty" firestore:"language,omitempty"`
	DateFormat                      string `json:"dateFormat,omitempty" gorm:"column:date_format" bson:"dateFormat,omitempty" dynamodbav:"dateFormat,omitempty" firestore:"dateFormat,omitempty"`
	DateTimeFormat                  string `json:"dateTimeFormat,omitempty" gorm:"column:date_time_format" bson:"dateTimeFormat,omitempty" dynamodbav:"dateTimeFormat,omitempty" firestore:"dateTimeFormat,omitempty"`
	TimeFormat                      string `json:"timeFormat,omitempty" gorm:"column:time_format" bson:"timeFormat,omitempty" dynamodbav:"timeFormat,omitempty" firestore:"timeFormat,omitempty"`
	Notification                    bool   `json:"notification" gorm:"column:notification" bson:"notification" dynamodbav:"notification" firestore:"notification"`
	EmailCommentsOfYourPosts        bool   `json:"emailCommentsOfYourPosts" gorm:"column:emailCommentsOfYourPosts" bson:"emailCommentsOfYourPosts" dynamodbav:"emailCommentsOfYourPosts" firestore:"emailCommentsOfYourPosts"`
	EmailEventInvitations           bool   `json:"emailEventInvitations" gorm:"column:emailEventInvitations" bson:"emailEventInvitations" dynamodbav:"emailEventInvitations" firestore:"emailEventInvitations"`
	EmailFeedUpdates                bool   `json:"emailFeedUpdates" gorm:"column:emailFeedUpdates" bson:"emailFeedUpdates" dynamodbav:"emailFeedUpdates" firestore:"emailFeedUpdates"`
	EmailPostMentions               bool   `json:"emailPostMentions" gorm:"column:emailPostMentions" bson:"emailPostMentions" dynamodbav:"emailPostMentions" firestore:"emailPostMentions"`
	EmailWhenNewEventsAround        bool   `json:"emailWhenNewEventsAround" gorm:"column:emailWhenNewEventsAround" bson:"emailWhenNewEventsAround" dynamodbav:"emailWhenNewEventsAround" firestore:"emailWhenNewEventsAround"`
	FollowingListPublicOnMyProfile  bool   `json:"followingListPublicOnMyProfile" gorm:"column:followingListPublicOnMyProfile" bson:"followingListPublicOnMyProfile" dynamodbav:"followingListPublicOnMyProfile" firestore:"followingListPublicOnMyProfile"`
	SearchEnginesLinksToMyProfile   bool   `json:"searchEnginesLinksToMyProfile" gorm:"column:searchEnginesLinksToMyProfile" bson:"searchEnginesLinksToMyProfile" dynamodbav:"searchEnginesLinksToMyProfile" firestore:"searchEnginesLinksToMyProfile"`
	NotifyFeedUpdates               bool   `json:"notifyFeedUpdates" gorm:"column:notifyFeedUpdates" bson:"notifyFeedUpdates" dynamodbav:"notifyFeedUpdates" firestore:"notifyFeedUpdates"`
	NotifyPostMentions              bool   `json:"notifyPostMentions" gorm:"column:notifyPostMentions" bson:"notifyPostMentions" dynamodbav:"notifyPostMentions" firestore:"notifyPostMentions"`
	NotifyCommentsOfYourPosts       bool   `json:"notifyCommentsOfYourPosts" gorm:"column:notifyCommentsOfYourPosts" bson:"notifyCommentsOfYourPosts" dynamodbav:"notifyCommentsOfYourPosts" firestore:"notifyCommentsOfYourPosts"`
	NotifyEventInvitations          bool   `json:"notifyEventInvitations" gorm:"column:notifyEventInvitations" bson:"notifyEventInvitations" dynamodbav:"notifyEventInvitations" firestore:"notifyEventInvitations"`
	NotifyWhenNewEventsAround       bool   `json:"notifyWhenNewEventsAround" gorm:"column:notifyWhenNewEventsAround" bson:"notifyWhenNewEventsAround" dynamodbav:"notifyWhenNewEventsAround" firestore:"notifyWhenNewEventsAround"`
	ShowMyProfileInSpacesAroundMe   bool   `json:"showMyProfileInSpacesAroundMe" gorm:"column:showMyProfileInSpacesAroundMe" bson:"showMyProfileInSpacesAroundMe" dynamodbav:"showMyProfileInSpacesAroundMe" firestore:"showMyProfileInSpacesAroundMe"`
	ShowAroundMeResultsInMemberFeed bool   `json:"showAroundMeResultsInMemberFeed" gorm:"column:showAroundMeResultsInMemberFeed" bson:"showAroundMeResultsInMemberFeed" dynamodbav:"showAroundMeResultsInMemberFeed" firestore:"showAroundMeResultsInMemberFeed"`
}

func (c Settings) Value() (driver.Value, error) {
	return json.Marshal(c)
}
func (c *Settings) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &c)
}
