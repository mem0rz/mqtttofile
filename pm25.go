package main

import (
    "encoding/json"
    mqtt "github.com/eclipse/paho.mqtt.golang"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"
)

var (
    // 目标地址
    destUrl = "http://localhost:1234"
)

func main() {
    chn := make(chan os.Signal, 1)
    signal.Notify(chn, os.Interrupt, syscall.SIGTERM)
    // mq地址
    opts := mqtt.NewClientOptions().AddBroker("tcp://192.168.0.16:1883")
    opts.SetUsername("yourusername").SetPassword("youpassword")
    opts.SetKeepAlive(2 * time.Second)

    opts.SetDefaultPublishHandler(func(client mqtt.Client, message mqtt.Message) {
        msg := &SensorMessage{}
        if err := json.Unmarshal(message.Payload(), msg); err != nil {
            log.Println("recv message parse error")
            return
        }
        os.WriteFile("/var/datatoecowitt/weatherstation/phicommtoecowitt/pm25.data", []byte(msg.PM25), 0666)
    })
    opts.SetPingTimeout(1 * time.Second)

    c := mqtt.NewClient(opts)
    if token := c.Connect(); token.Wait() && token.Error() != nil {
        log.Panicln(token.Error())
    }
    // topic 地址
    token := c.Subscribe("device/zm1/b0f8931e97d5/sensor", 0, nil)
    if token.Wait() && token.Error() != nil {
        log.Panicln(token.Error())
    }
    <-chn
    log.Println("client close")
}

type SensorMessage struct {
    Mac          string `json:"mac"`
    Temperature  string `json:"temperature"`
    Humidity     string `json:"humidity"`
    Formaldehyde string `json:"formaldehyde"`
    PM25         string `json:"PM25"`
}
