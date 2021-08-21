package restaurantlisting

type RestaurantModel struct {
	Name                 string      `json:"name"`
	Id                   string      `json:"id"`
	CityId               string      `json:"city_id"`
	EstimatedArrivalTime string      `json:"estimated_arrival_time"`
	MinimumPrice         string      `json:"minimum_price"`
	ImageUrl             string      `json:"image_url"`
	CategoryId           string      `json:"category_id"`
	Menu                 []MunuModel `json:"menu"`
}

type MunuModel struct {
	Name        string  `json:"name"`
	Id          string  `json:"id"`
	ImageUrl    string  `json:"image_url"`
	Price       string  `json:"price"`
	Description string  `json:"description"`
	DoublePrice float32 `json:"doublePrice"`
}
