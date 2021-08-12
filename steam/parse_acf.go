package steam

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var (
	mu     sync.RWMutex
	states = map[int]string{
		0:       "Invalid",        //0
		1 << 0:  "Uninstalled",    //1
		1 << 1:  "UpdateRequired", //2
		1 << 2:  "FullyInstalled", //4
		1 << 3:  "Encrypted",      //8
		1 << 4:  "Locked",         //16
		1 << 5:  "FilesMissing",   //32
		1 << 6:  "AppRunning",     //64
		1 << 7:  "FilesCorrupt",   //128
		1 << 8:  "UpdateRunning",  //256
		1 << 9:  "UpdatePaused",   //512
		1 << 10: "UpdateStarted",  //1024
		1 << 11: "Uninstalling",   //2048
		1 << 12: "BackupRunning",  //4096
		1 << 13: "Reconfiguring",  //65536
		1 << 14: "Validating",     //131072
		1 << 15: "AddingFiles",    //262144
		1 << 16: "Preallocating",  //524288
		1 << 17: "Downloading",    //1048576
		1 << 18: "Staging",        //2097152
		1 << 19: "Committing",     //4194304
		1 << 20: "UpdateStopping", //8388608
	}
)

// Acf is a metadata type for storing game information
type Acf struct {
	AppID                           string
	Universe                        string
	Name                            string
	StateFlags                      string
	Installdir                      string
	LastUpdated                     string
	UpdateResult                    string
	SizeOnDisk                      string
	Buildid                         string
	LastOwner                       string
	BytesToDownload                 string
	BytesDownloaded                 string
	BytesToStage                    string
	BytesStaged                     string
	AutoUpdateBehavior              string
	AllowOtherDownloadsWhileRunning string
	ScheduledAutoUpdate             string
	LauncherPath                    string
	StateName                       string
	GamePath                        string
}

func parseAcfDir(ctx context.Context, path string) ([]*Acf, error) {
	acfs := []*Acf{}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("readDir: %w", err)
	}
	for _, info := range files {
		name := info.Name()
		if info.IsDir() {
			continue
		}
		if filepath.Ext(name) != ".acf" {
			continue
		}

		//log.Debug().Msgf("parsing %s", name)
		a, err := parseAcf(ctx, fmt.Sprintf("%s/%s", path, name))
		if err != nil {
			return nil, fmt.Errorf("parseAcf: %w", err)
		}

		acfs = append(acfs, a)
	}
	if err != nil {
		return nil, fmt.Errorf("walkDir: %w", err)
	}
	return acfs, nil
}

func parseAcf(ctx context.Context, path string) (*Acf, error) {
	a := &Acf{}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("readFile: %w", err)
	}

	var entries = []struct {
		pattern string
		value   *string
	}{
		{pattern: `"appid"`, value: &a.AppID},
		{pattern: `"Universe"`, value: &a.Universe},
		{pattern: `"LauncherPath"`, value: &a.LauncherPath},
		{pattern: `"name"`, value: &a.Name},
		{pattern: `"StateFlags"`, value: &a.StateFlags},
		{pattern: `"installdir"`, value: &a.Installdir},
		{pattern: `"LastUpdated"`, value: &a.LastUpdated},
		{pattern: `"UpdateResult"`, value: &a.UpdateResult},
		{pattern: `"SizeOnDisk"`, value: &a.SizeOnDisk},
		{pattern: `"buildid"`, value: &a.Buildid},
		{pattern: `"LastOwner"`, value: &a.LastOwner},
		{pattern: `"BytesToDownload"`, value: &a.BytesToDownload},
		{pattern: `"BytesDownloaded"`, value: &a.BytesDownloaded},
		{pattern: `"BytesToStage"`, value: &a.BytesToStage},
		{pattern: `"BytesStaged"`, value: &a.BytesStaged},
		{pattern: `"AutoUpdateBehavior"`, value: &a.AutoUpdateBehavior},
		{pattern: `"AllowOtherDownloadsWhileRunning"`, value: &a.AllowOtherDownloadsWhileRunning},
		{pattern: `"ScheduledAutoUpdate"`, value: &a.ScheduledAutoUpdate},
	}
	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		lineNumber := i + 1
		if lineNumber == 1 && line != `"AppState"` {
			return nil, fmt.Errorf(`path: %s, invalid acf, first line is not "AppState"`, path)
		}
		line = strings.TrimSpace(line)
		for _, e := range entries {
			if !strings.HasPrefix(line, e.pattern) {
				continue
			}

			line = strings.TrimSpace(strings.TrimPrefix(line, e.pattern))
			line = strings.ReplaceAll(line, `"`, "")
			*e.value = line
		}
	}

	flag, err := strconv.Atoi(a.StateFlags)
	if err != nil {
		return nil, fmt.Errorf("atoid %s: %w", a.StateFlags, err)
	}
	mu.RLock()
	if flag == 0 {
		a.StateName = states[0]
	}
	for i := 1; i <= 21; i++ {
		if flag&i != i {
			continue
		}
		entry := states[i]
		if entry == "" {
			continue
		}
		if len(a.StateName) > 0 {
			a.StateName += ", " + entry
			continue
		}
		a.StateName = entry
	}
	mu.RUnlock()
	a.GamePath = fmt.Sprintf("%s/common/%s", filepath.Dir(path), a.Installdir)

	return a, nil
}
