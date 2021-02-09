package service

import (
	"context"
	"log"
	"microservices/authentication/models"
	"microservices/authentication/repository"
	"microservices/authentication/validators"
	"microservices/pb"
	"microservices/security"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type authService struct {
	usersRepository repository.UsersRepository
}

func NewAuthService(usersRepository repository.UsersRepository) pb.AuthServiceServer {
	return &authService{usersRepository: usersRepository}
}

func (s *authService) SignUp(ctx context.Context, req *pb.User) (*pb.User, error) {
	err := validators.ValidateSignUp(req)
	if err != nil {
		return nil, err
	}

	req.Password, err = security.EncryptPassword(req.Password)
	if err != nil {
		return nil, err
	}
	req.Name = strings.TrimSpace(req.Name)
	req.Email = validators.NormalizeEmail(req.Email)

	found, err := s.usersRepository.GetByEmail(req.Email)
	if err == mgo.ErrNotFound {
		user := new(models.User)
		user.FromProtoBuffer(req)
		err := s.usersRepository.Save(user)
		if err != nil {
			return nil, err
		}
		return user.ToProtoBuffer(), nil
	}

	if found == nil {
		return nil, err
	}

	return nil, validators.ErrEmailAlreadyExists
}

func (s *authService) SignIn(ctx context.Context, req *pb.SignInRequest) (*pb.SignInResponse, error) {
	req.Email = validators.NormalizeEmail(req.Email)

	user, err := s.usersRepository.GetByEmail(req.Email)
	if err != nil {
		log.Println("signin failed:", err.Error())
		return nil, validators.ErrSignInFailed
	}

	err = security.VerifyPassword(user.Password, req.Password)
	if err != nil {
		log.Println("signin failed:", err.Error())
		return nil, validators.ErrSignInFailed
	}

	token, err := security.NewToken(user.Id.Hex())
	if err != nil {
		log.Println("signin failed:", err.Error())
		return nil, validators.ErrSignInFailed
	}

	return &pb.SignInResponse{User: user.ToProtoBuffer(), Token: token}, nil
}

func (s *authService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	if !bson.IsObjectIdHex(req.Id) {
		return nil, validators.ErrInvalidUserId
	}
	found, err := s.usersRepository.GetById(req.Id)
	if err != nil {
		return nil, err
	}
	return found.ToProtoBuffer(), nil
}

func (s *authService) ListUsers(req *pb.ListUsersRequest, stream pb.AuthService_ListUsersServer) error {
	users, err := s.usersRepository.GetAll()
	if err != nil {
		return err
	}
	for _, user := range users {
		err := stream.Send(user.ToProtoBuffer())
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *authService) UpdateUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	if !bson.IsObjectIdHex(req.Id) {
		return nil, validators.ErrInvalidUserId
	}
	user, err := s.usersRepository.GetById(req.Id)
	if err != nil {
		return nil, err
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return nil, validators.ErrEmptyName
	}
	if req.Name == user.Name {
		return user.ToProtoBuffer(), nil
	}

	user.Name = req.Name
	user.Updated = time.Now()
	err = s.usersRepository.Update(user)
	return user.ToProtoBuffer(), err
}

func (s *authService) DeleteUser(ctx context.Context, req *pb.GetUserRequest) (*pb.DeleteUserResponse, error) {
	if !bson.IsObjectIdHex(req.Id) {
		return nil, validators.ErrInvalidUserId
	}
	err := s.usersRepository.Delete(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteUserResponse{Id: req.Id}, nil
}
