# arduino_rain_gauge

control a rain gauge using an arduino

### Getting a Dev Environment

Workflow is largely based around `arduino-cli` and `make`. Make sure all
required packages are installed.

Make sure you're a member of the `dialout` and `plugdev` groups and reboot your
machine.

```bash
user@host:~$ sudo apt-get install -y build-essentials clang-format
user@host:~$ sudo usermod -aG dialout $USER && sudo usermod -aG plugdev $USER
user@host:~$ sudo reboot 0
```

##### arduino-cli

See [install guide](https://arduino.github.io/arduino-cli/latest/installation/)
to get binary installed.

Plug in arduino and update the core libraries:

```bash
user@host:~$ arduino-cli core update-index
```

Locate the board for your arduino with `arduino-cli board list` and install the
appropriate board:

```bash
user@host:~$ arduino-cli board list
user@host:~$ arduino-cli core install <`Core` from previous command>
  ...
user@host:~$ arduino-cli core list # should show your core
```

Install the proper third-party libraries

```bash
user@host:~$ arduino-cli lib install LiquidCrystal
```

##### update udev

See comments in `udev/` directory for how to find and update udev rules for your
board. Symlink the file to `/lib/udev/rules.d/` or `/etc/udev/rules.d`

Reload udev and find the device:

```bash
user@host:~$ sudo udevadm control --reload && sudo udevadm trigger
user@host:~$ ls -al /dev/ttyACM99  # should see your board listed
```

#### put code on your arduino

Use `make` to compile and upload your code

```bash
user@host:~/arduino-rain-gauge $ make build   # compiles the code
user@host:~/arduino-rain-gauge $ make upload  # uploads to the code to the arduino
user@host:~/arduino-rain-gauge $ make all     # runs build and then upload
```
