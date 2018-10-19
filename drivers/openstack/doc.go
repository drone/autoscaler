/*
Package openstack contains a autoscaler driver for OpenStack
Configuration:

Authenticate with the usual OpenStack environment variables.
(Not all of these may be necessary:
see https://github.com/gophercloud/gophercloud/blob/master/openstack/auth_env.go)

OS_AUTH_URL=https://my.openstack.cloud:5000
OS_ENDPOINT_TYPE=publicURL
OS_IDENTITY_API_VERSION=2
OS_PASSWORD=<mypassword>
OS_DOMAIN_ID=default
OS_REGION_NAME=my-region
OS_TENANT_ID=my-tenant-id
OS_TENANT_NAME=my-tenant-name
OS_USERNAME=my-username

Configure driver with:
DRONE_OPENSTACK_SSHKEY=drone-key-name
DRONE_OPENSTACK_SECURITY_GROUP=my-security-group
# Pool for floating ips
DRONE_OPENSTACK_IP_POOL=my-ip-pool
DRONE_OPENSTACK_FLAVOR=v1-standard-2
DRONE_OPENSTACK_IMAGE=ubuntu-16.04-server-latest
DRONE_OPENSTACK_METADATA=name:agent,owner:drone-ci

*/
package openstack
