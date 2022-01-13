package client

import (
	"fmt"
	"log"
	"uav_client/src/common"
	"uav_client/src/http/post"
	"uav_client/src/proxy"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	port              int
	host              string
	subscribed_topics []string
	client            mqtt.Client
}

func (c *Client) TelemetryCallback(client mqtt.Client, msg mqtt.Message) {
	common.BuildTelemetryRequest(&common.TelemReq, string(msg.Payload()))
	status := post.Post(common.PostSendTelemetry, &common.TelemReq, &common.TelemResp)
	handled := proxy.HandleStatus(status)
	if !handled {
		// TODO:: Decide what happens if response from POST method is not SUCCESS
		return
	}
	c.Publish("raptor/telemetry/response", common.BuildTelemetryResponse(&common.TelemResp))
}

func contains(slice []string, value string) bool {
	for _, n := range slice {
		if value == n {
			return true
		}
	}
	return false
}

func (c *Client) Init(port int, host string) {
	c.host = host
	c.port = port
	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%d", host, port))
	opts.SetClientID("temp")
	c.client = mqtt.NewClient(opts)
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func (c *Client) Subscribe(topic string) {
	if topic == "raptor/telemetry" {
		if token := c.client.Subscribe(topic, 1, c.TelemetryCallback); token.Wait() && token.Error() != nil {
			log.Fatalln(token.Error())
			return
		}
	} else if topic == "raptor/land" {
		if token := c.client.Subscribe(topic, 1, c.TelemetryCallback); token.Wait() && token.Error() != nil {
			log.Fatalln(token.Error())
			return
		}
	} else {
		log.Fatalln("MQTT subscribe invalid topic name")
		return
	}
	c.subscribed_topics = append(c.subscribed_topics, topic)
}

func (c *Client) Publish(topic, message string) {
	token := c.client.Publish(topic, 1, false, message)
	token.Wait()
}

func (c *Client) Unsubscribe(topic string) {
	if contains := contains(c.subscribed_topics, topic); !contains {
		log.Println("Topic is not in subscribed topic list, returning from function")
		return
	}
	if token := c.client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		log.Fatalln(token.Error())
		return
	}
}

func (c *Client) Disconnect() {
	c.client.Disconnect(250)
}
