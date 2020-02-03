package news

// BiggerQuery does stuff
type BiggerQuery struct {
	Rss interface{} `xml:"rss"`
}

// // Channel stores stuff
// type Channel struct {
// 	Titles string   `xml:"title"`
// 	Links  []string `xml:"link"`
// 	Hello  string   `xml:"hello"`
// }

// Query stores query info
type Query struct {
	Channel Channel `xml:"channel"`
}

// Channel stores channel info
type Channel struct {
	Items []Item `xml:"item"`
}

// Item stores items
type Item struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
}
