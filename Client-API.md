# API

Resolute uses websockets and JSON to communicate. When you establish a connection, you can send the following commands...


## Rooms

Rooms represent conversations between 1 or more users and are secured with two types of keys, one-time-use keys and forever keys, both of which the creator of the room must request the server to generate as they are not created by default.

### Create Room

**Request**
```json
{
   "cmd":"create-room",
   "data":{
      "name":"test room",
      "username":"alice"
   }
}
```

**OK Response**
```json
{
   "cmd":"create-room",
   "room":{
      "room_id":"0uA4ujPmSOuHzYecbDe2ysGFmzeIsM5j",
      "name":"test room",
      "one_time_join_keys":[],
      "forever_join_key":""
   }
}
```

### Generate One-time Key for Room

These keys can only be used once. Once they are used to join a room, they are deleted.

**Request**
```json
{
   "cmd":"room-key-onetime",
   "data":{
      "room_id":"0uA4ujPmSOuHzYecbDe2ysGFmzeIsM5j"
   }
}
```

**Response**
```json
{
   "cmd":"room-key-onetime",
   "data":{
      "room_id":"0uA4ujPmSOuHzYecbDe2ysGFmzeIsM5j",
      "one_time_key":"oe82PE1bdSAqofRokah2pNChGm4LvunV"
   }
}
```

### Generate Forever Key for Room

These keys can be used to join a room repeatedly with no limit.

**Request**
```json
{
  "cmd":"room-key-forever",
  "data":{
    "room_id":"0uA4ujPmSOuHzYecbDe2ysGFmzeIsM5j"
  }
}
```

**Response**
```json
{
  "cmd":"room-key-forever",
  "data":{
    "room_id":"0uA4ujPmSOuHzYecbDe2ysGFmzeIsM5j",
    "forever_join_key":"O7lCoY5GvmBWVBvBT4egMXLkjhaYKdFA"
  }
}
```

### Join Room

First you need the Room ID and a one-time or forever key to join. The owner of the room must share this with you.

#### Join with one-time key

**Request**
```json
{
   "cmd":"join-room-onetime",
   "data":{
      "room_id":"iL6oB1RBw1AK1Odg6UoU397eDozV3fek",
      "one_time_key":"kQQ8ZKpz0PmoqQQbGILvgoMc72etKIbI",
      "username":"bob"
   }
}
```

**OK Response**
```json
{
   "cmd":"join-room-onetime",
   "data":{
      "room_id":"1cziERIkWDI9kKoeB4rytoaEijQUSfCG",
      "room_name":"test room"
   }
}
```

#### Join with forever key

**Request**
```json
{
   "cmd":"join-room-forever",
   "data":{
      "room_id":"iL6oB1RBw1AK1Odg6UoU397eDozV3fek",
      "forever_key":"kQQ8ZKpz0PmoqQQbGILvgoMc72etKIbI",
      "username":"bob"
   }
}
```

**OK Response**
```json
{
   "cmd":"join-room-forever",
   "data":{
      "room_id":"1cziERIkWDI9kKoeB4rytoaEijQUSfCG",
      "room_name":"test room"
   }
}
```


### Sending Messages

**Request**
```json
{
   "cmd":"send-message",
   "data":{
      "room_id":"1cziERIkWDI9kKoeB4rytoaEijQUSfCG",
      "content":"hello world"
   }
}
```

**OK Response**
```json
{
  "cmd":"send-message",
  "data":{
    "room_id":"1cziERIkWDI9kKoeB4rytoaEijQUSfCG",
    "user_id":"f2LySPQW8NVacLEKhZlS7QggT1wINjnf",
    "username":"alice",
    "content":"hello world"
  }
}
```


more coming soon...
