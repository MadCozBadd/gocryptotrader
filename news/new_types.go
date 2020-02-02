package news

// Query stores stuff
type Query struct {
	Chan Channel `xml:"Channel"`
}

// Link stores data in link
type Link struct {
	Link string `xml:"link"`
}

// Channel stores stuff
type Channel struct {
	Title string `xml:"title"`
	Links []Link
	Hello string `xml:"hello"`
}
