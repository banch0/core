package core

// Service ...
type Service struct {
	ID    int64
	Name  string `json:"-" xml:"-"`
	Price string `json:"-" xml:"-"`
}

// ATM ...
type ATM struct {
	ID      int64  `json:"-" xml:"-"`
	Name    string `json:"name" xml:"name"`
	Address string `json:"addres" xml:"addres"`
}

// UserType ...
type UserType struct {
	ID       int64  `json:"-" xml:"-"`
	Name     string `json:"name" xml:"name"`
	Phone    int    `json:"-" xml:"-"`
	Account  int    `json:"account" xml:"account"`
	Login    string `json:"-" xml:"-"`
	Password string `json:"-" xml:"-"`
	Balance  int64  `json:"-" xml:"-"`
}

// Manager ...
type Manager struct {
	ID       int64  `json:"-" xml:"-"`
	Name     string `json:"-" xml:"-"`
	Login    string `json:"-" xml:"-"`
	Password string `json:"-" xml:"-"`
}

// Sales ...
type Sales struct {
	ID        int64 `json:"-" xml:"-"`
	UserID    int64 `json:"-" xml:"-"`
	ServiceID int64 `json:"-" xml:"-"`
	Price     int64 `json:"-" xml:"-"`
}
