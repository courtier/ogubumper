package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/go-rod/bypass"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/pelletier/go-toml"
)

var (
	config *toml.Tree
)

func main() {

	config = loadConfig()

	if len(config.GetArray("threads").([]string)) >= 10 {
		fmt.Println("bumping >= 10 threads might end up with your post being deleted, warning points, and/or a ban")
	}

	fmt.Println("username", config.Get("username").(string))

	if config.Get("automode").(bool) {
		interval := config.Get("autointerval").(int64)
		fmt.Println("bumping auto every", interval, "minutes")
		intervalDur := time.Duration(interval) * time.Minute
		ticker := time.NewTicker(intervalDur)
		for range ticker.C {
			t, _ := time.Now().MarshalText()
			fmt.Println("ticker time", string(t))
			startBump()
		}
	} else {
		fmt.Println("bumping once")
		startBump()
	}

}

func startBump() {

	browser := rod.New().Timeout(5 * time.Minute).MustConnect()
	defer browser.MustClose()

	page := bypass.MustPage(browser)

	page.MustNavigate("https://ogusers.com/member.php?action=login")

	//log in

	//input username
	el := page.MustElementX("/html/body/div[4]/div/form[1]/div/div[1]/div/div[2]/div/span/div[1]/span/label/input")
	el.Focus()
	el.MustInput(config.Get("username").(string))

	//input password
	el = page.MustElementX("/html/body/div[4]/div/form[1]/div/div[1]/div/div[2]/div/span/div[2]/label/input")
	el.Focus()
	el.MustInput(config.Get("password").(string))

	//left click login
	el = page.MustElementX("/html/body/div[4]/div/form[1]/div/div[1]/div/div[2]/div/span/button")
	el.Click(proto.InputMouseButtonLeft)

	t, _ := time.Now().MarshalText()
	fmt.Println("logged in", string(t))

	//manually go to profile then threads created by user
	el = page.MustElementX("/html/body/div[7]/div/b/div[3]/div/div[1]/div[1]/div[1]/a")
	el.Click(proto.InputMouseButtonLeft)
	//click on threads created by user
	el = page.MustElementX("/html/body/div[7]/div/div[3]/div[2]/div/div/div/div/div[2]/div[1]/div[1]/a")
	el.Click(proto.InputMouseButtonLeft)

	//start bumping

	//load threads and messages
	threadsToBump := config.GetArray("threads").([]string)
	messages := config.GetArray("messages").([]string)

	//wait this many seconnds between each delay
	delay := config.Get("bumpdelay").(int64)
	bumpDelay := time.Duration(delay) * time.Second

	//loop through threads
	for index, threadToBump := range threadsToBump {

		//set seed for rand so its actually random
		rand.Seed(time.Now().UnixNano())
		//pick random message
		msgNum := rand.Intn(len(messages))
		//go to thread automatically, ogu dumb and they put two elements with same id on the page
		threadLink := "https://ogusers.com/" + threadToBump
		page.MustNavigate(threadLink)

		//write message in textarea
		el = page.MustElementX("/html/body/div[7]/div[1]/div[3]/form/div/div[1]/div[1]/textarea")
		el.Focus()
		el.Input(messages[msgNum])

		//send message
		el = page.MustElementX("/html/body/div[7]/div[1]/div[3]/form/div/div[2]/input[1]")
		el.Click(proto.InputMouseButtonLeft)

		t, _ := time.Now().MarshalText()
		fmt.Println("bumped thread", threadToBump, string(t))

		if index != (len(threadsToBump)-1) && delay != 0 {
			time.Sleep(bumpDelay)
		}

	}

	//log out for cleaner experience?

	//click on profile dropdown button
	el = page.MustElementX("/html/body/div[1]/div/div[4]/div/button")
	el.Click(proto.InputMouseButtonLeft)
	//log out
	el = page.MustElementX("/html/body/div[1]/div/div[4]/div/div/a")
	el.Click(proto.InputMouseButtonLeft)

	t, _ = time.Now().MarshalText()
	fmt.Println("logged out", string(t))

}

func contains(arr []string, el string) bool {
	for _, e := range arr {
		if e == el {
			return true
		}
	}
	return false
}

func loadConfig() (config *toml.Tree) {
	path := "./config.toml"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("couldnt find", path, ", quitting")
		os.Exit(1)
	} else {
		config, err := toml.LoadFile(path)
		if err != nil {
			fmt.Println("error while loading", path, " -> ", err.Error())
			os.Exit(2)
		}
		return config
	}
	return
}
