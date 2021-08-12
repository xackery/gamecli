package steam

import (
	"strings"

	"github.com/agext/levenshtein"
	"github.com/rs/zerolog/log"
)

func find(appID string, name string, acfs []*Acf) *Acf {
	lName := strings.ToLower(name)
	var acf *Acf
	var acfFuzzy *Acf
	var fuzzyCount float64
	for _, a := range acfs {
		if a.AppID == appID {
			acf = a
			break
		}
		if lName == a.Name {
			acf = a
			break
		}
		fc := levenshtein.Match(lName, strings.ToLower(a.Name), nil)
		//fmt.Println(lName, "vs", a.Name, "fc", fc)
		if fuzzyCount > fc {
			continue
		}
		fuzzyCount = fc
		acfFuzzy = a
	}
	if acf == nil {
		if fuzzyCount == 0 {
			return nil
		}
		acf = acfFuzzy
		log.Debug().Msgf("no exact search, fuzzy deduced %s with %0.2f precision", acf.Name, fuzzyCount)
	}
	return acf
}
