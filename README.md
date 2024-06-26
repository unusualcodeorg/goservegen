# goservegen - Go Backend Architecture Generator using goserve framework
Project generator for go backend architecture using goserve framework

Check out goserve framework [github.com/unusualcodeorg/goserve](https://github.com/unusualcodeorg/goserve)

## How To Use goservegen
1. Download the goservegen binary for your operating system from the goservegen latest release: [github.com/unusualcodeorg/goservegen/releases](https://github.com/unusualcodeorg/goservegen/releases)

2. Expand the compressed file (Example: Apple Mac M2: goservegen_Darwin_arm64.tar.gz)

3. Run the binary 
```bash
cd ~/Downloads/goservegen_Darwin_arm64

# ./goservegen [project directory path] [project module]
./goservegen ~/Downloads/example github.com/yourusername/example
```
Note: `./goservegen ~/Downloads/example github.com/yourusername/example` will generate project named `example` located at `~/Downloads` and module `github.com/yourusername/example`

4. Open the generated project in your IDE/editor of choice

5. Have fun developing your REST API server!

## Troubleshoot
Sometimes your operating system will block the binary from execution, you will have to provide permission to run it. 

In Mac you have to go System Settings > Privacy & Security > Allow goservegen