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

// HTTPResponse Structure used to define response object of every route request
type HTTPResponse struct {
	Status  bool   `json:"status"`
	Content string `json:"content"`
}

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
	res := new(HTTPResponse)

	if httpStatus {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(403)
	}

	res.Status = httpStatus
	res.Content = contents
	response, _ := json.Marshal(res)
	fmt.Fprintf(w, "%s", response)
}

// RouteRundeckNodes ...
func RouteRundeckNodes(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s", r.URL.Query().Get("project"))

	log.Print("Reading RunDeck nodes: project ...")

	projectName := "/secret/rundeck/ops/nodes"

	vaultClient, err := initVault()
	if err != nil {
		return
	}

	nodesList, err := vaultClient.Read(projectName)
	if err != nil {
		routeResponse(w, false, "")
	}
	routeResponse(w, true, fmt.Sprintf("%#v", nodesList.Data))
}
