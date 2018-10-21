import glob
import json
import os
import re
import requests
import sys
import selenium.webdriver.support.ui as ui

from re import match
from selenium import webdriver
from selenium.common import exceptions
from selenium.webdriver.support import expected_conditions
from selenium.webdriver.common.by import By
from selenium.webdriver.common.desired_capabilities import DesiredCapabilities
from typing import List


class Course:
    Open                 = False

    AcademicLevel        = ""
    CourseCode           = ""
    CourseDescription    = ""
    CourseName           = ""
    DateStart            = ""
    DateEnd              = ""
    Location             = ""
    MeetingInformation   = ""
    Supplies             = ""
    
    Credits              = 0
    SlotsAvailable       = 0
    SlotsCapacity        = 0
    SlotsWaitlist        = 0
    TimeEnd              = 0
    TimeStart            = 0

    ProfessorEmails      = []
    PrereqNonCourse      = []
    RecConcurrentCourses = []
    ReqConcurrentCourses = []

    PrereqCourses        = {}
    InstructionalMethods = {}

    def __init__(self):
        self.Open                 = False

        self.AcademicLevel        = ""
        self.CourseCode           = ""
        self.CourseDescription    = ""
        self.CourseName           = ""
        self.DateStart            = ""
        self.DateEnd              = ""
        self.Location             = ""
        self.MeetingInformation   = ""
        self.Supplies             = ""
        
        self.Credits              = 0
        self.SlotsAvailable       = 0
        self.SlotsCapacity        = 0
        self.SlotsWaitlist        = 0
        self.TimeEnd              = 0
        self.TimeStart            = 0

        self.ProfessorEmails      = []
        self.PrereqNonCourse      = []
        self.RecConcurrentCourses = []
        self.ReqConcurrentCourses = []

        self.PrereqCourses        = {}
        self.InstructionalMethods = {}


class Teacher:
    Email = ""
    Name  = ""
    Phone = ""
    def __init__(self, e, n, p):
        self.Email = e
        self.Name = n
        self.Phone = p


""" JSON Model for Data
    {
        "Teachers": {"{Teacher Email}": {Teacher Object}},
        "TotalTime": "{Total time it took for the data for all semesters to be collected}",
        "{Semester}": {
            "Time": "{Time it took for the data to be collected for the given semester}",
            "Courses": {
                "{Subject}": {"Course Code": [{"{Section Codes}": {Course Object}}]}
            }
        }
    }
"""

# with open(glob.glob("*_config.json")[0]) as fi:
#     config = json.load(fi)
# executable_path = {'executable_path': config['browser']}
# # b = Browser(config['type'], headless=config['headless'], incognito=True, **executable_path, fullscreen=True)
# teachers = {}
# totalData = {"Teachers": {}}

"""Handles WebDriver operations, runs up to {attempts} number of times
as exceptions occur, if number of attempts exceeds {attempts} then the last
exception raise from {func} will be raised

Raises:
    return_err -- Last exception raised"""
def handle_driver_operation(driver: webdriver, attempts: int, func, *args):
    return_err: Exception = Exception("Error with driver operation")

    # Attempt to run the function at most {attempt} times, will return
    # after function excecutes and no exceptions occur
    for attempt in range(attempts):
        try:
            func(*args, driver)
            return
        except Exception as err:
            return_err = err

    # Reraise from the last exception which was raised from {func}
    raise return_err

