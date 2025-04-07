package mongodate

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Now() primitive.DateTime {
	return primitive.NewDateTimeFromTime(time.Now())
}