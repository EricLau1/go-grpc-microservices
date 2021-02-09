package models

import (
	"microservices/pb"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id       bson.ObjectId `bson:"_id"`
	Name     string        `bson:"name"`
	Email    string        `bson:"email"`
	Password string        `bson:"password"`
	Created  time.Time     `bson:"created"`
	Updated  time.Time     `bson:"updated"`
}

func (u *User) ToProtoBuffer() *pb.User {
	return &pb.User{
		Id:       u.Id.Hex(),
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
		Created:  u.Created.Unix(),
		Updated:  u.Updated.Unix(),
	}
}

func (u *User) FromProtoBuffer(user *pb.User) {
	u.Id = bson.ObjectIdHex(user.GetId())
	u.Name = user.GetName()
	u.Email = user.GetEmail()
	u.Password = user.GetPassword()
	u.Created = time.Unix(user.Created, 0)
	u.Updated = time.Unix(user.Updated, 0)
}
