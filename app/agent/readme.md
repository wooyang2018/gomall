# *** Project

## introduce

- Use the [Kitex](https://github.com/cloudwego/kitex/) framework
- Generating the base code for unit tests.
- Provides basic config functions
- Provides the most basic MVC code hierarchy.

## Directory structure

| catalog     | introduce                                       |
|-------------|-------------------------------------------------|
| conf        | Configuration files                             |
| main.go     | Startup file                                    |
| handler.go  | Used for request processing return of response. |
| kitex_gen   | kitex generated code                            |
| biz/service | The actual business logic.                      |
| biz/dal     | Logic for operating the storage layer           |

## How to run

```shell
sh build.sh
sh output/bootstrap.sh
```
