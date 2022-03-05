#ifndef _CONFIG_HPP_
#define _CONFIG_HPP_

/* how frequently to take a temperature measurement in seconds */
constexpr unsigned long TEMP_INTERVAL_SEC = 3;

/* how long the LEDs should be on in milliseconds */
constexpr unsigned long LED_FLASH_DURATION_MS = 500;

/* how long to debounce regular buttons */
constexpr unsigned long BUTTON_DEBOUNCE_MS = 50;

/* how long to debounce raingauge */
constexpr unsigned long RAINGAUGE_DEBOUNCE_MS = 50;

/* how much does the gauge measure in millimeters */
constexpr float RAINGAUGE_AMT_MM = 0.2794;

#endif
