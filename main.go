package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

var base BaseSetupImpl

type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func defaultGOPATH() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		env = "home"
	}
	if home := os.Getenv(env); home != "" {
		def := filepath.Join(home, "go")
		if filepath.Clean(def) == filepath.Clean(runtime.GOROOT()) {
			// Don't set the default GOPATH to GOROOT,
			// as that will trigger warnings from the go tool.
			return ""
		}
		return def
	}
	return ""
}

func init() {
	// Setup correctly the GOPATH in the environment
	if goPath := os.Getenv("GOPATH"); goPath == "" {
		os.Setenv("GOPATH", defaultGOPATH())
	}

	base = BaseSetupImpl{
		ConfigFile:      "./config.yaml",
		ChannelID:       "suningchannel",
		OrgID:           "Org1",
		ChannelConfig:   "./suning-env/kafka/channel-artifacts/suningchannel.tx",
		ChainCodeID:     "suningCC",
		ConnectEventHub: true,
	}

	log.Printf("Start to Initialize the Fabric SDK")
	if err := base.Initialize(); err != nil {
		fmt.Printf("Initialize: %v", err)
		os.Exit(-1)
	}

}

func OutputJson(w http.ResponseWriter, code int, reason string, data interface{}) {
	out := &Result{code, reason, data}
	b, err := json.Marshal(out)
	if err != nil {
		fmt.Println("OutputJson fail:" + err.Error())
		return
	}

	w.Write(b)
}

func main() {
	http.HandleFunc("/block", Block)
	http.HandleFunc("/createOrg", CreateOrg)
	http.HandleFunc("/submitRecord", SubmitRecord)
	http.HandleFunc("/deleteRecord", DeleteRecord)
	http.HandleFunc("/queryRecord", QueryRecord)
	http.HandleFunc("/queryTransaction", QueryTransaction)
	http.HandleFunc("/issueCredit", IssueCredit)
	http.HandleFunc("/issueCreditToOrg", IssueCreditToOrg)
	http.HandleFunc("/transfer", Transfer)
	http.HandleFunc("/queryOrg", QueryOrg)
	http.HandleFunc("/queryAgency", QueryAgency)

	log.Printf("Start to listen the 8080 port")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
