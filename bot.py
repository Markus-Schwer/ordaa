import asyncio
import time
import os
import configparser 

from nio import AsyncClient, MatrixRoom, RoomMessageText

config_present = os.path.isfile("config.ini") 

config = configparser.ConfigParser(interpolation=None) 
if config_present:
    config.read("config.ini")

# get matrix settings
matrix_homeserver = config.get("MATRIX", "homeserver") if config_present else os.environ.get("MATRIX_HOMESERVER")
matrix_user = config.get("MATRIX", "user") if config_present else os.environ.get("MATRIX_USER")
matrix_password = config.get("MATRIX", "password") if config_present else os.environ.get("MATRIX_PASSWORD")
matrix_room_id = config.get("MATRIX", "room_id") if config_present else os.environ.get("MATRIX_ROOM_ID")

#print(matrix_homeserver, matrix_user, matrix_password, matrix_room_id)

async def parse_message(event: RoomMessageText) -> None:
    global order_mode_active
    global last_message

    msg = event.body    

    if msg == ".inder" and not order_mode_active:
        await send_message("Order mode active!")
        order_mode_active = True
    elif msg == ".inder" and order_mode_active:
        await send_message("Order mode is already active!")
    elif msg == ".order":
        await send_message("Are you sure you want to send the order?")
    elif msg == ".cancel" and order_mode_active:
        await send_message("Are you sure you want to cancel the order? [Y/N]")
    elif msg.upper() == "Y" and last_message == ".cancel":
        await send_message("Order mode disabled")
        order_mode_active = True 

    last_message = event.body

async def message_callback(room: MatrixRoom, event: RoomMessageText) -> None:
    if room.room_id == matrix_room_id and event.sender != matrix_user and (int(time.time() * 1000) - event.server_timestamp) <= 30 * 1000:
        print(f"{room.user_name(event.sender)} | {event.body}")
        await parse_message(event)
        #await send_message(event.body)

async def send_message(msg) -> None:
    await client.room_send(
        # Watch out! If you join an old room you'll see lots of old messages
        room_id=matrix_room_id,
        message_type="m.room.message",
        content={"msgtype": "m.text", "body": msg},
    )

async def main() -> None:
    global client
    client = AsyncClient(matrix_homeserver, matrix_user)
    client.add_event_callback(message_callback, RoomMessageText)

    print(await client.login(matrix_password))
    # "Logged in as @alice:example.org device id: RANDOMDID"

    # If you made a new room and haven't joined as that user, you can use
    # await client.join("your-room-id")

    await send_message("[.inder] Hey, I'm back!")

    await client.sync_forever(timeout=30000)  # milliseconds

global order_mode_active
order_mode_active = False

global last_message
last_message = ""

asyncio.get_event_loop().run_until_complete(main())

