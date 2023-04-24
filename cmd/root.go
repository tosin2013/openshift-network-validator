/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cobra-cli",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

type DNSRecord struct {
	Name   string `yaml:"name"`
	Port   string `yaml:"port"`
	Record string `yaml:"record"`
}

type NetworkingConfig struct {
	InternetAccess  bool     `yaml:"internet_access"`
	TestDHCP        bool     `yaml:"test_dhcp"`
	SubnetConflicts []string `yaml:"subnet_conflicts"`
}

type Config struct {
	Name        string           `yaml:"name"`
	DNS         []DNSRecord      `yaml:"dns"`
	ClusterName string           `yaml:"cluster_name"`
	BaseDomain  string           `yaml:"base_domain"`
	URLAccess   []string         `yaml:"url_access"`
	Networking  NetworkingConfig `yaml:"networking"`
}

var configFile string
var clusterName string
var baseDomain string

var dnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "Test DNS records",
	Run: func(cmd *cobra.Command, args []string) {
		configYaml, err := readConfig(configFile)
		if err != nil {
			fmt.Println("Error reading YAML file:", err)
			return
		}

		var config Config
		err = yaml.Unmarshal(configYaml, &config)
		if err != nil {
			fmt.Println("Error parsing YAML:", err)
			return
		}

		for _, dnsRecord := range config.DNS {
			record := strings.ReplaceAll(dnsRecord.Record, "<cluster_name>", clusterName)
			record = strings.ReplaceAll(record, "<base_domain>", baseDomain)

			ips, err := net.LookupIP(record)
			if err != nil {
				fmt.Println("Error resolving DNS record:", err)
				fmt.Println("Check that the DNS record is configured correctly in your DNS server.")
				fmt.Println("********************************************************************")
				const s = "\tAllocate two IP addresses, outside the DHCP range, and configure them with reverse DNS records.\n" +
					"\tA DNS record for api.<cluster_name>.<base_domain> pointing to the allocated IP address.\n" +
					"\tA DNS record for *.apps.<cluster_name>.<base_domain> pointing to the allocated IP address."
				fmt.Println(s)
				return
			}

			fmt.Println("IP addresses for", dnsRecord.Port, ":", ips)
		}
	},
}

var urlAccessCmd = &cobra.Command{
	Use:   "url-access",
	Short: "Test URL access and generate a CSV report",
	Run: func(cmd *cobra.Command, args []string) {
		configYaml, err := readConfig(configFile)
		if err != nil {
			fmt.Println("Error reading YAML file:", err)
			return
		}

		var config Config
		err = yaml.Unmarshal(configYaml, &config)
		if err != nil {
			fmt.Println("Error parsing YAML:", err)
			return
		}

		results := [][]string{{"URL", "Status Code"}}

		for _, url := range config.URLAccess {
			if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
				url = "https://" + url
			}

			resp, err := http.Get(url)
			if err != nil {
				fmt.Println("Error testing URL:", url, err)
				results = append(results, []string{url, err.Error()})
			} else {
				statusCode := fmt.Sprintf("%d", resp.StatusCode)
				fmt.Println("Tested URL:", url, "Status Code:", statusCode)
				results = append(results, []string{url, statusCode})
			}
		}

		err = writeResultsToFile(results, "url_access_report.csv")
		if err != nil {
			fmt.Println("Error writing CSV report:", err)
			return
		}

		fmt.Println("CSV report written to url_access_report.csv")
	},
}

func writeResultsToFile(results [][]string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, row := range results {
		err := writer.Write(row)
		if err != nil {
			return err
		}
	}

	return nil
}

func readConfig(file string) ([]byte, error) {
	configYaml, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return configYaml, nil
}

var testNetworkingCmd = &cobra.Command{
	Use:   "test-networking",
	Short: "Test networking configuration",
	Run: func(cmd *cobra.Command, args []string) {
		configYaml, err := readConfig(configFile)
		if err != nil {
			fmt.Println("Error reading YAML file:", err)
			return
		}

		var config Config
		err = yaml.Unmarshal(configYaml, &config)
		if err != nil {
			fmt.Println("Error parsing YAML:", err)
			return
		}

		if config.Networking.InternetAccess {
			fmt.Println("Testing internet access...")
			err = testInternetAccess()
			if err != nil {
				fmt.Println("Internet access test failed:", err)
			} else {
				fmt.Println("Internet access test succeeded")
			}
		}

		if config.Networking.TestDHCP {
			fmt.Println("Testing DHCP configuration...")
			err = testDHCP()
			if err != nil {
				fmt.Println("DHCP configuration test failed:", err)
			} else {
				fmt.Println("DHCP configuration test succeeded")
			}
		}

		if len(config.Networking.SubnetConflicts) > 0 {
			fmt.Println("Testing subnet conflicts...")
			err = testSubnetConflicts(config.Networking.SubnetConflicts)
			if err != nil {
				fmt.Println("Subnet conflicts test failed:", err)
			} else {
				fmt.Println("Subnet conflicts test succeeded")
			}
		}
	},
}

func testInternetAccess() error {
	_, err := net.LookupHost("google.com")
	if err != nil {
		return err
	}
	return nil
}

func testDHCP() error {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				// We found an IPv4 address
				if ipnet.IP[0] == 169 && ipnet.IP[1] == 254 {
					return fmt.Errorf("found DHCP-assigned address: %s", ipnet.IP.String())
				}
			}
		}
	}

	return nil
}

func testSubnetConflicts(subnets []string) error {
	for _, subnetStr := range subnets {
		_, subnet, err := net.ParseCIDR(subnetStr)
		if err != nil {
			return err
		}

		addrs, err := net.InterfaceAddrs()
		if err != nil {
			return err
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if subnet.Contains(ipnet.IP) {
					return fmt.Errorf("subnet conflict: %s", subnetStr)
				}
			}
		}
	}
	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(dnsCmd)

	dnsCmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to YAML configuration file")
	dnsCmd.MarkFlagRequired("config")

	dnsCmd.Flags().StringVarP(&clusterName, "cluster-name", "n", "", "Name of the Kubernetes cluster")
	dnsCmd.MarkFlagRequired("cluster-name")

	dnsCmd.Flags().StringVarP(&baseDomain, "base-domain", "d", "", "Base domain for the Kubernetes cluster")
	dnsCmd.MarkFlagRequired("base-domain")

	rootCmd.AddCommand(urlAccessCmd)

	urlAccessCmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to YAML configuration file")
	urlAccessCmd.MarkFlagRequired("config")

	rootCmd.AddCommand(testNetworkingCmd)
	testNetworkingCmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to YAML configuration file")
	testNetworkingCmd.MarkFlagRequired("config")

}
