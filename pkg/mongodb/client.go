package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const timeout = 10 * time.Second

// NewClient established connection to a mongoDb instance using provided URI and auth credentials.
func NewClient(uri, username, password string) (*mongo.Client, error) {

	clientOptions := options.Client().ApplyURI(uri)
	err := clientOptions.Validate()
	if err != nil {
		return nil, err
	}

	if username != "" && password != "" {
		clientOptions.SetAuth(options.Credential{
			Username: username,
			Password: password,
		})
	}

	//creating context for connection
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}
