package podcast

import (
	"PoDownloader/util"
	"github.com/mmcdole/gofeed"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
	"net/http"
	"sync"
)

type Parser struct {
	*gofeed.Parser
}

func NewPodcastParser(httpClient *http.Client) *Parser {
	rssParser := gofeed.NewParser()
	rssParser.Client = httpClient
	return &Parser{rssParser}
}

func (p *Parser) ParsePodcastRSS(RSS string) (*Podcast, error) {
	feed, err := p.ParseURL(RSS)
	if err != nil {
		return nil, err
	}
	var podcastCategories []*Category
	for _, category := range feed.ITunesExt.Categories {
		podcastCategory := &Category{
			Category: category.Text,
		}
		if category.Subcategory != nil {
			podcastCategory.SubCategory = category.Subcategory.Text
		}
		podcastCategories = append(podcastCategories, podcastCategory)
	}
	iTunesExt := &ITunesFeedExtension{
		Author:     feed.ITunesExt.Author,
		Categories: podcastCategories,
		Owner: &ITunesOwner{
			Email: feed.ITunesExt.Owner.Email,
			Name:  feed.ITunesExt.Owner.Name,
		},
		Subtitle: feed.ITunesExt.Subtitle,
		Summary:  feed.ITunesExt.Summary,
		Image:    feed.ITunesExt.Image,
		Explicit: feed.ITunesExt.Explicit,
	}
	var podcastItems []*Item
	for _, item := range feed.Items {
		var enclosures []*Enclosure
		for _, enclosure := range item.Enclosures {
			enclosures = append(enclosures, &Enclosure{
				URL:    enclosure.URL,
				Length: enclosure.Length,
				Type:   enclosure.Type,
			})
		}
		podcastItems = append(podcastItems, &Item{
			Title:       item.Title,
			SafeTitle:   util.SanitizeFileName(item.Title),
			Description: item.Description,
			PubDate:     item.PublishedParsed,
			GUID:        item.GUID,
			Enclosures:  enclosures,
			ITunesExt: &ITunesItemExtension{
				Author:   item.ITunesExt.Author,
				Subtitle: item.ITunesExt.Subtitle,
				Image:    item.ITunesExt.Image,
				Duration: item.ITunesExt.Duration,
				Order:    item.ITunesExt.Order,
			},
		})
	}
	return &Podcast{
		RSS:         RSS,
		Title:       feed.Title,
		SafeTitle:   util.SanitizeFileName(feed.Title),
		Description: feed.Description,
		ITunesExt:   iTunesExt,
		Items:       podcastItems,
	}, nil
}

func (p *Parser) ParsePodcastsFromRSSList(rssList []string) ([]*Podcast, []string) {
	var (
		podcasts []*Podcast
		failed   []string
	)
	downWg := new(sync.WaitGroup)
	downWg.Add(1)
	progressBar := mpb.New(mpb.WithWaitGroup(downWg))
	task := "[Parse Podcast RSS]"
	bar := progressBar.AddBar(
		int64(len(rssList)),
		mpb.PrependDecorators(
			decor.Name(task, decor.WC{W: len(task) + 1, C: decor.DidentRight}),
			decor.CountersNoUnit("%d / %d", decor.WCSyncWidth),
		),
		mpb.AppendDecorators(decor.Percentage(decor.WC{W: 5})),
	)
	go func() {
		defer downWg.Done()
		for _, rss := range rssList {
			podcast, err := p.ParsePodcastRSS(rss)
			if err != nil {
				failed = append(failed, rss)
			} else {
				podcasts = append(podcasts, podcast)
			}
			bar.IncrBy(1)
		}
	}()
	progressBar.Wait()
	return podcasts, failed
}
