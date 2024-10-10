### JDK Version Switcher

---

# JDK Version Switcher

<<<<<<< HEAD
**JDK Version Switcher** 是一个用 Go 语言编写的工具，旨在帮助 Windows 用户快速切换不同版本的 JDK（Java Development Kit）。该工具通过读取 Windows 注册表中已安装的 JDK 版本，并结合用户自定义的环境变量，提供一个交互式的终端界面，方便用户选择并设置所需的 JDK 版本。

## 特性

- **自动检测 JDK 版本**：从 Windows 注册表中读取已安装的 JDK 版本（支持 32 位和 64 位）。
- **手动添加 JDK 版本**：通过环境变量 `JAVA_SDK_{version}` 手动添加 JDK 版本。
- **交互式选择菜单**：使用 `promptui` 库提供的交互式菜单，支持键盘上的上下箭头键导航。
- **颜色化终端输出**：通过 `fatih/color` 库美化终端输出，使重要信息更易辨识。
- **修改环境变量**：自动设置系统级别的 `JAVA_HOME` 和更新 `PATH` 环境变量。

## 环境要求

- **操作系统**：Windows 10 或更高版本
- **依赖**：
  - [Go](https://golang.org/dl/) 1.16 或更高版本
  - [promptui](https://github.com/manifoldco/promptui) 库
  - [fatih/color](https://github.com/fatih/color) 库
  - [golang.org/x/sys/windows/registry](https://pkg.go.dev/golang.org/x/sys/windows/registry) 库
- **权限**：程序需要以管理员权限运行，以修改系统环境变量

## 安装步骤

### 1. 克隆项目代码

```bash
git clone https://github.com/your-username/jdk-switcher.git
```

### 2. 进入项目目录

```bash
cd jdk-switcher
```

### 3. 安装依赖

使用以下命令安装所需的第三方库：

```bash
go get github.com/manifoldco/promptui
go get github.com/fatih/color
go get golang.org/x/sys/windows/registry
```

### 4. 编译项目

运行以下命令编译项目：

```bash
go build -o jdk-switcher.exe main.go
```

### 5. 以管理员权限运行

由于该工具需要修改系统环境变量，请确保以管理员权限运行。

- **方法一**：通过文件资源管理器，右键点击 `jdk-switcher.exe`，选择 **“以管理员身份运行”**。
- **方法二**：在已提升权限的命令提示符或 PowerShell 中运行：

  ```bash
  .\jdk-switcher.exe
  ```

## 使用方法

1. **启动程序**：

   运行编译后的 `jdk-switcher.exe`，程序将自动检测系统中安装的 JDK 版本，并显示一个交互式选择菜单。

   ```bash
   **************************************************
   1. 该软件会自动搜索通过标准方式安装的 JDK 版本。
   2. 若需手动添加，请在系统环境变量中创建类似 `JAVA_SDK_{版本号}` 的变量，并指向相应的 JDK 安装目录。
   3. 在系统环境变量的 `PATH` 中，顶部添加 `%JAVA_HOME%\bin`。
   **************************************************
   Current JDK Version: C:\Program Files\Java\jdk-17.0.2

   ? 请选择一个 JDK 版本: (Use arrow keys)
    ❯ Exit
      (x86) 1.8.0_281	Path: C:\Program Files\Java\jdk1.8.0_281
      (x64) 11.0.10	Path: C:\Program Files\Java\jdk-11.0.10
      (env) JAVA_SDK_17	Path: C:\Java\jdk-17.0.2
   ```

2. **选择 JDK 版本**：

   使用键盘的上下箭头键导航菜单，选择您希望设置为 `JAVA_HOME` 的 JDK 版本。按 `Enter` 键确认选择。

3. **程序执行**：

   选择完成后，程序将自动更新 `JAVA_HOME` 和 `PATH` 环境变量。

   ```bash
   使用 C:\Program Files\Java\jdk-17.0.2
   成功设置 JAVA_HOME 为 C:\Program Files\Java\jdk-17.0.2
   成功更新 PATH。
   更改已成功应用。您可能需要重新启动终端或重新登录系统以使更改生效。
   ```

4. **退出程序**：

   如果选择“Exit”选项，程序将退出。

## 终端支持

该工具的输出包含颜色信息，确保您的终端支持 ANSI 转义码。推荐使用以下终端以获得最佳体验：

- [Windows Terminal](https://aka.ms/terminal)
- [PowerShell](https://docs.microsoft.com/powershell/)
- 支持颜色输出的其他终端模拟器

**注意**：旧版的 Windows 命令提示符（CMD）可能不完全支持颜色化输出，建议使用上述推荐的终端。

## 注意事项

1. **管理员权限**：
   - 由于工具需要修改系统级别的环境变量，因此请确保以管理员权限运行程序。
   - 运行程序前，可以右键点击终端程序图标，选择“以管理员身份运行”。

2. **重启终端**：
   - 更新 `JAVA_HOME` 和 `PATH` 后，某些终端可能不会立即反映更改，建议重启终端或重新登录系统。

3. **注册表访问**：
   - 该工具会从注册表中读取已安装的 JDK 信息，因此需要访问 Windows 注册表的权限。

4. **环境变量冲突**：
   - 确保环境变量中的 `JAVA_HOME` 和 `PATH` 不存在冲突，以避免潜在的问题。

## 许可证

本项目采用 [MIT 许可证](LICENSE)，详细信息请参阅 [LICENSE](LICENSE) 文件。
=======
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
>>>>>>> e38ca5ec0bd44f68ccc40bb555a0fe2b846883e8
