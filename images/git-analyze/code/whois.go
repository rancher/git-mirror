// From https://github.com/aaronhackney/gowhois
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"
)

type ReturnJSON struct {
	WhoisRecord   *Whois         `json:"whoIsRecord"`
	ContactRecord *ContactRecord `json:"ContactRecord,omitempty"`
}

type WhoisList struct {
	WhoisRecords []*Whois `json:"records"`
}

type Whois struct {
	StartAddress     string              `json:"startAddress"`
	EndAddress       string              `json:"endAddress"`
	Handle           string              `json:"handle"`
	Name             string              `json:"name"`
	RegistrationDate string              `json:"registrationDate,omitempty"`
	UpdateDate       string              `json:"updateDate,omitempty"`
	Version          string              `json:"version"`
	OriginASes       []string            `json:"originASes,omitempty"`
	ParentRefUrl     map[string]string   `json:"parentRefUrl,omitempty"`
	ContactRef       map[string]string   `json:"ContactRef,omitempty"`
	Comments         []string            `json:"comments,omitempty"`
	NetBlocks        []map[string]string `json:"netBlocks"`
}

type ContactRecord struct {
	Handle        string   `json:"handle"`
	Name          string   `json:"name"`
	StreetAddress []string `json:"address"`
	City          string   `json:"city"`
	State         string   `json:"state"`
	PostalCode    string   `json:"postalCode"`
	Country       string   `json:"country"`
	ContactType   string   `json:"type"`
	Reference     string   `json:"reference"`
}

func whois(ip string, wg *sync.WaitGroup) (*Whois, error) {
	defer wg.Done()

	url := fmt.Sprintf("https://whois.arin.net/rest/ip/%s", ip)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		log.Warnf("ip: %s whois: %v", ip, resp)
		return nil, errors.New(fmt.Sprintf("response code: %d", resp.StatusCode))
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return unmarshalResponse(data)
	//contactRecord, _ := getContactRecord(whois.ContactRef["url"], data)
	//jsonOutput, _ := whois.generateJson(whois, contactRecord)
	//log.Info(string(jsonOutput))
}

func getContactRecord(url string, data []byte) (*ContactRecord, error) {
	var contactRecord ContactRecord
	var jsonMap map[string]interface{}

	if err := json.Unmarshal(data, &jsonMap); err != nil { // unmarshall into a map of interfaces
		return nil, err
	}

	var prefix interface{}
	if org, exists := jsonMap["org"]; exists {
		contactRecord.ContactType = "org"
		prefix = org
	} else if cust, exists := jsonMap["customer"]; exists {
		contactRecord.ContactType = "customer"
		prefix = cust
	}

	for key, value := range prefix.(map[string]interface{}) {
		switch key {
		case "handle":
			contactRecord.Handle = value.(map[string]interface{})["$"].(string)
		case "name":
			contactRecord.Name = value.(map[string]interface{})["$"].(string)
		case "city":
			contactRecord.City = value.(map[string]interface{})["$"].(string)
		case "iso3166-2":
			contactRecord.State = value.(map[string]interface{})["$"].(string)
		case "postalCode":
			contactRecord.PostalCode = value.(map[string]interface{})["$"].(string)
		case "iso3166-1":
			contactRecord.Country = value.(map[string]interface{})["code2"].(map[string]interface{})["$"].(string)
		case "ref":
			contactRecord.Reference = value.(map[string]interface{})["$"].(string)
		}
	}

	contactRecord.StreetAddress, _ = getAddressLines(prefix)

	return &contactRecord, nil

}

