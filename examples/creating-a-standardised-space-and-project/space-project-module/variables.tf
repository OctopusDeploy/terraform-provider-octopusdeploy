variable "space" {
  description = "Space details"
  type = object({
    name        = string
    description = optional(string, "")
  })
}

variable "project" {
  description = "Application details"
  type = object({
    name        = string
    description = optional(string, "")
  })
}