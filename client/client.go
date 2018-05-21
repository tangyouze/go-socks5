package main

import (
	"github.com/tangyouze/go-socks5"
	"net"
	"golang.org/x/net/context"
	"golang.org/x/net/proxy"
	"fmt"
	"os"
	"net/http"
	"io/ioutil"
	"time"
)

var (
	proxyPool      = []string{"1080", "1081", "1082"}
	proxyAvailable = make(map[string]bool)
)

func checkProxy(proxyAddr, url string) error {
	// create a socks5 dialer
	dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
	if err != nil {
		return err
	}
	// setup a http client
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	// set our socks5 as the dialer
	httpTransport.Dial = dialer.Dial
	// create a request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	// use the http client to fetch the page
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading body:", err)
		return err
	}
	//fmt.Println(string(b))
	return nil
}

func newDial(ctx context.Context, network, addr string) (net.Conn, error) {

	fmt.Println(ctx, network, addr)
	dial, err := proxy.SOCKS5(network, "ubuntu.urwork.qbtrade.org:1080", nil, proxy.Direct)
	if err != nil {
		return nil, err
	}
	conn, err := dial.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func mainServe() {
	// Create a SOCKS5 server
	conf := &socks5.Config{
		Dial: newDial,
	}
	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	// Create SOCKS5 proxy on localhost port 8000
	if err := server.ListenAndServe("tcp", "127.0.0.1:8000"); err != nil {
		panic(err)
	}
}
func checkAllProxy() {
	for {
		for _, element := range proxyPool {
			err := checkProxy("ubuntu.urwork.qbtrade.org:"+element, "https://google.com")
			if err != nil {
				delete(proxyAvailable, element)
			} else {
				proxyAvailable[element] = true
			}
		}
		fmt.Println(proxyAvailable)
		time.Sleep(10 * time.Second)
	}

}
func main() {

	go checkAllProxy()
	mainServe()

}
