package proxy

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/proxy"
)

func init() {
	if socks5Url := os.Getenv("SOCKS5_PROXY"); socks5Url != "" {
		socks5(socks5Url)
	}
}

func socks5(socks5Url string) {
	var (
		auth    *proxy.Auth
		pUrl, _ = url.Parse(socks5Url)
	)

	if pUrl.User != nil {
		auth = &proxy.Auth{}
		auth.User = pUrl.User.Username()
		auth.Password, _ = pUrl.User.Password()
	}

	dialer, err := proxy.SOCKS5("tcp", pUrl.Host, auth, proxy.Direct)
	if err != nil {
		log.Fatalf("socks5 proxy err: %v\n", err)
	}
	http.DefaultClient.Transport = &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		}}
}
