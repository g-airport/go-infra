## mq



### kafka

```bash
go build -o kafkacli kafka.go
```

```bash
./kafkacli -h
Usage of ./client:
-ca string
CA Certificate (default "ca.pem")
-cert string
Client Certificate (default "cert.pem")
-command string
consumer|producer (default "consumer")
-host string
Common separated kafka hosts (default "localhost:9093")
-key string
Client Key (default "key.pem")
-partition int
Kafka topic partition
-password string
SASL password
-sasl
SASL enable
-tls
TLS enable
-topic string
Kafka topic (default "test--topic")
-username string
SASL username
```

**kafkacli usage**

- 作为producer启动

```bash
$ ./kafkacli -command producer \
-host kafka1:9092,kafka2:9092
```
```bash
# SASL/PLAIN enable
$ ./kafkacli -command producer \
-sasl -username kafkacli -password kafkapassword \
-host kafka1:9092,kafka2:9092
```
```bash
# TLS-enabled
$ ./kafkacli -command producer \
-tls -cert client.pem -key client.key -ca ca.pem \
-host kafka1:9093,kafka2:9093
```

> producer发送消息给kafka：\
> 1 \
2018/12/15 07:11:21 Produced message: [1] \
> 2 \
2018/12/15 07:11:30 Produced message: [2] \
> quit


- 作为consumer启动

```bash
$ ./kafkacli -command consumer \
-host kafka1:9092,kafka2:9092
```

```bash
# SASL/PLAIN enabled
$ ./kafkacli -command consumer \
-sasl -username kafkacli -password kafkapassword \
-host kafka1:9092,kafka2:9092
```

```bash
# TLS-enabled
$ ./kafkacli -command consumer \
-tls -cert client.pem -key client.key -ca ca.pem \
-host kafka1:9093,kafka2:9093
```

> consumer从kafka接受消息：\
2018/12/15 07:11:21 Consumed message: [aaa], offset: [4] \
2018/12/15 07:11:30 Consumed message: [bbb], offset: [5]
