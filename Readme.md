Attempting to Build JS Bindings to Ikemen Go
- Ikemen-Go - https://github.com/ikemen-engine/Ikemen-GO
- FileSystem - https://www.npmjs.com/package/browserfs
- Gamepad - https://developer.mozilla.org/en-US/docs/Web/API/Gamepad_API/Using_the_Gamepad_API
- Audio
- Graphics

# Build

- Zip up necessary Mugen Files to `./static/hidden.mugen_base.zip`
  - **This includes the directories**
    - `/chars`
    - `/data`
    - `/external`
    - `/font`
    - `/save`
    - `/sound`
    - `/stages`
- [Install Go](https://go.dev/doc/install)
- [Install Docker Compose](https://docs.docker.com/compose/install/)
- [Install Node](https://nodejs.org/en)
- terminal - `cd path/to/this/repo`
- To Build Go
  - terminal - `OOS=js GOARCH=wasm go build -o ./static/dist.wasm ./go/src/main.go ./go/src/testing.go`
- To Build Bindings
  - terminal - `cd ./web`
  - terminal - `npm run build`

# Run
- terminal - `cd path/to/this/repo`
- terminal - `cd build`
- terminal - `docker-compose up`
- web browser - go to `http://localhost:8080`