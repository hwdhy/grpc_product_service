package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc_tools/common"
	"grpc_tools/pb/category_pb"
	"product_service/db"
	"product_service/models"
)

var CategoryPermission = map[string]int{
	"/grpc_hwdhy.Category/List":   common.NotLogged,
	"/grpc_hwdhy.Category/Create": common.NotLogged,
}

type Category struct {
	category_pb.UnimplementedCategoryServer
}

func (c *Category) List(ctx context.Context, in *category_pb.CategoryListRequest) (*category_pb.CategoryListResponse, error) {
	Offset := (in.Page - 1) * in.PageSize

	var categorys []models.Category
	if err := db.PgsqlDB.Model(models.Category{}).Offset(int(Offset)).Limit(int(in.PageSize)).Find(&categorys).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "select category err: %v", err)
	}

	dateCount := len(categorys)

	res := make([]*category_pb.CategoryData, dateCount)
	for i, category := range categorys {
		res[i] = &category_pb.CategoryData{
			Id:    int64(category.ID),
			Name:  category.Name,
			Pid:   category.Pid,
			Image: category.Image,
		}
	}

	return &category_pb.CategoryListResponse{
		Page:         in.Page,
		PageSize:     in.PageSize,
		Count:        int64(dateCount),
		CategoryData: res,
	}, nil
}

func (c *Category) Create(ctx context.Context, in *category_pb.CategoryCreateRequest) (*category_pb.Response, error) {
	// 判断分类名称是否存在
	var findCategory models.Category
	db.PgsqlDB.Model(models.Category{}).Where("name = ?", in.GetName()).First(&findCategory)
	if findCategory.ID != 0 {
		return nil, status.Errorf(codes.AlreadyExists, "category(%s) is already exists", in.GetName())
	}

	categoryData := models.Category{
		Name:  in.GetName(),
		Pid:   in.GetPid(),
		Image: in.GetImage(),
	}

	if err := db.PgsqlDB.Model(models.Category{}).Create(&categoryData).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "create category err: %v", err)
	}
	logrus.Printf("careate category(%s) success", categoryData.Name)
	return &category_pb.Response{Code: int64(codes.OK)}, nil
}
