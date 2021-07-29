# tlv schema

tag and length are `unsigned char`

## TLV packets sent by arduino

| DESCRIPTION      | TAG | LENGTH | VALUE                | VALUE TYPE    |
| ---------------- | --- | ------ | -------------------- | ------------- |
| rain event       | 0   | 1      | 1                    | unsigned char |
| temperature      | 1   | 4      | temp in C            | int           |
| soft reset event | 2   | 1      | 1                    | unsigned char |
| hard reset event | 3   | 1      | 1, sent after reboot | unsigned char |
| pause            | 4   | 1      | 1                    | unsigned char |
| unpause          | 5   | 1      | 1                    | unsigned char |
| reserved         | 6   | n/a    | n/a                  | n/a           |
| reserved         | 7   | n/a    | n/a                  | n/a           |

## TLV packets received by arduino

Serial communication is currently one-way from the arduino to the host computer. No support currently for receiving UART
communication from the computer.
