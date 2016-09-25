package utils

// Initializer is a simple interface to handle extra customization on generated codes
type Initializer interface {
	// Initialize is used to handle the code generator extra works
	Initialize()
}

// DoInitialize try to initialize the object if the object is initializer
func DoInitialize(in interface{}) {
	if i, ok := in.(Initializer); ok {
		i.Initialize()
	}
}
