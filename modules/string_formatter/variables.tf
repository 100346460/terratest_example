variable "name" {
    description = "input string to append the random suffix to"
    type        = string
}

variable "add_random_suffix" {
    description = "Whether to add a random suffix or just pass through"
    type        = bool
    default     = false
}
