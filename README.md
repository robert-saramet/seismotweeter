# Seismotweeter

### Non-technical introduction

This project is a demonstration of a basic "seismograph" implemented using a microcontroller and an accelerometer module. While the device is rather accurate, it is not a sensible substitute for an actual seismometer, which are often implanted deep underground to isolate them from everyday mechanical shocks. As my approach considers only acceleration, it will treat e.g. nearby footsteps the same as earthquakes.

With this in mind, feel free to check it out yourself with the `go-plot` program or take a look at the **[seismotweeter](https://twitter.com/seismotweeter)** bot's Twitter page.

### Operating modes

There are multiple programs in these repository, divided as such:

- **seismograph**: reads accelerometer on Pico W and sends data over serial as fast as possible to the PC, which runs the `go-plot` program to receive data and plot it in realtime using custom raylib code (#graphsIn200FPS). To run this program, copy the `picow/seismograph` folder to the `CIRCUITPY` drive and run `go-plot` on your PC; this mode is useful for monitoring the seismic activity live
  
- **seismologger**: when connected to a computer, it behaves like the `seismograph` program, but with significantly lower throughput which results in unbearable framerates; on the other hand, when usb data is not connected, it logs all data to a csv file on the `CIRCUITPY` drive, which can later be copied to a computer for interpretation; this mode is useful for headless operation, where immediate analysis of the data is not necessary
  
- **seismotweeter**: this one is one-of-a-kind; it connects to a local Wi-Fi network and sends web requests to IFTTT when it detects anomalous seismic activity, with IFTTT then creating a tweet about the event; it has an impressively low latency of under 2 seconds, but I've decided to add 5s cooldowns between requests to avoid spamming the Twitter API
  

For more details about each program I encourage you to read the code, I've tried to make it as readable as possible.

### Reproducing this project

- The specific devices I have used are a Raspberry Pi Pico W and an NXP 9-DoF IMU from Adafruit. However, you can use any accelerometer and microcontroller with Circuitpython support with minimal changes to your code.
  
- For **all** of these projects to work out-of-the-box you will need to be running Circuitpython 8 on the Pico W. Then for each program just copy its files to the `CIRCUITPY` drive, erasing anything else.
  
- For the **seismograph** program, on PC you might need to change the serial port used by `go-plot`. You can do so by editing the `port0` and `port1` values at the top of its source code.
  
- For the **seismologger**, keep in mind that data logging mode will not be available when connected via usb for data. This is caused by a feature of Circuitpython which prevents writes to the `CIRCUITPY` drive when usb is connected, as this might cause filesystem corruption
  
- For the **seismotweeter**, you will first need to follow [this tutorial](https://www.tomshardware.com/how-to/connect-raspberry-pi-pico-w-to-twitter-via-ifttt) in order to set up IFTTT. Then make sure to edit the `settings.toml` file on the Pico W with your own Wi-Fi and IFTTT details.
