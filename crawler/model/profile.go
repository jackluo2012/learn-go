package model

type Profile struct {
	Name       string `json:"nickname"`
	Genter     string `json:"genderString"`
	Age        int    `json:"age"`
	Height     string `json:"heightString"`
	Income     string `json:"salaryString"`
	Marriage   string `json:"marriageString"`
	Education  string `json:"educationString"`
	AvatarURL string	`json:"avatarURL"`
	Hokou      string `json:"workCityString"`
	Xinzuo     string
	House      string
	Car        string
}
