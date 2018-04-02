package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	vault_api "github.com/hashicorp/vault/api"
)

var (
	listenAddress = flag.String("address", "0.0.0.0", "Domain or IP address where server is listening")
	listenPort    = flag.Int("port", 10000, "Webserver listening port")
	nodesListPath = flag.String("nodes-path", "/confs", "")
)

func main() {

	fmt.Println("")

	router := http.NewServeMux()
	router.HandleFunc("/rundeck/nodes", RouteRundeckNodes)

	log.Printf("Listening on port %d ...", *listenPort)

	//err := http.ListenAndServe(":"+listenPort, LogRequests(CheckURL(router)))
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", *listenAddress, *listenPort), router)
	if err != nil {
		log.Fatal(err)
	}
}

func initVault() (*vault_api.Logical, error) {
	config := &vault_api.Config{
		Address: os.Getenv("VAULT_ADDR"),
	}

	client, err := vault_api.NewClient(config)
	if err != nil {
		return nil, err
	}

	client.SetToken(os.Getenv("VAULT_TOKEN"))

	return client.Logical(), nil
}

// routeResponse Used to build response to API requests
func routeResponse(w http.ResponseWriter, httpStatus bool, contents string) {
	if httpStatus {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(403)
	}

	response, _ := json.Marshal(contents)
	fmt.Fprintf(w, "%s", response)
}

// RouteRundeckNodes is an endpoint that exposes all RunDeck nodes by a given
// project (it receives project as a querystring parameter)
// E.g.: `GET /rundeck/nodes?project=my-project`
func RouteRundeckNodes(w http.ResponseWriter, r *http.Request) {
	projectName := r.URL.Query().Get("project")
	vaultPath := "/secret/rundeck/ops/nodes/" + projectName

	log.Printf("Reading RunDeck nodes: project '%s'...", projectName)

	vaultClient, err := initVault()
	if err != nil {
		log.Panic(err)
		return
	}

	nodesList, err := vaultClient.Read(vaultPath)
	if err != nil {
		log.Panic(err)
		routeResponse(w, false, "")
	}

	var contents string
	fmt.Printf("%#v", nodesList.Data)

	routeResponse(w, true, contents)
}
