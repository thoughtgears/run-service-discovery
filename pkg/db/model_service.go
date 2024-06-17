package db

// Service represents a service in the database
// ID is a sha256 hash of the service name
// Name is the name of the service
// URL is the URL of the service
type Service struct {
	ID   string `json:"-" firestore:"id"`
	Name string `json:"name" firestore:"name"`
	URL  string `json:"url" firestore:"url"`
}
