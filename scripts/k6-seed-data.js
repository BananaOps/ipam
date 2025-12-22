import http from 'k6/http';
import { check, sleep } from 'k6';

// Configuration
const BASE_URL = __ENV.BASE_URL || 'http://localhost:8082'; // Updated port
const API_URL = `${BASE_URL}/api/v1`;

// Sample data for seeding with VPC/Subnet relationships
const sampleData = {
  // VPCs first (parent networks)
  vpcs: [
    {
      cidr: '10.0.0.0/16',
      name: 'Production VPC',
      description: 'Main production VPC in AWS us-east-1',
      location: 'us-east-1',
      location_type: 'CLOUD',
      cloud_info: {
        provider: 'aws',
        region: 'us-east-1',
        account_id: '123456789012',
        resource_type: 'vpc',
        vpc_id: 'vpc-prod-12345678'
      }
    },
    {
      cidr: '10.1.0.0/16',
      name: 'Development VPC',
      description: 'Development VPC in AWS eu-west-1',
      location: 'eu-west-1',
      location_type: 'CLOUD',
      cloud_info: {
        provider: 'aws',
        region: 'eu-west-1',
        account_id: '123456789012',
        resource_type: 'vpc',
        vpc_id: 'vpc-dev-87654321'
      }
    },
    {
      cidr: '10.2.0.0/16',
      name: 'Staging VPC',
      description: 'Staging environment VPC in Azure West Europe',
      location: 'westeurope',
      location_type: 'CLOUD',
      cloud_info: {
        provider: 'azure',
        region: 'westeurope',
        account_id: 'sub-azure-staging-001',
        resource_type: 'vpc',
        vpc_id: 'vnet-staging-abcdef'
      }
    },
    {
      cidr: '10.3.0.0/16',
      name: 'Analytics VPC',
      description: 'Analytics and ML workloads VPC in GCP',
      location: 'us-central1',
      location_type: 'CLOUD',
      cloud_info: {
        provider: 'gcp',
        region: 'us-central1',
        account_id: 'project-analytics-123',
        resource_type: 'vpc',
        vpc_id: 'vpc-analytics-xyz789'
      }
    }
  ],

  // Subnets (child networks) - will be linked to VPCs after creation
  subnets: [
    // Production VPC Subnets (10.0.0.0/16)
    {
      cidr: '10.0.1.0/24',
      name: 'Production Web Tier',
      description: 'Public subnet for web servers',
      location: 'us-east-1a',
      location_type: 'CLOUD',
      parent_vpc_cidr: '10.0.0.0/16',
      cloud_info: {
        provider: 'aws',
        region: 'us-east-1',
        account_id: '123456789012',
        resource_type: 'subnet',
        vpc_id: 'vpc-prod-12345678',
        subnet_id: 'subnet-web-11111111'
      }
    },
    {
      cidr: '10.0.2.0/24',
      name: 'Production App Tier',
      description: 'Private subnet for application servers',
      location: 'us-east-1a',
      location_type: 'CLOUD',
      parent_vpc_cidr: '10.0.0.0/16',
      cloud_info: {
        provider: 'aws',
        region: 'us-east-1',
        account_id: '123456789012',
        resource_type: 'subnet',
        vpc_id: 'vpc-prod-12345678',
        subnet_id: 'subnet-app-22222222'
      }
    },
    {
      cidr: '10.0.3.0/24',
      name: 'Production DB Tier',
      description: 'Private subnet for databases',
      location: 'us-east-1b',
      location_type: 'CLOUD',
      parent_vpc_cidr: '10.0.0.0/16',
      cloud_info: {
        provider: 'aws',
        region: 'us-east-1',
        account_id: '123456789012',
        resource_type: 'subnet',
        vpc_id: 'vpc-prod-12345678',
        subnet_id: 'subnet-db-33333333'
      }
    },
    {
      cidr: '10.0.10.0/24',
      name: 'Production Management',
      description: 'Management and monitoring subnet',
      location: 'us-east-1c',
      location_type: 'CLOUD',
      parent_vpc_cidr: '10.0.0.0/16',
      cloud_info: {
        provider: 'aws',
        region: 'us-east-1',
        account_id: '123456789012',
        resource_type: 'subnet',
        vpc_id: 'vpc-prod-12345678',
        subnet_id: 'subnet-mgmt-44444444'
      }
    },

    // Development VPC Subnets (10.1.0.0/16)
    {
      cidr: '10.1.1.0/24',
      name: 'Development Web',
      description: 'Development web servers',
      location: 'eu-west-1a',
      location_type: 'CLOUD',
      parent_vpc_cidr: '10.1.0.0/16',
      cloud_info: {
        provider: 'aws',
        region: 'eu-west-1',
        account_id: '123456789012',
        resource_type: 'subnet',
        vpc_id: 'vpc-dev-87654321',
        subnet_id: 'subnet-dev-web-55555555'
      }
    },
    {
      cidr: '10.1.2.0/24',
      name: 'Development API',
      description: 'Development API services',
      location: 'eu-west-1b',
      location_type: 'CLOUD',
      parent_vpc_cidr: '10.1.0.0/16',
      cloud_info: {
        provider: 'aws',
        region: 'eu-west-1',
        account_id: '123456789012',
        resource_type: 'subnet',
        vpc_id: 'vpc-dev-87654321',
        subnet_id: 'subnet-dev-api-66666666'
      }
    },

    // Staging VPC Subnets (10.2.0.0/16) - Azure
    {
      cidr: '10.2.1.0/24',
      name: 'Staging Frontend',
      description: 'Staging frontend services',
      location: 'westeurope',
      location_type: 'CLOUD',
      parent_vpc_cidr: '10.2.0.0/16',
      cloud_info: {
        provider: 'azure',
        region: 'westeurope',
        account_id: 'sub-azure-staging-001',
        resource_type: 'subnet',
        vpc_id: 'vnet-staging-abcdef',
        subnet_id: 'subnet-staging-frontend'
      }
    },
    {
      cidr: '10.2.2.0/24',
      name: 'Staging Backend',
      description: 'Staging backend services',
      location: 'westeurope',
      location_type: 'CLOUD',
      parent_vpc_cidr: '10.2.0.0/16',
      cloud_info: {
        provider: 'azure',
        region: 'westeurope',
        account_id: 'sub-azure-staging-001',
        resource_type: 'subnet',
        vpc_id: 'vnet-staging-abcdef',
        subnet_id: 'subnet-staging-backend'
      }
    },

    // Analytics VPC Subnets (10.3.0.0/16) - GCP
    {
      cidr: '10.3.1.0/24',
      name: 'Analytics Compute',
      description: 'Analytics compute instances',
      location: 'us-central1-a',
      location_type: 'CLOUD',
      parent_vpc_cidr: '10.3.0.0/16',
      cloud_info: {
        provider: 'gcp',
        region: 'us-central1',
        account_id: 'project-analytics-123',
        resource_type: 'subnet',
        vpc_id: 'vpc-analytics-xyz789',
        subnet_id: 'subnet-analytics-compute'
      }
    },
    {
      cidr: '10.3.2.0/24',
      name: 'Analytics Storage',
      description: 'Analytics data storage network',
      location: 'us-central1-b',
      location_type: 'CLOUD',
      parent_vpc_cidr: '10.3.0.0/16',
      cloud_info: {
        provider: 'gcp',
        region: 'us-central1',
        account_id: 'project-analytics-123',
        resource_type: 'subnet',
        vpc_id: 'vpc-analytics-xyz789',
        subnet_id: 'subnet-analytics-storage'
      }
    }
  ],

  // Additional cloud subnets without VPC relationships
  standaloneSubnets: [
    // Scaleway
    {
      cidr: '10.4.1.0/24',
      name: 'Scaleway Production API',
      description: 'Production API servers in Scaleway Paris',
      location: 'fr-par-1',
      location_type: 'CLOUD',
      cloud_info: {
        provider: 'scaleway',
        region: 'fr-par-1',
        account_id: 'scw-prod-789',
        resource_type: 'subnet',
        subnet_id: 'scw-subnet-api-001'
      }
    },
    {
      cidr: '10.4.2.0/24',
      name: 'Scaleway Production Storage',
      description: 'Production storage network in Scaleway Paris',
      location: 'fr-par-1',
      location_type: 'CLOUD',
      cloud_info: {
        provider: 'scaleway',
        region: 'fr-par-1',
        account_id: 'scw-prod-789',
        resource_type: 'subnet',
        subnet_id: 'scw-subnet-storage-001'
      }
    },

    // OVH
    {
      cidr: '10.5.1.0/24',
      name: 'OVH Production Web',
      description: 'Production web servers in OVH Gravelines',
      location: 'gra',
      location_type: 'CLOUD',
      cloud_info: {
        provider: 'ovh',
        region: 'gra',
        account_id: 'ovh-prod-abc',
        resource_type: 'subnet',
        subnet_id: 'ovh-subnet-web-001'
      }
    }
  ],

  // On-premise networks
  onPremiseNetworks: [
    {
      cidr: '192.168.1.0/24',
      name: 'Paris DC1 Management',
      description: 'Management network for Paris datacenter',
      location: 'Paris DC1',
      location_type: 'DATACENTER'
    },
    {
      cidr: '192.168.2.0/24',
      name: 'Paris DC1 Production',
      description: 'Production servers in Paris datacenter',
      location: 'Paris DC1',
      location_type: 'DATACENTER'
    },
    {
      cidr: '192.168.10.0/24',
      name: 'London DC1 Management',
      description: 'Management network for London datacenter',
      location: 'London DC1',
      location_type: 'DATACENTER'
    },
    {
      cidr: '192.168.11.0/24',
      name: 'London DC1 Production',
      description: 'Production servers in London datacenter',
      location: 'London DC1',
      location_type: 'DATACENTER'
    },

    // Site Networks
    {
      cidr: '172.16.1.0/24',
      name: 'New York Office',
      description: 'New York office network',
      location: 'New York Office',
      location_type: 'SITE'
    },
    {
      cidr: '172.16.2.0/24',
      name: 'San Francisco Office',
      description: 'San Francisco office network',
      location: 'San Francisco Office',
      location_type: 'SITE'
    },
    {
      cidr: '172.16.3.0/24',
      name: 'Tokyo Office',
      description: 'Tokyo office network',
      location: 'Tokyo Office',
      location_type: 'SITE'
    },
    {
      cidr: '172.16.4.0/24',
      name: 'Berlin Office',
      description: 'Berlin office network',
      location: 'Berlin Office',
      location_type: 'SITE'
    }
  ]
};

