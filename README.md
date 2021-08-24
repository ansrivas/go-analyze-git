go-analyze-git:
---

This project can be used to analyze git-data from data files exported from git.


## Installation:
```bash
go get -u gitlab.com/ansrivas/go-analyze-git
```
Local Build ( check the `build` folder after running following command)
```
make release
```

## Usage:

```bash
❯ ./go-analyze-git 
                                                  _                                        _   _
   __ _    ___             __ _   _ __     __ _  | |  _   _   ____   ___            __ _  (_) | |_
  / _` |  / _ \   _____   / _` | | '_ \   / _` | | | | | | | |_  /  / _ \  _____   / _` | | | | __|
 | (_| | | (_) | |_____| | (_| | | | | | | (_| | | | | |_| |  / /  |  __/ |_____| | (_| | | | | |_
  \__, |  \___/           \__,_| |_| |_|  \__,_| |_|  \__, | /___|  \___|          \__, | |_|  \__|
  |___/                                               |___/                        |___/
NAME:
   go-analyze-git - A new cli application

USAGE:
   go-analyze-git [global options] command [command options] [arguments...]

DESCRIPTION:
   Run analytics on git data set

AUTHOR:
   Ankur Srivastava <ankur.srivastava@email.de>

COMMANDS:
   repository, r  Commands related to repository operations
   user, u        Commands related to user operations
   help, h        Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
```


## Examples:

### NOTE: 
If you are running the examples from `data` folder, make sure to run `git lfs pull` first. CSV files are
committed as `large files`.


1. To print in tabular format
    ```bash
    ./go-analyze-git repository topk-by-events --events-file ./data/events.csv --repos-file ./data/repos.csv 
    +-------------------------------------+-------+
    |               REPOID                | COUNT |
    +-------------------------------------+-------+
    | victorqribeiro/isocity              |    44 |
    | neutraltone/awesome-stock-resources |    11 |
    | GitHubDaily/GitHubDaily             |    11 |
    | sw-yx/spark-joy                     |    10 |
    | imsnif/bandwhich                    |     8 |
    | Chakazul/Lenia                      |     7 |
    | BurntSushi/xsv                      |     7 |
    | FiloSottile/age                     |     6 |
    | neeru1207/AI_Sudoku                 |     6 |
    | ErikCH/DevYouTubeList               |     6 |
    +-------------------------------------+-------+
    ```

2. To print in json format
    ```bash
    ./go-analyze-git repository tw --events-file ./data/events.csv --repos-file ./data/repos.csv  --json | jq
    ```

3. Top-k by commits:
   ```
    ./go-analyze-git --debug repository topk-by-commits --events-file ./data/events.csv --repos-file ./data/repos.csv --commits-file ./data/commits.csv
    ```

4. User operations
   ```
   ./go-analyze-git --debug user topk-by-pc --events-file=./data/events.csv --commits-file=./data/commits.csv --actors-file=./data/actors.csv --count 10
   ```

## Tests
To run tests:
   `make test`
## Benchmarks:
   To run benchmarks simply use:
   ```
   go test -cpuprofile cpu.prof -memprofile mem.prof -bench BenchmarkTop ./pkg/repository/
   go tool pprof mem.prof
   # Inside the prompt type `web`
   go tool pprof cpu.prof
   # Inside the prompt type `web`


   go test -cpuprofile cpu.prof -memprofile mem.prof -bench BenchmarkTop ./pkg/user/
   go tool pprof mem.prof
   # Inside the prompt type `web`
   go tool pprof cpu.prof
   # Inside the prompt type `web`
   ```
