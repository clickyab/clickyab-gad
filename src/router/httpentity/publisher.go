package httpentity

import (
	"config"
	"entity"
	"models"
)

// Publisher is the publisher interface
type website struct {
	ws            *models.Website
	publisherType entity.PublisherType
}

// GetID return the publisher id
func (ws *website) ID() int64 {
	return ws.ws.WID
}

// FloorCPM is the floor cpm for publisher
func (ws *website) FloorCPM() int64 {
	if ws.ws.WFloorCpm.Int64 < config.Config.Clickyab.MinCPMFloorWeb {
		ws.ws.WFloorCpm.Int64 = config.Config.Clickyab.MinCPMFloorWeb
		ws.ws.WFloorCpm.Valid = true
	}

	return ws.ws.WFloorCpm.Int64
}

// Name of publisher
func (ws *website) Name() string {
	return ws.ws.WDomain.String
}

// Active is the publisher active?
func (ws *website) Active() bool {
	return ws.ws.WStatus == 0 || ws.ws.WStatus == 1
}

// Type return the publisher type
func (ws *website) Type() entity.PublisherType {
	return ws.publisherType
}

// Attributes is the generic attribute system
func (ws *website) Attributes(entity.PublisherAttributes) interface{} {
	// TODO : implement if needed
	return nil
}

// NewHTTPublisherByPublicID return a website by its public id
func NewHTTPublisherByPublicID(publicID string, ptype entity.PublisherType) (entity.Publisher, error) {
	w, err := models.NewManager().FetchWebsiteByPublicID(publicID)
	if err != nil {
		return nil, err
	}

	return &website{
		ws:            w,
		publisherType: ptype,
	}, nil
}
