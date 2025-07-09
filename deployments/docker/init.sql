-- Create database if not exists
SELECT 'CREATE DATABASE golang_service'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'golang_service')\gexec

-- Connect to the database
\c golang_service;

-- Create extensions if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create AWS EC2 instances table
CREATE TABLE IF NOT EXISTS aws_ec2_instances (
    _cq_sync_time timestamp without time zone,
    _cq_source_name text,
    _cq_id uuid NOT NULL,
    _cq_parent_id uuid,
    account_id text,
    region text,
    arn text NOT NULL,
    state_transition_reason_time timestamp without time zone,
    tags jsonb,
    ami_launch_index bigint,
    architecture text,
    block_device_mappings jsonb,
    boot_mode text,
    capacity_reservation_id text,
    capacity_reservation_specification jsonb,
    client_token text,
    cpu_options jsonb,
    current_instance_boot_mode text,
    ebs_optimized boolean,
    elastic_gpu_associations jsonb,
    elastic_inference_accelerator_associations jsonb,
    ena_support boolean,
    enclave_options jsonb,
    hibernation_options jsonb,
    hypervisor text,
    iam_instance_profile jsonb,
    image_id text,
    instance_id text,
    instance_lifecycle text,
    instance_type text,
    ipv6_address text,
    kernel_id text,
    key_name text,
    launch_time timestamp without time zone,
    licenses jsonb,
    maintenance_options jsonb,
    metadata_options jsonb,
    monitoring jsonb,
    network_interfaces jsonb,
    outpost_arn text,
    placement jsonb,
    platform text,
    platform_details text,
    private_dns_name text,
    private_dns_name_options jsonb,
    private_ip_address text,
    product_codes jsonb,
    public_dns_name text,
    public_ip_address text,
    ramdisk_id text,
    root_device_name text,
    root_device_type text,
    security_groups jsonb,
    source_dest_check boolean,
    spot_instance_request_id text,
    sriov_net_support text,
    state jsonb,
    state_reason jsonb,
    state_transition_reason text,
    subnet_id text,
    tpm_support text,
    usage_operation text,
    usage_operation_update_time timestamp without time zone,
    virtualization_type text,
    vpc_id text,
    PRIMARY KEY (arn)
);

-- Create Azure VM instances table
CREATE TABLE IF NOT EXISTS azure_compute_virtual_machines (
    _cq_sync_time timestamp without time zone,
    _cq_source_name text,
    _cq_id uuid NOT NULL,
    _cq_parent_id uuid,
    subscription_id text,
    instance_view jsonb,
    location text,
    extended_location jsonb,
    identity jsonb,
    plan jsonb,
    properties jsonb,
    tags jsonb,
    zones text[],
    id text NOT NULL,
    name text,
    resources jsonb,
    type text,
    PRIMARY KEY (id)
);

-- Create GCP Compute instances table
CREATE TABLE IF NOT EXISTS gcp_compute_instances (
    _cq_sync_time timestamp without time zone,
    _cq_source_name text,
    _cq_id uuid NOT NULL,
    _cq_parent_id uuid,
    project_id text,
    advanced_machine_features jsonb,
    can_ip_forward boolean,
    confidential_instance_config jsonb,
    cpu_platform text,
    creation_timestamp text,
    deletion_protection boolean,
    description text,
    disks jsonb,
    display_device jsonb,
    fingerprint text,
    guest_accelerators jsonb,
    hostname text,
    id bigint,
    instance_encryption_key jsonb,
    key_revocation_action_type text,
    kind text,
    label_fingerprint text,
    labels jsonb,
    last_start_timestamp text,
    last_stop_timestamp text,
    last_suspended_timestamp text,
    machine_type text,
    metadata jsonb,
    min_cpu_platform text,
    name text,
    network_interfaces jsonb,
    network_performance_config jsonb,
    params jsonb,
    private_ipv6_google_access text,
    reservation_affinity jsonb,
    resource_policies text[],
    resource_status jsonb,
    satisfies_pzs boolean,
    scheduling jsonb,
    self_link text NOT NULL,
    service_accounts jsonb,
    shielded_instance_config jsonb,
    shielded_instance_integrity_policy jsonb,
    source_machine_image text,
    source_machine_image_encryption_key jsonb,
    start_restricted boolean,
    status text,
    status_message text,
    tags jsonb,
    zone text,
    PRIMARY KEY (self_link)
);

