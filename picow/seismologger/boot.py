import supervisor
import storage

if __name__ == 'main':
    new_name = "Seismograph"
    storage.remount("/", readonly=False)
    mnt = storage.getmount("/")
    mnt.label = new_name

if supervisor.runtime.usb_connected:
    mode = "serial"
else:
    mode = "logger"
    storage.remount("/", readonly=True)