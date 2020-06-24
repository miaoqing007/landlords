package obj

type User struct {
	Id         string
	Name       string
	PlatformId string
}

type PlatformData struct {
	Id     string
	UserId string
}
