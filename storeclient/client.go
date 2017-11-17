package storeclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/vice-registry/vice-util/communication"
	"github.com/vice-registry/vice-util/persistence"

	"github.com/vice-registry/vice-util/models"
)

// defines timeout in seconds for waiting for store component
const fileTransferTimeout = 300

// StoreRequest defines a store request from worker/api to store component
type StoreRequest struct {
	ImageID    string `json:"imageid"`
	Connection string `json:"connection"`
	AuthToken  string `json:"authtoken"`
	Action     string `json:"action"`
}

// NewStoreRequest issues a new store request: pushes to message queue, opens local tcp socket for file transfer
func NewStoreRequest(image *models.Image, reader io.Reader) error {

	authToken := generateToken()
	tcpPort := generateTCPPort()
	ipAddress, _ := getIPAddress()
	connection := ipAddress + ":" + strconv.Itoa(tcpPort)

	var wg sync.WaitGroup
	var serverError error

	// start tcp server for file transfer
	log.Printf("listen on %s", connection)
	wg.Add(1)
	go func() {
		defer wg.Done()
		serverError = storeRequestServer(authToken, reader, connection)
	}()

	// publish store request to message queue
	request := StoreRequest{
		ImageID:    image.ID,
		Connection: connection,
		AuthToken:  authToken,
		Action:     "store",
	}
	b, err := json.Marshal(request)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	message := string(b)
	communication.SendMessage("store", message)

	wg.Wait()

	return serverError
}

// NewRetrieveRequest issues a new retrieve request: pushes to message queue, opens local tcp socket for file transfer
func NewRetrieveRequest(image *models.Image, writer io.Writer) error {

	authToken := generateToken()
	tcpPort := generateTCPPort()
	ipAddress, _ := getIPAddress()
	connection := ipAddress + ":" + strconv.Itoa(tcpPort)

	var wg sync.WaitGroup
	var serverError error

	// start tcp server for file transfer
	log.Printf("listen on %s", connection)
	wg.Add(1)
	go func() {
		defer wg.Done()
		serverError = retrieveRequestServer(authToken, writer, connection)
	}()

	// publish store request to message queue
	request := StoreRequest{
		ImageID:    image.ID,
		Connection: connection,
		AuthToken:  authToken,
		Action:     "retrieve",
	}
	b, err := json.Marshal(request)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	message := string(b)
	communication.SendMessage("store", message)

	wg.Wait()

	return serverError
}

func storeRequestServer(authToken string, reader io.Reader, connectionURL string) error {
	addr, err := net.ResolveTCPAddr("tcp", connectionURL)
	if err != nil {
		log.Printf("Error parsind address from url %s: %s", connectionURL, err)
		return err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Printf("Error listening on %s for file transfer: %s", connectionURL, err)
		return err
	}
	defer listener.Close()

	err = listener.SetDeadline(time.Now().Add(fileTransferTimeout * time.Second))
	if err != nil {
		log.Printf("Error setting deadline: %s", err)
		return err
	}

	connection, err := listener.Accept()
	if err != nil {
		log.Printf("Error accepting incoming request on file transfer: %s", err)
		return err
	}
	defer connection.Close()
	log.Printf("A client connected to file transfer server.")

	buffer := make([]byte, 1024)
	for {
		// read a chunk
		n, err := reader.Read(buffer)
		if err != nil && err != io.EOF {
			log.Printf("Error in storage: %s", err)
			return err
		}
		if n == 0 {
			break
		}

		// write a chunk
		if _, err := connection.Write(buffer[:n]); err != nil {
			if err != nil {
				log.Printf("Error in storage: %s", err)
				return err
			}
		}
	}
	log.Printf("File transfer server finished for %s", connectionURL)

	return nil
}

func retrieveRequestServer(authToken string, writer io.Writer, connectionURL string) error {
	addr, err := net.ResolveTCPAddr("tcp", connectionURL)
	if err != nil {
		log.Printf("Error parsind address from url %s: %s", connectionURL, err)
		return err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Printf("Error listening on %s for file transfer: %s", connectionURL, err)
		return err
	}
	defer listener.Close()

	err = listener.SetDeadline(time.Now().Add(fileTransferTimeout * time.Second))
	if err != nil {
		log.Printf("Error setting deadline: %s", err)
		return err
	}

	connection, err := listener.Accept()
	if err != nil {
		log.Printf("Error accepting incoming request on file transfer: %s", err)
		return err
	}
	defer connection.Close()
	log.Printf("A client connected to file transfer server.")

	buffer := make([]byte, 1024)
	for {
		// read a chunk
		n, err := connection.Read(buffer)
		if err != nil && err != io.EOF {
			log.Printf("Error in storage: %s", err)
			return err
		}
		if n == 0 {
			break
		}

		// write a chunk
		if _, err := writer.Write(buffer); err != nil {
			if err != nil {
				log.Printf("Error in storage: %s", err)
				return err
			}
		}
	}
	log.Printf("File retrieve server finished for %s", connectionURL)

	return nil
}

func generateToken() string {
	return persistence.GenerateID(12)
}

func generateTCPPort() int {
	rand.Seed(time.Now().Unix())
	port := rand.Intn(65535-65000) + 65000
	return port
}

func getIPAddress() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}
