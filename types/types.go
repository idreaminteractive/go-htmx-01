package types

type IndexData struct {
	Name  string
	Posts []PostData
}

type PostData struct {
	Id      int
	Name    string
	Content string
}
