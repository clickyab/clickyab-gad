package routes

// Publisher is a publisher object, app, website
type Publisher interface {
	// GetID return the id of object
	GetID() int64
	// GetName return the name of object
	GetName() string

	// FloorCPM is the floor value for this site
	FloorCPM() int64

	// GetActive return if the app is acive
	GetActive() bool
}
