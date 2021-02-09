package resthandlers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"microservices/api/restutil"
	"microservices/pb"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type AuthHandlers interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	SignIn(w http.ResponseWriter, r *http.Request)
	PutUser(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	GetUsers(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}

type authHandlers struct {
	authSvcClient pb.AuthServiceClient
}

func NewAuthHandlers(authSvcClient pb.AuthServiceClient) AuthHandlers {
	return &authHandlers{authSvcClient: authSvcClient}
}

func (h *authHandlers) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		restutil.WriteError(w, http.StatusBadRequest, restutil.ErrEmptyBody)
		return
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		restutil.WriteError(w, http.StatusBadRequest, err)
		return
	}
	user := new(pb.User)
	err = json.Unmarshal(body, user)
	if err != nil {
		restutil.WriteError(w, http.StatusBadRequest, err)
		return
	}
	user.Created = time.Now().Unix()
	user.Updated = user.Created
	user.Id = bson.NewObjectId().Hex()
	resp, err := h.authSvcClient.SignUp(r.Context(), user)
	if err != nil {
		restutil.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}
	restutil.WriteAsJson(w, http.StatusCreated, resp)
}

func (h *authHandlers) SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		restutil.WriteError(w, http.StatusBadRequest, restutil.ErrEmptyBody)
		return
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		restutil.WriteError(w, http.StatusBadRequest, err)
		return
	}
	input := new(pb.SignInRequest)
	err = json.Unmarshal(body, input)
	if err != nil {
		restutil.WriteError(w, http.StatusBadRequest, err)
		return
	}
	resp, err := h.authSvcClient.SignIn(r.Context(), input)
	if err != nil {
		restutil.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}
	restutil.WriteAsJson(w, http.StatusOK, resp)
}

func (h *authHandlers) PutUser(w http.ResponseWriter, r *http.Request) {
	tokenPayload, err := restutil.AuthRequestWithId(r)
	if err != nil {
		restutil.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if r.Body == nil {
		restutil.WriteError(w, http.StatusBadRequest, restutil.ErrEmptyBody)
		return
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		restutil.WriteError(w, http.StatusBadRequest, err)
		return
	}
	user := new(pb.User)
	err = json.Unmarshal(body, user)
	if err != nil {
		restutil.WriteError(w, http.StatusBadRequest, err)
		return
	}
	user.Id = tokenPayload.UserId
	resp, err := h.authSvcClient.UpdateUser(r.Context(), user)
	if err != nil {
		restutil.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}
	restutil.WriteAsJson(w, http.StatusOK, resp)
}

func (h *authHandlers) GetUser(w http.ResponseWriter, r *http.Request) {
	tokenPayload, err := restutil.AuthRequestWithId(r)
	if err != nil {
		restutil.WriteError(w, http.StatusBadRequest, err)
		return
	}
	resp, err := h.authSvcClient.GetUser(r.Context(), &pb.GetUserRequest{Id: tokenPayload.UserId})
	if err != nil {
		restutil.WriteError(w, http.StatusBadRequest, err)
		return
	}
	restutil.WriteAsJson(w, http.StatusOK, resp)
}

func (h *authHandlers) GetUsers(w http.ResponseWriter, r *http.Request) {
	stream, err := h.authSvcClient.ListUsers(r.Context(), &pb.ListUsersRequest{})
	if err != nil {
		restutil.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}
	var users []*pb.User
	for {
		user, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			restutil.WriteError(w, http.StatusBadRequest, err)
			return
		}
		users = append(users, user)
	}
	restutil.WriteAsJson(w, http.StatusOK, users)
}

func (h *authHandlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	tokenPayload, err := restutil.AuthRequestWithId(r)
	if err != nil {
		restutil.WriteError(w, http.StatusBadRequest, err)
		return
	}
	resp, err := h.authSvcClient.DeleteUser(r.Context(), &pb.GetUserRequest{Id: tokenPayload.UserId})
	if err != nil {
		restutil.WriteError(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", resp.Id)
	restutil.WriteAsJson(w, http.StatusNoContent, nil)
}
