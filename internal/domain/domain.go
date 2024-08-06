package domain

type Employee struct {
	UserID    int    `json:"userID"`
	Password  string `json:"password"`
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Day       int    `json:"day"`
	Month     int    `json:"month"`
	Year      int    `json:"year"`
}

type ResponseEmployee struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Username  string `json:"username"`
	Day       int    `json:"day"`
	Month     int    `json:"month"`
	Year      int    `json:"year"`
}

type Notify struct {
	Subscriber int `json:"subscriber"`
	Publisher  int `json:"publisher"`
}
