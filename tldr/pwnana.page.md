# pwnana

> Pronounceable, keyboard-alternating password generator.
> Generates passwords that alternate QWERTY hands and can include leet digits, symbols, or Japanese goroawase hints.
> More information: <https://github.com/kusanaginoturugi/pwnana>.

- Generate passwords with the default length, filling about half of the terminal height:

`pwnana`

- Generate passwords with a specific length:

`pwnana {{20}}`

- Generate a specific count of passwords (odd counts are rounded up):

`pwnana {{20}} {{5}}`

- Generate passwords that include one digit:

`pwnana -d`

- Generate passwords that include one symbol:

`pwnana -s`

- Generate passwords that include both a digit and a symbol:

`pwnana {{24}} -d -s`

- Generate passwords with embedded goroawase digits and hints:

`pwnana -g`
