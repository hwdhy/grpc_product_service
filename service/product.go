package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc_tools/common"
	"grpc_tools/pb/product_pb"
	"product_service/db"
	"product_service/models"
)

var ProductPermission = map[string]int{
	"/grpc_hwdhy.Product/Add": common.NotLogged,
}

type Product struct {
	product_pb.UnimplementedProductServer
}

func (p *Product) Add(ctx context.Context, input *product_pb.ProductAddRequest) (*product_pb.ProductAddResponse, error) {

	productData := models.Product{
		Name:  input.Name,
		Price: input.Price,
		Image: input.Image,
	}
	if err := db.PgsqlDB.Model(models.Product{}).Create(&productData).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "create product err: %v", err)
	}
	logrus.Printf("create product(%+v) success", productData)

	return &product_pb.ProductAddResponse{
		Status: "success",
	}, nil
}
