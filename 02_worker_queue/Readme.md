## What should you notice

### Durable Queue
set durable true, queue will persist on server restart

```
q, err := ch.QueueDeclare(
		"worker_queue_name", // name
		true,                // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
```

### Fair Dispatch
set prefetch count to 1, queue will dispatch message one by one/first come first serve.
it better message distribution.

it also prevent a new "consummer" worker don't get old message. 
```
ch.Qos(
    1,    // prefetch count
    0,    // prefetch size
    false, // global
)
```


### Publish Message Persistence
set amqp.Publishing.DeliveryMode to amqp.Persistent
```
ch.Publish(
    "",     // exchange
    q.Name, // routing key
    false,  // mandatory
    false,  // immediate
    amqp.Publishing{
        DeliveryMode: amqp.Persistent,
        ContentType: "text/plain",
        Body:        []byte("Hello World!"),
    })
```