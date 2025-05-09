package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StatusValue string

type status struct {
	Draft      StatusValue
	Processing StatusValue
	Complete   StatusValue
}

var AnalyticsJobStatus = status{
	Draft:      "DRAFT",
	Processing: "PROCESSING",
	Complete:   "COMPLETE",
}

type AnalyticsJob struct {
	ID              *primitive.ObjectID `bson:"_id,omitempty"`
	CreatedTime     *time.Time          `bson:"created_time,omitempty"`
	LastUpdatedTime *time.Time          `bson:"last_updated_time,omitempty"`

	Latitude  *float64    `bson:"latitude,omitempty"`
	Longitude *float64    `bson:"longitude,omitempty"`
	FileUrl   string      `bson:"file_url,omitempty"`
	VideoUrl  string      `bson:"video_url,omitempty"`
	Status    StatusValue `bson:"status,omitempty"`
	Result    string      `bson:"result,omitempty"`
}
