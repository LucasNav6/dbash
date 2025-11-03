# 🧘‍♂️ Hush — Shell Scripting in Human Language

**Hush** is a **Domain-Specific Language (DSL)** designed to replace the complexity of traditional Bash scripts with a clean, readable, and declarative syntax. Inspired by Dockerfile, but built specifically for **local automation workflows**, Hush allows you to express shell tasks as if you were writing human-readable instructions — while still executing safely and efficiently.

> **Hush = Humanized Bash**  
Read it like documentation, run it like a validated and automated shell script.

---

## ✨ Key Features

- 🧠 Human-friendly, declarative syntax  
- 🔍 Automatic pre-execution validations  
- 🔒 Early requirement detection (sudo, binaries, versions)  
- 🧵 String-first variable system with simple interpolation  
- 🧰 Built-in commands for Git, directories, execution, and more  
- 🧪 Reproducible, predictable scripting — no more cryptic Bash files  

---

## 📦 Installation

```bash
# Coming soon…
hush install
```
> Installation instructions and package manager support will be added once the project is released

## 🤝 Open Source Project

Hush is built as an open-source and community-driven project.
The goal is to create a scripting language that anyone can read and understand, even without Bash expertise.

Contributions are welcome, including:
- 🐛 Bug Reports
- 💡 Feature Requests
- 🧪 Test Scripts & Examples
- 🧩 Language Design Feedback

If you’d like to shape the evolution of Hush, join the discussion and contribute!

---

## 🧱 Hush Language — Core Instructions
Below is a quick reference to the main instructions supported by Hush.

### REQUIRE_SUDO

Declares that the script requires elevated privileges.
Hush will prompt for sudo before executing anything.

```bash
REQUIRE_SUDO
```

### REQUIRE <binary[@version]>

Validates system requirements.
If a binary is missing or its version is lower than required, execution stops early with a clear message.

```bash
#Example	        Meaning
REQUIRE git@2.0.0	#Requires Git ≥ 2.0.0
REQUIRE go	        #Requires Go (any version)
REQUIRE make	    #Requires Make (any version)
REQUIRE git@2.0.0
REQUIRE go
REQUIRE make
```

### ARGS $var = value

Declares a variable.
>All variables are always strings.

```bash
ARGS $version = 12.5.0
```

### SET "<prompt>" $var

Prompts the user for input.
> If left empty, the previous value remains unchanged.

Supports interpolation with {{$variable}}.

```bash
SET "Which version do you want to install? (Default {{$version}}): " $version
```

### CLONE <repo-url>

Clones a public GitHub repository.
*Private repos are not supported yet. (comming soon)*

```bash
CLONE https://github.com/gravitational/teleport.git
```

### WORKDIR <path>

Changes the working directory.

```bash
WORKDIR folder
```

### CHECKOUT <tag|branch>

Checks out a Git tag or branch.
Supports ${} variable interpolation.

```bash
CHECKOUT v${version}
```

### RUN <command>

Executes a terminal command.

```bash
RUN make build/tsh
RUN ./build/tsh --version
```