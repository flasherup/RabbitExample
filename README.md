<H1>RabbitMQ RPC example</H1>

This is an example of how to use RabbitMQ to implement RPC in GO services.

<H3>How to run</H3>
Before running the services, you need to install RabbitMQ and start it.
Here is an example of how to install RabbitMQ in docker:
> docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 -e RABBITMQ_DEFAULT_USER=admin -e RABBITMQ_DEFAULT_PASS=admin  rabbitmq:3.12-management

Then you can run the services:

Counter service:
> go run services/counter/cmd/main.go

Observer service:
> go run services/observer/cmd/main.go

Once you run both services, you can see the next output.

Counter service:
- 2023/08/31 12:12:22 Counter service is running
- 2023/08/31 12:12:29 Increment 1
- 2023/08/31 12:12:31 Increment 3
- 2023/08/31 12:13:18 Increment 0
- 2023/08/31 12:13:19 Increment 1
- 2023/08/31 12:13:20 Increment 2
- 2023/08/31 12:13:21 Increment 3
- 2023/08/31 12:13:22 Increment 4
- 2023/08/31 12:13:23 Increment 5



Observer service:
- 2023/08/31 12:11:19 Observer service is running
- 2023/08/31 12:11:20 Request 0
- 2023/08/31 12:11:20 Response 1
- 2023/08/31 12:11:21 Request 1
- 2023/08/31 12:11:21 Response 2
- 2023/08/31 12:11:22 Request 2
- 2023/08/31 12:11:22 Response 3
- 2023/08/31 12:11:23 Request 3
- 2023/08/31 12:11:23 Response 4
- 2023/08/31 12:11:24 Request 4
- 2023/08/31 12:11:24 Response 5


