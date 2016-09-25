# -*- mode: ruby -*-
# vi: set ft=ruby :

ENV['VAGRANT_DEFAULT_PROVIDER'] = 'docker'

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.ssh.username = "develop"
  config.ssh.password = "bita123"

  config.vm.provision "shell" do |s|
    s.path   = "bin/provision.sh"
    s.args   = [%x(ip addr | grep inet | grep docker0 | awk -F" " '{print $2}'| sed -e 's/\\/.*$//')]
  end

  config.vm.network "forwarded_port", guest: 80,    host: 80    # nginx
  config.vm.synced_folder ".", "/home/develop/gad", owner: "develop", group: "develop", create: true

  config.vm.provider "docker" do |d|
    d.image = "docker.clickyab.com/clickyab/baseimage-go"
    d.has_ssh = true
    d.cmd = ["/bin/bash", "/home/develop/gad/bin/init.sh"]
    d.expose= [5432]
  end
end
