package services

import (
	"context"
	"io"
	"time"
)

type UserService struct{}

// GetUserScore 普通方法一次性返回
func (*UserService) GetUserScore(ctx context.Context, in *UserScoreRequest) (*UserScoreResponse, error) {
	var score int32 = 101
	users := make([]*UserInfo, 0)
	for _, user := range in.Users {
		user.UserScore = score
		score++
		users = append(users, user)
	}
	return &UserScoreResponse{Users: users}, nil
}

// GetUserScoreByServerStream 服务端流
func (this *UserService) GetUserScoreByServerStream(in *UserScoreRequest, stream UserService_GetUserScoreByServerStreamServer) error {
	var score int32 = 101
	users := make([]*UserInfo, 0)
	for index, user := range in.Users {
		user.UserScore = score
		score++
		users = append(users, user)

		if (index+1)%2 == 0 && index > 0 {
			// 每隔2条发送
			err := stream.Send(&UserScoreResponse{Users: users})
			if err != nil {
				return err
			}
			users = (users)[0:0] // 发送完成之后清空切片，方便下次发送
		}
		time.Sleep(time.Second * 1) // 模拟这里处理比较耗时
	}
	if len(users) > 0 { // 发送剩余残留的数据
		err := stream.Send(&UserScoreResponse{Users: users})
		if err != nil {
			return err
		}
	}
	return nil
}

// GetUserScoreByClientStream 客户端流
func (this *UserService) GetUserScoreByClientStream(stream UserService_GetUserScoreByClientStreamServer) error {
	var score int32 = 101
	users := make([]*UserInfo, 0)
	for {
		req, err := stream.Recv()
		if err == io.EOF { // 接受结束
			return stream.SendAndClose(&UserScoreResponse{Users: users}) // 发送并关闭
		}
		if err != nil { // 接受出错
			return err
		}
		for _, user := range req.Users {
			user.UserScore = score // 服务端做的业务处理
			score++
			users = append(users, user) // 把处理的结果放入 users
		}
	}
}
