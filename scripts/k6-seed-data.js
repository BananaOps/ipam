import http from 'k6/http';
import { check, sleep } from 'k6';

// Configuration
const BASE_URL = __ENV.BASE_URL || 'http://localhost:8081';
const API_URL = `${BASE_URL}/api/v1`;

// Sample data for seeding
const sampleSubnets = [
  // AWS Subnets
  {
    cidr: '10.0.1.0/24',
    name: 'AWS Production Web Tier',
    description: 'Production web servers in AWS us-east-1',
    location: 'AWS US East 1',
    location_type: 'CLOUD',
    cloud_info: {
      provider: 'aws',
      region: 'us-east-1',
      account_id: '123456789012'
    }
  },
  {
    cidr: '10.0.2.0/24',
    name: 'AWS Production App Tier',
    description: 'Productionlication servers in AWS us-east-1',
    location: 'AWS US East 1',
    location_type: 'CLOUD',
    cloud_info: {
      provider: 'aws',
      region: 'us-east-1',
      account_id: '123456789012'
    }
  },
  {
    cidr: '10.1.0.0/16',
    name: 'AWS Development VPC',
    description: 'Development environment in AWS eu-west-1',
    location: 'AWS EU West 1',
    location_type: 'CLOUD',
    cloud_info: {
      provider: 'aws',
      region: 'eu-west-1',
      account_id: '123456789012'
    }
  },

  // Azure Subnets
  {
    cidr: '10.2.1.0/24',
    name: 'Azure Production Frontend',
    description: 'Production frontend services in Azure West Europe',
    location: 'Azure West Europe',
    location_type: 'CLOUD',
    cloud_info: {
      provider: 'azure',
      region: 'westeurope',
      account_id: 'sub-azure-prod-001'
    }
  },
  {
    cidr: '10.2.2.0/24',
    name: 'Azure Production Backend',
    description: 'Production backend services in Azure West Europe',
    location: 'Azure West Europe',
    location_type: 'CLOUD',
    cloud_info: {
      provider: 'azure',
      region: 'westeurope',
      account_id: 'sub-azure-prod-001'
    }
  },
  {
    cidr: '10.3.0.0/16',
    name: 'Azure Staging Environment',
    description: 'Staging environment in Azure East US',
    location: 'Azure East US',
    location_type: 'CLOUD',
    cloud_info: {
      provider: 'azure',
      region: 'eastus',
      account_id: 'sub-azure-staging-001'
    }
  },

  // GCP Subnets
  {
    cidr: '10.4.1.0/24',
    name: 'GCP Production Compute',
    description: 'Production compute instances in GCP us-central1',
    location: 'GCP US Central 1',
    location_type: 'CLOUD',
    cloud_info: {
      provider: 'gcp',
      region: 'us-central1',
      account_id: 'project-gcp-prod-123'
    }
  },
  {
    cidr: '10.4.2.0/24',
    name: 'GCP Production Database',
    description: 'Production database subnet in GCP us-central1',
    location: 'GCP US Central 1',
    location_type: 'CLOUD',
    cloud_info: {
      provider: 'gcp',
      region: 'us-central1',
      account_id: 'project-gcp-prod-123'
    }
  },
  {
    cidr: '10.5.0.0/16',
    name: 'GCP Development',
    description: 'Development environment in GCP europe-west1',
    location: 'GCP Europe West 1',
    location_type: 'CLOUD',
    cloud_info: {
      provider: 'gcp',
      region: 'europe-west1',
      account_id: 'project-gcp-dev-456'
    }
  },

  // Scaleway Subnets
  {
    cidr: '10.6.1.0/24',
    name: 'Scaleway Production API',
    description: 'Production API servers in Scaleway Paris',
    location: 'Scaleway Paris',
    location_type: 'CLOUD',
    cloud_info: {
      provider: 'scaleway',
      region: 'fr-par-1',
      account_id: 'scw-prod-789'
    }
  },
  {
    cidr: '10.6.2.0/24',
    name: 'Scaleway Production Storage',
    description: 'Production storage network in Scaleway Paris',
    location: 'Scaleway Paris',
    location_type: 'CLOUD',
    cloud_info: {
      provider: 'scaleway',
      region: 'fr-par-1',
      account_id: 'scw-prod-789'
    }
  },

  // OVH Subnets
  {
    cidr: '10.7.1.0/24',
    name: 'OVH Production Web',
    description: 'Production web servers in OVH Gravelines',
    location: 'OVH Gravelines',
    location_type: 'CLOUD',
    cloud_info: {
      provider: 'ovh',
      region: 'gra',
      account_id: 'ovh-prod-abc'
    }
  },

  // On-premise Subnets
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
];

export let options = {
  stages: [
    { duration: '10s', target: 1 }, // Single user for seeding
  ],
};

export default function () {
  console.log('🌱 Starting IPAM data seeding...');
  
  // Check if API is available
  let healthCheck = http.get(`${API_URL}/subnets`);
  if (!check(healthCheck, {
    'API is available': (r) => r.status === 200,
  })) {
    console.error('❌ API is not available. Make sure the backend is running.');
    return;
  }

  console.log('✅ API is available, starting to seed data...');

  let successCount = 0;
  let errorCount = 0;

  // Create each subnet
  sampleSubnets.forEach((subnet, index) => {
    console.log(`📡 Creating subnet ${index + 1}/${sampleSubnets.length}: ${subnet.name}`);
    
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
      'subnet created successfully': (r) => r.status === 201,
    })) {
      successCount++;
      console.log(`✅ Created: ${subnet.name} (${subnet.cidr})`);
    } else {
      errorCount++;
      console.log(`❌ Failed to create: ${subnet.name} - Status: ${response.status}`);
      if (response.body) {
        console.log(`   Error: ${response.body}`);
      }
    }

    // Small delay between requests
    sleep(0.1);
  });

  console.log('\n📊 Seeding Summary:');
  console.log(`✅ Successfully created: ${successCount} subnets`);
  console.log(`❌ Failed to create: ${errorCount} subnets`);
  console.log(`📈 Total processed: ${sampleSubnets.length} subnets`);
  
  if (successCount > 0) {
    console.log('\n🎉 Data seeding completed! You can now view the subnets in your IPAM interface.');
  }
}
