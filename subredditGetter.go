package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/vartanbeno/go-reddit/v2/reddit"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var url string = "www.reddit.com"
var progName string = "subredditGetter.go"
var progVersionString string = "v0.1.5"
var progVersion = [2]string{progName + progVersionString, progVersionString}

// for file reading
func check(e error) {
	if e != nil {
		panic(e)
	}
}
func envFill() {
	fmt.Print("Please enter your reddit Client ID: ")
	var clientID string
	fmt.Scanln(&clientID)

	fmt.Print("Please enter your reddit Client Secret: ")
	var clientSecret string
	fmt.Scanln(&clientSecret)

	fmt.Print("Please enter your reddit Username: ")
	var username string
	fmt.Scanln(&username)

	fmt.Print("Please enter your reddit Password: ")
	var password string
	fmt.Scanln(&password)
	fmt.Println("")
	f, err := os.Create(".env")
	check(err)
	defer f.Close()
	_, err2 := f.WriteString("client_id=" + clientID + "\n" + "client_secret=" + clientSecret + "\n" + "username=" + username + "\n" + "password=" + password)
	check(err2)
	log.Println("Successfully created .env file")
	godotenv.Load()
}

// checks env file
func init() {
	err := godotenv.Load()
	if err != nil {
		// if no env file found, create new
		log.Println("Error loading .env file")
		envFill()

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
		Time: "all", // sorts by all time
	})
	check(err)

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
