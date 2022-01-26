package main

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"

	"os"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

func main() {

	opts := []selenium.ServiceOption{
		selenium.Output(os.Stderr), // Output debug information to STDERR
	}
	service, err := selenium.NewChromeDriverService("/home/nexdata/chromedriver", 9515, opts...)
	if err != nil {
		fmt.Printf("Error starting the ChromeDriver server: %v", err)
	}
	defer service.Stop()

	// call browser
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	// set chrome arguments
	chromeCaps := chrome.Capabilities{
		Args: []string{
			"--headless",   // do not open the browser (run in background)
			"--no-sandbox", //  allow non-root to execute chrome
			"--disable-deb-shm-usage",
			"--window-size=1400,1500",
			//"--start-maximized",
		},
	}
	caps.AddChrome(chromeCaps)

	// connect to the webdriver instance which running locally
	wd, err := selenium.NewRemote(caps, "http://127.0.0.1:9515/wd/hub")
	if err != nil {
		fmt.Printf("connect to the webDriver faild: %v", err)
	}
	// delay closing Chrome
	defer wd.Quit()

	// connect to the target website
	if err := wd.Get("https://www.selenium.dev/projects/"); err != nil {
		fmt.Printf("connect to the reflow server failed: %v", err)
	}

	time.Sleep(time.Duration(1) * time.Second)

	ele, _ := wd.FindElement(selenium.ByXPATH, "/html/body/div/main/div[1]/div")
	scrnsht, _ := wd.Screenshot()
	ioutil.WriteFile("scrnsht.png", scrnsht, 0666)
	loc, _ := ele.Location()
	sz, _ := ele.Size()
	// fmt.Println(loc)
	// fmt.Println(sz)
	file, _ := os.Open("./scrnsht.png")
	defer file.Close()
	img, _ := png.Decode(file)
	sub_image := img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(image.Rect(loc.X, loc.Y, loc.X+sz.Width, loc.Y+sz.Height))
	file, _ = os.Create("./crop.png")
	png.Encode(file, sub_image)

	fmt.Println("爬取完成")

}
