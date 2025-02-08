This is a simple implementation of rabbitmq.
there is a publisher and a consumer.

## What Publisher do:
- open new connection and open new channel
- create new "named" queue if not exist
- send message to the "named" queue (main job)


## What Consumer do:
- open new connection and open new channel
- create new "named" queue if not exist (queue name same with queue name that created by publisher)
- waiting message to the "named" queue (main job)



## What should you notice?

### Create new Queue
- name: make sure publisher and consumer has same queue name
- durable: queue always exist eventhough RabbitMQ down
- autoDelete: if true, then the queue will be deleted after there is no consumer.
- exclusive: if true, then the queue can be only accessed by one consummer, and will be deleted after the connection closed.
- noWait: if true, then creating new queue will not wait response from server
- args: other configuration

### Create New Consumer:
- auto-ack: if false, you need to send ack (acknowledge) manually

### When Consumming Message:
- msg.Ack(bool): if false, you only ack this current message. If true 


