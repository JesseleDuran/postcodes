package postcodesio

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
)

var client *resty.Client
var api API

func TestMain(m *testing.M) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stdout)

	client = resty.New()
	httpmock.ActivateNonDefault(client.GetClient())
	api = API{HttpClient: client}

	httpmock.RegisterResponder(http.MethodGet, PostcodesURI,
		httpmock.NewStringResponder(http.StatusOK, getServiceResponse("postcodes")))

	code := m.Run()
	os.Exit(code)
}

func TestAPI_PostCode(t *testing.T) {
	pc, err := api.PostCode("51.58647", "-0.112639", "1")
	expected := "N8 7EA"
	if err != nil {
		t.Fatalf("status code should be 200, got %s", err.Error())
	}
	if pc != expected {
		t.Fatalf("postcode not expected, got: %v, expected: %s", pc, expected)
	}
}

func getServiceResponse(state string) string {
	// Read entire file content, giving us little control but
	// making it very simple. No need to close the file.
	content, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s-response.json",
		state))
	if err != nil {
		log.Fatal(err)
	}
	// Convert []byte to string
	return string(content)
}
