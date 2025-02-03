package main

import (
  "fmt"
  "net"
  "flag"
  "os"
  "sync"
  "time"
)

func isValidHost(host string) bool {
  if net.ParseIP(host) != nil {
    return true
  }
  
  _, err := net.LookupHost(host)
  return err == nil
}

func scanPort(host string, port int, results chan int, wg *sync.WaitGroup) {
  defer wg.Done()

  address := fmt.Sprintf("%s:%d", host, port)
  conn, err := net.DialTimeout("tcp", address, 2*time.Second)
  if err != nil {
    return
  }
  conn.Close()
  results <- port
}

func main() {
  banner := `
   ███▄    █  ███▄ ▄███▓ ▄▄▄       ██▓███       ▄████  ▒█████  
 ██ ▀█   █ ▓██▒▀█▀ ██▒▒████▄    ▓██░  ██▒    ██▒ ▀█▒▒██▒  ██▒
▓██  ▀█ ██▒▓██    ▓██░▒██  ▀█▄  ▓██░ ██▓▒   ▒██░▄▄▄░▒██░  ██▒
▓██▒  ▐▌██▒▒██    ▒██ ░██▄▄▄▄██ ▒██▄█▓▒ ▒   ░▓█  ██▓▒██   ██░
▒██░   ▓██░▒██▒   ░██▒ ▓█   ▓██▒▒██▒ ░  ░   ░▒▓███▀▒░ ████▓▒░
░ ▒░   ▒ ▒ ░ ▒░   ░  ░ ▒▒   ▓▒█░▒▓▒░ ░  ░    ░▒   ▒ ░ ▒░▒░▒░ 
░ ░░   ░ ▒░░  ░      ░  ▒   ▒▒ ░░▒ ░          ░   ░   ░ ▒ ▒░ 
   ░   ░ ░ ░      ░     ░   ▒   ░░          ░ ░   ░ ░ ░ ░ ▒  
         ░        ░         ░  ░                  ░     ░ ░  

  `


  fmt.Println(banner)

  var host string
  fmt.Print("Enter a Target IP or Hostname: ")
  fmt.Scan(&host)

  if !isValidHost(host) {
    fmt.Printf("Invalid host: %s\n", host)
    os.Exit(1)
  }

  startPort := flag.Int("start", 1, "Start port")
  endPort := flag.Int("end", 1024, "End port")
  flag.Parse()

  fmt.Printf("Scanning %s from port %d to %d...\n", host, *startPort, *endPort)

  results := make(chan int, 100)
  var wg sync.WaitGroup

  for port := *startPort; port <= *endPort; port++ {
    wg.Add(1)
    go scanPort(host, port, results, &wg)
  }

  go func() {
    wg.Wait()
    close(results)
  }()

  for port := range results {
    fmt.Printf("[+] Open Port: %d\n", port)
  }
}
