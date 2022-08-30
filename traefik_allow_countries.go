// Package "allow countries" is a Traefik plugin to allow requests based on their country of origin and block everything else.
// Thanks to https://github.com/PascalMinder/GeoBlock for the initial idea.
package traefik_allow_countries

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

/**********************************
 *        Define constants        *
 **********************************/

const (
	HoursInMillis         = 60 * 60 * 1000
	PrivateIpAddressesTag = "private"
)

/**********************************
 *          Define types          *
 **********************************/

type traefik_allow_countries struct {
	next               http.Handler
	name               string
	allowedIPRanges    []*IpRangesTimestamp
	allowLocalRequests bool
	cidrFileFolder     string
	cidrFileUpdate     bool
	countries          []string
	logAllowedRequests bool
	logDetails         bool
	logLocalRequests   bool
}

type Config struct {
	AllowLocalRequests bool     `yaml:"allowLocalRequests"`
	CidrFileFolder     string   `yaml:"cidrFileFolder"`
	CidrFileUpdate     bool     `yaml:"cidrFileUpdate"`
	Countries          []string `yaml:"countries,omitempty"`
	LogAllowedRequests bool     `yaml:"logAllowedRequests"`
	LogDetails         bool     `yaml:"logDetails"`
	LogLocalRequests   bool     `yaml:"logLocalRequests"`
}

type IpRangesTimestamp struct {
	Country   string
	IpRanges  []*net.IPNet
	Timestamp time.Time
}

/**********************************
 * Define traefik related methods *
 **********************************/

// CreateConfig creates the default plugin configuration.
// Returns a empty config object.
func CreateConfig() *Config {
	return &Config{
		AllowLocalRequests: false,
		CidrFileUpdate:     true,
		LogAllowedRequests: false,
		LogDetails:         true,
		LogLocalRequests:   true,
	}
}

// New creates a new plugin.
// Returns the configured AllowCountries plugin object.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.CidrFileFolder) == 0 {
		return nil, fmt.Errorf("the CIDR file folder is empty")
	}

	if len(config.Countries) == 0 {
		return nil, fmt.Errorf("the list of allowed countries is empty")
	}

	// log.Println("Allow local IPs: ", config.AllowLocalRequests)
	// log.Println("Allowed countries: ", config.Countries)
	// log.Println("CIDR file folder: ", config.CidrFileFolder)
	// log.Println("CIDR file update: ", config.CidrFileUpdate)
	// log.Println("Log allowed requests: ", config.LogAllowedRequests)
	// log.Println("Log details: ", config.LogDetails)
	// log.Println("Log local requests: ", config.LogLocalRequests)

	return &traefik_allow_countries{
		next:               next,
		name:               name,
		allowedIPRanges:    InitializeAllowedIPRanges(config.Countries, config.CidrFileFolder),
		allowLocalRequests: config.AllowLocalRequests,
		cidrFileFolder:     config.CidrFileFolder,
		cidrFileUpdate:     config.CidrFileUpdate,
		countries:          config.Countries,
		logAllowedRequests: config.LogAllowedRequests,
		logDetails:         config.LogDetails,
		logLocalRequests:   config.LogLocalRequests,
	}, nil
}

// This method is the middleware called during runtime and handling middleware actions.
func (allowCountries *traefik_allow_countries) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	// Collect the IP addresses from the HTTP request.
	requestIPAddressList, err := allowCountries.CollectRemoteIP(request)
	if err != nil {
		// If one of the ip addresses could not be parsed, return status forbidden.
		log.Println(err)
		responseWriter.WriteHeader(http.StatusForbidden)
		return
	}

	// Find the block of private ip addresses
	var privateIPBlocks []*net.IPNet
	for index := range allowCountries.allowedIPRanges {
		if allowCountries.allowedIPRanges[index].Country == PrivateIpAddressesTag {
			privateIPBlocks = allowCountries.allowedIPRanges[index].IpRanges
			break
		}
	}

	// Iterate over the addresses.
	for _, ipAddress := range requestIPAddressList {
		// Check whether the current IP address is a private one.
		isPrivateIp := IsPrivateIP(*ipAddress, privateIPBlocks)
		if isPrivateIp {
			// If local requests are allowed everything is fine.
			if allowCountries.allowLocalRequests {
				if allowCountries.logLocalRequests {
					log.Println("Local IP allowed: ", ipAddress, request.URL)
				}
				allowCountries.next.ServeHTTP(responseWriter, request)
			} else {
				// If local requests are prohibited write StatusForbidden.
				if allowCountries.logLocalRequests {
					log.Println("Local IP denied: ", ipAddress, request.URL)
				}
				responseWriter.WriteHeader(http.StatusForbidden)
			}

			// We handled a private IP address here, so we can safely return here.
			return
		}

		// Check country ip ranges.
		var found bool = false
		for index := range allowCountries.allowedIPRanges {
			if allowCountries.allowedIPRanges[index].Country != PrivateIpAddressesTag {
				// Check whether an update is needed.
				if allowCountries.cidrFileUpdate && time.Since(allowCountries.allowedIPRanges[index].Timestamp).Hours() >= HoursInMillis {
					allowCountries.allowedIPRanges[index] = CreateCountryIPBlocks(allowCountries.allowedIPRanges[index].Country, allowCountries.cidrFileFolder)
				}

				found = IsIpInList(*ipAddress, allowCountries.allowedIPRanges[index].IpRanges)
				// If IP was found we can break the current cycle.
				if found {
					if allowCountries.logAllowedRequests {
						log.Printf("%s: Request (%s) allowed for IP [%s]", allowCountries.name, request.URL, ipAddress)
					}
					break
				}
			}
		}

		if !found {
			log.Printf("%s: Request (%s) denied for IP [%s]", allowCountries.name, request.URL, ipAddress)
			responseWriter.WriteHeader(http.StatusForbidden)

			return
		}
	}

	allowCountries.next.ServeHTTP(responseWriter, request)
}

