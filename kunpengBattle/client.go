package kunpengBattle

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"strconv"
)

type KunPengBattleClient struct {
	TeamID     int
	TeamName   string
	ServerIP   string
	ServerPort int
	conn       *battleConnection
	reader     *bufio.Reader
	writer     *bufio.Writer
	recvChan   chan []byte
	replyChan  chan []byte
	errChan    chan error
}

func NewKunPengBattleClient(teamID int, teamName string) KunPengBattleClient {
	client := KunPengBattleClient{TeamID: teamID, TeamName: teamName}
	client.conn = new(battleConnection)
	client.reader = bufio.NewReaderSize(client.conn, 1024*10)
	client.writer = bufio.NewWriterSize(client.conn, 1024*10)

	return client
}

func (c *KunPengBattleClient) Connect(ip string, port int) error {
	c.ServerIP = ip
	c.ServerPort = port
	address := ip + ":" + strconv.Itoa(port)
	err := c.conn.connect(address)

	return err
}

func (c *KunPengBattleClient) Start() error {

	if err := c.registrate(); err != nil {
		return err
	}

	log.Println("registrate finish!!!")

	c.receive()

	return nil
}

func (c *KunPengBattleClient) receive() {
	scanner := bufio.NewScanner(c.reader)
	scanner.Split(splitPackage)
	for scanner.Scan() {
		msgBytes := scanner.Bytes()
		// log.Printf(string(msgBytes))
		msg := new(KunPengMsg)
		err := json.Unmarshal(msgBytes, msg)
		if err != nil {
			log.Printf("Unmarshal err: %v", err)
		}

		log.Println(msg)

		// switch msg.Name {
		// case "leg_start":
		// 	leg_start := new(KunPengLegStart)
		// 	err := json.Unmarshal(msg.Data, leg_start)
		// case "round":
		// 	round := new(KunPengLegStart)
		// 	err := json.Unmarshal(msg.Data, round)
		// case "leg_end":
		// 	leg_end := new(KunPengLegStart)
		// 	err := json.Unmarshal(msg.Data, leg_end)
		// case "game_over":
		// 	log.Println("game_over!!!")
		// default:
		// 	log.Printf("Unknown msg Name: %v", msg.Name)
		// }

	}
}

func (c *KunPengBattleClient) registrate() error {
	kpMsg := new(KunPengMsg)
	kpr := KunPengRegistration{TeamID: c.TeamID, TeamName: c.TeamName}
	kpMsg.Name = "registration"
	kpMsg.Data = kpr

	msgBytes, _ := json.Marshal(kpMsg)

	log.Println(string(msgBytes))
	if err := c.send(msgBytes); err != nil {
		return err
	}

	err := c.writer.Flush()
	return err
}

func (c *KunPengBattleClient) send(msgBytes []byte) error {
	msgLen := len(msgBytes)

	if msgLen > 99999 {
		return errors.New("msgLen shoule not greater than 99999")
	}

	msgLenStr := strconv.FormatInt(int64(msgLen), 10)
	sizeBytes := make([]byte, 0, msgLen+5)
	for i := 0; i < 5-len(msgLenStr); i++ {
		sizeBytes = append(sizeBytes, '0')
	}
	sizeBytes = append(sizeBytes, []byte(msgLenStr)...)

	sendBytes := append(sizeBytes, msgBytes...)

	log.Println(string(sendBytes))

	if _, err := c.writer.Write(sendBytes); err != nil {
		return err
	}

	if err := c.writer.Flush(); err != nil {
		return err
	}

	return nil
}
