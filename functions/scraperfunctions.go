package functions

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	serrors "github.com/the-rileyj/WebAdvisorScraper/errors"
	"github.com/the-rileyj/WebAdvisorScraper/structs"
)

const (
	loginURL     = "https://loginshib.sdbor.edu/idp/Authn/UserPassword"
	landingURL   = "https://portal.sdbor.edu/all-student/Pages/default.aspx"
	postLoginURL = "https://portal.sdbor.edu/dsu-student/Pages/default.aspx"
)

// compostActions adds Actions/Tasks into a uniform Task list
func composeActions(action chromedp.Action, actions ...chromedp.Action) chromedp.Tasks {
	tasks := make(chromedp.Tasks, 0)

	switch action.(type) {
	case chromedp.Tasks:
		actionTasks := action.(chromedp.Tasks)

		tasks = append(tasks, actionTasks...)
	default:
		tasks = append(tasks, action)
	}

	for _, extraAction := range actions {
		switch extraAction.(type) {
		case chromedp.Tasks:
			extraTasks := extraAction.(chromedp.Tasks)

			tasks = append(tasks, extraTasks...)
		default:
			tasks = append(tasks, extraAction)
		}
	}

	return tasks
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

func GetSubjectInformation(ctx context.Context, c *chromedp.CDP, semester, subject string, semesterIndex, subjectIndex int) ([]structs.CourseInfomation, error) {
	err := ReinitializeSearch(ctx, c)

	if err != nil {
		return nil, err
	}

	tasks := composeActions(
		SelectDropDownOption("VAR1", semester, semesterIndex),
		SelectDropDownOption("LIST_VAR1_1", subject, subjectIndex),
		SelectDropDownOption("VAR6", "DSU", 2),
		chromedp.Click(`//input[@id="WASubmit"]`, chromedp.BySearch),
		chromedp.Sleep(8*time.Second),
	)

	err = c.Run(ctx, tasks)

	if err != nil {
		return nil, err
	}

	return []structs.CourseInfomation{}, nil
}

func GetSubjectList(ctx context.Context, c *chromedp.CDP) ([]string, error) {
	return getNodeValueAttr(ctx, c, `//*[@id="LIST_VAR1_1"]/option[not(@value="")]`)
}

func KillWindows(ctx context.Context, c *chromedp.CDP) {
	var nodes []*cdp.Node

	// Return if there are no nodes found
	if c.Run(ctx, chromedp.Nodes(`//button[@type="button"][@class="btn btn-sm btn-danger"]`, &nodes, chromedp.BySearch)); len(nodes) == 0 {
		return
	}

	// Exit all the extra windows
	for _, node := range nodes[1:] {
		c.Run(ctx, chromedp.Click(node.FullXPath()))
	}
}

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

// func NavigateToSearch(username, password string, auth bool) chromedp.Tasks {
// 	tasks := chromedp.Tasks{
// 		chromedp.Navigate(loginURL),
// 		chromedp.WaitVisible(`#username`, chromedp.ByID),
// 		chromedp.SendKeys(`#username`, username, chromedp.ByID),
// 		chromedp.SendKeys(`#password`, password, chromedp.ByID),
// 		chromedp.Click(`#submitthis`, chromedp.ByID),
// 		chromedp.WaitVisible(`#SearchBox`, chromedp.ByID),
// 		chromedp.Navigate(postLoginURL),
// 		chromedp.WaitVisible(`//*[text()="WebAdvisor for Prospective Students"]`, chromedp.BySearch),
// 		chromedp.Click(`//*[text()="WebAdvisor for Prospective Students"]`, chromedp.BySearch),
// 		chromedp.WaitVisible(`//*[text()="Admission Information"]`, chromedp.BySearch),
// 		chromedp.Click(`//*[text()="Admission Information"]`, chromedp.BySearch),
// 		chromedp.WaitVisible(`//*[text()="Search for Class Sections"]`, chromedp.BySearch),
// 		chromedp.Click(`//*[text()="Search for Class Sections"]`, chromedp.BySearch),
// 		chromedp.WaitVisible(`//div/a/span[text()="Search for Class Sections"]`, chromedp.BySearch),
// 	}

// 	// if auth {
// 	// 	tasks = append(tasks, NavigateToSearchAuthActions(username, password)...)
// 	// }

// 	// tasks = append(tasks, NavigateToSearchActions()...)

// 	return tasks
// }

func NavigateToSearchActions() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.WaitVisible(`#SearchBox`, chromedp.ByID),
		chromedp.Navigate(postLoginURL),
		chromedp.WaitVisible(`//*[text()="WebAdvisor for Prospective Students"]`, chromedp.BySearch),
		chromedp.Click(`//*[text()="WebAdvisor for Prospective Students"]`, chromedp.BySearch),
		chromedp.WaitVisible(`//*[text()="Admission Information"]`, chromedp.BySearch),
		chromedp.Click(`//*[text()="Admission Information"]`, chromedp.BySearch),
		chromedp.WaitVisible(`//*[text()="Search for Class Sections"]`, chromedp.BySearch),
		chromedp.Click(`//*[text()="Search for Class Sections"]`, chromedp.BySearch),
		chromedp.WaitVisible(`//div/a/span[text()="Search for Class Sections"]`, chromedp.BySearch),
	}
}

