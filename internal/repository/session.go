package repository

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sso/pkg/utils"

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
	const op = "repository.session..create"
	l := r.log.With(slog.String(utils.Operation, op))

	ttlIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "createdAt", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(int32(session.TTL)),
	}
	_, err := r.mongo.Indexes().CreateOne(ctx, ttlIndex)
	if err != nil {
		l.Error("r.mongo.Indexes.CreateOne: can't create index",
			slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	info, err := r.mongo.InsertOne(ctx, session)
	if err != nil {
		l.Error("r.mongo.InsertOne: can't create session",
			slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	l.Debug("insert session to mongodb", slog.Any("info", info))

	return nil
}

func (r *sessionRepo) FindByID(ctx context.Context, sid string) (domain.Session, error) {
	const op = "repository.session.findById"
	l := r.log.With(slog.String(utils.Operation, op))

	var session domain.Session

	if err := r.mongo.FindOne(ctx, bson.M{"_id": sid}).Decode(&session); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			l.Error("findOne: no documents found", slog.String("error", err.Error()))
			return domain.Session{}, fmt.Errorf("%s: %w", op, err)
		}
		l.Error("findOne: can't find session",
			slog.String("error", err.Error()))
		return domain.Session{}, fmt.Errorf("%s: %w", op, err)
	}

	l.Debug("find session from mongodb",
		slog.Any("session", session),
	)

	return session, nil
}

func (r *sessionRepo) FindAll(ctx context.Context, aid string) ([]domain.Session, error) {
	const op = "repository.session.findAll"
	l := r.log.With(slog.String(utils.Operation, op))

	cursor, err := r.mongo.Find(ctx, bson.M{"accountId": bson.M{"$eq": aid}})
	if err != nil {
		l.Error("r.mongo.FindAll: can't find sessions",
			slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer cursor.Close(ctx) //todo ?

	var sessions []domain.Session

	if err = cursor.All(ctx, &sessions); err != nil {
		l.Error("cursor.All: can't find sessions",
			slog.String("error", err.Error()))
		return sessions, fmt.Errorf("%s: %w", op, err)
	}
	return sessions, nil
}

func (r *sessionRepo) Delete(ctx context.Context, sid string) error {
	const op = "repository.session.Delete"
	l := r.log.With(slog.String(utils.Operation, op))

	res, err := r.mongo.DeleteOne(ctx, bson.M{"_id": sid})
	if err != nil {
		l.Error("r.deleteOne",
			slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	l.Debug("deleted document", slog.Int64("count", res.DeletedCount))
	return nil
}

func (r *sessionRepo) DeleteAll(ctx context.Context, aid, currSid string) error {
	const op = "repository.session.deleteAll"
	l := r.log.With(slog.String(utils.Operation, op))

	_, err := r.mongo.DeleteMany(ctx,
		bson.M{
			"_id":       bson.M{"$ne": currSid},
			"accountId": aid,
		})
	if err != nil {
		l.Error("r.deleteMany",
			slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
