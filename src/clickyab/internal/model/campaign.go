package model

import (
	"database/sql"
	"entity"
)

type Campaign struct {
	CampaignID        sql.NullInt64  `db:"cp_id"`
	CampaignName      sql.NullString `db:"cp_name"`
	CampaignMaxBid    sql.NullInt64  `db:"cp_maxbid"`
	CampaignFrequency sql.NullInt64  `db:"cp_frequency"`
	capping           entity.Capping `db:"-"`
}

func (cp *Campaign) ID() int64 {
	return cp.CampaignID.Int64
}

func (cp *Campaign) Name() string {
	return cp.CampaignName.String
}

func (cp *Campaign) MaxBID() int64 {
	return cp.CampaignMaxBid.Int64
}

func (cp *Campaign) Frequency() int {
	return int(cp.CampaignFrequency.Int64)
}

func (cp *Campaign) SetCapping(c entity.Capping) {
	cp.capping = c
}
