package functions

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

func NavigateToSearch(username, password string) chromedp.Tasks {
	sdborURL := `https://portal.sdbor.edu`

	return chromedp.Tasks{
		chromedp.Sleep(1 * time.Second),
		chromedp.Navigate(sdborURL),
		chromedp.WaitVisible(`#username`, chromedp.ByID),
		chromedp.SendKeys(`#username`, username, chromedp.ByID),
		chromedp.SendKeys(`#password`, password, chromedp.ByID),
		chromedp.Click(`#submitthis`, chromedp.ByID),
		chromedp.WaitVisible(`//div[@id="SearchBox"]`, chromedp.BySearch),
		chromedp.Navigate(fmt.Sprintf("%s/dsu-student/Pages/default.aspx", sdborURL)),
		chromedp.WaitVisible(`//*[text()="WebAdvisor for Prospective Students"]`, chromedp.BySearch),
		chromedp.Click(`//*[text()="WebAdvisor for Prospective Students"]`, chromedp.BySearch),
		chromedp.WaitVisible(`//*[text()="Admission Information"]`, chromedp.BySearch),
		chromedp.Click(`//*[text()="Admission Information"]`, chromedp.BySearch),
		chromedp.WaitVisible(`//*[text()="Search for Class Sections"]`, chromedp.BySearch),
		chromedp.Click(`//*[text()="Search for Class Sections"]`, chromedp.BySearch),
	}
}

func getNodeValueAttr(ctx context.Context, c *chromedp.CDP, xpath string) ([]string, error) {
	var nodes []*cdp.Node

	if err := c.Run(ctx, chromedp.Nodes(xpath, &nodes, chromedp.BySearch)); err != nil {
		return []string{}, err
	}

	var ok bool
	var attrValue string

	valuesList := make([]string, 0)

	for _, node := range nodes {
		c.Run(ctx, chromedp.AttributeValue(node.FullXPath(), "value", &attrValue, &ok))

		if !ok {
			return []string{}, nil
		}

		valuesList = append(valuesList, attrValue)
	}

	return valuesList, nil
}

func GetSemesterList(ctx context.Context, c *chromedp.CDP) ([]string, error) {
	return getNodeValueAttr(ctx, c, `//*[@id="VAR1"]/option[not(@value="")]`)
}

func GetSubjectInformation(ctx context.Context, c *chromedp.CDP) (structs.SubjectInfomation, error) {

}

func GetSubjectList(ctx context.Context, c *chromedp.CDP) ([]string, error) {
	return getNodeValueAttr(ctx, c, `//*[@id="LIST_VAR1_1"]/option[not(@value="")]`)
}