-- Insert dummy data for AWS EC2 instances
INSERT INTO aws_ec2_instances (_cq_id, account_id, region, arn, instance_id, instance_type, private_ip_address, public_ip_address, vpc_id, subnet_id, launch_time, state) VALUES
(gen_random_uuid(), '123456789012', 'us-east-1', 'arn:aws:ec2:us-east-1:123456789012:instance/i-1234567890abcdef0', 'i-1234567890abcdef0', 't2.micro', '10.0.1.100', '54.123.45.67', 'vpc-12345678', 'subnet-12345678', NOW() - INTERVAL '30 days', '{"name": "running", "code": 16}'),
(gen_random_uuid(), '123456789012', 'us-east-1', 'arn:aws:ec2:us-east-1:123456789012:instance/i-1234567890abcdef1', 'i-1234567890abcdef1', 't2.small', '10.0.2.100', '54.123.45.68', 'vpc-12345678', 'subnet-87654321', NOW() - INTERVAL '25 days', '{"name": "running", "code": 16}'),
(gen_random_uuid(), '123456789012', 'us-west-2', 'arn:aws:ec2:us-west-2:123456789012:instance/i-1234567890abcdef2', 'i-1234567890abcdef2', 't2.medium', '10.0.3.100', '54.123.45.69', 'vpc-87654321', 'subnet-11111111', NOW() - INTERVAL '20 days', '{"name": "running", "code": 16}'),
(gen_random_uuid(), '123456789012', 'us-west-2', 'arn:aws:ec2:us-west-2:123456789012:instance/i-1234567890abcdef3', 'i-1234567890abcdef3', 't2.large', '10.0.4.100', NULL, 'vpc-87654321', 'subnet-22222222', NOW() - INTERVAL '15 days', '{"name": "stopped", "code": 80}'),
(gen_random_uuid(), '123456789012', 'eu-west-1', 'arn:aws:ec2:eu-west-1:123456789012:instance/i-1234567890abcdef4', 'i-1234567890abcdef4', 't2.xlarge', '10.0.5.100', '54.123.45.70', 'vpc-33333333', 'subnet-33333333', NOW() - INTERVAL '10 days', '{"name": "running", "code": 16}'),
(gen_random_uuid(), '123456789012', 'eu-west-1', 'arn:aws:ec2:eu-west-1:123456789012:instance/i-1234567890abcdef5', 'i-1234567890abcdef5', 't2.2xlarge', '10.0.6.100', '54.123.45.71', 'vpc-33333333', 'subnet-44444444', NOW() - INTERVAL '5 days', '{"name": "running", "code": 16}'),
(gen_random_uuid(), '123456789012', 'ap-southeast-1', 'arn:aws:ec2:ap-southeast-1:123456789012:instance/i-1234567890abcdef6', 'i-1234567890abcdef6', 't3.micro', '10.0.7.100', '54.123.45.72', 'vpc-44444444', 'subnet-55555555', NOW() - INTERVAL '3 days', '{"name": "running", "code": 16}'),
(gen_random_uuid(), '123456789012', 'ap-southeast-1', 'arn:aws:ec2:ap-southeast-1:123456789012:instance/i-1234567890abcdef7', 'i-1234567890abcdef7', 't3.small', '10.0.8.100', '54.123.45.73', 'vpc-44444444', 'subnet-66666666', NOW() - INTERVAL '2 days', '{"name": "running", "code": 16}'),
(gen_random_uuid(), '123456789012', 'sa-east-1', 'arn:aws:ec2:sa-east-1:123456789012:instance/i-1234567890abcdef8', 'i-1234567890abcdef8', 't3.medium', '10.0.9.100', NULL, 'vpc-55555555', 'subnet-77777777', NOW() - INTERVAL '1 day', '{"name": "stopped", "code": 80}'),
(gen_random_uuid(), '123456789012', 'sa-east-1', 'arn:aws:ec2:sa-east-1:123456789012:instance/i-1234567890abcdef9', 'i-1234567890abcdef9', 't3.large', '10.0.10.100', '54.123.45.74', 'vpc-55555555', 'subnet-88888888', NOW(), '{"name": "running", "code": 16}');

