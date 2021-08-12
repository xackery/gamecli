package steam

import (
	"context"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func parseVdf(ctx context.Context, path string) ([]string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("readFile: %w", err)
	}
	libs := []string{}

	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		lineNumber := i + 1
		if lineNumber == 1 && line != `"LibraryFolders"` {
			return libs, fmt.Errorf(`path: %s, invalid vdf, first line is not "LibraryFolders"`, path)
		}
		if lineNumber == 1 {
			continue
		}
		line = strings.TrimSpace(line)
		line = strings.ReplaceAll(line, `"`, "")

		records := strings.Split(line, "\t")
		if len(records) < 3 {
			continue
		}
		_, err = strconv.Atoi(records[0])
		if err != nil {
			continue
		}
		libs = append(libs, fmt.Sprintf("%s/SteamApps", strings.TrimSpace(records[2])))
	}

	return libs, nil
}
