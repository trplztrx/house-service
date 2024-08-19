package transport

type ApartmentCreateRequest struct {
	ApartmentID int `json:"flat_id"`
	HouseID 	int `json:"house_id"`
	Price   	int `json:"price"`
	Rooms   	int `json:"rooms"`
}

type ApartmentCreateResponse struct {
	ApartmentID int		`json:"flat_id"`
	HouseID 	int		`json:"house_id"`
	Price 		int		`json:"price"`
	Rooms 		int		`json:"rooms"`
	Status		string 	`json:"status"`
}

type ApartmentUpdateRequest struct {
	ApartmentID		int    `json:"id"`
	HouseID 		int    `json:"house_id"`
	Status			string `json:"status,omitempty"`
}

type ApartmentInfo struct {
    ApartmentID string `json:"apartment_id"`
    HouseID		int	   `json:"house_id"`
    Price       int    `json:"price"`
    Rooms   	int    `json:"rooms"`
    Status      string `json:"status"`
}