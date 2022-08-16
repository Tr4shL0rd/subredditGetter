package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/vartanbeno/go-reddit/v2/reddit"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var url string = "www.reddit.com"
var progVersion = [2]string{"subredditGetter.go v0.0.5", "v0.0.5"}

// for file reading
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// checks env file
func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func main() {
	title := flag.Bool("title", false, "get title of posts")
	link := flag.Bool("url", false, "get url of posts")
	both := flag.Bool("both", false, "get title and url of posts")
	sreddit := flag.String("subreddit", "golang", "searches in the specified subreddit")
	old := flag.Bool("old", false, "uses old.reddit.com")
	limit := flag.Int("limit", 10, "limits the amount of posts")
	help := flag.Bool("help", false, "Outputs this help message")
	version := flag.Bool("version", false, "outputs version")
	flag.Parse()

	postLimit := 10
	subreddit := "all"
	if *sreddit != "" {
		subreddit = *sreddit
	}
	if *old {
		url = "old.reddit.com"
	}
	if *limit >= 1 {
		postLimit = *limit
	}

	ctx := context.Background()
	clientID := os.Getenv("client_id")
	clientSecret := os.Getenv("client_secret")
	username := os.Getenv("username")
	password := os.Getenv("password")

	Credentials := reddit.Credentials{
		ID:       clientID,
		Secret:   clientSecret,
		Username: username,
		Password: password,
	}
	//client, err
	client, _ := reddit.NewClient(Credentials)

	//posts, resp, err
	posts, _, err := client.Subreddit.TopPosts(ctx, subreddit, &reddit.ListPostOptions{
		ListOptions: reddit.ListOptions{
			Limit: postLimit,
		},
		Time: "all",
	})
	if err != nil {
		fmt.Println(err)
	}

	for _, post := range posts {
		if *title {
			postTitle(post)
		} else if *link {
			postURL(post)
		} else if *both {
			postTitleUrl(post)
		} else if *help {
			helpMessage()
		} else if *version {
			fmt.Println(progVersion[1])
			os.Exit(0)
		} else {
			helpMessage()
		}

	}

}
func postUpvotes(post *reddit.Post) string {
	p := message.NewPrinter(language.German)
	withComSep := p.Sprintf("%d", post.Score)
	return withComSep
}
func postTitle(post *reddit.Post) {
	fmt.Printf("[%s] %s\n", postUpvotes(post), post.Title)

}

func postURL(post *reddit.Post) {
	url := url + post.Permalink
	fmt.Println(url)
}

func postTitleUrl(post *reddit.Post) {
	fmt.Printf("[%s] | %s\n%s%s\n\n", postUpvotes(post), post.Title, url, post.Permalink)
}

func helpMessage() {
	dat, err := os.ReadFile("help.txt")
	check(err)
	fmt.Println(string(dat))
	fmt.Println(progVersion[1])

	os.Exit(0)
}
