package modules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

// Services contains only one attribute called Services, which 
// is an array of the Service struct
type Services struct {
	Services []Service `json:"services"`
}

// Method to convert Services struct to string
func (services *Services) ToStr() string {
    
	var str_val strings.Builder

	str_val.WriteString("[")
	for index, service := range services.Services {

		// Add service
		str_val.WriteString("{" + service.ToStr() + "}")
		// Skip seperator for last instance
		if index == len(services.Services) - 1 { 
			continue 
		}
		// Add seperator
		str_val.WriteString(", ")
	}
	str_val.WriteString("]")

    return str_val.String()
}

// Service is the service which needs to be monitored. It consists of an endpoint
// and a description. Endpoint should be in the form `http(s)://<domain>`. Description is
// optional and may be left empty
type Service struct {
	Endpoint 	string `json:"endpoint"`
	Description string `json:"description"`
}

// Method to convert Service struct to string
func (service *Service) ToStr() string {
    
    // Convert array attrs to string
    endpoint     := &service.Endpoint
    description  := &service.Description

    // Convert all attrs to string
    str_val := fmt.Sprintf("endpoint: \"%s\", description: \"%s\"", 
		*endpoint,
		*description,
	)

    return str_val
}

// Read file and process it into Services
func ReadServices(file_path string) *Services {

	service_file, err := os.Open(file_path)
	if err != nil {
		slog.Error("cannot read services.json file")
		panic(err)
	}
	defer service_file.Close()

	// create services slice
	var services Services

	buf := make([]byte, 1024) // create buffer
	for {
		// read chunk
		read_len, err := service_file.Read(buf) // read 1024 bytes
		buf = bytes.Trim(buf, "\x00")  			// trim null values
		if err != nil && err != io.EOF {
			slog.Error("an error occured while reading the services file")
			panic(err)
		}
		// break if no lines were read
		if read_len == 0 {
			break
		}
		// write chunk
		if err := json.Unmarshal(buf, &services); err != nil {
			slog.Error("an error occured while processing the services file")
			panic(err)
		}
	}

	return &services
}