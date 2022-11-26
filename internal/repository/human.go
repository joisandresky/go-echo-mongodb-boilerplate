package repository

import (
	"context"

	paginate "github.com/gobeam/mongo-go-pagination"
	"github.com/joisandresky/go-echo-mongodb-boilerplate/database"
	"github.com/joisandresky/go-echo-mongodb-boilerplate/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const CollectionName = "humans"

type HumanRepository interface {
	Paginate(ctx context.Context, page int64, limit int64, filter bson.M) ([]entity.Human, *paginate.PaginationData, error)
	FindAll(ctx context.Context) ([]entity.Human, error)
	FindById(ctx context.Context, id string) (*entity.Human, error)
	Store(ctx context.Context, human entity.Human) (*mongo.InsertOneResult, error)
	UpdateById(ctx context.Context, id string, human entity.Human) (*mongo.UpdateResult, error)
	Delete(ctx context.Context, id string) (*mongo.DeleteResult, error)
}

type humanRepository struct {
	c *mongo.Collection
}

func NewHumanRepository(conn database.Connection) HumanRepository {
	return &humanRepository{conn.DB().Collection(CollectionName)}
}

func (r *humanRepository) Paginate(ctx context.Context, page int64, limit int64, filter bson.M) ([]entity.Human, *paginate.PaginationData, error) {
	var humans []entity.Human
	paginated, err := paginate.New(r.c).Context(ctx).Limit(limit).Page(page).Filter(filter).Decode(&humans).Find()
	if err != nil {
		return nil, nil, err
	}

	return humans, &paginated.Pagination, err
}

func (r *humanRepository) FindAll(ctx context.Context) ([]entity.Human, error) {
	var humans []entity.Human
	cursor, err := r.c.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var human entity.Human
		err := cursor.Decode(&human)
		if err != nil {
			return nil, err
		}

		humans = append(humans, human)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return humans, nil
}

func (r *humanRepository) FindById(ctx context.Context, id string) (*entity.Human, error) {
	var human *entity.Human
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	if err := r.c.FindOne(ctx, bson.M{"_id": objId}).Decode(&human); err != nil {
		return nil, err
	}

	return human, nil
}

func (r *humanRepository) Store(ctx context.Context, human entity.Human) (*mongo.InsertOneResult, error) {
	result, err := r.c.InsertOne(ctx, &human)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *humanRepository) UpdateById(ctx context.Context, id string, human entity.Human) (*mongo.UpdateResult, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return r.c.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": human})
}

func (r *humanRepository) Delete(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return r.c.DeleteOne(ctx, bson.M{"_id": objId})
}
