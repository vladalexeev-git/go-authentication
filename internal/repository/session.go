package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"log/slog"

	"sso/internal/domain"
)

type sessionRepo struct {
	log   *slog.Logger
	mongo *mongo.Collection
}

func NewSessionRepo(mongo *mongo.Database, logger *slog.Logger) *sessionRepo {
	return &sessionRepo{mongo: mongo.Collection("session"), log: logger}
}

func (r *sessionRepo) Create(ctx context.Context, session domain.Session) error {
	ttlIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "createdAt", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(int32(session.TTL)),
	}
	_, err := r.mongo.Indexes().CreateOne(ctx, ttlIndex)
	if err != nil {
		r.log.Error("r.mongo.Indexes.CreateOne: can't create index",
			slog.String("error", err.Error()))
		return err
	}

	_, err = r.mongo.InsertOne(ctx, session)
	if err != nil {
		r.log.Error("r.mongo.InsertOne: can't create session",
			slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (r *sessionRepo) FindByID(ctx context.Context, sid string) (domain.Session, error) {
	var session domain.Session

	if err := r.mongo.FindOne(ctx, bson.M{"_id": sid}).Decode(&session); err != nil {
		r.log.Error("findOne: can't find session",
			slog.String("error", err.Error()))
		return domain.Session{}, err
	}
	return session, nil
}

func (r *sessionRepo) FindAll(ctx context.Context, aid string) ([]domain.Session, error) {
	cursor, err := r.mongo.Find(ctx, bson.M{"accountId": bson.M{"$eq": aid}})
	if err != nil {
		r.log.Error("r.mongo.FindAll: can't find sessions",
			slog.String("error", err.Error()))
		return nil, err
	}
	defer cursor.Close(ctx) //todo ?

	var sessions []domain.Session

	if err = cursor.All(ctx, &sessions); err != nil {
		r.log.Error("cursor.All: can't find sessions",
			slog.String("error", err.Error()))
		return sessions, err
	}
	return sessions, nil
}

func (r *sessionRepo) Delete(ctx context.Context, sid string) error {
	_, err := r.mongo.DeleteOne(ctx, bson.M{"_id": sid})
	if err != nil {
		r.log.Error("r.deleteOne",
			slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (r *sessionRepo) DeleteAll(ctx context.Context, aid, currSid string) error {
	_, err := r.mongo.DeleteMany(ctx,
		bson.M{
			"_id":       bson.M{"$ne": currSid},
			"accountId": aid,
		})
	if err != nil {
		r.log.Error("r.deleteMany",
			slog.String("error", err.Error()))
		return err
	}

	return nil
}