func NavigateToSearchAuthActions(username, password string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.WaitVisible(`#username`, chromedp.ByID),
		chromedp.SendKeys(`#username`, username, chromedp.ByID),
		chromedp.SendKeys(`#password`, password, chromedp.ByID),
		chromedp.Click(`#submitthis`, chromedp.ByID),
	}
}

// func NavigateToSearchWithCaution(ctx context.Context, c *chromedp.CDP, username, password string) error {
// 	err := c.Run(ctx, chromedp.Tasks{
// 		chromedp.Sleep(1 * time.Second),
// 		chromedp.Navigate(loginURL),
// 		chromedp.WaitVisible(`#username`, chromedp.BySearch)
// 	})

// 	if err != nil {
// 		return err
// 	}

// 	var nodes []*cdp.Node
// 	var currentURL string

// 	err = c.Run(ctx, chromedp.Nodes(`username`, &nodes, chromedp.ByID))

// 	for c.Run(ctx, chromedp.Location(&currentURL)); len(nodes) == 0 && currentURL != landingURL; {
// 		fmt.Println(currentURL)
// 		err = c.Run(ctx, chromedp.Nodes(`username`, &nodes, chromedp.ByID))
// 		fmt.Println(currentURL)
// 	}

// 	if err != nil {
// 		return err
// 	}

// 	tasks := make(chromedp.Tasks, 0)

// 	if len(nodes) == 1 {
// 		tasks = append(tasks, NavigateToSearchAuthActions(username, password)...)
// 	}

// 	tasks = append(tasks, NavigateToSearchActions()...)

// 	err = c.Run(ctx, tasks)

// 	return err
// }

func SelectDropDownOption(id, value string, valueIndex int) chromedp.Tasks {
	var res string

	return chromedp.Tasks{
		chromedp.Evaluate(fmt.Sprintf(`document.getElementById("%s").selectedIndex = %d`, id, valueIndex), &res),
	}
}

func ReinitializeSearch(ctx context.Context, c *chromedp.CDP) error {
	KillWindows(ctx, c)

	err := c.Run(ctx, chromedp.Tasks{
		chromedp.Click(`//li/a/span[text()="Search for Class Sections"]`, chromedp.BySearch),
		chromedp.WaitVisible(`//div/a/span[text()="Search for Class Sections"]`, chromedp.BySearch),
	})

	if err != nil {
		return serrors.WrapFullReinitializationError(err, "could not reinitialize correctly")
	}

	return nil
}
