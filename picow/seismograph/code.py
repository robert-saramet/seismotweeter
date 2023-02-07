import time
import board
import busio
import digitalio
import adafruit_fxos8700

led = digitalio.DigitalInOut(board.LED)
led.direction = digitalio.Direction.OUTPUT
led.value = True

scl = board.GP1
sda = board.GP0
i2c = busio.I2C(scl, sda)
fxos = adafruit_fxos8700.FXOS8700(i2c)

while True:
    print(f'{fxos.accelerometer[2]:.2f}')

'''
The code is intentionally as barebones as possible to
maximise throughput when used with the plotter program
'''
