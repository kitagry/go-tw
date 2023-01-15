## go-tw

normal output

```bash
$ go test ./...                                                                                                                                                                                                             (raspi/kube-system)
?       github.com/kitagry/repository       [no test files]
--- FAIL: TestYourMethod (0.00s)
    --- FAIL: TestYourMethod/run_a (0.00s)
        huga_test.go:209: error message
    --- FAIL: TestYourMethod/run_b (0.00s)
        huga_test.go:209: error message
FAIL
FAIL    github.com/kitagry/repository/hoge    0.006s
FAIL
```

go-tw output

```bash
$ go-tw ./...                                                                                                                                                                                                             (raspi/kube-system)
?       github.com/kitagry/repository       [no test files]
--- FAIL: TestYourMethod (0.00s)
    --- FAIL: TestYourMethod/run_a (0.00s)
        /home/path/to/repository/hoge/huga_test.go:209: error message
    --- FAIL: TestYourMethod/run_b (0.00s)
        /home/path/to/repository/hoge/huga_test.go:209: error message
FAIL
FAIL    github.com/kitagry/repository/hoge    0.006s
FAIL
```

diff

```diff
$ go-tw ./...                                                                                                                                                                                                             (raspi/kube-system)
?       github.com/kitagry/repository       [no test files]
--- FAIL: TestYourMethod (0.00s)
    --- FAIL: TestYourMethod/run_a (0.00s)
-        huga_test.go:209: error message
+        /home/path/to/repository/hoge/huga_test.go:209: error message
    --- FAIL: TestYourMethod/run_b (0.00s)
-        huga_test.go:209: error message
+        /home/path/to/repository/hoge/huga_test.go:209: error message
FAIL
FAIL    github.com/kitagry/repository/hoge    0.006s
FAIL
```
