## Getting Started
* Install Android SDK and NDK
* `go install golang.org/x/mobile/cmd/gomobile@latest`
* `gomobile init`
* `./build.sh`

## Project Structure
The `mobile` folder contains some basic interfaces and definitions that help form the structure for the shape the business logic of an app should take. This is all very lightweight and only provides loose structure. The intent is that a more experienced engineer would help define codebase pattern and would guide and coach other engineers as opposed to having a rigid framework to work within.

The `app` folder contains app specific needs and implementation. Any functionality not providable by Go needs an interface that would require a platform-specific implementation. An example of this is Firebase Remote Config. `app/firebase/remote_config.go` defines an interface to provide a means for Go to retrieve remote config values and a method to fetch and activate the configuration itself. This definition is very basic and could be expanded on, but it is all this small app needs and would suffice for many larger applications as well.

Due to limitations in what types can be shared between Go and the target platform with Go Mobile, notice that the FetchAndActivate interface method receives a struct with a callback function defined on it rather than a direct function parameter.

Business / domain logic is kept in the `logic` folder. In this case, some structures for managing array and map values are defined here. These are necessary because no map or dictionary values can be shared between Go and iOS or Android and the only array-like values allowed are byte slices (`[]byte`).

The view models are kept in the `viewmodel` folder. They simply define actions the UI can provide. Each action may affect state which then notifies any observers which may update the UI.

## Conclusions
I found Go Mobile to be an effective way of keeping a separation of business / domain logic and UI. While I have not yet introduced any testing, this approach would allow for much easier testing. One could even create a CLI tool where the tool acts as a platform and is scriptable to test app flows. Thus each platform only needs to have its UI tested.

This approach can also be an effective way of developing business logic once while utilizing it on more than one platform. The result would be more parity between iOS and Android implementations and potentially faster feature iteration since the logic is implemented once instead of once per platform.

The main negatives lie with debugging. I have only just started looking into how to debug the Go logic, but it seems like it may be a bit tricky.
