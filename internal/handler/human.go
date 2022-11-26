package handler

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/joisandresky/go-echo-mongodb-boilerplate/internal/entity"
	"github.com/joisandresky/go-echo-mongodb-boilerplate/internal/helper"
	"github.com/joisandresky/go-echo-mongodb-boilerplate/internal/repository"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type HumanHandler interface {
	Paginate(c echo.Context) error
	GetAllHuman(c echo.Context) error
	FindHumanById(c echo.Context) error
	CreateHuman(c echo.Context) error
	UpdateHumanById(c echo.Context) error
	DeleteHuman(c echo.Context) error
}

type humanHandler struct {
	humanRepo repository.HumanRepository
}

func NewHumanHandler(humanRepo repository.HumanRepository) HumanHandler {
	return &humanHandler{humanRepo}
}

func (h *humanHandler) Paginate(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pageParam := c.QueryParam("page")
	if pageParam == "" {
		pageParam = "1"
	}
	page, _ := strconv.Atoi(pageParam)

	limitParam := c.QueryParam("limit")
	if limitParam == "" {
		limitParam = "15"
	}
	limit, _ := strconv.Atoi(limitParam)
	filter := bson.M{
		"name": bson.M{"$regex": c.QueryParam("search"), "$options": "i"},
	}

	humans, pagination, err := h.humanRepo.Paginate(ctx, int64(page), int64(limit), filter)
	if err != nil {
		return helper.BadRequestResponse(c, helper.Response{
			Message: "Failed to Get Human List",
			Errors:  err.Error(),
		})
	}

	return helper.OkResponse(c, helper.Response{
		Data:       humans,
		Pagination: pagination,
	})
}

func (h *humanHandler) GetAllHuman(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	humans, err := h.humanRepo.FindAll(ctx)
	if err != nil {
		return helper.BadRequestResponse(c, helper.Response{
			Message: "Failed to Get Human List",
			Errors:  err.Error(),
		})
	}

	return helper.OkResponse(c, helper.Response{
		Data: humans,
	})
}

func (h *humanHandler) FindHumanById(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	paramId := c.Param("id")

	if paramId == "" {
		return helper.BadRequestResponse(c, helper.Response{
			Message: "Please Provide Valid Human ID!",
		})
	}

	human, err := h.humanRepo.FindById(ctx, paramId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return helper.NotFoundResponse(c, helper.Response{
				Message: fmt.Sprintf("Human With [ID: %s] Not Found", paramId),
			})
		}

		return helper.BadRequestResponse(c, helper.Response{
			Message: fmt.Sprintf("Failed to Get Human With [ID: %s]", paramId),
			Errors:  err.Error(),
		})
	}

	return helper.OkResponse(c, helper.Response{
		Data: human,
	})
}

func (h *humanHandler) CreateHuman(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var human entity.Human
	if err := c.Bind(&human); err != nil {
		return helper.BadRequestResponse(c, helper.Response{
			Message: "Failed to Get Value from Body",
			Errors:  err.Error(),
		})
	}

	result, err := h.humanRepo.Store(ctx, human)
	if err != nil {
		return helper.UnprocResponse(c, helper.Response{
			Message: "Failed to Save Human!",
			Errors:  err.Error(),
		})
	}

	return helper.CreatedResponse(c, helper.Response{
		Message: "Human Saved!",
		Data:    result.InsertedID,
	})
}

func (h *humanHandler) UpdateHumanById(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var human entity.Human
	paramId := c.Param("id")
	if paramId == "" {
		return helper.BadRequestResponse(c, helper.Response{
			Message: "Please Provide Valid Human ID!",
		})
	}

	if err := c.Bind(&human); err != nil {
		return helper.BadRequestResponse(c, helper.Response{
			Message: "Failed to Get Value from Body",
			Errors:  err.Error(),
		})
	}

	result, err := h.humanRepo.UpdateById(ctx, paramId, human)
	if err != nil {
		return helper.UnprocResponse(c, helper.Response{
			Message: fmt.Sprintf("Failed to Update Human With [ID: %s]", paramId),
			Errors:  err.Error(),
		})
	}

	if result.MatchedCount == 0 {
		return helper.NotFoundResponse(c, helper.Response{
			Message: fmt.Sprintf("Human ID [ID: %s] Not Found", paramId),
		})
	}

	return helper.OkResponse(c, helper.Response{
		Message: "Human Updated!",
	})
}

func (h *humanHandler) DeleteHuman(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	paramId := c.Param("id")
	if paramId == "" {
		return helper.BadRequestResponse(c, helper.Response{
			Message: "Please Provide Valid Human ID!",
		})
	}

	result, err := h.humanRepo.Delete(ctx, paramId)
	if err != nil {
		return helper.UnprocResponse(c, helper.Response{
			Message: fmt.Sprintf("Failed to Delete Human With [ID: %s]", paramId),
			Errors:  err.Error(),
		})
	}

	if result.DeletedCount == 0 {
		return helper.NotFoundResponse(c, helper.Response{
			Message: fmt.Sprintf("Human With [ID: %s] Not Found", paramId),
		})
	}

	return helper.OkResponse(c, helper.Response{
		Message: "Human Deleted!",
	})
}
