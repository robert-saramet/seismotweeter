# Accelerometer imports
import time
import board
import busio
import digitalio
import adafruit_fxos8700

# Wifi imports
import os
import rtc
import ssl
import wifi
import ipaddress
import socketpool
import microcontroller
import adafruit_ntp
import adafruit_requests

# Seismometer config
samples = 4
delay = 0.1
min_deviation = 1.5
g = 9.8
cooldown = 5

# Fetch data from toml file
ssid = os.getenv('CIRCUITPY_WIFI_SSID')
passwd = os.getenv('CIRCUITPY_WIFI_PASSWORD')
url = os.getenv('IFTTT_URL')
api_key = os.getenv('IFTTT_KEY')

# Connect to Wifi
wifi.radio.connect(ssid, passwd)
pool = socketpool.SocketPool(wifi.radio)

# Display connection status
print(f'Succesfully connected to network \'{ssid}\'')
print('IP address is', wifi.radio.ipv4_address)
print()

# Sync local time to NTP
ntp = adafruit_ntp.NTP(pool, tz_offset=2)
rtc.RTC().datetime = ntp.datetime

# Set up requests
requests = adafruit_requests.Session(pool, ssl.create_default_context())
base_url = url + api_key + '?value1='

# Set up accelerometer
scl = board.GP1
sda = board.GP0
i2c = busio.I2C(scl, sda)
fxos = adafruit_fxos8700.FXOS8700(i2c)

# Initialize LED
led = digitalio.DigitalInOut(board.LED)
led.direction = digitalio.Direction.OUTPUT

while True:
    # Measure acceleration
    accel = 0
    for i in range(samples):
        accel += fxos.accelerometer[2]
        time.sleep(delay / samples)
    avg = accel / samples
    rel_accel = abs(avg - g)
    # If acceleration higher than threshold
    if  abs(avg - g) >= min_deviation:
        led.value = True
        et = time.localtime()
        print(f'Earthquake of relative acceleration {rel_accel:.2f} detected at {et.tm_hour}:{et.tm_min}:{et.tm_sec}')
        try:
            # Send request to IFTTT for tweeting
            message = base_url + f'{rel_accel:.2f}'
            response = requests.get(message)
            print(response.text)
            print()
            time.sleep(cooldown)
            led.value = False
        except Exception as e:
            # Start from scratch
            print("Error:\n", str(e))
            print("Resetting in 10 seconds")
            led.value = True
            for i in range(5):
                led.value = not led.value
                time.sleep(2)
            microcontroller.reset()
