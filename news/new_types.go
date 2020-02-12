package news

import (
	"time"
)

// BiggerQuery does stuff
type BiggerQuery struct {
	Rss interface{} `xml:"rss"`
}

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
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	PubTime string `xml:"pubDate"`
}

// Storage stores items from a site
type Storage struct {
	Path       string
	LastUpdate time.Time
	Items      []Item
}

// AuthParams has auth request params
type AuthParams struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}
