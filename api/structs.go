package api

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// UserLoged struct for resp logged user
type UserLoged struct {
	ID       bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Username string        `bson:"username" json:"username"`
	Avatar   string        `bson:"avatar" json:"avatar,omitempty"`
}

// User struct fro mongodb model
type User struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Username  string        `bson:"username" json:"username" valid:"required,alphanumunicode2,min=5,max=60"`
	Email     string        `bson:"email" json:"email,omitempty" valid:"required,email"`
	Password  string        `bson:"-" json:"password,omitempty" valid:"required,alphanumunicode2,min=5,max=60"`
	HPassword string        `bson:"hPassword" json:"-"`
	Avatar    string        `bson:"avatar" json:"avatar,omitempty"`
	IsActive  bool          `bson:"isActive" json:"isActive"`
	LastLog   *time.Time    `bson:"lastLog,omitempty" json:"lastLog,omitempty"`
	Dated     *time.Time    `bson:"dated,omitempty" json:"dated,omitempty"`
	Updated   *time.Time    `bson:"updated,omitempty" json:"updated,omitempty"`
}
