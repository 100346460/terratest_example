locals {
    appended = join("-", compact(concat([var.name], [random_string.suffix.result])))
}

resource "random_string" "suffix" {
    length    = 4
    min_lower = 4
    special   = false
}

output "suffix" {
    value = random_string.suffix.result
}
output "result" {
    value = var.add_random_suffix ? local.appended : var.name
}