// This file will run in the background as our server runs
// Long running job
package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/reinhardbuyabo/RSS-Aggregator/internal/database"
)

// This function shouldn't return ath because it's going to be a long running job
func startScraping(
	db *database.Queries, // connection to the database
	concurrency int, // indicate to the start scraping function how many go routines we want to use to go fetch all of those different feeds ... the whole point is we can fetch all of them at the same time // number of concurrency units, i.e, how many go routines we want to do the scraping on
	timeBetweenRequest time.Duration, // how much time we want in between each request to go scrape a new RSS Feed
) {
	log.Printf("Scraping on %v go routines every %s duration", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	// executing the body of the for loop everytime a new value comes across the tickers Channel
	// ticker Struct has a field called C which is a Channel, where every 1 minute will be sent across the channel ...
	// i.e run this for loop every 1 minute
	// we have an empty initialization and update because we want it to execute immediately the first time, the first time we get to line 17, there's immediate execution
	for ; ; <-ticker.C {
		// 1. Context
		// 2. Limit
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(), // this is the GLOBAL context ... what you use if you don't have access to a scoped context like we do for our individual http Request
			int32(concurrency),
		)
		if err != nil {
			log.Println("error fetching feeds:", err)
			continue
		}
		// now we have a slice of feeds
		// LOGIC that fetches it feed individually, and fetching each individually at the SAME TIME
		// We are going to need a SYNCHRONIZATION MECHANISM [wait group]
		wg := &sync.WaitGroup{} // anytime you want to spawn a new go routine within the context of a wait group, you do a wait group.add ... and add some number to it
		// iterating over all the feeds we want to fetch and we're going to add 1 to the wait group // iterating over all the feeds on the same go routine as the startScraping function i.e on the main go routine
		for _, feed := range feeds {
			wg.Add(1) // on the main go routine, we are adding 1 to the wait group, for every feed ... say we had a concurrency of 30, then we would be adding 30 to the wait group

			go scrapeFeed(db, wg, feed) // spawning separate go routines as we do [that?]
		}
		// when they are all done, line 48 will execute and move past that, before they're all done, we'll be blocking on line 35, which is what we want because we don't want to continue to the next iteration of the loop until we're sure we've actually scraped all of the feed
		wg.Wait() // at the end of the loop ... we're going to be waiting on the wait group for 30 distinct calls to Done() ... Done() effectively decrements the counter by 1
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done() // this will always be called at the end of this function

	// 1. Mark that we've fetched this feed
	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Error marking feed as fetched:", err)
		return
	}

	// 2. Scraping the feed(Heavy Lifting)
	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed:", err)
	}

	// iterating over the items
	for _, item := range rssFeed.Channel.Item {
		log.Println("Found post", item.Title, "on feed", feed.Name) // logging each individual post, or that we found a post
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item)) // and logging how many posts we found
}
