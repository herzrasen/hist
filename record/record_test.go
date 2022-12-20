package record

import (
	"fmt"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"testing"
)

func TestFuzzy(t *testing.T) {
	ranks := fuzzy.RankFind("saml", []string{"ll", "saml2aws login -p prod-admin -a prod --skip-prompt"})
	fmt.Printf("%+v", ranks)
}
