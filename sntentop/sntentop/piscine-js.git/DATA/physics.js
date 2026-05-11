function getAcceleration(obj) {
  if (obj.f != null && obj.m != null) {
    return obj.f / obj.m
  }
  
  if (obj.Δv != null && obj.Δt != null) {
    return obj.Δv / obj.Δt
  }
  
  if (obj.d != null && obj.t != null) {
    return (2 * obj.d) / (obj.t ** 2)
  }
  
  return "impossible"
}
