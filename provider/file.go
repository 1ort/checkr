package provider

import (
	"bufio"
	"log"
	"os"

	"github.com/1ort/checkr/proxy"
)

func FromFiles(files ...string) ProxyProvider {
	if len(files) == 0 {
		return nil
	}
	ps := make([]ProxyProvider, len(files))
	for i, file := range files {
		ps[i] = &fileProvider{
			file,
		}
	}
	return Group(ps...)
}

type fileProvider struct {
	path string
}

func (p *fileProvider) Fetch(done chan int) <-chan proxy.Proxy {
	ch := make(chan proxy.Proxy)
	go func() {
		defer close(ch)
		file, err := os.Open(p.path)
		if err != nil {
			log.Println(err)
			return
		}
		defer file.Close()
		fileScanner := bufio.NewScanner(file)
		for fileScanner.Scan() {
			proxy, err := proxy.FromHostPort(fileScanner.Text())
			if err != nil {
				log.Println(err)
				continue
			}
			select {
			case <-done:
				return
			default:
				ch <- proxy
			}
		}
	}()
	return ch
}
