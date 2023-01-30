// Library is organized in a clean or onion architecture. The library package
// holds common entities: models, errors and interfaces that modules share with
// each other.
//
// Submodules typically only import from this project and not other submodules
// directly (utilities are an exception). Dependencies are specified as
// interfaces, usually in the submodule itself but can be added to the root
// library package if impractical.
//
// Testing mocks are generated using [gomock](github.com/golang/mock/gomock).
//
// Utilities typically use global state and are set up using [compile-time
// plugins](https://eli.thegreenplace.net/2021/plugins-in-go/).
package library
