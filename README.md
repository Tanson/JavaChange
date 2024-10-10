### README 文件示例

---

# JDK Version Switcher

这是一个用 Go 编写的工具，用于在 Windows 系统上快速切换不同版本的 JDK（Java Development Kit）。该工具会自动检测通过标准方式安装的 JDK 版本，并允许用户选择需要设置为 `JAVA_HOME` 和 `PATH` 环境变量的 JDK 版本。

## 特性

- 自动从 Windows 注册表读取安装的 JDK 版本（32 位和 64 位）。
- 支持通过环境变量 (`JAVA_SDK_{version}`) 手动添加 JDK 版本。
- 修改系统级别的 `JAVA_HOME` 和 `PATH` 环境变量。
- 终端输出带颜色的提示信息，方便识别当前 JDK 版本（需要支持 ANSI 的终端）。
  
## 环境要求

- **操作系统**：Windows 10 或更新版本
- **依赖**：需要安装 Go 语言开发环境
- **权限**：程序需要管理员权限以修改系统环境变量

## 安装步骤

1. **克隆项目代码**：

    ```bash
    git clone https://github.com/your-username/jdk-switcher.git
    ```

2. **进入项目目录**：

    ```bash
    cd jdk-switcher
    ```

3. **安装依赖**：

    项目依赖 `golang.org/x/sys/windows/registry` 包。如果要使用带颜色输出的版本，还需要安装 `github.com/fatih/color` 包。

    使用以下命令安装依赖：

    ```bash
    go get golang.org/x/sys/windows/registry
    go get github.com/fatih/color
    ```

4. **编译项目**：

    运行以下命令编译项目：

    ```bash
    go build -o jdk-switcher main.go
    ```

5. **以管理员权限运行**：

    由于该工具需要修改系统环境变量，请确保以管理员权限运行。

    ```bash
    ./jdk-switcher
    ```

## 使用方法

1. 启动程序后，它会自动检测系统中安装的 JDK 版本，并显示一个菜单供用户选择。

    ```bash
    **************************************************
    1. The software automatically searches for JDK versions installed via standard methods.
    2. To add manually, create system environment variables for all JDK versions with names like: JAVA_SDK_{version}
    3. Add `%JAVA_HOME%\bin` at the top of the `PATH` system environment variable.
    **************************************************

    Current JDK Version: C:\Program Files\Java\jdk1.8.0_271

    Please select a JDK version:
    0    : Exit
    1    : (x64) 1.8.0_271   Path: C:\Program Files\Java\jdk1.8.0_271
    2    : (x64) 11.0.9.1    Path: C:\Program Files\Java\jdk-11.0.9.1
    3    : (env) JAVA_SDK_1.8    Path: C:\Java\jdk1.8
    ```

2. 输入您想要切换到的 JDK 版本编号，程序会自动更新 `JAVA_HOME` 和 `PATH`。

3. 如果需要手动添加 JDK 版本，可以在系统环境变量中创建类似 `JAVA_SDK_{version}` 的变量名，并将其路径指向所需的 JDK 安装目录。

4. 选择完成后，程序将更新环境变量，您可能需要重新启动终端或重新登录系统以使更改生效。

## 终端支持

该工具的输出包含颜色信息，确保您的终端支持 ANSI 转义码。如果您使用的是旧版 Windows 命令提示符（CMD），建议切换到 Windows Terminal 或使用支持颜色输出的终端。

## 注意事项

1. **管理员权限**：由于工具需要修改系统级别的环境变量，因此请确保以管理员权限运行程序。
2. **重启终端**：更新 `JAVA_HOME` 和 `PATH` 后，某些终端可能不会立即反映更改，建议重启终端或重新登录系统。
3. **注册表访问**：该工具会从注册表中读取已安装的 JDK 信息，因此需要访问 Windows 注册表的权限。

## 开发与贡献

如有任何问题或建议，欢迎提交 Issue 或 Pull Request。

## 许可证

本项目采用 MIT 许可证，详细信息请参阅 [LICENSE](LICENSE) 文件。