-- Insert dummy data for Azure VM instances
INSERT INTO azure_compute_virtual_machines (_cq_id, subscription_id, id, name, location, properties) VALUES
(gen_random_uuid(), 'subscription-12345678', '/subscriptions/subscription-12345678/resourceGroups/rg-web/providers/Microsoft.Compute/virtualMachines/vm-web-01', 'vm-web-01', 'eastus', '{"hardwareProfile": {"vmSize": "Standard_D2s_v3"}, "provisioningState": "Succeeded"}'),
(gen_random_uuid(), 'subscription-12345678', '/subscriptions/subscription-12345678/resourceGroups/rg-app/providers/Microsoft.Compute/virtualMachines/vm-app-01', 'vm-app-01', 'eastus2', '{"hardwareProfile": {"vmSize": "Standard_D4s_v3"}, "provisioningState": "Succeeded"}'),
(gen_random_uuid(), 'subscription-12345678', '/subscriptions/subscription-12345678/resourceGroups/rg-db/providers/Microsoft.Compute/virtualMachines/vm-db-01', 'vm-db-01', 'westus', '{"hardwareProfile": {"vmSize": "Standard_D8s_v3"}, "provisioningState": "Succeeded"}'),
(gen_random_uuid(), 'subscription-12345678', '/subscriptions/subscription-12345678/resourceGroups/rg-cache/providers/Microsoft.Compute/virtualMachines/vm-cache-01', 'vm-cache-01', 'westus2', '{"hardwareProfile": {"vmSize": "Standard_D16s_v3"}, "provisioningState": "Stopped"}'),
(gen_random_uuid(), 'subscription-12345678', '/subscriptions/subscription-12345678/resourceGroups/rg-monitoring/providers/Microsoft.Compute/virtualMachines/vm-monitoring-01', 'vm-monitoring-01', 'centralus', '{"hardwareProfile": {"vmSize": "Standard_E2s_v3"}, "provisioningState": "Succeeded"}'),
(gen_random_uuid(), 'subscription-12345678', '/subscriptions/subscription-12345678/resourceGroups/rg-lb/providers/Microsoft.Compute/virtualMachines/vm-lb-01', 'vm-lb-01', 'northcentralus', '{"hardwareProfile": {"vmSize": "Standard_E4s_v3"}, "provisioningState": "Succeeded"}'),
(gen_random_uuid(), 'subscription-12345678', '/subscriptions/subscription-12345678/resourceGroups/rg-api/providers/Microsoft.Compute/virtualMachines/vm-api-01', 'vm-api-01', 'southcentralus', '{"hardwareProfile": {"vmSize": "Standard_E8s_v3"}, "provisioningState": "Succeeded"}'),
(gen_random_uuid(), 'subscription-12345678', '/subscriptions/subscription-12345678/resourceGroups/rg-worker/providers/Microsoft.Compute/virtualMachines/vm-worker-01', 'vm-worker-01', 'westcentralus', '{"hardwareProfile": {"vmSize": "Standard_F2s_v2"}, "provisioningState": "Succeeded"}'),
(gen_random_uuid(), 'subscription-12345678', '/subscriptions/subscription-12345678/resourceGroups/rg-backup/providers/Microsoft.Compute/virtualMachines/vm-backup-01', 'vm-backup-01', 'canadacentral', '{"hardwareProfile": {"vmSize": "Standard_F4s_v2"}, "provisioningState": "Stopped"}'),
(gen_random_uuid(), 'subscription-12345678', '/subscriptions/subscription-12345678/resourceGroups/rg-test/providers/Microsoft.Compute/virtualMachines/vm-test-01', 'vm-test-01', 'canadaeast', '{"hardwareProfile": {"vmSize": "Standard_F8s_v2"}, "provisioningState": "Succeeded"}');

