# Ketikode Executor

## Building the Source

```
cargo build
```

## Adding a New Programming Language

Currently, the programming language available are

- C++ - `g++`
- Java - `javac`
- Python3 - `python`

You can configure the compiler/intepreter by yourself on [Language.toml](Language.toml)

## Time limit

Implemented using `timeout`

```
timeout 3s ./exec
```

## Sandboxing

The program execution was sandboxed using [firejail](https://github.com/netblue30/firejail/)

## License

[Mozilla Public License Version 2.0](LICENSE)

## Author

- **Andra Antariksa**