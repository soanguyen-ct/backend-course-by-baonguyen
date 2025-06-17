package mongostore

import (
	"ct-backend-course-baonguyen/internal/entity"
	"ct-backend-course-baonguyen/pkg/hashpass"
	"github.com/google/uuid"
	"time"
)

var (
	_ IDocument = &ImageDoc{}
	_ IDocument = &UserDoc{}
)

type IDocument interface {
	GetDocId() string
}

func NewDoc() Doc {
	docId := uuid.New().String()
	return Doc{
		DocId:     docId,
		Version:   1,
		CreatedDt: time.Now(),
		UpdatedDt: time.Now(),
	}
}

type Doc struct {
	DocId     string    `bson:"docId"`
	Version   int64     `bson:"version"`
	CreatedDt time.Time `bson:"createdDt"`
	UpdatedDt time.Time `bson:"createdDt"`
}

func (d *Doc) GetDocId() string {
	return d.DocId
}

type ImageDoc struct {
	Doc  `bson:",inline"`
	User string `bson:"user"`
	Name string `bson:"name"`
	Path string `bson:"path"`
}

type UserDoc struct {
	Doc      `bson:",inline"`
	Username string `json:"username"`
	HashPass string `json:"hashPass"`
	FullName string `json:"fullName"`
	Address  string `json:"address"`
}

func NewImageDocument(info entity.ImageInfo) *ImageDoc {
	return &ImageDoc{
		Doc:  NewDoc(),
		User: info.UserName,
		Name: info.FileName,
		Path: info.ImagePath,
	}
}

func NewUserDocument(info entity.UserInfo) *UserDoc {
	return &UserDoc{
		Doc:      NewDoc(),
		Username: info.Username,
		HashPass: hashpass.HashPassword(info.Password),
		FullName: info.FullName,
		Address:  info.Address,
	}
}
