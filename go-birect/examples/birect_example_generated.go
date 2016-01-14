package test_birect

import (
	"github.com/marcuswestin/birect/go-birect"
	"github.com/marcuswestin/go-errs"
)

// AUTO-GEN
///////////

func Listen() {
	// reflect all fields
	// for field in range fields
	//		if handlers[field] == nil
	//			panic()
}

type SendMessageReq struct {
	ToUserId uint32 `protobuf:"varint,1,opt,name=req_id" json:"req_id,omitempty"`
	Text  string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
}
type SendMessageRes struct {
	Timestamp uint32 `protobuf:"varint,1,opt,name=req_id" json:"req_id,omitempty"`
}
type userAPIHandlers struct {
	SendMessage func(req SendMessageReq) (res SendMessageRes)
}
var handlers userAPIHandlers
func init() {
	handlers = userAPIHandlers{}
}

// MANUALAs
///////////

func init() {
	handle("userAPI/SendMessage", func(req ProtoReq) (res ProtoRes, err error) {
		var par SendMessageParams
		req.ParseParams(&par)
		
		timestamp, err := sendMessage(req, par.ToUserId, par.Text)
		if err != nil {
			return
		}
		
		info, reader, err := req.NextAttachment()
		if err != nil {
			return
		}
		
		res = SendMessageRes{timestamp}
		return
	})
}
// or
func init() {
	userAPI.SendMessage = func(req SendMessageReq, par SendMessageParams) (res SendMessageRes, err error) {
		timestamp, err := sendMessage(req, par.ToUserId, par.Text)
		if err != nil {
			return
		}
		
		info, reader, err := req.NextAttachment()
		if err != nil {
			return
		}
		
		
		res.Timestamp = timestamp
		return
	}

	userAPIHandlers.SendMessage = func(req SendMessageReq) (res SendMessageRes) {
		
	}
}

type authAPI struct {
	AuthenticateWithPassword func(req AuthReq) (res AuthRes)
}
//....

// Autoge


// Auth package
///////////////
type AuthWithPasswordParams struct {
	Username string
	Password string
}
type AuthWithPasswordResponse struct {
	Secret string
}
func init() {
	handleJSONReq("auth/AuthWithPassword", func(req JSONReq) (res JSONRes, err error) {
		var par AuthWithPasswordParams
		req.ParseParams(&par)
		
		secret, err := authenticateUser(req, par.Username, par.Password)
		if err != nil {
			return
		}
		
		res = AuthWithPasswordResponse{secret}
		return
	})
}


// userAPI package
//////////////////

type UserReq struct {
	birect.JSONReq
}
func (u *UserReq) IsAuthenticatedUser() bool {
	_, found := connUserIdMap[u.JSONReq.Conn]
	return found
}

func (u *UserRep) AssertIsAuthenticatedUser() {
	if !u.IsAuthenticatedUser {
		panic(u.JSONReq.Error("Not authenticated user"))
	}
}

// UserId
func (u *UserReq) UserId() tp.UserId {
	u.AssertIsAuthenticatedUser()
	return connUserIdMap[u.JSONReq.Conn]
}

var connUserIdMap = make(map[birect.Conn]tp.UserId, 10000)

type JSONReqHandler birect.JSONReqHandler
type JSONReq birect.JSONReq
type JSONRes birect.JSONRes
func handleJSONReq(reqName string, handler JSONHandler) {
	if !strings.HasPrefix(reqName, "users/") {
		panic("Bad user handler request reqName")
	}
	birect.HandleJSONReq(reqName, func(req birect.JSONReq) (birect.JSONRes, error) {
		jsonReq := JSONReq{req}
		jsonReq.AssertIsAuthenticatedUser()
		res, err := handler(jsonReq)
		return birect.JSONRes(res), err
	})
})


// userAPI/SendMessage.go

func init() {
	handleJSONReq("user/SendMessage", func(req JSONReq) (res JSONRes, err error) {
		var par SendMessageResponse
		req.ParseParams(&par)
		
		timestamp, err := sendMessage(req, req.UserId(), params.ToUserId, params.Text)
		if err != nil {
			return
		}
		
		res = SendMessageResponse{timestamp}
		return
	})
}


// ctx/userContexts.go
//////////////////////

type UserContext interface {
	IsAuthenticatedUser() bool
	UserId() tp.UserId
}

// tp/dataTypes.go
//////////////////

type UserId uint64


