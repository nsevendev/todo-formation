package task

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Label	 string             `bson:"label" json:"label"`
	Done	 bool               `bson:"done" json:"done"`
	CreatedAt primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt primitive.DateTime `bson:"updated_at" json:"updated_at"`

	IdUser primitive.ObjectID `bson:"id_user" json:"id_user"`
}

type TaskCreateDto struct {
	Label string `json:"label" validate:"required"`
}

type TaskUpdateLabelDto struct {
	Label string `json:"label" validate:"required"`
}

type TaskDeleteManyDto struct {
	Ids []string `json:"ids" validate:"required"`
}

func (u *Task) SetTimeStamps() {
	now := primitive.NewDateTimeFromTime(time.Now())
	if u.CreatedAt == 0 {
		u.CreatedAt = now
	}
	u.UpdatedAt = now
}