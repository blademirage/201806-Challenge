package main

import (
	"testing"

	"github.com/coreos/go-semver/semver"
)

func stringToVersionSlice(stringSlice []string) []*semver.Version {
	versionSlice := make([]*semver.Version, len(stringSlice))
	for i, versionString := range stringSlice {
		versionSlice[i] = semver.New(versionString)
	}
	return versionSlice
}

func versionToStringSlice(versionSlice []*semver.Version) []string {
	stringSlice := make([]string, len(versionSlice))
	for i, version := range versionSlice {
		stringSlice[i] = version.String()
	}
	return stringSlice
}

func TestLatestVersions(t *testing.T) {
	testCases := []struct {
		versionSlice   []string
		expectedResult []string
		minVersion     *semver.Version
	}{
		{
			versionSlice:   []string{"1.8.11", "1.9.6", "1.10.1", "1.9.5", "1.8.10", "1.10.0", "1.7.14", "1.8.9", "1.9.5"},
			expectedResult: []string{"1.10.1", "1.9.6", "1.8.11"},
			minVersion:     semver.New("1.8.0"),
		},
		{
			versionSlice:   []string{"1.8.11", "1.9.6", "1.10.1", "1.9.5", "1.8.10", "1.10.0", "1.7.14", "1.8.9", "1.9.5"},
			expectedResult: []string{"1.10.1", "1.9.6"},
			minVersion:     semver.New("1.8.12"),
		},
		{
			versionSlice:   []string{"1.10.1", "1.9.5", "1.8.10", "1.10.0", "1.7.14", "1.8.9", "1.9.5"},
			expectedResult: []string{"1.10.1"},
			minVersion:     semver.New("1.10.0"),
		},
		{
			versionSlice:   []string{"2.2.1", "2.2.0"},
			expectedResult: []string{"2.2.1"},
			minVersion:     semver.New("2.2.1"),
		},
		// New test cases
		{
			versionSlice:   []string{"1.8.11", "1.7.6", "1.10.1", "1.9.6", "1.6.10", "1.10.0", "1.7.14", "1.8.9", "1.8.11", "1.9.5"},
			expectedResult: []string{"1.10.1", "1.9.6", "1.8.11", "1.7.14", "1.6.10"},
			minVersion:     semver.New("1.5.0"),
		},
		{
			versionSlice:   []string{"1.11.1-beta.2", "1.11.0", "1.10.1", "1.9.5", "1.8.10", "1.10.0", "1.7.14", "1.8.9", "1.9.5"},
			expectedResult: []string{"1.11.1-beta.2", "1.10.1"},
			minVersion:     semver.New("1.10.0"),
		},
		{
			versionSlice:   []string{"10.4.1", "9.11.2", "8.11.3", "6.14.3", "10.4.0", "10.3.0", "10.2.1", "10.2.0", "8.11.2", "10.1.0", "6.14.2", "10.0.0", "9.11.1", "9.11.0", "9.10.1", "8.11.1", "6.14.1", "4.9.1", "9.10.0", "8.11.0", "6.14.0", "4.9.0", "9.9.0", "9.8.0", "8.10.0", "6.13.1", "9.7.1", "9.7.0", "9.6.1", "9.6.0", "6.13.0", "9.5.0", "9.4.0", "8.9.4", "6.12.3", "9.3.0", "9.2.1", "8.9.3", "6.12.2", "4.8.7", "8.9.2", "6.12.1", "9.2.0", "9.1.0", "8.9.1", "6.12.0", "4.8.6", "9.0.0", "8.9.0", "8.8.1", "8.8.0", "6.11.5", "4.8.5", "8.7.0", "6.11.4", "8.6.0", "8.5.0", "6.11.3"},
			expectedResult: []string{"10.4.1", "10.3.0", "10.2.1", "10.1.0", "10.0.0", "9.11.2", "9.10.1", "9.9.0", "9.8.0", "9.7.1", "9.6.1", "9.5.0", "9.4.0", "9.3.0", "9.2.1", "9.1.0", "9.0.0", "8.11.3", "8.10.0", "8.9.4", "8.8.1", "8.7.0", "8.6.0", "8.5.0"},
			minVersion:     semver.New("8.0.0"),
		},
	}

	test := func(versionData []string, expectedResult []string, minVersion *semver.Version) {
		stringSlice := versionToStringSlice(LatestVersions(stringToVersionSlice(versionData), minVersion))
		for i, versionString := range stringSlice {
			if versionString != expectedResult[i] {
				t.Errorf("Received %s, expected %s", stringSlice, expectedResult)
				return
			}
		}
	}

	for _, testValues := range testCases {
		test(testValues.versionSlice, testValues.expectedResult, testValues.minVersion)
	}
}