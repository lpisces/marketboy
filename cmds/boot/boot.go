package boot

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
	"net/url"
	"os"
	"os/signal"
	"time"
)

const (
	AuthExpire = 7 * 24 * 60 * 60
)

type CMD struct {
	Command string        `json:"op"`
	Args    []interface{} `json:"args"`
}

func Run(c *cli.Context) (err error) {
	if err := Conf.Load(c); err != nil {
		return err
	}

	for _, v := range Conf.Trading.Symbol {
		params := make(map[string]interface{})
		params["symbol"] = v
		params["leverage"] = Conf.Trading.Leverage
		if err := setLeverage(params); err != nil {
			log.Info(err)
		}
	}

	Ping := time.Duration(Conf.Trading.Watch)

	log.SetLevel(log.InfoLevel)
	if Conf.Debug {
		log.SetLevel(log.DebugLevel)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: Conf.WSConfig.Scheme, Host: Conf.WSConfig.Host, Path: Conf.WSConfig.Path}
	log.Infof("connecting to %s", u.String())

	// connection
	conn, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	log.Info(resp.Header)
	if err != nil {
		return err
	}
	defer conn.Close()

	done := make(chan struct{})

	// receive message
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Error("read:", err)
				return
			}
			if err := dispatch(message); err != nil {
				log.Error(err)
			}
		}
	}()

	// ping
	ticker := time.NewTicker(time.Second * Ping)
	defer ticker.Stop()

	// Auth
	expires := time.Now().Unix() + int64(AuthExpire)
	sign := HmacSha256([]byte(Conf.AuthConfig.Secret), []byte(fmt.Sprintf("%s%d", "GET/realtime", expires)))
	cmd := CMD{
		Command: "authKeyExpires",
		Args:    []interface{}{Conf.AuthConfig.Key, expires, sign},
	}
	msg, _ := json.Marshal(cmd)
	err = conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		return err
	}

	// subscribe topics
	subscribe := func(topic interface{}) error {
		cmd := CMD{
			Command: "subscribe",
			Args:    []interface{}{topic},
		}
		msg, _ := json.Marshal(cmd)
		return conn.WriteMessage(websocket.TextMessage, msg)
	}

	for _, topic := range Conf.Subscribe.Topic {
		retryLimit := 3
	Retry:
		if retryLimit <= 0 {
			return fmt.Errorf("subscribe %v failed", topic)
		}
		if err := subscribe(topic); err != nil {
			log.Error(err)
			time.Sleep(time.Second * 1)
			retryLimit--
			goto Retry
		}
	}

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			err := conn.WriteMessage(websocket.TextMessage, []byte("ping"))
			if err != nil {
				log.Error("write:", err)
				return err
			}
		case <-interrupt:
			log.Info("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Error("write close:", err)
				return err
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return nil
		}
	}

	return
}