-- Insert dummy data for GCP Compute instances
INSERT INTO gcp_compute_instances (_cq_id, project_id, self_link, name, status, zone, machine_type) VALUES
(gen_random_uuid(), 'project-123456', 'https://www.googleapis.com/compute/v1/projects/project-123456/zones/us-central1-a/instances/gcp-web-01', 'gcp-web-01', 'RUNNING', 'us-central1-a', 'e2-micro'),
(gen_random_uuid(), 'project-123456', 'https://www.googleapis.com/compute/v1/projects/project-123456/zones/us-central1-b/instances/gcp-app-01', 'gcp-app-01', 'RUNNING', 'us-central1-b', 'e2-small'),
(gen_random_uuid(), 'project-123456', 'https://www.googleapis.com/compute/v1/projects/project-123456/zones/us-central1-c/instances/gcp-db-01', 'gcp-db-01', 'RUNNING', 'us-central1-c', 'e2-medium'),
(gen_random_uuid(), 'project-123456', 'https://www.googleapis.com/compute/v1/projects/project-123456/zones/us-central1-a/instances/gcp-cache-01', 'gcp-cache-01', 'STOPPED', 'us-central1-a', 'e2-standard-2'),
(gen_random_uuid(), 'project-123456', 'https://www.googleapis.com/compute/v1/projects/project-123456/zones/us-central1-b/instances/gcp-monitoring-01', 'gcp-monitoring-01', 'RUNNING', 'us-central1-b', 'e2-standard-4'),
(gen_random_uuid(), 'project-123456', 'https://www.googleapis.com/compute/v1/projects/project-123456/zones/us-central1-c/instances/gcp-lb-01', 'gcp-lb-01', 'RUNNING', 'us-central1-c', 'e2-standard-8'),
(gen_random_uuid(), 'project-123456', 'https://www.googleapis.com/compute/v1/projects/project-123456/zones/us-central1-a/instances/gcp-api-01', 'gcp-api-01', 'RUNNING', 'us-central1-a', 'n1-standard-1'),
(gen_random_uuid(), 'project-123456', 'https://www.googleapis.com/compute/v1/projects/project-123456/zones/us-central1-b/instances/gcp-worker-01', 'gcp-worker-01', 'RUNNING', 'us-central1-b', 'n1-standard-2'),
(gen_random_uuid(), 'project-123456', 'https://www.googleapis.com/compute/v1/projects/project-123456/zones/us-central1-c/instances/gcp-backup-01', 'gcp-backup-01', 'STOPPED', 'us-central1-c', 'n1-standard-4'),
(gen_random_uuid(), 'project-123456', 'https://www.googleapis.com/compute/v1/projects/project-123456/zones/us-central1-a/instances/gcp-test-01', 'gcp-test-01', 'RUNNING', 'us-central1-a', 'n1-standard-8');

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_aws_ec2_instances_account_id ON aws_ec2_instances(account_id);
CREATE INDEX IF NOT EXISTS idx_aws_ec2_instances_region ON aws_ec2_instances(region);
CREATE INDEX IF NOT EXISTS idx_aws_ec2_instances_instance_type ON aws_ec2_instances(instance_type);

CREATE INDEX IF NOT EXISTS idx_azure_vm_subscription_id ON azure_compute_virtual_machines(subscription_id);
CREATE INDEX IF NOT EXISTS idx_azure_vm_location ON azure_compute_virtual_machines(location);

CREATE INDEX IF NOT EXISTS idx_gcp_compute_project_id ON gcp_compute_instances(project_id);
CREATE INDEX IF NOT EXISTS idx_gcp_compute_zone ON gcp_compute_instances(zone);
CREATE INDEX IF NOT EXISTS idx_gcp_compute_machine_type ON gcp_compute_instances(machine_type);