package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"grpc_tools/etcd"
	"grpc_tools/pb/product_pb"
	"strconv"
	"time"
)

func main() {
	conn := etcd.ClientConn("productService", 0, "")
	if conn == nil {
		logrus.Fatalf("get grpc client err")
	}
	defer conn.Close()

	// 获取客户端对象
	c := product_pb.NewProductClient(conn)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	for i := 0; i < 10; i++ {
		response, err := c.Add(ctx, &product_pb.ProductAddRequest{
			Name:  "product" + strconv.Itoa(i),
			Price: uint64(i),
			Image: "image url",
		})
		if err != nil {
			logrus.Fatalf("add product err: %v", err)
		}
		logrus.Printf("add product success, count : %d, status: %v", i, response.GetStatus())
	}
}
