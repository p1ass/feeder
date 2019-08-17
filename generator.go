package feeder

import "github.com/p1ass/feeds"

// ToRSS generates RSS feed
func (f *Feed) ToRSS() (string, error) {
	return f.convert().ToRss()
}

// ToAtom generates Atom feed
func (f *Feed) ToAtom() (string, error) {
	return f.convert().ToAtom()
}

// ToJSON generates JSON feed
func (f *Feed) ToJSON() (string, error) {
	return f.convert().ToJSON()
}

func (l *Link) convert() *feeds.Link {
	return &feeds.Link{
		l.Href,
		l.Rel,
		l.Type,
		l.Length,
	}
}

func (a *Author) convert() *feeds.Author {
	return &feeds.Author{
		a.Name,
		a.Email,
	}
}
func (i *Image) convert() *feeds.Image {
	return &feeds.Image{
		i.URL,
		i.Title,
		i.Link,
		i.Width,
		i.Height,
	}
}

func (e *Enclosure) convert() *feeds.Enclosure {
	return &feeds.Enclosure{
		e.URL,
		e.Length,
		e.Type,
	}
}

func (i *Item) convert() *feeds.Item {
	feedsItem := &feeds.Item{
		Title:       i.Title,
		Description: i.Description,
		Id:          i.ID,
		Content:     i.Content,
	}

	if i.Link != nil {
		feedsItem.Link = i.Link.convert()
	}
	if i.Updated != nil {
		feedsItem.Updated = *i.Updated
	}
	if i.Created != nil {
		feedsItem.Created = *i.Created
	}

	if i.Source != nil {
		feedsItem.Source = i.Source.convert()
	}

	if i.Author != nil {
		feedsItem.Author = i.Author.convert()
	}

	if i.Enclosure != nil {
		feedsItem.Enclosure = i.Enclosure.convert()
	}

	return feedsItem
}

func convert(items []*Item) []*feeds.Item {
	convertedItems := []*feeds.Item{}

	for _, i := range items {
		convertedItems = append(convertedItems, i.convert())
	}
	return convertedItems
}

func (f *Feed) convert() *feeds.Feed {
	feed := &feeds.Feed{
		Title:       f.Title,
		Description: f.Description,
		Updated:     f.Updated,
		Created:     f.Created,
		Id:          f.Id,
		Subtitle:    f.Subtitle,
		Items:       convert(f.Items),
		Copyright:   f.Copyright,
	}

	if f.Link != nil {
		feed.Link = f.Link.convert()
	}

	if f.Author != nil {
		feed.Author = f.Author.convert()
	}

	if f.Image != nil {
		feed.Image = f.Image.convert()
	}

	return feed
}
