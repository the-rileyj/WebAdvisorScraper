// Command click is a chromedp example demonstrating how to use a selector to
// click on an element.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/chromedp/chromedp"
	"github.com/pkg/errors"
	serrors "github.com/the-rileyj/WebAdvisorScraper/errors"
	"github.com/the-rileyj/WebAdvisorScraper/functions"
	"github.com/the-rileyj/WebAdvisorScraper/structs"
)

// func click(str *string) chromedp.Tasks { // You can use xPaths and CSS selectors!

// 	tasks := chromedp.Tasks{
// 		chromedp.Sleep(1 * time.Second),
// 		chromedp.Navigate(`https://portal.sdbor.edu`),
// 		// chromedp.ActionFunc(func(ctext context.Context, ex chromedp.)),
// 		// chromedp.ActionFunc(func(con context.Context, ex cdp.Executor) error {
// 		// 	cdp.Node.
// 		// 	return nil
// 		// }),
// 		chromedp.WaitVisible(`#username`, chromedp.ByID),
// 		// chromedp.WaitVisible(`#submitthis`, chromedp.ByID),
// 		chromedp.WaitVisible(`//div[@id="SearchBox"]`, chromedp.BySearch),
// 		// chromedp.Text(`//title`, str),
// 	}

// 	return tasks
// }

func main() {
	killChannel := make(chan os.Signal, 1)

	signal.Notify(killChannel,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	var err error
	var operationCounter int

	// Create context and defer cancelation
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create chrome instance
	c, err := chromedp.New(ctxt)
	// c, err := chromedp.New(ctxt, chromedp.WithLog(log.Printf))
	// c, err := chromedp.New(ctxt, chromedp.WithLog(log.Printf), chromedp.WithTargets(client.New().WatchPageTargets(ctxt)))
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		signal := <-killChannel

		if signal.String() == "RJ" {
			return
		}

		cancel()

		// Shutdown chrome
		err = c.Shutdown(ctxt)
		if err != nil {
			log.Fatal(err)
		}

		// Wait for chrome to finish
		err = c.Wait()
		if err != nil {
			log.Fatal(err)
		}

		os.Exit(1)
	}()

	defer func() {
		killChannel <- structs.RJSignal{}
	}()

	credentials, err := structs.GetCredentials("../info.json")

	if err != nil {
		fmt.Println(errors.Wrap(err, "could not get credentials"))
		os.Exit(1)
	}

	// c.Run(ctxt, chromedp.Tasks{
	// 	chromedp.Navigate("https://portal.sdbor.edu"),
	// 	chromedp.WaitVisible(`#username`, chromedp.ByID),
	// 	chromedp.SendKeys(`#username`, credentials.Username, chromedp.ByID),
	// 	chromedp.SendKeys(`#password`, credentials.Password, chromedp.ByID),
	// 	chromedp.Click(`#submitthis`, chromedp.ByID),
	// })

	// chromedp.Location

	// Run task list
	// for err = c.Run(ctxt, chromedp.Tasks{
	// 	chromedp.WaitVisible(`#SearchBox`, chromedp.ByID),
	// 	chromedp.Navigate("https://portal.sdbor.edu/dsu-student/Pages/default.aspx"),
	// 	chromedp.WaitVisible(`//*[text()="WebAdvisor for Prospective Students"]`, chromedp.BySearch),
	// 	chromedp.Click(`//*[text()="WebAdvisor for Prospective Students"]`, chromedp.BySearch),
	// 	chromedp.WaitVisible(`//*[text()="Admission Information"]`, chromedp.BySearch),
	// 	chromedp.Click(`//*[text()="Admission Information"]`, chromedp.BySearch),
	// 	chromedp.WaitVisible(`//*[text()="Search for Class Sections"]`, chromedp.BySearch),
	// 	chromedp.Click(`//*[text()="Search for Class Sections"]`, chromedp.BySearch),
	// 	chromedp.WaitVisible(`//div/a/span[text()="Search for Class Sections"]`, chromedp.BySearch),
	// });
	// for err = functions.NavigateToSearch(ctxt, c, credentials.Username, credentials.Password, true); err != nil; {
	// 	if operationCounter == 10 {
	// 		log.Fatal(errors.Wrap(err, "scraper failed to reach initial search position"))
	// 	}

	// 	operationCounter++
	// }

	err = c.Run(ctxt, functions.NavigateToSearch(credentials.Username, credentials.Password))

	if err != nil {
		log.Fatal(errors.Wrap(err, "scraper failed to reach initial search position"))
	}

	// Get list of semesters
	var semesterList []string
	operationCounter = 0
	for semesterList, err = functions.GetSemesterList(ctxt, c); err != nil; {
		if operationCounter == 10 {
			log.Fatal(errors.Wrap(err, "scraper failed to get semester list"))
		}

		semesterList, err = functions.GetSemesterList(ctxt, c)
		operationCounter++
	}

	// Get list of subjects
	var subjectList []string
	operationCounter = 0
	for subjectList, err = functions.GetSubjectList(ctxt, c); err != nil; {
		if operationCounter == 10 {
			log.Fatal(errors.Wrap(err, "scraper failed to get subject list"))
		}

		subjectList, err = functions.GetSubjectList(ctxt, c)
		operationCounter++
	}

	for semesterIndex, semester := range semesterList {
		for subjectIndex, subject := range subjectList {
			var subjectInformation []structs.CourseInfomation
			operationCounter = 0

			// Get information for the given subject
			for subjectInformation, err = functions.GetSubjectInformation(ctxt, c, semester, subject, semesterIndex+1, subjectIndex+1); err != nil; {
				fmt.Println(err)

				if _, ok := err.(*serrors.FullReinitializationError); ok {
					functions.NavigateToSearch("", "")
				}

				if operationCounter == 10 {
					fmt.Println(errors.Wrapf(err, "scraper failed to get information for subject %s in the %s semester", subject, semester))
					break
				}

				subjectInformation, err = functions.GetSubjectInformation(ctxt, c, semester, subject, semesterIndex+1, subjectIndex+1)
				operationCounter++
			}

			fmt.Println(subjectInformation)
		}
	}

	// Shutdown chrome
	err = c.Shutdown(ctxt)
	if err != nil {
		log.Fatal(err)
	}

	// Wait for chrome to finish
	err = c.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
