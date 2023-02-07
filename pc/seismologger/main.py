import time
import serial
import matplotlib.pyplot as plt
import matplotlib.animation as animation

# Parameters
x_len = 300         # Number of points to display
y_range = [8, 12]   # Range of possible Y values to display

# Create figure for plotting
fig = plt.figure()
ax = fig.add_subplot(1, 1, 1)
xs = list(range(0, 300))
ys = [0] * x_len
ax.set_ylim(y_range)

# Initialize serial communication
ser = serial.Serial('/dev/ttyACM0', 115200, timeout=1)
time.sleep(1)
ser.read_all()

# Create a blank line to be updated in animate()
line, = ax.plot(xs, ys)

# Add labels
plt.title('Seismograf')
plt.xlabel('Samples')
plt.ylabel('Acceleration (m/s^2)')


# This function is called periodically from FuncAnimation
def animate(i, ys):
    # Read next packet of data
    data = ser.read_until(b'\r\n').decode().splitlines()[0]
    val = float(data)
    # Add y to list
    ys.append(val)
    # Limit y list to set number of items
    ys = ys[-x_len:]
    # Update line with new Y values
    line.set_ydata(ys)
    return line,


# Set up plot to call animate() function periodically
a = animation.FuncAnimation(fig, animate, fargs=(ys,), interval=40, blit=True)
plt.show()
