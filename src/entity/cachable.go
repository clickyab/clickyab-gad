package entity

import (
	"fmt"
	"io"
)

// Serializable represent the object that can be serialized
type Serializable interface {
	// Decode is the decoder of this function
	Decode(io.Writer) error
	// Encode is the encoder function
	Encode(io.Reader) error
}

// Cacheable is the object that can be cached into
type Cacheable interface {
	Serializable
	fmt.Stringer
}
