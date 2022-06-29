package model

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var UserAll = []User{{Id: "1", Username: "A", Password: "A"}, {Id: "2", Username: "B", Password: "B"}, {Id: "3", Username: "C", Password: "C"}, {Id: "4", Username: "D", Password: "D"}}
var GroupMap = map[string][]User{"game": {UserAll[0], UserAll[1], UserAll[2]}, "work": {UserAll[1]}}
var UserGroupMap = map[string][]string{"B": {"game", "work"}}