def scrape_information(username: str, passsord: str, driver: webdriver):
    handle_driver_operation(driver, 10, initToQuery, username, passsord)

    subjects  = get_values_from_select_options("LIST_VAR1_1", driver)
    semesters = get_values_from_select_options("VAR1", driver)
    semester = semesters[1]

    print(subjects, semesters)
    return
    subjectCourses = {}
    for subject in subjects:
        unsuccessfulRun = True
        run = 0
        while unsuccessfulRun:
            run += 1
            print("Run", run, "For", subject)
            try:
                handleExtraExits()
                if b.is_element_present_by_text("Section Selection Results"):
                    unsuccessfulBackOut = True
                    while unsuccessfulBackOut:
                        try:
                            b.find_by_text("Go back").first.click()
                            unsuccessfulBackOut = False
                        except Exception as e:
                            print("Failed to go back: ", e)
                    while b.is_element_not_present_by_id("VAR1", 1):
                        pass
                elif not b.is_element_present_by_text("Search for Class Sections"):
                    b.execute_script("window.location.reload()")
                    while b.is_element_not_present_by_text("Search for Class Sections", 1):
                        pass
                if b.is_element_not_present_by_id("VAR1", 5):
                    b.execute_script("window.location.reload()")
                    attempts, maximumAttempts = 0, 7
                    while b.is_element_not_present_by_text("Search for Class Sections", 1) and attempts < maximumAttempts:
                        attempts += 1
                    continue
                selectDropdown("VAR1", semester)
                selectDropdown("LIST_VAR1_1", subject)
                selectDropdown("VAR6", "DSU")
                b.find_by_id("WASubmit").first.click()
                attempts, maximumAttempts = 0, 7
                while b.is_element_not_present_by_text("Section Selection Results", 1) and attempts < maximumAttempts:
                    attempts += 1
                if b.is_text_present("No classes meeting the search criteria have been found."):
                    unsuccessfulRun = False
                    continue
                elif b.is_element_not_present_by_text("Section Selection Results", 1):
                    raise Exception("Bad result, trying again")
                currentPage, totalPages = "0", "1"
                courses = []
                while currentPage != totalPages:
                    while b.is_element_not_present_by_css('table[summary="Paged Table Navigation Area"]'):
                        pass
                    m = match(r"Page (\d+) of (\d+)", b.find_by_css('table[summary="Paged Table Navigation Area"]').first.find_by_css('td[align="right"]').first.text)
                    currentPage, totalPages = m.groups(1)
                    courses.extend(scrapeTable(b.find_by_css('table[summary="Sections"]')))
                    handleExtraExits()
                    if currentPage != totalPages:
                        b.find_by_css('input[value="NEXT"]').first.click()
                        attempts, maximumAttempts = 0, 12
                        while not b.is_text_present("Page {0} of {1}".format(int(currentPage) + 1, totalPages)) and attempts < maximumAttempts:
                            attempts += 1
                        if not b.is_text_present("Page {0} of {1}".format(int(currentPage) + 1, totalPages)):
                            raise Exception("Bad result, trying again")
                subjectCourses[subject] = {}
                for c in courses:
                    if c.CourseCode in subjectCourses[subject]:
                        subjectCourses[subject][c.CourseCode].append({c.CourseCodeSection: objectToDict(c)})
                    else:
                        subjectCourses[subject][c.CourseCode] = [{c.CourseCodeSection: objectToDict(c)}]
                unsuccessfulRun = False
            except Exception as e:
                print("Trying again after error: \n{0}".format(e))
                while 1:
                    b.execute_script("window.location.reload()")
                    attempts, maximumAttempts = 0, 10
                    while b.is_element_not_present_by_text("Search for Class Sections", 1) and attempts < maximumAttempts:
                        attempts += 1
                    if b.is_element_not_present_by_text("Search for Class Sections", 1):
                        continue
                    break
    totalData[semester] = subjectCourses
    for t in teachers:
        totalData["Teachers"][teachers[t].Email] = objectToDict(teachers[t])

    with open("{0}.json".format(semester), "w") as fi:
        json.dump(totalData, fi)
    b.quit()

