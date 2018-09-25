# Pending

Special thanks to external contributors on this release:

BREAKING CHANGES:

* CLI/RPC/Config

* Apps

* Go API
- [node] Remove node.RunForever
- [config] \#2232 timeouts as time.Duration, not ints

FEATURES:

IMPROVEMENTS:
- [config] \#2232 added ValidateBasic method, which performs basic checks

BUG FIXES:
- [node] \#2434 Make node respond to signal interrupts while sleeping for genesis time
