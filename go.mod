module github.com/gregorgebhardt/redblack/v2

go 1.23.0

toolchain go1.23.2

retract (
  // let's stay at v0
  v2.0.0
  v2.0.1
)

require golang.org/x/exp v0.0.0-20241009180824-f66d83c29e7c
