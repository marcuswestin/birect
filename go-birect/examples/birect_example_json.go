package test_birect

import (
	"github.com/marcuswestin/birect/go-birect"
	"github.com/marcuswestin/go-errs"
)

// Echo handler
///////////////

type SendMessageParams struct {
	ToUserId tp.UserId
	Text string
}
type SendMessageResponse struct {
	Timestamp uint32
}
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

// Auth handler
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


// User shorthands
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
		res, err := handler(JSONReq(req))
		return birect.JSONRes(res), err
	})
})



// ctx/userContexts.go
//////////////////////

type UserContext interface {
	IsAuthenticatedUser() bool
	UserId() tp.UserId
}

// tp/dataTypes.go
//////////////////

type UserId uint64
