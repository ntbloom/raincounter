# tlv schema

tag and length are `unsigned char`

## TLV packets sent by arduino

| DESCRIPTION      | TAG | LENGTH | VALUE                | VALUE TYPE    |
| ---------------- | --- | ------ | -------------------- | ------------- |
| rain event       | 0   | 1      | 1                    | unsigned char |
| temperature      | 1   | 4      | temp in C            | int           |
| RESERVED         | 2   | NA     | NA                   | NA            |
| hard reset event | 3   | 1      | 1, sent after reboot | unsigned char |
| start pause      | 4   | 1      | 1                    | unsigned char |
| stop pause       | 5   | 1      | 1                    | unsigned char |
| RESERVED         | 6   | NA     | NA                   | NA            |
| RESERVED         | 7   | NA     | NA                   | NA            |

## TLV packets received by arduino

Serial communication is currently one-way from the arduino to the host computer.
No support currently for receiving serial communication from the computer.
