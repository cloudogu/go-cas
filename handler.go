package cas

import (
	"bytes"
	"fmt"
	"github.com/golang/glog"
	"io"
	"net/http"
)

const (
	sessionCookieName = "_cas_session"
)

// clientHandler handles CAS Protocol HTTP requests
type clientHandler struct {
	c               *Client
	h               http.Handler
	isLogoutRequest func(r *http.Request) bool
}

// ServeHTTP handles HTTP requests, processes CAS requests
// and passes requests up to its child http.Handler.
func (ch *clientHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if glog.V(2) {
		glog.Infof("cas: handling %v request for %v", r.Method, r.URL)
	}

	fmt.Println("================\n\n\n")
	fmt.Println(getCookie(w, r))
	fmt.Println("=!=!=!=!=!")
	fmt.Println(ch.c.sessions)
	fmt.Println("================\n\n\n")

	setClient(r, ch.c)

	if ch.isSingleLogoutRequest(r) {
		ch.performSingleLogout(w, r)
		ch.c.RedirectToLogout(w, r)
		return
	}

	ch.c.getSession(w, r)
	ch.h.ServeHTTP(w, r)
	return
}

// isSingleLogoutRequest determines if the http.Request is a CAS Single Logout Request.
//
// The rules for a SLO request are, HTTP POST urlencoded form with a logoutRequest parameter.
func (ch *clientHandler) isSingleLogoutRequest(r *http.Request) bool {
	content, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	defer func() {
		r.Body = io.NopCloser(bytes.NewBuffer(content))
	}()

	if r.Method != "POST" {
		return false
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/x-www-form-urlencoded" {
		return false
	}

	if v := r.FormValue("logoutRequest"); v == "" {
		return false
	}

	return true
}

// performSingleLogout processes a single logout request
func (ch *clientHandler) performSingleLogout(w http.ResponseWriter, r *http.Request) {
	rawXML := r.FormValue("logoutRequest")
	logoutRequest, err := parseLogoutRequest([]byte(rawXML))

	if err != nil {
		glog.Info(err.Error())
		return
	}

	if err := ch.c.tickets.Delete(logoutRequest.SessionIndex); err != nil {
		glog.Info(err.Error())
		return
	}

	ch.c.deleteSession(logoutRequest.SessionIndex)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}
