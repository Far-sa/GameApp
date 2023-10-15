package entity

type Category struct {
	ID          uint
	Title       string
	Description string
}

// *
// type Category string

// const (
// 	CategorySport   Category = "sport"
// 	CategoryHistory Category = "history"
// )

//*
// type Category uint

// const (
// 	CategorySport   Category = iota + 1
// 	CategoryHistory
//  CategoryTechonolgy
// )

// func (c Category) String() string {
// 	switch c {
// 	case 1:
// 		return "sport"
// 	case 2:
// 		return "techonolgy"
// 	}

// 	return ""
// }
