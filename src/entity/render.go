package entity

import "io"

// HTMLRenderer is the html renderer
type HTMLRenderer interface {
	// HTMLRender render into html
	HTMLRender(io.Writer) error
}

//VASTRenderer is the vast renderer
type VASTRenderer interface {
	// VASTRender is the render interface
	VASTRender(io.Writer) error
}

// APPRenderer is the app renderer
type APPRenderer interface {
	// APPRender render into app
	APPRender(io.Writer) error
}
