# GALACTUS aka order manager

## Message queue

> run with `podman run --rm -p 5672:5672 -p 8081:15672 rabbitmq:3-management`

- `order/actions`: OrderAction JSON strings
- `order/finalized`: Map of users to their orders

## Endpoints

POST: `/new`
> create a new order, returns the order number

POST: `/{orderNo}/{action}`
> change the order in progress

Actions:
- add
- remove
- finalize
- arrived
- cancel

Add and remove need a JSON body with `user` and `item`.

GET: `/{orderNo}/status`
> status of the order

