package main

import (
	"context"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/ab-elhaddad/rssagg/internal/database"
	"github.com/google/uuid"
)

func startScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Scraping %v feeds in %v", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Printf("Error getting feeds to fetch: %v", err)
			continue
		}

		wg := sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, feed, &wg)
		}
		wg.Wait()
	}
}

func prepareBatchInsertData(items []RSSItem, feedID uuid.UUID) ([]string, []string, []string, []time.Time, []uuid.UUID) {
	var titles []string
	var urls []string
	var descriptions []string
	var publishedAts []time.Time
	var feedIDs []uuid.UUID

	for _, item := range items {
		publishedAt, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {
			log.Printf("Error parsing published at: %v", err)
			continue
		}

		titles = append(titles, item.Title)
		urls = append(urls, item.Link)
		descriptions = append(descriptions, item.Description)
		publishedAts = append(publishedAts, publishedAt)
		feedIDs = append(feedIDs, feedID)
	}

	return titles, urls, descriptions, publishedAts, feedIDs
}

func scrapeFeed(db *database.Queries, feed database.GetNextFeedsToFetchRow, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("Scraping feed %v", feed.Name)
	feedData, err := fetchFeed(feed.Url)
	if err != nil {
		log.Printf("Error fetching feed: %v", err)
		return
	}

	log.Printf("Feed %v: %v posts found", feed.Name, len(feedData.Channel.Item))

	// Prepare batch insert data
	titles, urls, descriptions, publishedAts, feedIDs := prepareBatchInsertData(feedData.Channel.Item, feed.ID)

	// Batch insert all posts
	if len(titles) > 0 {
		_, err = db.CreatePosts(context.Background(), database.CreatePostsParams{
			Column1: titles,
			Column2: urls,
			Column3: descriptions,
			Column4: publishedAts,
			Column5: feedIDs,
		})
		if err != nil {
			log.Printf("Error batch creating posts: %v", err)
		}
	}

	db.MarkFeedAsFetched(context.Background(), feed.ID)
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(feedURL string) (*RSSFeed, error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := httpClient.Get(feedURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rssFeed RSSFeed
	err = xml.Unmarshal(dat, &rssFeed)
	if err != nil {
		return nil, err
	}

	return &rssFeed, nil
}
