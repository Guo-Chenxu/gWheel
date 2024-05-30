go test                          # 运行所有测试
go test -v                       # 运行所有测试并输出详细信息, 可输出t.Log()内容
go test -run TestName            # 运行指定测试函数
go test name_test.go             # 运行指定测试文件
go test -coverprofile=c          # 运行测试并输出覆盖率信息到c.out文件
go tool cover -html=c            # 生成HTML格式的覆盖率报告
go test -bench=. -cpuprofile=cpu # 运行性能测试并输出CPU性能分析报告
go tool pprof cpu                # 分析CPU性能分析报告
