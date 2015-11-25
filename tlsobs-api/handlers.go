package main

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/context"

	pg "github.com/mozilla/tls-observatory/database"
	"github.com/mozilla/tls-observatory/logger"
)

func ScanHandler(w http.ResponseWriter, r *http.Request) {
	var (
		status int
		err    error
	)
	defer func() {
		if nil != err {
			http.Error(w, err.Error(), status)
		}
	}()
	log := logger.GetLogger()
	log.WithFields(logrus.Fields{
		"form values": r.Form,
		"headers":     r.Header,
	}).Debug("Received request")

	val, ok := context.GetOk(r, dbKey)
	if !ok {
		log.Error("Could not find db in request context")
		status = http.StatusInternalServerError
	}

	db := val.(*pg.DB)

	domain := r.FormValue("target")
	if validateDomain(domain) {

		scan, err := db.NewScan(domain, -1) //no replay
		if err != nil {
			log.WithFields(logrus.Fields{
				"domain": domain,
				"error":  err.Error(),
			}).Error("Could not create new scan")
			status = http.StatusInternalServerError
		}

		resp := fmt.Sprintf(`{"scan_id":"%d"}`, scan.ID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, resp)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func ResultHandler(w http.ResponseWriter, r *http.Request) {

	var (
		status int
		err    error
	)

	defer func() {
		if nil != err {
			http.Error(w, err.Error(), status)
		}
	}()

	domain := r.FormValue("id")

	if validateDomain(domain) {

		status = http.StatusOK

	} else {
		status = http.StatusBadRequest
		return
	}

}

func CertificateHandler(w http.ResponseWriter, r *http.Request) {

	var (
		status int
		err    error
	)

	defer func() {
		if nil != err {
			http.Error(w, err.Error(), status)
		}
	}()

	domain := r.FormValue("target")

	if validateDomain(domain) {

		//		raw := r.FormValue("raw")

		//		rawCert := false

		//		if raw == "true" {
		//			rawCert = true
		//		}

		status = http.StatusOK

	} else {
		status = http.StatusBadRequest
		return
	}

}

func validateDomain(domain string) bool {

	// TODO
	// Need to validate the domain, in a way,
	// before passing it to the retriever queue

	return true
}