package db

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/firestore"
)

// FirestoreDB represents a service in the database
type FirestoreDB struct {
	client     *firestore.Client
	collection string
}

// NewFirestoreDB creates a new FirestoreDB instance
// It requires a context and a projectID
// The projectID is the GCP project ID where Firestore is located
func NewFirestoreDB(ctx context.Context, projectID string) (*FirestoreDB, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("firestore.NewClient: %v", err)
	}

	return &FirestoreDB{
		client:     client,
		collection: "services",
	}, nil
}

// GetService retrieves a service from Firestore by its name
// If the service does not exist, an error will be returned
// The ID is a sha256 hash of the service name
func (db *FirestoreDB) GetService(ctx context.Context, ID string) (*Service, error) {
	var service Service

	docRef := db.client.Collection(db.collection).Doc(ID)
	doc, err := docRef.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting document in firestore")
	}

	if err := doc.DataTo(&service); err != nil {
		return nil, fmt.Errorf("doc.DataTo: %v", err)
	}

	return &service, nil
}

// GetServices retrieves a list of services from Firestore
func (db *FirestoreDB) GetServices(ctx context.Context) ([]*Service, error) {
	var services []*Service

	iter := db.client.Collection(db.collection).Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error getting document: %v", err)
		}
		var service *Service

		if err := doc.DataTo(&service); err != nil {
			return nil, fmt.Errorf("error parsing document: %v", err)
		}

		services = append(services, service)
	}

	return services, nil
}

// SetService sets a service in Firestore, it will be used both for updating and creating services
// If the service does not exist, it will be created
func (db *FirestoreDB) SetService(ctx context.Context, ID string, data Service) error {
	_, err := db.client.Collection(db.collection).Doc(ID).Set(ctx, data)
	return err
}

// UpdateServiceURL updates the URL of a service in Firestore
// It requires the ID of the service and the new URL
func (db *FirestoreDB) UpdateServiceURL(ctx context.Context, ID string, url string) error {
	_, err := db.client.Collection(db.collection).Doc(ID).Update(ctx, []firestore.Update{
		{Path: "url", Value: url},
	})
	return err
}

// SetID generates a unique ID for a service based on its name
// This is useful for Firestore, which requires a unique ID for each document
func SetID(name string) string {
	hash := sha256.Sum256([]byte(name))
	return hex.EncodeToString(hash[:])
}

// DeleteService deletes a service from Firestore by its ID
func (db *FirestoreDB) DeleteService(ctx context.Context, ID string) error {
	_, err := db.client.Collection(db.collection).Doc(ID).Delete(ctx)
	return err
}
