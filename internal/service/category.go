package service

import (
	"context"
	"example-grpc/internal/database"
	"example-grpc/internal/pb"
	"io"
)

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	CategoryDB database.Category
}

func NewCategoryService(categoryDB database.Category) *CategoryService {
	return &CategoryService{
		CategoryDB: categoryDB,
	}
}

func (c *CategoryService) GetCategory(ctx context.Context, in *pb.GetCategoryRequest) (*pb.Category, error) {
	category, err := c.CategoryDB.Find(in.Id)

	if err != nil {
		return nil, err
	}

	categoryModel := &pb.Category{
		Id: category.ID,
		Name: category.Name,
		Description: category.Description,
	}

	return categoryModel, nil
}

func (c *CategoryService) ListCategories(ctx context.Context, in *pb.Blank) (*pb.CategoryList, error) {
	categories, err := c.CategoryDB.FindAll()

	if err != nil {
		return nil, err
	}

	var categoriesModel []*pb.Category

	for _, category := range categories {
		categoryModel := &pb.Category{
			Id: category.ID,
			Name: category.Name,
			Description: category.Description,
		}

		categoriesModel = append(categoriesModel, categoryModel)
	}

	return &pb.CategoryList{Categories: categoriesModel}, nil
}

func (c *CategoryService) CreateCategoryStream(stream pb.CategoryService_CreateCategoryStreamServer) (error) {
	categoryList := &pb.CategoryList{}

	for {
		categoryRequest, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(categoryList)
		}

		if err != nil {
			return err
		}

		category, err := c.CategoryDB.Create(categoryRequest.Name, categoryRequest.Description)

		if err != nil {
			return err
		}

		categoryList.Categories = append(categoryList.Categories, &pb.Category{
			Id: category.ID,
			Name: category.Name,
			Description: category.Description,
		})
	}
}

func (c *CategoryService) CreateCategory(ctx context.Context, in *pb.CreateCategoryRequest) (*pb.CategoryResponse, error) {
	category, err := c.CategoryDB.Create(in.Name, in.Description)

	if err != nil {
		return nil, err
	}

	categoryResponse := &pb.CategoryResponse{
		Category: &pb.Category{
			Id: category.ID,
			Name: category.Name,
			Description: category.Description,
		},
	}

	return categoryResponse, nil
}
