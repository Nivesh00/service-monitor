package modules

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"

)

// Return type for goroutines
type result struct {
	SvcResp *ServiceResponse
	Err *error
}


// Wrapper for function ProbeService, which loop over an array of structs
// and send HTTP requests to the endpoint attribute
func ProbeAllServices(services *Services) *ServicesResponse {

	var servicesResponse ServicesResponse

	// Loop array and probe each endpoint
	for _, service := range services.Services {
		
		// Create cannel for execution
		execChannel := make(chan *result, 1)
		// Exec goroutine
		go ProbeService(&service, &execChannel)
		resultObj := <-execChannel
		// If error, continue
		if resultObj.Err != nil {
			slog.Error("could not probe endpoint " + service.Endpoint, slog.Any("error", *resultObj.Err))
			continue
		}

		// Log success
		slog.Info("successfully probed endpoint " + service.Endpoint)

		// Lock before writing and unlock again afterwards
		servicesResponse.mu.Lock()
		servicesResponse.ServicesResponse = append(servicesResponse.ServicesResponse, *resultObj.SvcResp)
		servicesResponse.mu.Unlock()
	}

	return &servicesResponse
}

// Send HTTP request to an endpoint
func ProbeService(service *Service, execChannel *chan *result) {
	fmt.Println("Start ", service.Endpoint)

	var serviceResp ServiceResponse
	var probeResult result

	// Probe service
    resp, err := http.Get(service.Endpoint)
    if err != nil {
		// Set Err of result to err ptr
		probeResult.Err = &err
        *execChannel <- &probeResult
		return
    }
    defer resp.Body.Close()

	// Assign attributes to ServiceResponse struct
	serviceResp.ContentLength = int(resp.ContentLength)
	serviceResp.ContentType   = strings.Join(resp.Header["Content-Type"], ", ")
	serviceResp.RequestUrl    = service.Endpoint
	serviceResp.StatusCode    = resp.StatusCode

	// TODO: parse resp.Body
	//

	// Set response for service
	probeResult.SvcResp = &serviceResp

	*execChannel <- &probeResult
}

type ServicesResponse struct {
	mu sync.Mutex
	ServicesResponse []ServiceResponse `json:"servicesResponse"`
}

// Method to convert ServicesResponse struct to string
func (servicesResponse *ServicesResponse) ToStr() string {
    
	var str_val strings.Builder

	str_val.WriteString("[")
	for index, serviceResponse := range servicesResponse.ServicesResponse {

		// Add service
		str_val.WriteString("{" + serviceResponse.ToStr() + "}")
		// Skip seperator for last instance
		if index == len(servicesResponse.ServicesResponse) - 1 { 
			continue 
		}
		// Add seperator
		str_val.WriteString(", ")
	}
	str_val.WriteString("]")

    return str_val.String()
}

type ServiceResponse struct {
	ContentLength int    `json:"contentLength"`
	ContentType   string `json:"contentType"`
	RequestUrl    string `json:"requestUrl"`
	StatusCode    int 	 `json:"status"`
}

// Method to convert ServiceResponse struct to string
func (serviceResponse *ServiceResponse) ToStr() string {
    
    // Convert array attrs to string
	ContentLength := &serviceResponse.ContentLength
	ContentType   := &serviceResponse.ContentType
	RequestUrl 	  := &serviceResponse.RequestUrl
	StatusCode 	  := &serviceResponse.StatusCode

    // Convert all attrs to string
    str_val := fmt.Sprintf("Content-Length: \"%d\", Content-Type: \"%s\", Request-URL: \"%s\", Status-Code: \"%d\"", 
		*ContentLength,
		*ContentType,
		*RequestUrl,
		*StatusCode,
	)

    return str_val
}