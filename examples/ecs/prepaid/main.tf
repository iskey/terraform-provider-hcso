data "hcso_availability_zones" "myaz" {}

data "hcso_compute_flavors" "myflavor" {
  availability_zone = data.hcso_availability_zones.myaz.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "hcso_vpc_subnet" "mynet" {
  name = "subnet-default"
}

data "hcso_images_image" "myimage" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

resource "hcso_compute_instance" "basic" {
  name              = "basic"
  image_id          = data.hcso_images_image.myimage.id
  flavor_id         = data.hcso_compute_flavors.myflavor.ids[0]
  security_groups   = ["default"]
  availability_zone = data.hcso_availability_zones.myaz.names[0]
  admin_pass        = "Test@123"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1

  network {
    uuid = data.hcso_vpc_subnet.mynet.id
  }
}
