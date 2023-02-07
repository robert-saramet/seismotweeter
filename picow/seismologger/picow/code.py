import time
import board
import busio
import storage
import digitalio
import adafruit_fxos8700

from boot import mode

samples = 20
delay = 0.2

led = digitalio.DigitalInOut(board.LED)
led.direction = digitalio.Direction.OUTPUT
led.value = True

scl = board.GP1
sda = board.GP0
i2c = busio.I2C(scl, sda)
fxos = adafruit_fxos8700.FXOS8700(i2c)

if mode == "logger":
    file = open('data.csv', 'w')
    file.write('time accel')

while True:
    accel = 0
    for i in range(samples):
        accel += fxos.accelerometer[2]
        time.sleep(delay / samples)
    avg = accel / samples
    
    if mode == 'serial':
        print(f'{avg:.2f}')
    elif mode == 'logger':
        try:
            file.write(f'{time.monotonic()} {avg}')
        except:
            file.close()