export let options = {
  stages: [
    { duration: '10s', target: 1 }, // Single user for seeding
  ],
};

export default function () {
  console.log('🌱 Starting IPAM data seeding with VPC/Subnet relationships...');
  
  // Check if API is available
  let healthCheck = http.get(`${API_URL}/subnets`);
  if (!check(healthCheck, {
    'API is available': (r) => r.status === 200,
  })) {
    console.error('❌ API is not available. Make sure the backend is running.');
    return;
  }

  console.log('✅ API is available, starting to seed data...');

  // Optional: Delete existing subnets first (if you want to start fresh)
  const CLEAN_FIRST = __ENV.CLEAN_FIRST === 'true';
  if (CLEAN_FIRST) {
    console.log('🧹 Cleaning existing data...');
    let existingSubnets = http.get(`${API_URL}/subnets`);
    if (existingSubnets.status === 200) {
      try {
        let data = JSON.parse(existingSubnets.body);
        if (data.subnets && data.subnets.length > 0) {
          console.log(`   Found ${data.subnets.length} existing subnets, deleting...`);
          data.subnets.forEach((subnet) => {
            http.del(`${API_URL}/subnets/${subnet.id}`);
          });
          console.log('   ✅ Existing data cleaned');
          sleep(1); // Wait a bit for deletions to complete
        }
      } catch (e) {
        console.log('   ⚠️  Could not parse existing subnets');
      }
    }
  }

  let successCount = 0;
  let errorCount = 0;
  let vpcIdMap = {}; // Map to store VPC CIDR -> created subnet ID

  // Step 1: Create VPCs first
  console.log('\n🏗️  Creating VPCs...');
  sampleData.vpcs.forEach((vpc, index) => {
    console.log(`📡 Creating VPC ${index + 1}/${sampleData.vpcs.length}: ${vpc.name}`);
    
    let response = http.post(
      `${API_URL}/subnets`,
      JSON.stringify(vpc),
      {
        headers: {
          'Content-Type': 'application/json',
        },
      }
    );

    if (check(response, {
      'VPC created successfully': (r) => r.status === 201,
    })) {
      successCount++;
      console.log(`✅ Created VPC: ${vpc.name} (${vpc.cidr})`);
      
      // Store VPC ID for later subnet linking
      try {
        let responseData = JSON.parse(response.body);
        vpcIdMap[vpc.cidr] = responseData.id;
      } catch (e) {
        console.log(`⚠️  Could not parse response for VPC ${vpc.name}`);
      }
    } else if (response.status === 500 && response.body && response.body.includes('UNIQUE constraint failed')) {
      // Handle duplicate CIDR - try to find existing subnet
      console.log(`⚠️  VPC ${vpc.name} (${vpc.cidr}) already exists, trying to find it...`);
      let existingResponse = http.get(`${API_URL}/subnets`);
      if (existingResponse.status === 200) {
        try {
          let data = JSON.parse(existingResponse.body);
          let existingVpc = data.subnets.find(s => s.cidr === vpc.cidr);
          if (existingVpc) {
            vpcIdMap[vpc.cidr] = existingVpc.id;
            console.log(`   ✅ Found existing VPC with ID: ${existingVpc.id}`);
            successCount++; // Count as success since we found it
          }
        } catch (e) {
          console.log(`   ❌ Could not find existing VPC: ${e.message}`);
          errorCount++;
        }
      } else {
        errorCount++;
      }
    } else {
      errorCount++;
      console.log(`❌ Failed to create VPC: ${vpc.name} - Status: ${response.status}`);
      console.log(`   Request: ${JSON.stringify(vpc)}`);
      if (response.body) {
        console.log(`   Response: ${response.body}`);
      }
    }

    sleep(0.1);
  });

  // Step 2: Create subnets with parent relationships
  console.log('\n🔗 Creating subnets with VPC relationships...');
  sampleData.subnets.forEach((subnet, index) => {
    console.log(`📡 Creating subnet ${index + 1}/${sampleData.subnets.length}: ${subnet.name}`);
    
    // Create a copy of the subnet data
    let subnetData = { ...subnet };
    
    // Add parent_id if we have the VPC ID
    if (subnet.parent_vpc_cidr && vpcIdMap[subnet.parent_vpc_cidr]) {
      subnetData.parent_id = vpcIdMap[subnet.parent_vpc_cidr];
      console.log(`   🔗 Linking to VPC: ${subnet.parent_vpc_cidr} (ID: ${subnetData.parent_id})`);
    }
    
    // Remove the temporary field
    delete subnetData.parent_vpc_cidr;
    
    let response = http.post(
      `${API_URL}/subnets`,
      JSON.stringify(subnetData),
      {
        headers: {
          'Content-Type': 'application/json',
        },
      }
    );

    if (check(response, {
      'subnet created successfully': (r) => r.status === 201,
    })) {
      successCount++;
      console.log(`✅ Created subnet: ${subnet.name} (${subnet.cidr})`);
    } else if (response.status === 500 && response.body && response.body.includes('UNIQUE constraint failed')) {
      // Handle duplicate CIDR
      console.log(`⚠️  Subnet ${subnet.name} (${subnet.cidr}) already exists, skipping...`);
      successCount++; // Count as success since it exists
    } else {
      errorCount++;
      console.log(`❌ Failed to create subnet: ${subnet.name} - Status: ${response.status}`);
      if (response.body) {
        console.log(`   Error: ${response.body}`);
      }
    }

    sleep(0.1);
  });

  // Step 3: Create standalone cloud subnets
  console.log('\n☁️  Creating standalone cloud subnets...');
  sampleData.standaloneSubnets.forEach((subnet, index) => {
    console.log(`📡 Creating standalone subnet ${index + 1}/${sampleData.standaloneSubnets.length}: ${subnet.name}`);
    
    let response = http.post(
      `${API_URL}/subnets`,
      JSON.stringify(subnet),
      {
        headers: {
          'Content-Type': 'application/json',
        },
      }
    );

    if (check(response, {
      'standalone subnet created successfully': (r) => r.status === 201,
    })) {
      successCount++;
      console.log(`✅ Created standalone subnet: ${subnet.name} (${subnet.cidr})`);
    } else if (response.status === 500 && response.body && response.body.includes('UNIQUE constraint failed')) {
      // Handle duplicate CIDR
      console.log(`⚠️  Standalone subnet ${subnet.name} (${subnet.cidr}) already exists, skipping...`);
      successCount++; // Count as success since it exists
    } else {
      errorCount++;
      console.log(`❌ Failed to create standalone subnet: ${subnet.name} - Status: ${response.status}`);
      if (response.body) {
        console.log(`   Error: ${response.body}`);
      }
    }

    sleep(0.1);
  });

  // Step 4: Create on-premise networks
  console.log('\n🏢 Creating on-premise networks...');
  sampleData.onPremiseNetworks.forEach((network, index) => {
    console.log(`📡 Creating on-premise network ${index + 1}/${sampleData.onPremiseNetworks.length}: ${network.name}`);
    
    let response = http.post(
      `${API_URL}/subnets`,
      JSON.stringify(network),
      {
        headers: {
          'Content-Type': 'application/json',
        },
      }
    );

    if (check(response, {
      'on-premise network created successfully': (r) => r.status === 201,
    })) {
      successCount++;
      console.log(`✅ Created on-premise network: ${network.name} (${network.cidr})`);
    } else if (response.status === 500 && response.body && response.body.includes('UNIQUE constraint failed')) {
      // Handle duplicate CIDR
      console.log(`⚠️  On-premise network ${network.name} (${network.cidr}) already exists, skipping...`);
      successCount++; // Count as success since it exists
    } else {
      errorCount++;
      console.log(`❌ Failed to create on-premise network: ${network.name} - Status: ${response.status}`);
      if (response.body) {
        console.log(`   Error: ${response.body}`);
      }
    }

    sleep(0.1);
  });

  // Summary
  let totalNetworks = sampleData.vpcs.length + sampleData.subnets.length + 
                     sampleData.standaloneSubnets.length + sampleData.onPremiseNetworks.length;

  console.log('\n📊 Seeding Summary:');
  console.log(`🏗️  VPCs created: ${sampleData.vpcs.length}`);
  console.log(`🔗 Subnets with VPC relationships: ${sampleData.subnets.length}`);
  console.log(`☁️  Standalone cloud subnets: ${sampleData.standaloneSubnets.length}`);
  console.log(`🏢 On-premise networks: ${sampleData.onPremiseNetworks.length}`);
  console.log(`✅ Successfully created: ${successCount} networks`);
  console.log(`❌ Failed to create: ${errorCount} networks`);
  console.log(`📈 Total processed: ${totalNetworks} networks`);
  
  if (successCount > 0) {
    console.log('\n🎉 Data seeding completed!');
    console.log('📋 You can now view:');
    console.log('   • VPCs with VPC badges');
    console.log('   • Subnets with SUBNET badges');
    console.log('   • Parent-child relationships');
    console.log('   • Multi-cloud provider support');
    console.log('   • Click on VPCs to see their child subnets');
  }

  // Test the children endpoint for one of the VPCs
  if (Object.keys(vpcIdMap).length > 0) {
    console.log('\n🧪 Testing children endpoint...');
    let firstVpcId = Object.values(vpcIdMap)[0];
    let childrenResponse = http.get(`${API_URL}/subnets/${firstVpcId}/children`);
    
    if (check(childrenResponse, {
      'children endpoint works': (r) => r.status === 200,
    })) {
      try {
        let childrenData = JSON.parse(childrenResponse.body);
        console.log(`✅ Children endpoint working: Found ${childrenData.count} child subnets`);
      } catch (e) {
        console.log('✅ Children endpoint responding (could not parse JSON)');
      }
    } else {
      console.log(`⚠️  Children endpoint test failed: Status ${childrenResponse.status}`);
    }
  }
}
