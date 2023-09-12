# GALACTUS aka order manager

## Endpoints

POST: `/new`
> create a new order

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

