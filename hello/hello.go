package hello

type User struct {
	ID   int64
	Name string
	Addr *Address
}

type Address struct {
	City string
	ZIP  int
}

func Hello() string {
	address := Address{"Tokyo", 1234567}
	user := User{1, "Kaori", &address}
	return "HELLO, " + user.Name
}
