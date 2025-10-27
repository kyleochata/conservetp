# conservetp
dockerized microservice conservetp


## Docker

### Compose
`develop.watch{action: rebuild, path: ./, target: /app}`

Command to rebuild the container after changes to the path are made. Target path is the location within the docker container rebuild is to take place.

- Use rebuild to recreate the docker image and re-initialize the go server. 
- sync+restart is available but doesn't build the container again (issue with go needing to be compiled.)
    - Must use a two stage rebuild (build binary then cmd ["./main"])


## Services
Users, Orders, Inventory, EventBroker

### Users
Store user(s) with addresses.
    - id, name, email, created_on, last_updated, active_status, active_status
    - Addresses: is_primary
User-id will be the primary key for orders.
    - Call to orders service to pull past and current orders

* As of now no event creations (except maybe a local one for logging purposes)

### Orders
Store(s) & Create(s) orders. Create order event tickets for use in the eventbroker to check for inventory item(s).
    - id, user_id, number of items, items_ordered, total_price, paid
    Item table
        - id, name, description, size, notes, price, materials
    Materials table
        - id, name, description, qty

### EventBroker


### Inventory
Create and store product & material(s).
Product Table
    - id, name, description, in_stock, total_sold, price_per_unit, price_per_build, specials, appx_build_time (hr)
    Product_size table
        - id, product_id, size, in_stock, amt, discontinued
    Product-Mat table
        - id, mat_name, product_id, desc, items_id, amount_needed
Items table
    - id, name, description, in_stock, amount_onHand, amount_onOrder, total_Price
    
    