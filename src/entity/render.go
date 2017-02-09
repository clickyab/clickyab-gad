package entity

import "io"

// Renderer is the app renderer
type Renderer interface {
	// Render render into app
	Render(Advertise, io.Writer) error
}
