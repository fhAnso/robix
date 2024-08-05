package lib

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func FormatUrl(url string) string {
	re := regexp.MustCompile(`^https?://`)
	remLastSlash := strings.TrimSuffix(url, "/")
	remUrlProto := re.ReplaceAllString(remLastSlash, "")
	return remUrlProto + "-robots.txt"
}

func IsRoot(entry string) bool {
	return entry == "Allow: /" || entry == "Disallow: /"
}

func ProcessEntry(url string, entry string) string {
	if entry == "" {
		fmt.Println()
		return ""
	}
	if strings.HasPrefix(entry, "#") {
		fmt.Println(entry)
		return ""
	}
	if strings.Contains(entry, "Allow") || strings.Contains(entry, "Disallow") {
		if !strings.ContainsAny(entry, "$?*=") && !IsRoot(entry) {
			return url + entry
		} else {
			fmt.Println(entry + " # check manually")
			return ""
		}
	}
	fmt.Println(entry)
	return ""
}

func ReadRobots(url string) error {
	var e string
	target := Target{
		Url: url,
	}
	session := target.SessionInit()
	robots, err := target.GetRobotsFile(session)
	if err != nil {
		e = fmt.Sprintf("unable to get contents of robots.txt from %s: %s", url, err.Error())
		return errors.New(e)
	}
	output := FormatUrl(url)
	// Write content to file: format: domain.xyz-robots.txt
	if err := os.WriteFile(output, []byte(robots), 0666); err != nil {
		e = fmt.Sprintf("failed to create output file \"%s\": %s", output, err.Error())
		return errors.New(e)
	}
	fileStream, err := os.Open(output)
	if err != nil {
		return errors.New("could not open output file: " + err.Error())
	}
	defer fileStream.Close()
	scanner := bufio.NewScanner(fileStream)
	url = strings.TrimSuffix(url, "/")
	for scanner.Scan() {
		entry := scanner.Text()
		requestCode := ProcessEntry(url, entry)
		if requestCode != "" {
			statusCode := target.HttpStatusCode(session)
			result := fmt.Sprintf("%s # http status code: %d", entry, statusCode)
			fmt.Println(result)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
