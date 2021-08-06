package opml

// OPML is the root node of an OPML document
type OPML struct {
	Head *Head `xml:"head"`
	Body *Body `xml:"body"`
}

// Head is the head section of OPML
type Head struct {
	Title        string `xml:"title"`
	DateCreated  string `xml:"dateCreated,omitempty"`
	DateModified string `xml:"dateModified,omitempty"`
	OwnerName    string `xml:"ownerName,omitempty"`
	OwnerEmail   string `xml:"ownerEmail,omitempty"`
	OwnerID      string `xml:"ownerId,omitempty"`
}

// Body is the body section of OPML
type Body struct {
	Outlines []*Outline `xml:"outline"`
}

// Outline contains all information in an outline tag
type Outline struct {
	Text        string `xml:"text,attr"`
	Title       string `xml:"title,attr,omitempty"`
	Type        string `xml:"type,attr,omitempty"`
	XMLUrl      string `xml:"xmlUrl,attr,omitempty"`
	HTMLUrl     string `xml:"htmlUrl,attr,omitempty"`
	Description string `xml:"description,attr,omitempty"`
}

// GetAllXMLUrl gets all outline XMLUrls from OPML
func (o *OPML) GetAllXMLUrl() []string {
	var XMLUrls []string
	for _, outline := range o.Body.Outlines {
		XMLUrls = append(XMLUrls, outline.XMLUrl)
	}
	return XMLUrls
}
