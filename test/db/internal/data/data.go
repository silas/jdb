package data

import (
	"time"
)

type Row struct {
	Kind            string
	ID              string
	ParentKind      string
	ParentID        string
	UniqueStringKey string
	StringKey       string
	NumericKey      float64
	TimeKey         time.Time
	Data            string
	CreateTime      time.Time
	UpdateTime      time.Time
}

var (
	UserKind = "user"

	UserDomain = "example.com"
	createTime = time.Date(2001, 2, 3, 4, 5, 6, 0, time.UTC)
	days       = 24 * time.Hour
	now        = time.Date(2018, 5, 1, 4, 2, 45, 0, time.UTC)

	User1ID          = "1"
	User1Email       = "jane@" + UserDomain
	User1GivenName   = "Jane"
	User1FamilyName  = "Doe"
	User1Aliases     = []string{"Janie", "Roe"}
	User1Age         = 34
	User1RefreshTime = time.Date(2018, 2, 1, 5, 23, 3, 0, time.UTC)
	User1Data        = `{"Email": "jane@example.com", "Name": {"GivenName": "Jane", "FamilyName": "Doe", "Aliases": ["Janie", "Roe"]}, "Age": 34}`
	User1CreateTime  = createTime
	User1UpdateTime  = time.Date(2002, 3, 4, 5, 6, 7, 0, time.UTC)

	User2ID          = "2"
	User2Email       = "john@" + UserDomain
	User2GivenName   = "John"
	User2FamilyName  = "Smith"
	User2Aliases     = []string{"Richard", "Johnny", "Roe"}
	User2Age         = 23
	User2RefreshTime = time.Date(2018, 3, 2, 1, 4, 43, 0, time.UTC)
	User2Data        = `{"Email": "john@example.com", "Name": {"GivenName": "John", "FamilyName": "Smith", "Aliases": ["Richard", "Johnny", "Roe"]}, "Age": 23}`
	User2CreateTime  = User1CreateTime.Add(time.Millisecond)
	User2UpdateTime  = now

	User3ID         = "3"
	User3CreateTime = User2CreateTime.Add(30 * days)
	User3UpdateTime = now.Add(-24 * time.Hour)

	UserRows = []Row{
		{
			Kind:            UserKind,
			ID:              User1ID,
			UniqueStringKey: User1Email,
			StringKey:       UserDomain,
			NumericKey:      50,
			TimeKey:         User1RefreshTime,
			Data:            User1Data,
			CreateTime:      User1CreateTime,
			UpdateTime:      User1UpdateTime,
		},
		{
			Kind:            UserKind,
			ID:              User2ID,
			UniqueStringKey: User2Email,
			StringKey:       UserDomain,
			NumericKey:      10,
			TimeKey:         User2RefreshTime,
			Data:            User2Data,
			CreateTime:      User2CreateTime,
			UpdateTime:      User2UpdateTime,
		},
		{
			Kind:       UserKind,
			ID:         User3ID,
			NumericKey: 3,
			CreateTime: User3CreateTime,
			UpdateTime: User3UpdateTime,
		},
	}
)
