package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/Shopify/sarama"
)

var (
	command    string
	hosts      string
	topic      string
	partition  int
	saslEnable bool
	username   string
	password   string
	tlsEnable  bool
	clientCert string
	clientKey  string
	caCert     string
)

func main() {
	flag.StringVar(&command, "command", "consumer", "consumer|producer")
	flag.StringVar(&hosts, "host", "localhost:9093", "Common separated kafka hosts")
	flag.StringVar(&topic, "topic", "test--topic", "Kafka topic")
	flag.IntVar(&partition, "partition", 0, "Kafka topic partition")

	flag.BoolVar(&saslEnable, "sasl", false, "SASL enable")
	flag.StringVar(&username, "username", "", "SASL Username")
	flag.StringVar(&password, "password", "", "SASL Password")

	flag.BoolVar(&tlsEnable, "tls", false, "TLS enable")
	flag.StringVar(&clientCert, "cert", "cert.pem", "Client Certificate")
	flag.StringVar(&clientKey, "key", "key.pem", "Client Key")
	flag.StringVar(&caCert, "ca", "ca.pem", "CA Certificate")
	flag.Parse()

	config := sarama.NewConfig()
	if saslEnable {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = username
		config.Net.SASL.Password = password
	}

	if tlsEnable {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
		tlsConfig, err := genTLSConfig(clientCert, clientKey, caCert)
		if err != nil {
			log.Fatal(err)
		}

		config.Net.TLS.Enable = true
		config.Net.TLS.Config = tlsConfig
	}

	client, err := sarama.NewClient(strings.Split(hosts, ","), config)
	if err != nil {
		log.Fatalf("unable to create kafka client: %q", err)
	}

	if command == "consumer" {
		consumer, err := sarama.NewConsumerFromClient(client)
		if err != nil {
			log.Fatal(err)
		}
		defer consumer.Close()
		loopConsumer(consumer, topic, partition)
	} else {
		producer, err := sarama.NewAsyncProducerFromClient(client)
		if err != nil {
			log.Fatal(err)
		}
		defer producer.Close()
		loopProducer(producer, topic, partition)
	}
}

func genTLSConfig(clientCertFile, clientKeyFile, caCertFile string) (*tls.Config, error) {
	// load client cert
	clientCert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		return nil, err
	}

	// load ca cert pool
	caCert, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// generate tls config
	tlsConfig := tls.Config{}
	tlsConfig.RootCAs = caCertPool
	tlsConfig.Certificates = []tls.Certificate{clientCert}
	tlsConfig.BuildNameToCertificate()
	// tlsConfig.InsecureSkipVerify = true // This can be used on test server if domain does not match cert
	return &tlsConfig, err
}

func loopProducer(producer sarama.AsyncProducer, topic string, partition int) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
		} else if text == "exit" || text == "quit" {
			break
		} else {
			producer.Input() <- &sarama.ProducerMessage{Topic: topic, Key: nil, Value: sarama.StringEncoder(text)}
			log.Printf("Produced message: [%s]\n", text)
		}
		fmt.Print("> ")
	}
}

func loopConsumer(consumer sarama.Consumer, topic string, partition int) {
	partitionConsumer, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
	if err != nil {
		log.Println(err)
		return
	}
	defer partitionConsumer.Close()

	for {
		msg := <-partitionConsumer.Messages()
		log.Printf("Consumed message: [%s], offset: [%d]\n", msg.Value, msg.Offset)
	}
}
