package handler

import (
	context "context"
	"crypto/sha512"
	"economic_srvs/user_srv/global"
	"economic_srvs/user_srv/model"
	"economic_srvs/user_srv/proto"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"strings"
	"time"
)

type UserServer struct {
	proto.UnimplementedUserServer
}

func (u *UserServer) mustEmbedUnimplementedUserServer() {
	panic("implement me")
}

func ModelToResponse(user model.User) proto.UserInfoResponse {
	userInfoResponse := proto.UserInfoResponse{
		Id:       int32(user.ID),
		Password: user.Password,
		Mobile:   user.Mobile,
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     int32(user.Role),
	}
	if user.Birthday != nil {
		userInfoResponse.Birthday = uint64(user.Birthday.Unix())
	}
	return userInfoResponse
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func (u *UserServer) GetUserList(ctx context.Context, request *proto.PageInfo) (*proto.UserListResponse, error) {
	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	response := &proto.UserListResponse{}
	response.Total = int32(result.RowsAffected)

	global.DB.Scopes(Paginate(int(request.Pn), int(request.PSize))).Find(&users)
	for _, user := range users {
		userInfoResponse := ModelToResponse(user)
		response.Data = append(response.Data, &userInfoResponse)
	}
	return response, nil
}

func (u *UserServer) GetUserByMobile(ctx context.Context, request *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: request.Mobile}).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	userInfoResponse := ModelToResponse(user)
	return &userInfoResponse, nil
}

func (u *UserServer) GetUserById(ctx context.Context, request *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.First(&user, request.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	userInfoResponse := ModelToResponse(user)
	return &userInfoResponse, nil
}

func (u *UserServer) CreateUser(ctx context.Context, request *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: request.Mobile}).First(&user)
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已经存在")
	}
	user.Mobile = request.Mobile
	user.NickName = request.NickName

	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	salt, encodedPwd := password.Encode(request.Password, options)
	passwordSaltedStr := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	user.Password = passwordSaltedStr
	result = global.DB.Create(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	userInfoResponse := ModelToResponse(user)
	return &userInfoResponse, nil
}

func (u *UserServer) UpdateUser(ctx context.Context, request *proto.UpdateUserInfo) (*emptypb.Empty, error) {
	var user model.User
	result := global.DB.First(&user, request.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	birthday := time.Unix(int64(request.Birthday), 0)
	user.NickName = request.NickName
	user.Birthday = &birthday
	user.Gender = request.Gender

	result = global.DB.Save(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return &emptypb.Empty{}, nil
}

func (u *UserServer) CheckPassWord(ctx context.Context, request *proto.PasswordCheckInfo) (*proto.CheckResponse, error) {
	passwordRawStr := strings.Split(request.EncryptedPassword, "$")
	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	check := password.Verify(request.Password, passwordRawStr[2], passwordRawStr[3], options)
	return &proto.CheckResponse{
		Success: check,
	}, nil
}
