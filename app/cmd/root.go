// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/btcsuite/btcutil/base58"
	"github.com/spf13/cobra"
	"net"
	"strconv"
	"strings"

	"encoding/base64"

	"github.com/kprc/basclient/dnsclient"
	//"github.com/kprc/basserver/dns/server"
	"encoding/hex"
	"github.com/BASChain/go-bas-dns-server/lib/dns"
)

//var cfgFile string
var encodetyp string
var domainnametyp string
var remotehost string
var querystring string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bas",
	Short: "Send A DNS Request And Get the Result",
	Long:  `Send A DNS Request And Get the Result`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(querystring) == 0 {
			fmt.Println("Please input domain for search")
			return
		}
		qs := querystring

		if domainnametyp != "dn" {
			var b []byte
			var err error

			switch encodetyp {
			case "base64":
				b, err = base64.StdEncoding.DecodeString(qs)
				if err != nil {
					fmt.Println("error:", err)
					return
				}
			case "base16":
				if strings.ToLower(qs[:2]) == "0x" {
					//b, err = hexutil.Decode(qs)
					//if err != nil {
					//	fmt.Println(err)
					//	return
					//}
					b := make([]byte, 128)
					n, err := hex.Decode(b, []byte(qs[:2]))
					if err != nil {
						fmt.Println(err)
						return
					}
					b = b[:n]
				} else {
					fmt.Println("base16 address must have 0x prefix")
					return
				}
			case "base58":
				b = base58.Decode(qs)
				if len(b) == 0 {
					fmt.Println("base58 address error")
					return
				}
			default:
				fmt.Println("no enodetyp input")
				return

			}
			if len(b) > 0 {
				qs = base58.Encode(b)
			}

		}

		rh := strings.Split(remotehost, ":")
		if len(rh) > 2 {
			fmt.Println("please input correct host")
			return
		}

		rhost := rh[0]
		ip := net.ParseIP(rh[0])
		if ip == nil {
			fmt.Println("please input correct host")
			return
		}

		if len(rh) == 1 {
			rhost += ":53"
		} else {
			port, err := strconv.Atoi(rh[1])
			if err != nil || (port == 0 || port > 65535) {
				fmt.Println("please input correct host")
				return
			}
			rhost += ":" + rh[1]
		}

		typ := dns.TypeA

		if domainnametyp != "dn" {
			typ = 65
		}
		qs1 := qs + "."
		msg := dnsclient.SendAndRcv(rhost, qs1, typ)
		if msg == nil {
			fmt.Println("command line failed, please try again")
		} else {
			fmt.Println("Src Address:", querystring)
			if msg.Rcode == dns.RcodeBadKey {
				fmt.Println("Not Found in Blockchain Address System")
			} else if len(msg.Answer) > 0 {
				for i := 0; i < len(msg.Answer); i++ {
					rr := msg.Answer[i]
					switch rr.(type) {
					case *dns.A:
						a := rr.(*dns.A)
						fmt.Println("Ip  Address:", a.A.String())
					case *dns.CNAME:
						cn := rr.(*dns.CNAME)
						fmt.Println("CName:", cn.Target)
					case *dns.NULL:
						alias := rr.(*dns.NULL)
						fmt.Println("Null:", alias.Data)
					default:
						fmt.Println("Not support")
					}
				}

			} else {
				fmt.Println("No Answer")
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

//
func init() {
	//cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.app.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.Flags().StringVarP(&encodetyp, "encode-type", "e", "base16", "encoding type [base16,base58,base64]")
	rootCmd.Flags().StringVarP(&domainnametyp, "domain-type", "t", "dn", "domain name type [dn:domain name,eth: ethereum address,ipv4,other: other address]")
	rootCmd.Flags().StringVarP(&remotehost, "remote-bas-server", "r", "103.45.98.72", "remote bas server (also you can input port like as 103.45.98.72:53)")
	rootCmd.Flags().StringVarP(&querystring, "query-string", "q", "", "domain name or ethereum address or other address")
}

//
//// initConfig reads in config file and ENV variables if set.
//func initConfig() {
//	if cfgFile != "" {
//		// Use config file from the flag.
//		viper.SetConfigFile(cfgFile)
//	} else {
//		// Find home directory.
//		home, err := homedir.Dir()
//		if err != nil {
//			fmt.Println(err)
//			os.Exit(1)
//		}
//
//		// Search config in home directory with name ".app" (without extension).
//		viper.AddConfigPath(home)
//		viper.SetConfigName(".app")
//	}
//
//	viper.AutomaticEnv() // read in environment variables that match
//
//	// If a config file is found, read it in.
//	if err := viper.ReadInConfig(); err == nil {
//		fmt.Println("Using config file:", viper.ConfigFileUsed())
//	}
//}
