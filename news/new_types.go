package news

// Query stores stuff
type Query struct {
	Chan Channel `xml:"Channel"`
}

// Channel stores stuff
type Channel struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
	Hello string `xml:"hello"`
}
