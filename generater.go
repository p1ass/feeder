package feeder

import "github.com/naoki-kishi/feeds"

func (f *Feed) ToRSS() (string, error) {
	return f.convert().ToRss()
}

func (f *Feed) ToAtom() (string, error) {
	return f.convert().ToAtom()
}

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
		i.Url,
		i.Title,
		i.Link,
		i.Width,
		i.Height,
	}
}

func (e *Enclosure) convert() *feeds.Enclosure {
	return &feeds.Enclosure{
		e.Url,
		e.Length,
		e.Type,
	}
}

func (i *Item) convert() *feeds.Item {
	return &feeds.Item{
		i.Title,
		i.Link.convert(),
		i.Source.convert(),
		i.Author.convert(),
		i.Description,
		i.Id,
		*i.Updated,
		*i.Created,
		i.Enclosure.convert(),
		i.Content,
	}
}

func (items *Items) convert() []*feeds.Item {
	convertedItems := []*feeds.Item{}

	for _, i := range items.items {
		convertedItems = append(convertedItems, i.convert())
	}
	return convertedItems
}

func (f *Feed) convert() *feeds.Feed {
	return &feeds.Feed{
		f.Title,
		f.Link.convert(),
		f.Description,
		f.Author.convert(),
		f.Updated,
		f.Created,
		f.Id,
		f.Subtitle,
		f.Items.convert(),
		f.Copyright,
		f.Image.convert(),
	}
}