func unmarshalResponse(b []byte) (*Whois, error) {
	var whois Whois
	var jsonMap map[string]interface{}
	var returnNetBlocks []map[string]string

	if err := json.Unmarshal(b, &jsonMap); err != nil { // unmarshall into a map of interfaces
		return nil, err
	}

	for key, value := range jsonMap["net"].(map[string]interface{}) { // Extract the top level json nest []net
		switch key {
		case "startAddress":
			whois.StartAddress = value.(map[string]interface{})["$"].(string)
		case "endAddress":
			whois.EndAddress = value.(map[string]interface{})["$"].(string)
		case "handle":
			whois.Handle = value.(map[string]interface{})["$"].(string)
		case "name":
			whois.Name = value.(map[string]interface{})["$"].(string)
		case "version":
			whois.Version = value.(map[string]interface{})["$"].(string)
		case "orgRef", "customerRef":
			whois.ContactRef = map[string]string{
				"url":    value.(map[string]interface{})["$"].(string),
				"handle": value.(map[string]interface{})["@handle"].(string),
				"name":   value.(map[string]interface{})["@name"].(string),
			}
		case "comment":
			comments, _ := convertToSlice(value.(map[string]interface{})["line"])
			var returnComments []string
			for i := range comments {
				returnComments = append(returnComments, comments[i].(map[string]interface{})["$"].(string))
			}
			whois.Comments = returnComments
		case "originASes":
			originAS, _ := convertToSlice(value.(map[string]interface{})["originAS"])
			var originASes []string
			for i := range originAS {
				originASes = append(originASes, originAS[i].(map[string]interface{})["$"].(string))
			}
			whois.OriginASes = originASes
		case "parentNetRef":
			whois.ParentRefUrl = map[string]string{
				"url":    value.(map[string]interface{})["$"].(string),
				"handle": value.(map[string]interface{})["@handle"].(string),
				"name":   value.(map[string]interface{})["@name"].(string),
			}
		case "registrationDate":
			whois.RegistrationDate = value.(map[string]interface{})["$"].(string)
		case "updateDate":
			whois.UpdateDate = value.(map[string]interface{})["$"].(string)
		case "netBlocks":
			netBlockList, err := convertToSlice(value.(map[string]interface{})["netBlock"])
			if err != nil {
				fmt.Println("ERROR: ", err)
			}
			for i := range netBlockList {
				description := netBlockList[i].(map[string]interface{})["description"].(map[string]interface{})["$"].(string)
				endAddress := netBlockList[i].(map[string]interface{})["endAddress"].(map[string]interface{})["$"].(string)
				startAddress := netBlockList[i].(map[string]interface{})["startAddress"].(map[string]interface{})["$"].(string)
				blockType := netBlockList[i].(map[string]interface{})["type"].(map[string]interface{})["$"].(string)
				cidrLength := netBlockList[i].(map[string]interface{})["cidrLength"].(map[string]interface{})["$"].(string)
				netBlockObject := map[string]string{
					"description":  description,
					"startAddress": startAddress,
					"endAddress":   endAddress,
					"cidrLength":   cidrLength,
					"type":         blockType,
				}
				returnNetBlocks = append(returnNetBlocks, netBlockObject)
			}
			whois.NetBlocks = returnNetBlocks
		}
	}

	return &whois, nil
}

func (*Whois) generateJson(whoisRecord *Whois, contactRecord *ContactRecord) ([]byte, error) {
	var returnJson ReturnJSON
	returnJson.WhoisRecord = whoisRecord
	returnJson.ContactRecord = contactRecord

	jsonOutput, err := json.MarshalIndent(&returnJson, "", "\t")

	if err != nil {
		fmt.Println("err:", err.Error())
		return nil, err
	}

	return jsonOutput, nil
}

func convertToSlice(object interface{}) ([]interface{}, error) {
	switch v := object.(type) {
	case []interface{}:
		return object.([]interface{}), nil
	case interface{}:
		var returnInterfaceArray []interface{} = make([]interface{}, 1)
		returnInterfaceArray[0] = object
		return returnInterfaceArray, nil
	default:
		fmt.Println(v)
		return nil, nil
	}
}

func getAddressLines(rawJson interface{}) ([]string, error) {
	var streetAddress []string

	address, err := convertToSlice(rawJson.(map[string]interface{})["streetAddress"].(map[string]interface{})["line"])
	if err != nil {
		return nil, err
	}

	for line := range address {
		//fmt.Printf("\nADDRESS LINE: %+v\n", address[line])
		streetAddress = append(streetAddress, address[line].(map[string]interface{})["$"].(string))
	}
	return streetAddress, nil
}
