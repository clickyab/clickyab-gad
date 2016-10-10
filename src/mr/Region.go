package mr

type RegionData struct {
	LocationID	int64	`json:"location_id" db:"location_id"`
	LocationName	string	`json:"location_name" db:"location_name"`
	LocationNameFa	string	`json:"location_name_persian" db:"location_name_persian"`
	LocationMaster 	bool	`json:"location_master" db:"location_master"`
	LocationSelect	bool	`json:"location_select" db:"location_select"`
}
type RegionsData []RegionData