/**********************************
 *         Private methods        *
 **********************************/

// This method collects the remote IP address.
// It tries to parse the IP from the HTTP request.
// Returns the parsed IP and no error on success, otherwise the so far generated list and an error.
func (allowCountries *traefik_allow_countries) CollectRemoteIP(request *http.Request) ([]*net.IP, error) {
	var ipList []*net.IP

	// Helper method to split a string at char ','
	splitFn := func(c rune) bool {
		return c == ','
	}

	// Try to parse from header "X-Forwarded-For"
	xForwardedForValue := request.Header.Get("X-Forwarded-For")
	xForwardedForIPs := strings.FieldsFunc(xForwardedForValue, splitFn)
	for _, value := range xForwardedForIPs {
		ipAddress, err := ParseIP(value)
		if err != nil {
			return ipList, fmt.Errorf("parsing failed: %s", err)
		}

		ipList = append(ipList, &ipAddress)
	}

	// Try to parse from header "X-Real-IP"
	xRealIpValue := request.Header.Get("X-Real-IP")
	xRealIpIPs := strings.FieldsFunc(xRealIpValue, splitFn)
	for _, value := range xRealIpIPs {
		ipAddress, err := ParseIP(value)
		if err != nil {
			return ipList, fmt.Errorf("parsing failed: %s", err)
		}

		ipList = append(ipList, &ipAddress)
	}

	return ipList, nil
}

// Creates a new IP ranges timestamp entity for the provided country.
func CreateCountryIPBlocks(country string, cidrFileFolder string) *IpRangesTimestamp {
	var countryIPBlocks []*net.IPNet

	for _, ipType := range []string{
		"ipv4",
		"ipv6",
	} {
		lines, err := ReadFile(cidrFileFolder + "/" + ipType + "/" + strings.ToLower(country) + ".cidr")
		if err != nil {
			panic(fmt.Errorf("failed to read file for version %q and country %q: %v", ipType, country, err))
		}

		for _, value := range lines {
			_, block, err := net.ParseCIDR(value)
			if err != nil {
				panic(fmt.Errorf("parse error on %q: %v", value, err))
			}
			countryIPBlocks = append(countryIPBlocks, block)
		}
	}

	return &IpRangesTimestamp{
		Country:   country,
		IpRanges:  countryIPBlocks,
		Timestamp: time.Now(),
	}
}

// This method initializes the allowed IP ranges.
// It uses a predefined range of private CIDR addresses and reads cidr files based on the provided countries.
// Returns a list of IP ranges timestamp objects.
func InitializeAllowedIPRanges(countries []string, cidrFileFolder string) []*IpRangesTimestamp {
	var allowedIPBlocks []*IpRangesTimestamp

	// Append the private IP addresses first.
	allowedIPBlocks = append(allowedIPBlocks, &IpRangesTimestamp{
		Country:   PrivateIpAddressesTag,
		IpRanges:  InitializePrivateIPBlocks(),
		Timestamp: time.Now(),
	})

	// Now read the country files and append them.
	for _, country := range countries {
		allowedIPBlocks = append(allowedIPBlocks, CreateCountryIPBlocks(country, cidrFileFolder))
	}

	return allowedIPBlocks
}

// This method initializes a list of private IP addresses.
// It uses a predefined range of CIDR addresses.
// Returns a list of private IP blocks.
// https://stackoverflow.com/questions/41240761/check-if-ip-address-is-in-private-network-space
func InitializePrivateIPBlocks() []*net.IPNet {
	var privateIPBlocks []*net.IPNet

	for _, cidr := range []string{
		"127.0.0.0/8",    // IPv4 loopback
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
		"169.254.0.0/16", // RFC3927 link-local
		"::1/128",        // IPv6 loopback
		"fe80::/10",      // IPv6 link-local
		"fc00::/7",       // IPv6 unique local addr
	} {
		_, block, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(fmt.Errorf("parse error on %q: %v", cidr, err))
		}
		privateIPBlocks = append(privateIPBlocks, block)
	}

	return privateIPBlocks
}

// This method checks whether a provided IP is a private IP.
// If this is the case it returns true, otherwise false.
// https://stackoverflow.com/questions/41240761/check-if-ip-address-is-in-private-network-space
func IsPrivateIP(ip net.IP, privateIPBlocks []*net.IPNet) bool {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}

	return IsIpInList(ip, privateIPBlocks)
}

// Checks whether a string is in a list of strings.
// Returns true if this is the case, otherwise returns false.
func IsIpInList(ip net.IP, list []*net.IPNet) bool {
	for _, block := range list {
		if block.Contains(ip) {
			return true
		}
	}

	return false
}

// Tries to parse the IP from a provided address.
// Returns the ip and no error on success, otherwise returns nil and the occured error.
func ParseIP(address string) (net.IP, error) {
	ipAddress := net.ParseIP(address)

	if ipAddress == nil {
		return nil, fmt.Errorf("unable to parse IP from address [%s]", address)
	}

	return ipAddress, nil
}

// Reads a file based on the provided file path.
// Returns each line in a list on success, otherwise nil and an error.
func ReadFile(file string) ([]string, error) {
	var lines []string

	// Open file
	f, err := os.Open(file)
	if err != nil {
		return lines, err
	}
	// Remember to close the file at the end of the program.
	defer f.Close()

	// Read the file line by line using scanner.
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return lines, err
	} else {
		return lines, nil
	}
}
