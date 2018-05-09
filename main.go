package main

import (
    "encoding/json"
    "log"
    "net/http"
    "net/http/httputil"
    "io"
    "io/ioutil"
    "os"
    "fmt"
)

type CitiesResponse struct {
    Cities []string `json:"cities"`
}

type CityES struct {
  Name string `json:"Name"`
  State string `json:"State"`
}
type CitiesResponseES struct {
    Cities []CityES `json:"cities"`
}

func CityHandlerES(res http.ResponseWriter, req *http.Request) {
    esUrl := os.Getenv("ES_URL")
    esUser := os.Getenv("ES_USER")
    esPass := os.Getenv("ES_PASS")
    if len(esUrl) == 0 {
      io.WriteString(res, "Environment variable ES_URL is not set! Cannot retrieve cities from Elasticsearch.")
      return
    }

    esReq, err := http.NewRequest("GET", esUrl + "/cities/cities/_search", nil)
    if len(esUser) != 0 && len(esPass) != 0 {
      esReq.SetBasicAuth(esUser, esPass)
    }
    //dumpReq, _ := httputil.DumpRequest(esReq, true)
    //fmt.Println(string(dumpReq))
    client := &http.Client{}
    esResp, err := client.Do(esReq)
    //esResp, err := http.Get(esUrl + "/cities/cities/_search")
    if err != nil {
       //handle err
       io.WriteString(res, "An error occurred. "+err.Error()+"\n")
       return
    }
    if esResp.StatusCode != 200 {
       //handle err
       io.Copy(res, esResp.Body)
       return
    }
    defer esResp.Body.Close()

    //Unmarshal json data to a Map
    var data map[string]interface{}
    bodyStr,_ := ioutil.ReadAll(esResp.Body)
    err = json.Unmarshal([]byte(bodyStr), &data)
    if err != nil {
       io.WriteString(res, "An error occurred. "+err.Error()+"\n")
       return
    }
    
    //get es response[hits][hits]
    hits := data["hits"].(map[string]interface{})["hits"].([]interface{})

    //extract _source from hits and make the CitiesResponseES
    var cities []CityES
    for _,hit:= range hits {
      source := hit.(map[string]interface{})["_source"].(map[string]interface{})
//      log.Println("Got source ")
//      log.Println(source)
      
      city := CityES{Name:source["name"].(string), State:source["state"].(string)}
      cities = append(cities, city)
    }    
    writeJsonResponse(res, &CitiesResponseES{Cities: cities})
}

func writeJsonResponse(res http.ResponseWriter, myresp interface{}) {
    data, _ := json.MarshalIndent(myresp, "", "  ")
    res.Header().Set("Content-Type", "application/json; charset=utf-8")
    res.Write(data)
}

func CityHandler(res http.ResponseWriter, req *http.Request) {
    citiesResponse := &CitiesResponse{
        Cities: []string{
            "NYC",
            "LA",
            "Chicago",
            "Philly",
        },
    }
    writeJsonResponse(res, citiesResponse)
}

func logHandler(w http.ResponseWriter, res *http.Request) {
    res.Write(w)
    fmt.Fprintf(w, "RemoteAddr: %s", res.RemoteAddr)
    log.Println("Request Body: %s", res.Body)
}


func main() {
    log.Println("Listening on this host: http://localhost:5005")
    log.Println("Available Endpoints: /cities.json, /es_cities")

    http.HandleFunc("/", logHandler)
    http.HandleFunc("/cities.json", CityHandler)
    http.HandleFunc("/es_cities", CityHandlerES)
    err := http.ListenAndServe(":5005", nil)
    if err != nil {
        log.Fatal("Unable to listen on :5005: ", err)
    }
}
