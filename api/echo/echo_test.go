// generated by authenticaTed v1.0.0 at 2016-12-27 20:21:22.133402047 -0300 BRT

// generated by authenticaTed v1.0.0 at 2016-12-27 20:21:13.327886771 -0300 BRT

// generated by authenticaTed v1.0.0 at 2016-12-27 20:20:38.594604774 -0300 BRT

// generated by authenticaTed v1.0.0 at 2016-12-27 20:20:25.447135665 -0300 BRT

// generated by authenticaTed v1.0.0 at 2016-12-27 20:20:19.739561652 -0300 BRT

// generated by authenticaTed v1.0.0 at 2016-12-27 20:20:14.467113369 -0300 BRT

// generated by authenticaTed v1.0.0 at 2016-12-27 20:20:13.101972987 -0300 BRT

// generated by authenticaTed v1.0.0 at 2016-12-27 20:20:04.768958274 -0300 BRT

// generated by authenticaTed v1.0.0 at 2016-12-27 20:19:04.457914276 -0300 BRT

// generated by authenticaTed v1.0.0 at 2016-12-27 20:18:58.909807063 -0300 BRT

package echo

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	URL = "/api/v1/users"
)

var (
	id     string
	server *httptest.Server
	ex     *httpexpect.Expect
)

func insert(t *testing.T) {
	// create httpexpect instance
	ex = httpexpect.New(t, server.URL)
	httpexpect.NewDebugPrinter(t, true)
}

func TestSetup(t *testing.T) {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Pre(middleware.RemoveTrailingSlash())

	// insert rest api endpoints
	err := SetAPI(&e)
	if err != nil {
		panic(err)
	}

	//s := standard.New("")
	//s.SetHandler(e)
	handler := http.Handler(e)

	// run server using httptest
	server = httptest.NewServer(handler)
}

func TestGetNone(t *testing.T) {
	insert(t)
	expect := ex.GET(URL).Expect()

	if expect.Raw().StatusCode != http.StatusNoContent {
		obj := expect.JSON().Object().Raw()

		if usrs, ok := obj["users"]; ok {
			if len(usrs.(map[string]interface{})) > 0 {
				expect.Status(http.StatusOK)
			}
		}
	}
}

func TestPost(t *testing.T) {
	insert(t)

	u := map[string]interface{}{
		"username": "gopher",
		"password": "wood",
		"email":    "gopher@ufo.gov",
	}

	obj := ex.POST(URL).
		WithJSON(u).
		Expect().
		Status(http.StatusCreated).
		JSON().
		Object()

	obj.Keys().ContainsOnly("success", "user") // contains keys

	// user object
	uobj := obj.Value("user").Object()
	uobj.ValueEqual("username", "gopher") // username == "gopher"

	// id number
	idn := uobj.Value("id").String()
	id = idn.Raw()
}

func TestGetMany(t *testing.T) {
	insert(t)

	obj := ex.GET(URL).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	obj.Keys().ContainsOnly("success", "users")
	obj.ValueEqual("success", true)
	obj.Value("users").Array().Length().Gt(0) // len(users) > 0
}

func TestPut(t *testing.T) {
	insert(t)

	newUsername := "gophersour"

	u := map[string]interface{}{
		"username": newUsername,
	}

	obj := ex.PUT(URL + "/" + id).
		WithJSON(u).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()

	obj.Keys().ContainsOnly("success", "user")

	// user object
	uobj := obj.Value("user").Object()
	uobj.ValueEqual("username", newUsername) // user.username == newUsername
	uobj.Value("id").String().Equal(id)      // user.id == id
}

func TestGetSingle(t *testing.T) {
	insert(t)

	obj := ex.GET(URL + "/" + id).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	obj.Keys().ContainsOnly("success", "user")
	obj.ValueEqual("success", true)

	// user object
	uobj := obj.Value("user").Object()
	uobj.ValueEqual("username", "gophersour") // user.username == "gophersour"
	uobj.Value("id").String().Equal(id)       // user.id == id
}

func TestEnd(t *testing.T) {
	server.Close()
}
