package model

import "encoding/json"

type Profile struct {
	Url       string `json:"url"`
	Id        int    `json:"memberID"`
	Name      string `json:"nickname"`
	Genter    string `json:"genderString"`
	Age       int    `json:"age"`
	Height    string `json:"heightString"`
	Income    string `json:"salaryString"`
	Marriage  string `json:"marriageString"`
	Education string `json:"educationString"`
	AvatarURL string `json:"avatarURL"`
	Hokou     string `json:"workCityString"`
	Xinzuo    string
	House     string
	Car       string
}

func FromJsonObj(o interface{}) (Profile, error) {
	var profile Profile
	s, err := json.Marshal(o)
	if err != nil {
		return profile, err
	}
	err = json.Unmarshal(s, &profile)
	if err != nil {
		return profile, err
	}
	return profile, nil

}
