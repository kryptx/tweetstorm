package twitter

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	twit "github.com/dghubble/go-twitter/twitter"
)

// TweetWriter is a type that receives tweets and writes them somewhere
type TweetWriter interface {
	Write(*twit.Tweet) <-chan error
}

func writeTweet(tweet *twit.Tweet, writers []TweetWriter) {
	var errReceivers []<-chan error
	for _, writer := range writers {
		errReceivers = append(errReceivers, writer.Write(tweet))
	}
	for _, ch := range errReceivers {
		err := <-ch
		if err != nil {
			log.Output(0, err.Error())
		}
	}
}

// StreamTweets streams tweets and writes them to the provided TweetWriter
func StreamTweets(filterTerms []string, httpClient *http.Client, writers []TweetWriter) {
	client := twit.NewClient(httpClient)
	demux := twit.NewSwitchDemux()
	demux.Tweet = func(tweet *twit.Tweet) {
		go writeTweet(tweet, writers)
	}

	// Filter
	filterParams := &twit.StreamFilterParams{
		Track:         filterTerms,
		StallWarnings: twit.Bool(true),
		Language:      []string{"en"},
	}

	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}

	go demux.HandleChan(stream.Messages)

	// stop if the app receives SIGINT or SIGTERM
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	fmt.Println("Stopping Stream...")
	stream.Stop()
}
