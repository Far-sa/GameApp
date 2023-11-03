package entity

// type Category struct {
// 	ID          uint
// 	Title       string
// 	Description string
// }

// *
type Category string

const (
	SportCategory   = "football"
	HistoryCategory = "history"
)

func (c Category) IsValid() bool {
	switch c {
	case SportCategory:
		return true
	}
	return false
}

func CategoryList() []Category {
	return []Category{SportCategory, HistoryCategory}

}

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