"""Navigates to the search query, handles the case where already logged in"""
def initToQuery(username: str, passsord: str, driver: webdriver):
    wait5 = ui.WebDriverWait(driver, 5)
    wait3 = ui.WebDriverWait(driver, 3)

    url = "https://portal.sdbor.edu"

    # Navigate to login
    driver.get(url)

    try:
        # Wait until login page has loaded and log in, exception is thrown if already logged in
        wait3.until(lambda d: d.find_element_by_id("username") is not None)
        driver.find_element_by_id("username").send_keys(username)
        driver.find_element_by_id("password").send_keys(passsord + "\n")
        # Wait until post-login page has loaded
        wait5.until(lambda d: d.find_element_by_id("SearchBox") is not None)
    # Exception will be thrown if already logged in
    except exceptions.TimeoutException:
        pass

    # Navigate to DSU landing page
    driver.get(url + "/dsu-student/Pages/default.aspx")

    # If this isn't the first time on this page (have already logged in and been here once)
    # the dropdown will already be open and thus the "Admission Information" link can be
    # clicked right away
    try:
        # Try waiting until "Admission Information" link is visible, then click it
        element = wait3.until(expected_conditions.element_to_be_clickable((By.XPATH, """//a/span[text()="Admission Information"]""")))
        element.click()
    # Exception will be thrown if this is the first time landing on this page
    except (exceptions.TimeoutException, exceptions.WebDriverException):
        # Wait until DSU landing page has loaded and click on the "WebAdvisor for Prospective Students" link
        wait3.until(lambda d: d.find_element_by_xpath("""//a/h5[text()="WebAdvisor for Prospective Students"]""") is not None)
        driver.find_element_by_xpath("""//a/h5[text()="WebAdvisor for Prospective Students"]""").click()
        element = wait3.until(expected_conditions.element_to_be_clickable((By.XPATH, """//a/span[text()="Admission Information"]""")))
        element.click()

    # Wait until "Search for Class Sections" link is visible, then click it
    element = wait5.until(expected_conditions.element_to_be_clickable((By.XPATH, """//*[text()="Search for Class Sections"]""")))
    element.click()

    # Wait until search page has loaded
    wait5.until(lambda d: d.find_element_by_id("VAR1") is not None)

def get_values_from_select_options(ID: str, driver: webdriver) -> List[str]:
    # Assumes that you're already on the search page
    select_element = driver.find_element_by_id(ID)
    option_tags = select_element.find_elements_by_tag_name("option")

    return [option_tag.get_attribute("value") for option_tag in option_tags if option_tag.text != ""]

def handleExtraExits():
    for e, x in enumerate(b.find_by_css('button[type="button"][class="btn btn-sm btn-danger"]')):
        if e:
            x.click()
    
def objectToDict(obj):
    return {attr: getattr(obj, attr) for attr in type(obj).__dict__ if not attr.startswith("__")}

def scrapeTable(tab):
    courses = []
    for n in range(len(tab.find_by_tag("tr"))):
        if n == 0:
            continue
        course = Course()
        for e in range(len(tab.find_by_tag("tr")[n].find_by_tag("td"))):
            if e == 2:
                course.Open = "Open" in tab.find_by_tag("tr")[n].find_by_tag("td")[e].text
            elif e == 3:
                tab.find_by_tag("tr")[n].find_by_tag("td")[e].find_by_tag("a").first.click()
                while b.is_element_not_present_by_text("Section Information", 1):
                    print("Waiting on Section Information")
                #Start scraping for the bulk of the course data
                course.CourseName = b.find_by_id("VAR1").first.text
                course.CourseCodeSection = b.find_by_id("VAR2").first.text
                nm = match(r"(.+-.+)-.+", course.CourseCodeSection)
                course.CourseCode = nm.groups(1)[0] if nm is not None else course.CourseCodeSection
                course.CourseDescription = b.find_by_id("VAR3").first.text
                course.Location = b.find_by_id("VAR40").first.text
                course.Credits = float(b.find_by_id("VAR4").first.text)
                course.DateStart = b.find_by_id("VAR6").first.text
                course.DateEnd = b.find_by_id("VAR7").first.text
                course.MeetingInformation = b.find_by_id("LIST_VAR12_1").first.text
                course.TimeStart, course.TimeEnd = getTimes(course.MeetingInformation)
                course.Supplies = b.find_by_id("LIST_VAR14_1").first.text

                course.ReqConcurrentCourses = getConcurrent(b.find_by_id("VAR_LIST7_1").first.text)
                course.RecConcurrentCourses = getConcurrent(b.find_by_id("VAR_LIST8_1").first.text)
                course.PrereqCourses = getPrereqs(b.find_by_id("VAR_LIST1_1").first.text)
                PrereqNonCourse = b.find_by_id("VAR_LIST4_1").first.text
                course.PrereqNonCourse = PrereqNonCourse if PrereqNonCourse.lower() != "none" else ""
                teacherDicts = {}
                teacherDicts = scrapeTeachers(b.find_by_css('table[summary="Faculty Contact"]').first)
                for t in teacherDicts:
                    email = t["email"]
                    if not email: 
                        continue
                    name = t["name"]
                    phone = t["phone"]
                    course.InstructionalMethods[email] = t["lecture"]
                    if email:
                        course.ProfessorEmails.append(email)
                        if email not in teachers:
                            teachers[email] = Teacher(email, name, phone)
                        else:
                            if not teachers[email].Email:
                                teachers[email].Email = email
                            if not teachers[email].Name and name:
                                teachers[email].Name = name
                            if not teachers[email].Phone and phone:
                                teachers[email].Phone = phone
                handleExtraExits()
            elif e == 7:
                nm = match(r"(-?\d+) +\/ +(-?\d+) +\/ +(-?\d+)", tab.find_by_tag("tr")[n].find_by_tag("td")[e].text)
                if nm is None:
                    course.SlotsAvailable = course.SlotsCapacity = course.SlotsWaitlist = ""
                else:
                    course.SlotsAvailable = int(nm.group(1))
                    course.SlotsCapacity = int(nm.group(2))
                    course.SlotsWaitlist = int(nm.group(3))
            elif e == 9:
                course.AcademicLevel = tab.find_by_tag("tr")[n].find_by_tag("td")[e].text
        print(course.CourseCode)
        courses.append(course)
    return courses

