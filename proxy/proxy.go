package proxy

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/varun-muthanna/loadbalancer/balancer"
)

type handler struct {
	addr string
	lb   *balancer.Balancer
}

func StartReverseProxy(listenaddress string, lb *balancer.Balancer) {

	h := &handler{
		addr: listenaddress,
		lb:   lb,
	}

	s := &http.Server{
		Addr:    listenaddress,
		Handler: h,
	}

	ch := make(chan os.Signal, 1) //buffer as incase go routine is not listening and signal arrives it can be lost

	go func() {

		err := s.ListenAndServe()

		if err != nil {
			log.Printf("proxy listener stopped: %v", err)
		}

	}()

	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)

	sig := <-ch //blocking

	//default behaivour if no blocking that will is to terminate the process (no defer is run) , now it is handled manually  (defer is run)

	fmt.Println("Gracefull shutdown initiated", sig)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	defer s.Shutdown(ctx)
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	srv := h.lb.GetLeastConnections(r.Host)

	var url string

	if srv == nil { //public dns
		ipList, err := net.LookupIP(r.Host)

		if err != nil || len(ipList)==0 {
			log.Printf("Error in DNS resolution %s", err)
		}
		
		url = "http://[" + ipList[0].String() +"]"

	} else { //load balanced to one of localhost
		url = "http://" + srv.GetAddress() + r.RequestURI
		srv.IncrementConnection()
		defer srv.DecrementConnections()
	}

	//fmt.Printf("url , %s\n",url)
	req, err := http.NewRequest(r.Method, url, r.Body)
	
	if err != nil {
		log.Printf("Error in creating request %s", err)
		return
	}

	req.Host=r.Host

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Printf(" error making http request: %s", err)
		if srv != nil  {
			srv.SetHealth(false)
		}
		return
	}

	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body) //captures w.Write([]) from backend servers/public DNS

	if err != nil {
		log.Printf("error reading response: %s", err)
		return
	}

	w.Write(resBody) //essentially doing it twice

	
	if srv !=nil {
		fmt.Println(string(resBody))
	}
}

// architecture

// client ----------- proxy ---------  server
//      make request        select
// 							connect
//		transfer data		transfer data
