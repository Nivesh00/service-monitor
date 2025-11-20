package modules

import(
	// "bufio"
	// "errors"
	"fmt"
	"net/http"
	"log/slog"
	"strings"
)

// Wrapper for function ProbeService, which loop over an array of structs
// and send HTTP requests to the endpoint attribute
func ProbeAllServices(services *Services) *ServicesResponse {

	var services_resp ServicesResponse

	// Loop array and probe each endpoint
	for _, service := range services.Services {
		service_resp, err := ProbeService(&service)
		// Log error
		if err != nil {
			slog.Error("could not probe endpoint " + service.Endpoint, slog.Any("error", *err))
			continue
		}
		// Log success
		slog.Info("successfully probed endpoint " + service.Endpoint)
		services_resp.ServicesResponse = append(services_resp.ServicesResponse, *service_resp)
	}

	return &services_resp
}

// Send HTTP request to an endpoint
func ProbeService(service *Service) (*ServiceResponse, *error) {

	var service_resp ServiceResponse

	// Probe service
    resp, err := http.Get(service.Endpoint)
    if err != nil {
        return &service_resp, &err
    }
    defer resp.Body.Close()

	// Assign attributes to ServiceResponse struct
	service_resp.ContentLength = int(resp.ContentLength)
	service_resp.ContentType = strings.Join(resp.Header["Content-Type"], ", ")
	service_resp.RequestUrl = service.Endpoint
	service_resp.StatusCode = resp.StatusCode

	// TODO: parse resp.Body
	//

	return &service_resp, nil
}

type ServicesResponse struct {
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
	ContentLength int  `json:"contentLength"`
	ContentType string `json:"contentType"`
	RequestUrl string  `json:"requestUrl"`
	StatusCode int 	   `json:"status"`
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