def scrapeTeachers(tab):
    if tab is None:
        return []
    rl = []
    for e, tr in enumerate(tab.find_by_tag("tr")):
        if e:
            tdict = {}
            for n, td in enumerate(tr.find_by_tag("td")):
                if n == 1:
                    tdict["name"] = td.find_by_tag("div").first.text
                if n == 2:
                    tdict["phone"] = td.find_by_tag("div").first.text
                if n == 4:
                    tdict["email"] = td.find_by_tag("div").first.text
                if n == 5:
                    tdict["lecture"] = td.find_by_tag("div").first.text
            rl.append(tdict)
    return rl

def getPrereqs(s):
    AND = re.findall(r"\+ (\w+ \d+) ?", s) + re.findall(r"^(\w+ \d+) \+", s) + re.findall(r"^(\w+ \d+)$", s)
    OR  = re.findall(r"(\w+ \d+),", s) + re.findall(r"or (\w+ \d+)", s)
    return {"and": replaceSpaces(AND), "or": replaceSpaces(OR)}

def replaceSpaces(l):
    return list(map(lambda a: a.replace(" ", "-"), l))

def getConcurrent(s):
    return replaceSpaces(re.findall(r"(\w+ \d+)", s))

def getTimes(s):
	r = re.search(r"(\d\d\:\d\d)(\w\w) - (\d\d\:\d\d)(\w\w)", s)
	if r is None:
		return (0, 0)
	st, et = int(r.group(1).replace(":", "")), int(r.group(3).replace(":", ""))
	if r.group(2).lower() == "pm" and st < 1200:
		st += 1200
	if r.group(4).lower() == "pm" and et < 1200:
		et += 1200
	return (st, et)

def selectDropdown(ddt, t): #assumes that you're already on the Prospective students search page
    selects = b.find_by_id(ddt).first
    selects.select(t)

def get_config(path: str) -> dict:
    with open(path) as creds_file:
        return json.load(creds_file)

def main():
    creds: dict = get_config("./creds.json")
    username: str = creds["username"]
    password: str = creds["password"]
    remote = False

    if not remote:
        with webdriver.Chrome("C:/RJFiles/Assets/chromedriver_win32/chromedriver.exe") as driver:
            scrape_information(username, password, driver)
    else:
        with webdriver.Remote("http://127.0.0.1:4444/wd/hub", desired_capabilities=DesiredCapabilities.CHROME.copy()) as driver:
            scrape_information(username, password, driver)
            driver.quit()



if __name__ == "__main__":
    main()