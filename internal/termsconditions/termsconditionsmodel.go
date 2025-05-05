package termsconditions

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TermsConditions struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Subtitle    string             `bson:"subtitle" json:"subtitle"`
}

type TermsConditionsCreateDto struct {
	Title    string `json:"title" validate:"required"`
	Subtitle string `json:"subtitle" validate:"required"`
}

type TermsConditionsUpdateDto struct {
	Title    string `json:"title" validate:"required"`
	Subtitle string `json:"subtitle" validate:"required"`
}

type TermsConditionsDeleteManyDto struct {
	Ids []string `json:"ids" validate:"required"`
}

