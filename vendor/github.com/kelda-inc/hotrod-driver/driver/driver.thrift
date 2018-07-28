
struct DriverLocation {
  1: required string  driver_id
  2: required string  location
}

struct Result {
  1: optional bool  result
}

service Driver  {
    list<DriverLocation> findNearest(1: string location)
    Result lock(1: string id)
    Result unlock(1: string id)
}