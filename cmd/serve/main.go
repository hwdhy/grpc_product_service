package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"grpc_tools/common"
	"grpc_tools/etcd"
	interceptorTool "grpc_tools/interceptor"
	"grpc_tools/pb/category_pb"
	"grpc_tools/pb/product_pb"
	"net"
	"net/http"
	"product_service/db"
	"product_service/service"
	"strconv"
)

var (
	grpcPort          = flag.Int("grpcPort", 50061, "the grpc server port")
	restPort          = flag.Int("restPort", 8090, "the rest server port")
	host              = flag.String("host", "127.0.0.1", "the server host")
	serviceName       = "productService"
	etcdExpire  int64 = 5
)

func main() {
	flag.Parse()

	grpcListen, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcPort))
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}
	restListen, err := net.Listen("tcp", fmt.Sprintf(":%d", *restPort))
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}

	go func() {
		err := runGRPCServer(grpcListen)
		if err != nil {
			logrus.Fatal(err)
		}
	}()
	logrus.Fatal(runRESTServer(restListen))
}

// 启动grpc服务
func runGRPCServer(listen net.Listener) error {
	db.InitConnectionPgsql() // 初始化数据库连接
	//更新接口权限
	e := common.InitAdapter([]map[string]int{
		service.ProductPermission,
		service.CategoryPermission,
	})

	interceptor := interceptorTool.NewAuthInterceptor()
	serverOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptor.Unary(e)),
	}

	server := grpc.NewServer(serverOptions...)
	product_pb.RegisterProductServer(server, &service.Product{})
	category_pb.RegisterCategoryServer(server, &service.Category{})

	etcdRegister, err := etcd.NewEtcdRegister()
	if err != nil {
		logrus.Fatal(err)
	}
	defer etcdRegister.Close()
	addr := net.JoinHostPort(*host, strconv.Itoa(*grpcPort))

	err = etcdRegister.RegisterServer("/etcd/"+serviceName, addr, etcdExpire)
	if err != nil {
		logrus.Fatalf("register error %v ", err)
	}

	logrus.Printf("server listening at %v", listen.Addr())
	return server.Serve(listen)
}

// 启动rest服务
func runRESTServer(listen net.Listener) error {
	conn := etcd.ClientConn(serviceName, 0, "")
	if conn == nil {
		logrus.Fatalf("get grpc client err")
	}

	mux := runtime.NewServeMux()
	err := product_pb.RegisterProductHandler(context.Background(), mux, conn)
	if err != nil {
		logrus.Fatalf("register product handler err: %v", err)
	}
	err = category_pb.RegisterCategoryHandler(context.Background(), mux, conn)
	if err != nil {
		logrus.Fatalf("register category handler err: %v", err)
	}

	logrus.Printf("start REST server at %s", listen.Addr())
	return http.Serve(listen, mux)
}
