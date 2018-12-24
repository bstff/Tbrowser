package rchrome

import (
    "Tbrowser/common"
    "context"
    "fmt"

    "github.com/chromedp/chromedp"
    "github.com/chromedp/chromedp/runner"

    "io/ioutil"
    "log"
)

type WaitInputFunc func(ev common.Event)

func SetWaitInputFunc(f WaitInputFunc) {
    waitFunc = f
}

const (
    defaultViewWidth  = 1280
    defaultViewHeight = 960
)

var (
    waitFunc WaitInputFunc

    scrollX = 0
    scrollY = 0
)

type Chrome struct {
    ctxt       context.Context
    cancelFunc context.CancelFunc
    cdp        *chromedp.CDP
}

func NewChromeDP(port int) *Chrome {
    return NewChromeDPWithViewSize(port, defaultViewWidth, defaultViewHeight)
}

func NewChromeDPWithViewSize(port int, w, h int) *Chrome {
    ctxt, cancel := context.WithCancel(context.Background())

    log.SetOutput(ioutil.Discard)
    c, err := chromedp.New(ctxt, chromedp.WithRunnerOptions(
        // runner.Flag("headless", true),
        runner.Flag("disable-gpu", true),
        // runner.Flag("no-sandbox", true),
        runner.Flag("hide-scrollbars", true),
        runner.Flag("remote-debugging-port", port),
        runner.Flag("window-size", fmt.Sprintf("%d,%d", defaultViewWidth, defaultViewHeight)),
    ))
    if err != nil {
        fmt.Println(err)
        cancel()
        return nil
    }

    return &Chrome{
        ctxt,
        cancel,
        c,
    }
}

func (c *Chrome) Close() {

    c.cdp.Shutdown(c.ctxt)
    c.cancelFunc()

    // wait for chrome to finish
    err := c.cdp.Wait()
    if err != nil {
        log.Fatal(err)
    }
}

func (c *Chrome) maybeNavigate(x, y int) bool {
    cdp := c.cdp
    ctxt := c.ctxt

    var err error

    // var urlB string
    // err = cdp.Run(ctxt, chromedp.Location(&urlB))
    // if err != nil {
    //     fmt.Println(err)
    //     return false
    // }

    err = cdp.Run(ctxt, chromedp.MouseClickXY(int64(x), int64(y)))
    if err != nil {
        return false
    }

    // var urlA string
    // err = cdp.Run(ctxt, chromedp.Location(&urlA))
    // if err != nil {
    //     fmt.Println(err)
    //     return false
    // }
    // if urlA != urlB {
    //     // c.Navigate(urlA)
    //     return true
    // }
    return false
}

func (c *Chrome) Clicked(col, row, left, top, right, bottom int) int {
    x := (left + right) / 2
    y := (top + bottom) / 2

    if c.maybeNavigate(x, y) {
        return 1
    }
    // if c.maybeInput(left, top, right, bottom) {
    //     return 2
    // }

    return 0
}
