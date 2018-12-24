package rchrome

import (
	"context"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	// "io/ioutil"

	"time"
)

func (c *Chrome) RunScreenshot(quit chan struct{}, ch chan []byte, delay int) {
	ctxt := c.ctxt
	cdp := c.cdp

	go func() {
		for {
			select {
			case <-quit:
				return
			default:
				var buf []byte
				err := cdp.Run(ctxt, chromedp.Tasks{
					// chromedp.WaitEnabled(sel, chromedp.ByQuery),
					// setLayoutMetrics(),
					chromedp.CaptureScreenshot(&buf),
				})
				if err == nil {
					ch <- buf
				}
				// err = ioutil.WriteFile(name, buf, 0644)
				time.Sleep(time.Duration(delay) * time.Millisecond)
			}
		}
	}()
}

func setLayoutMetrics() chromedp.ActionFunc {
	return func(i context.Context, executor cdp.Executor) error {
		_, _, contentSize, err := page.GetLayoutMetrics().Do(i, executor)
		if err != nil {
			return err
		}
		w := int64(contentSize.Width)
		h := int64(contentSize.Height)
		scale := 1.0
		err = emulation.SetDeviceMetricsOverride(w, h, scale, false).
			WithScale(scale).Do(i, executor)
		if err != nil {
			return err
		}

		return nil
	}
}

func (c *Chrome) Navigate(url string) error {
	cdp := c.cdp
	ctxt := c.ctxt

	return cdp.Run(ctxt, chromedp.Navigate(url))
}
