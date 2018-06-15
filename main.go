package main

import (
	"fmt"
	"os"
	"log"
	"bufio"
	"strings"
	"context"
	"github.com/coreos/go-semver/semver"
	"github.com/google/go-github/github"
)

type Source struct {
	owner      string
	repo       string
	minVersion string
}

func main() {
	// Command line arguments
	path := GetFileFromCmdArgs()
	targets := ProcessInputFile(path)

	// Github
	client := github.NewClient(nil)
	ctx := context.Background()
	// Increased record returned per page in order
	// to reduce required number of API calls when
	// the the min version is very old
	opt := &github.ListOptions{PerPage: 100}

	// Retrieve the release information of each repository
	for _, target := range targets {
		allReleases := GetReleases(client, ctx, opt, target.owner, target.repo)
		minVersion := semver.New(target.minVersion)
		versionSlice := LatestVersions(allReleases, minVersion)
		fmt.Printf("latest versions of %s/%s: %s\n", target.owner, target.repo, versionSlice)
	}
}

// Function to process and generate desired release/version information of specific repository
func LatestVersions(releases []*semver.Version, minVersion *semver.Version) []*semver.Version {
	// Sort the release information in ascending order
	semver.Sort(releases)
	var versionSlice []*semver.Version
	maxVersion := releases[len(releases)-1]
	// Add the highest released version to the head of the result list
	versionSlice = append(versionSlice, maxVersion)
	idxMajor := maxVersion.Major
	idxMinor := maxVersion.Minor
	// Loop the entire list of released versions
	for i := len(releases) - 1; i >= 0; i-- {
		version := releases[i]
		curMajor := version.Major
		curMinor := version.Minor
		// Save all highest patch version to the list
		if curMajor == idxMajor && curMinor < idxMinor && version.Compare(*minVersion) >= 0 {
			idxMajor = version.Major
			idxMinor = version.Minor
			versionSlice = append(versionSlice, version)
		}
	}
	return versionSlice
}

// Function to retrieve release information of a specific repository
func GetReleases(client *github.Client, ctx context.Context, opt *github.ListOptions, owner string, repo string) []*semver.Version {
	releases, _, err := client.Repositories.ListReleases(ctx, owner, repo, opt)
	if err != nil {
		fmt.Printf("[Error] Failed to retrieve release of %s/%s \n Error message: %s", owner, repo, err)
		return nil
	}
	var allReleases []*semver.Version
	for _, release := range releases {
		versionString := *release.TagName
		if versionString[0] == 'v' {
			versionString = versionString[1:]
		}
		ver, verErr := semver.NewVersion(versionString)
		// Add versions that can be parsed into the list
		if verErr == nil {
			allReleases = append(allReleases, ver)
		}
	}
	return allReleases
}

// Function to process the file specified in the command line argument
func ProcessInputFile(path string) []Source {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatal("[Fatal] File " + path + " does not exist. Exiting...")
	}
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	// Defer close the file
	defer file.Close()
	var targets []Source
	scanner := bufio.NewScanner(file)
	firstLine := true
	for scanner.Scan() {
		// Skip the first line of the file
		if firstLine == false {
			line := scanner.Text()
			lineArray := strings.Split(line, ",")
			repo := strings.Split(lineArray[0], "/")
			target := Source{repo[0], repo[1], lineArray[1]}
			targets = append(targets, target)
		} else {
			firstLine = false
		}

	}
	if err := scanner.Err(); err != nil {
		// Use panic here to ensure all deferred function can be executed
		// and to release all opened resources
		panic(err)
	}
	return targets
}

// Function to check command line arguments
func GetFileFromCmdArgs() string {
	if len(os.Args) < 2 {
		log.Fatal("[Fatal] File name not specified. Exiting...")
	}
	path := os.Args[1]
	return path
}