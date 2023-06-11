package main

import (
	"context"
	"economic_srvs/user_srv/proto"
	"fmt"
	"google.golang.org/grpc"
)

var userClient proto.UserClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:8088", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	userClient = proto.NewUserClient(conn)
}

func TestGetUserList() {
	list, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    1,
		PSize: 2,
	})
	if err != nil {
		panic(err)
	}
	for _, user := range list.Data {
		fmt.Println(user.NickName, user.Password, user.Mobile)
		checked, err := userClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
			Password:          "yinlei19980505",
			EncryptedPassword: user.Password,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(checked.Success)
	}
}

func TestCreateUser() {
	resp, err := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		NickName: "张春秀",
		Password: "123456789",
		Mobile:   "13350011809",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Id)
}

func main() {
	Init()
	TestCreateUser()
	//TestGetUserList()
	conn.Close()
}
