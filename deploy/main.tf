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


	provisioner "remote-exec" {
		inline = [
			"docker run -d --name audition --restart unless-stopped -p 80:8080 -v $(pwd)/db:/db arbourd/audition:latest",
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
