provider "digitalocean" {
	token = "${var.do_token}"
}

resource "digitalocean_droplet" "app" {
	image              = "coreos-stable"
	region             = "tor1"
	size               = "512mb"
	private_networking = true
	backups            = false
	ipv6               = true
	name               = "app-1"
	ssh_keys		   = "${var.ssh_key_ids}"

	provisioner "file" {
		source      = "audition.service"
		destination = "/tmp/audition.service"

		connection {
			type        = "ssh"
			private_key = "${file("${var.private_key}")}"
			user        = "core"
			timeout     = "2m"
		}
	}

	provisioner "remote-exec" {
		inline = [
			"/usr/bin/sudo mv /tmp/audition.service /etc/systemd/system/audition.service",
			"/usr/bin/sudo systemctl enable /etc/systemd/system/audition.service",
			"/usr/bin/sudo systemctl start audition.service"
		]

		connection {
			type        = "ssh"
			private_key = "${file("${var.private_key}")}"
			user        = "core"
			timeout     = "2m"
		}
	}
}

variable "do_token" {
	type    = "string"
	default = ""
}

variable "ssh_key_ids" {
	type    = "list"
	default = []
}

variable "private_key" {
	type    = "string"
	default = "~/.ssh/id_rsa"
}

output "Public IP" {
	value = "${digitalocean_droplet.app.ipv4_address}"
}

output "Name" {
  value = "${digitalocean_droplet.app.name}"
}
