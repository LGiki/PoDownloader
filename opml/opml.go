package opml

type OPML struct {
	Head *Head `xml:"head"`
	Body *Body `xml:"body"`
}

type Head struct {
	Title        string `xml:"title"`
	DateCreated  string `xml:"dateCreated,omitempty"`
	DateModified string `xml:"dateModified,omitempty"`
	OwnerName    string `xml:"ownerName,omitempty"`
	OwnerEmail   string `xml:"ownerEmail,omitempty"`
	OwnerId      string `xml:"ownerId,omitempty"`
}

type Body struct {
	Outlines []*Outline `xml:"outline"`
}

type Outline struct {
	Text        string `xml:"text,attr"`
	Title       string `xml:"title,attr,omitempty"`
	Type        string `xml:"type,attr,omitempty"`
	XMLUrl      string `xml:"xmlUrl,attr,omitempty"`
	HtmlUrl     string `xml:"htmlUrl,attr,omitempty"`
	Description string `xml:"description,attr,omitempty"`
}

// GetAllXMLUrl
// Get all outline XML urls from OPML
func (o *OPML) GetAllXMLUrl() []string {
	var XMLUrls []string
	for _, outline := range o.Body.Outlines {
		XMLUrls = append(XMLUrls, outline.XMLUrl)
	}
	return XMLUrls
}
