# gosystract
gosystract returns the names and IDs of all system calls being called inside a go application.

Usage:

```sh
gosystrac goapp.dump
```

To generate a dump file from a go application use: 
```sh
go tool objdump goapp > goapp.dump
```


## License

This application is licensed under the MIT License, you may obtain a copy of it [here](LICENSE).